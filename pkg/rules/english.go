package rules

import (
	"fmt"
	"go/ast"
	"unicode"

	"github.com/hel1th/loglinter/pkg/loggers"
	"golang.org/x/tools/go/analysis"
)

type EnglishOnlyRule struct{}

func (r *EnglishOnlyRule) Name() string {
	return "english-only"
}

func (r *EnglishOnlyRule) Message() string {
	return "log message should contain only English characters"
}

func (r *EnglishOnlyRule) Check(pass *analysis.Pass, logCall loggers.LogCall) []analysis.Diagnostic {
	message, ok := loggers.ExtractStringLit(logCall.Message)
	if !ok {
		return nil
	}

	nonEnglishChars := r.findNonEnglishChars(message)

	if len(nonEnglishChars) > 0 {
		return []analysis.Diagnostic{
			{
				Pos:      logCall.Message.Pos(),
				End:      logCall.Message.End(),
				Message:  fmt.Sprintf("%s (found: %s)", r.Message(), formatNonEnglishChars(nonEnglishChars)),
				Category: r.Name(),
			},
		}
	}

	return nil
}

func (r *EnglishOnlyRule) findNonEnglishChars(message string) []rune {
	var nonEnglish []rune
	seen := make(map[rune]bool)

	for _, char := range message {
		if unicode.IsSpace(char) || unicode.IsPunct(char) || unicode.IsDigit(char) {
			continue
		}

		if unicode.IsLetter(char) && !isLatinLetter(char) {
			if !seen[char] {
				nonEnglish = append(nonEnglish, char)
				seen[char] = true
			}
		}
	}

	return nonEnglish
}

func isLatinLetter(r rune) bool {
	if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
		return true
	}

	return false
}

func formatNonEnglishChars(chars []rune) string {
	if len(chars) == 0 {
		return ""
	}

	result := ""
	for i, char := range chars {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("'%c'", char)

		if i >= 2 {
			if len(chars) > 3 {
				result += fmt.Sprintf(" and %d more", len(chars)-3)
			}
			break
		}
	}

	return result
}

func (r *EnglishOnlyRule) CheckExpr(expr ast.Expr) (bool, string) {
	message, ok := loggers.ExtractStringLit(expr)
	if !ok || len(message) == 0 {
		return true, ""
	}

	nonEnglishChars := r.findNonEnglishChars(message)

	if len(nonEnglishChars) > 0 {
		return false, fmt.Sprintf("contains non-English characters: %s", formatNonEnglishChars(nonEnglishChars))
	}

	return true, ""
}
