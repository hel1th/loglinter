package rules

import (
	"go/parser"
	"testing"
)

func TestEnglishOnlyRule(t *testing.T) {
	rule := &EnglishOnlyRule{}

	tests := []struct {
		name      string
		message   string
		wantValid bool
	}{
		{
			name:      "valid english",
			message:   `"starting server"`,
			wantValid: true,
		},
		{
			name:      "invalid cyrillic",
			message:   `"запуск сервера"`,
			wantValid: false,
		},
		{
			name:      "invalid mixed",
			message:   `"server запуск"`,
			wantValid: false,
		},
		{
			name:      "valid with numbers",
			message:   `"server 8080"`,
			wantValid: true,
		},
		{
			name:      "valid with punctuation",
			message:   `"server started successfully!"`,
			wantValid: true,
		},
		{
			name:      "invalid chinese",
			message:   `"服务器启动"`,
			wantValid: false,
		},
		{
			name:      "invalid arabic",
			message:   `"خادم بدأ"`,
			wantValid: false,
		},
		{
			name:      "invalid extended latin (café)",
			message:   `"café server"`,
			wantValid: false,
		},
		{
			name:      "empty string",
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
