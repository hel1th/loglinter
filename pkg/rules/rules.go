package rules

import (
	"github.com/hel1th/loglinter/pkg/loggers"
	"golang.org/x/tools/go/analysis"
)

type Rule interface {
	Check(pass *analysis.Pass, logCall loggers.LogCall) []analysis.Diagnostic

	Name() string
	Message() string
}

func (rs *RuleSet) CheckLogCall(pass *analysis.Pass, logCall loggers.LogCall) []analysis.Diagnostic {
	var diagnostics []analysis.Diagnostic

	for _, rule := range rs.rules {
		diags := rule.Check(pass, logCall)
		diagnostics = append(diagnostics, diags...)
	}

	return diagnostics
}

type RuleSet struct {
	rules []Rule
}

func NewRuleSet() *RuleSet {
	return &RuleSet{
		rules: make([]Rule, 0),
	}
}

func (rs *RuleSet) AddRule(rule Rule) {
	rs.rules = append(rs.rules, rule)
}

func (rs *RuleSet) GetRules() []Rule {
	return rs.rules
}

func DefaultRuleSet() *RuleSet {
	rs := NewRuleSet()

	rs.AddRule(&LowercaseRule{})

	return rs
}
