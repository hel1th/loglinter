package loggers

import (
	"go/ast"
	"go/token"
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

	ast.Inspect(file, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		logCall := d.analyzeCallExpr(callExpr)
		if logCall == nil {
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
	if _, ok := logMethods[method]; !ok {
		return false
	}

	return true
}

func (d *Detector) loggerType(expr ast.Expr) LoggerType {
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
			return str[1 : len(str)-1], true
		}

	case *ast.BinaryExpr:
		if left, ok := ExtractStringLit(v.X); ok {
			return left, true
		}
	}

	return "", false
}
