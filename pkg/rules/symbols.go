package rules

import (
	"fmt"
	"go/ast"
	"strings"
	"unicode"

	"github.com/hel1th/loglinter/pkg/loggers"
	"golang.org/x/tools/go/analysis"
)

// WHITELIST: разрешены только буквы цифры пробелы дефис и нижнее подчеркивание
type NoSpecialSymbolsRule struct{}

func (r *NoSpecialSymbolsRule) Name() string {
	return "no-special-symbols"
}

func (r *NoSpecialSymbolsRule) Message() string {
	return "log message should contain only letters, digits, spaces, hyphens and underscores"
}

func (r *NoSpecialSymbolsRule) Check(pass *analysis.Pass, logCall loggers.LogCall) []analysis.Diagnostic {
	message, ok := loggers.ExtractStringLit(logCall.Message)
	if !ok {
		return nil
	}

	invalidChars := r.findInvalidCharacters(message)

	if len(invalidChars) > 0 {
		cleanedMessage := r.CleanMessage(message)

		return []analysis.Diagnostic{
			{
				Pos:      logCall.Message.Pos(),
				End:      logCall.Message.End(),
				Message:  r.formatDiagnosticMessage(invalidChars),
				Category: r.Name(),
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: "Remove special symbols",
						TextEdits: []analysis.TextEdit{
							{
								Pos:     logCall.Message.Pos(),
								End:     logCall.Message.End(),
								NewText: []byte(fmt.Sprintf(`"%s"`, cleanedMessage)),
							},
						},
					},
				},
			},
		}
	}

	return nil
}

func (r *NoSpecialSymbolsRule) isAllowedChar(char rune) bool {
	if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
		return true
	}
	if char >= '0' && char <= '9' {
		return true
	}
	if char == ' ' || char == '-' || char == '_' {
		return true
	}
	return false
}

func (r *NoSpecialSymbolsRule) findInvalidCharacters(message string) []InvalidChar {
	var invalid []InvalidChar
	seen := make(map[rune]bool)

	for position, char := range message {
		if !r.isAllowedChar(char) && !seen[char] {
			invalid = append(invalid, InvalidChar{
				Char:     char,
				Position: position,
				Category: r.categorizeChar(char),
			})
			seen[char] = true

			if len(invalid) >= 5 {
				break
			}
		}
	}

	return invalid
}

type InvalidChar struct {
	Char     rune
	Position int
	Category string
}

func (r *NoSpecialSymbolsRule) categorizeChar(char rune) string {
	if unicode.Is(unicode.Cyrillic, char) {
		return "cyrillic"
	}
	if unicode.IsPunct(char) {
		return "punctuation"
	}
	if unicode.IsSymbol(char) {
		return "symbol"
	}
	if char > 0x1F000 {
		return "emoji"
	}
	return "special character"
}

func (r *NoSpecialSymbolsRule) formatDiagnosticMessage(invalidChars []InvalidChar) string {
	if len(invalidChars) == 0 {
		return r.Message()
	}

	categories := make(map[string][]string)
	for _, ic := range invalidChars {
		char := fmt.Sprintf("'%c' (U+%04X)", ic.Char, ic.Char)
		categories[ic.Category] = append(categories[ic.Category], char)
	}

	var parts []string
	for category, chars := range categories {
		part := fmt.Sprintf("%s: %s", category, strings.Join(chars, ", "))
		parts = append(parts, part)
	}

	return fmt.Sprintf("%s - found %s", r.Message(), strings.Join(parts, "; "))
}

func (r *NoSpecialSymbolsRule) CheckExpr(expr ast.Expr) (bool, string) {
	message, ok := loggers.ExtractStringLit(expr)
	if !ok || len(message) == 0 {
		return true, ""
	}

	invalidChars := r.findInvalidCharacters(message)

	if len(invalidChars) > 0 {
		return false, r.formatDiagnosticMessage(invalidChars)
	}

	return true, ""
}

func (r *NoSpecialSymbolsRule) CleanMessage(message string) string {
	var cleaned strings.Builder
	cleaned.Grow(len(message))

	for _, char := range message {
		if r.isAllowedChar(char) {
			cleaned.WriteRune(char)
		}
	}

	result := cleaned.String()
	result = strings.TrimSpace(result)
	result = r.normalizeSpaces(result)

	return result
}

func (r *NoSpecialSymbolsRule) normalizeSpaces(s string) string {
	var result strings.Builder
	result.Grow(len(s))

	prevSpace := false
	for _, char := range s {
		if char == ' ' {
			if !prevSpace {
				result.WriteRune(char)
			}
			prevSpace = true
		} else {
			result.WriteRune(char)
			prevSpace = false
		}
	}

	return result.String()
}
