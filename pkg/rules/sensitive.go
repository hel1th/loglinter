package rules

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/hel1th/loglinter/pkg/loggers"
	"golang.org/x/tools/go/analysis"
)

type SensitiveDataRule struct {
	// Можно добавить кастомные паттерны
	customPatterns []string
}

func (r *SensitiveDataRule) Name() string {
	return "no-sensitive-data"
}

func (r *SensitiveDataRule) Message() string {
	return "log message may contain sensitive data"
}

func (r *SensitiveDataRule) Check(pass *analysis.Pass, logCall loggers.LogCall) []analysis.Diagnostic {
	violations := r.analyzeMessageExpression(logCall.Message)

	if len(violations) > 0 {
		return []analysis.Diagnostic{
			{
				Pos:      logCall.Message.Pos(),
				End:      logCall.Message.End(),
				Message:  fmt.Sprintf("%s: %s", r.Message(), strings.Join(violations, ", ")),
				Category: r.Name(),
			},
		}
	}

	return nil
}

func (r *SensitiveDataRule) analyzeMessageExpression(expr ast.Expr) []string {
	var violations []string

	switch v := expr.(type) {
	case *ast.BasicLit:
		if message, ok := loggers.ExtractStringLit(expr); ok {
			if keywords := r.findSensitiveKeywords(message); len(keywords) > 0 {
				violations = append(violations, keywords...)
			}
		}

	case *ast.BinaryExpr:
		violations = append(violations, r.analyzeBinaryExpr(v)...)

	case *ast.CallExpr:
		violations = append(violations, r.analyzeCallExpr(v)...)
	}

	return violations
}

func (r *SensitiveDataRule) analyzeBinaryExpr(expr *ast.BinaryExpr) []string {
	var violations []string

	if message, ok := loggers.ExtractStringLit(expr.X); ok {
		if keywords := r.findSensitiveKeywords(message); len(keywords) > 0 {

			violations = append(violations, fmt.Sprintf("concatenating sensitive field: %s", strings.Join(keywords, ", ")))
		}
	}

	violations = append(violations, r.analyzeMessageExpression(expr.Y)...)

	return violations
}

func (r *SensitiveDataRule) analyzeCallExpr(expr *ast.CallExpr) []string {
	var violations []string

	if len(expr.Args) > 0 {
		if message, ok := loggers.ExtractStringLit(expr.Args[0]); ok {
			if keywords := r.findSensitiveKeywords(message); len(keywords) > 0 {
				violations = append(violations, keywords...)
			}
		}
	}

	return violations
}

func (r *SensitiveDataRule) findSensitiveKeywords(message string) []string {
	var found []string
	messageLower := strings.ToLower(message)
	sensitiveKeywords := []string{
		"password", "passwd", "pwd",
		"token", "access_token", "refresh_token", "bearer",
		"api_key", "apikey", "api-key",
		"secret", "client_secret",
		"private_key", "private key", "privatekey",
		"credit_card", "card_number", "cvv", "cvc",
		"ssn", "social security",
		"auth", "authorization",
		"credentials", "credential",
	}

	sensitiveKeywords = append(sensitiveKeywords, r.customPatterns...)

	for _, keyword := range sensitiveKeywords {
		if r.containsSensitivePattern(messageLower, keyword) {
			found = append(found, keyword)
		}
	}

	return found
}

func (r *SensitiveDataRule) containsSensitivePattern(message, keyword string) bool {
	patterns := []string{
		keyword + ":",
		keyword + " :",
		keyword + "=",
		keyword + " =",
		keyword + " is",
	}

	for _, pattern := range patterns {
		if strings.Contains(message, pattern) {
			return true
		}
	}

	words := strings.Fields(message)
	for _, word := range words {
		cleaned := strings.Trim(word, ".,;:!?\"'")
		if cleaned == keyword {
			return true
		}
	}

	return false
}

func (r *SensitiveDataRule) CheckExpr(expr ast.Expr) (bool, string) {
	violations := r.analyzeMessageExpression(expr)

	if len(violations) > 0 {
		return false, strings.Join(violations, ", ")
	}

	return true, ""
}

func (r *SensitiveDataRule) SetCustomPatterns(patterns []string) {
	r.customPatterns = patterns
}

func (r *SensitiveDataRule) GetDefaultPatterns() []string {
	return []string{
		"password", "passwd", "pwd",
		"token", "access_token", "refresh_token",
		"api_key", "apikey",
		"secret", "client_secret",
		"private_key", "privatekey",
		"credit_card", "cvv",
		"credentials",
	}
}
