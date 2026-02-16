package rules

import (
	"go/parser"
	"testing"
)

func TestLowercaseRule(t *testing.T) {
	rule := &LowercaseRule{}

	tests := []struct {
		name      string
		message   string
		wantValid bool
		wantError string
	}{
		{
			name:      "valid lowercase",
			message:   `"starting app"`,
			wantValid: true,
		},
		{
			name:      "valid lowercase with spaces",
			message:   `"  starting app"`,
			wantValid: true,
		},
		{
			name:      "valid  spaces",
			message:   `"       "`,
			wantValid: true,
		},
		{
			name:      "invalid uppercase",
			message:   `"Starting app"`,
			wantValid: false,
			wantError: "message starts with uppercase letter: 'S'",
		},
		{
			name:      "empty str",
			message:   `""`,
			wantValid: true,
		},
	}

	for _, tt := range tests {
		expr, err := parser.ParseExpr(tt.message)
		if err != nil {
			t.Fatalf("failed to parse expr: %s", tt.message)
		}

		errMsg, valid := rule.CheckExpr(expr)
		if valid != tt.wantValid {
			t.Errorf("TEST: %s - CheckExpr() valid = %v, want %v", tt.name, valid, tt.wantValid)
		}

		if !tt.wantValid && errMsg != tt.wantError {
			t.Errorf("TEST: %q - CheckExpr() error = %v, want %v", tt.name, errMsg, tt.wantError)
		}
	}

}
