package analyzer

import (
	"github.com/hel1th/loglinter/pkg/loggers"
	"github.com/hel1th/loglinter/pkg/rules"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name:             "loglinter",
	Doc:              "checks log messages for bad patterns",
	Run:              run,
	Requires:         []*analysis.Analyzer{},
	RunDespiteErrors: false,
}

type Config struct {
	EnabledRules            []string
	DisabledRules           []string
	CustomSensitivePatterns []string
}

var config = &Config{}

func run(pass *analysis.Pass) (any, error) {
	detector := loggers.NewDetector(pass)
	ruleSet := createRuleSet()

	for _, file := range pass.Files {
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
		&rules.EnglishOnlyRule{},
		&rules.NoSpecialSymbolsRule{},
		&rules.SensitiveDataRule{},
	}

	for _, rule := range rulesList {
		if shouldEnableRule(rule.Name()) {
			ruleSet.AddRule(rule)
		}
	}

	return ruleSet
}

func shouldEnableRule(ruleName string) bool {
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
