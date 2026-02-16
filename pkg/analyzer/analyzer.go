package analyzer

import (
	"go/ast"

	"github.com/hel1th/loglinter/pkg/loggers"
	"github.com/hel1th/loglinter/pkg/rules"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name:     "loglinter",
	Doc:      "corrects log messages",
	Run:      run,
	Requires: []*analysis.Analyzer{},
}

type Config struct {
	EnabledRules            []string
	DisabledRules           []string
	CustomSensitivePatterns []string
}

var config = &Config{}

func SetConfig(cfg *Config) {
	config = cfg
}

func run(pass *analysis.Pass) (any, error) {
	detector := loggers.NewDetector(pass)

	ruleSet := createRuleSet()

	for _, file := range pass.Files {
		if isTestFile(file) {
			continue
		}

		logCalls := detector.DetectLogCalls(file)

		for _, logCall := range logCalls {
			diagnos := ruleSet.CheckLogCall(pass, logCall)

			for _, diag := range diagnos {
				pass.Report(diag)
			}
		}
	}

	return nil, nil
}

func createRuleSet() *rules.RuleSet {
	ruleSet := rules.NewRuleSet()

	rulesList := []rules.Rule{
		&rules.LowercaseRule{},
	}

	for _, rule := range rulesList {
		if RuleEnabled(rule) {
			ruleSet.AddRule(rule)
		}
	}

	return ruleSet
}

func RuleEnabled(rule rules.Rule) bool {
	ruleName := rule.Name()
	for _, disabled := range config.DisabledRules {
		if disabled == ruleName {
			return false
		}
	}

	if len(config.EnabledRules) > 0 {
		for _, enabled := range config.EnabledRules {
			if enabled == ruleName {
				return true
			}
		}
		return false
	}

	return true
}

func isTestFile(file *ast.File) bool {
	if file.Name.Name == "test" || file.Name.Name == "testing" {
		return true
	}

	return false
}
