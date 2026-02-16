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
		},
		{
			name:      "empty str",
			message:   `""`,
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tt.message)
			if err != nil {
				t.Fatalf("failed to parse expression: %v", err)
			}

			valid, _ := rule.CheckExpr(expr)

			if valid != tt.wantValid {
				t.Errorf("CheckExpr() valid = %v, want %v", valid, tt.wantValid)
			}
		})
	}

}
