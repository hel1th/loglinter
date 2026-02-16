package rules

import (
	"fmt"
	"go/ast"
	"unicode"

	"github.com/hel1th/loglinter/pkg/loggers"
	"golang.org/x/tools/go/analysis"
)

type LowercaseRule struct{}

func (r *LowercaseRule) Check(pass *analysis.Pass, logCall loggers.LogCall) []analysis.Diagnostic {
	message, ok := loggers.ExtractStringLit(logCall.Message)
	if !ok {
		return nil
	}

	if len(message) == 0 {
		return nil
	}

	firstCharIndex := 0
	for firstCharIndex < len(message) && unicode.IsSpace(rune(message[firstCharIndex])) {
		firstCharIndex++
	}

	if firstCharIndex >= len(message) {
		return nil
	}

	firstChar := rune(message[firstCharIndex])

	if !unicode.IsLetter(firstChar) {
		return nil
	}

	if unicode.IsUpper(firstChar) {
		suggested := string(unicode.ToLower(firstChar)) + message[firstCharIndex+1:]
		logCallPos := logCall.Call.Pos()
		logCallEnd := logCall.Call.End()
		diag := []analysis.Diagnostic{{
			Pos:      logCallPos,
			End:      logCallEnd,
			Message:  r.Message(),
			Category: r.Name(),
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "Convert first letter to lowercase",
					TextEdits: []analysis.TextEdit{{
						Pos:     logCallPos,
						End:     logCallEnd,
						NewText: []byte(fmt.Sprintf(`"%s"`, suggested)),
					}},
				},
			},
		}}
		return diag
	}
	return nil
}

func (r *LowercaseRule) Name() string {
	return "lowercase-start"
}

func (r *LowercaseRule) Message() string {
	return "log message should start with a lowercase letter"
}

func (r *LowercaseRule) CheckExpr(expr ast.Expr) (bool, string) {
	message, ok := loggers.ExtractStringLit(expr)
	if !ok || len(message) == 0 {
		return true, ""
	}

	firstNonSpaceIdx := 0

	for firstNonSpaceIdx < len(message) && unicode.IsSpace(rune(message[firstNonSpaceIdx])) {
		firstNonSpaceIdx++
	}

	if firstNonSpaceIdx >= len(message) {
		return true, ""
	}
	firstNonSpace := rune(message[firstNonSpaceIdx])

	if unicode.IsUpper(firstNonSpace) {
		return false, fmt.Sprintf("message starts with uppercase letter: '%c'", firstNonSpace)
	}

	return true, ""
}
