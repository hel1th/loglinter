package loggers

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

type Detector struct {
	pass *analysis.Pass
}

func NewDetector(pass *analysis.Pass) *Detector {
	return &Detector{
		pass: pass,
	}
}

type LogCall struct {
	Call    *ast.CallExpr
	Message ast.Expr
	Logger  LoggerType
	Method  string
}

func (d *Detector) DetectLogCalls(file *ast.File) []LogCall {
	var logCalls []LogCall

	if d.pass == nil || d.pass.TypesInfo == nil {
		return logCalls
	}

	ast.Inspect(file, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		logCall := d.analyzeCallExpr(callExpr)
		if logCall != nil {
			logCalls = append(logCalls, *logCall)
		}

		return true
	})

	return logCalls
}

func (d *Detector) analyzeCallExpr(call *ast.CallExpr) *LogCall {
	selectorExp, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	method := selectorExp.Sel.Name

	if !isLogMethod(method) {
		return nil
	}

	loggerType := d.loggerType(selectorExp.X)
	if loggerType == UnknownLogger {
		return nil
	}

	if len(call.Args) == 0 {
		return nil
	}

	return &LogCall{
		Call:    call,
		Message: call.Args[0],
		Logger:  loggerType,
		Method:  method,
	}
}

func isLogMethod(method string) bool {
	_, ok := logMethods[method]
	return ok
}

func (d *Detector) loggerType(expr ast.Expr) LoggerType {
	if d.pass == nil || d.pass.TypesInfo == nil {
		return UnknownLogger
	}

	if ident, ok := expr.(*ast.Ident); ok {
		return d.identifyByIdent(ident)
	}

	return d.identifyByType(expr)
}

func (d *Detector) identifyByIdent(ident *ast.Ident) LoggerType {
	obj := d.pass.TypesInfo.ObjectOf(ident)
	if obj == nil {
		return UnknownLogger
	}

	if pkgName, ok := obj.(*types.PkgName); ok {
		pkgPath := pkgName.Imported().Path()

		switch pkgPath {
		case "log/slog":
			return SlogLogger
		case "go.uber.org/zap":
			return ZapLogger
		case "log":
			return LogLogger
		default:
			return UnknownLogger
		}
	}

	return d.identifyByType(ident)
}

func (d *Detector) identifyByType(expr ast.Expr) LoggerType {
	logType := d.pass.TypesInfo.TypeOf(expr)
	if logType == nil {
		return UnknownLogger
	}

	typStr := logType.String()

	for _, check := range loggerChecks {
		if strings.Contains(typStr, check.pattern) {
			return check.logType
		}
	}

	return UnknownLogger
}

func ExtractStringLit(expr ast.Expr) (string, bool) {
	switch v := expr.(type) {
	case *ast.BasicLit:
		if v.Kind == token.STRING {
			str := v.Value
			if len(str) >= 2 {
				return str[1 : len(str)-1], true
			}
		}

	case *ast.BinaryExpr:
		if left, ok := ExtractStringLit(v.X); ok {
			return left, true
		}
	}

	return "", false
}
