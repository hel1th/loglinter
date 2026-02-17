package analyzer

import (
	"os"
	"slices"

	"github.com/hel1th/loglinter/pkg/config"
	"github.com/hel1th/loglinter/pkg/loggers"
	"github.com/hel1th/loglinter/pkg/rules"
	"github.com/joho/godotenv"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name:             "loglinter",
	Doc:              "checks log messages for bad patterns",
	Run:              run,
	Requires:         []*analysis.Analyzer{},
	RunDespiteErrors: false,
}
var cfg *config.Config

func init() {
	godotenv.Load(".env")

	var err error
	cfg, err = config.LoadConfig(os.Getenv("CONFIG_PATH"))
	if err != nil {
		cfg = config.DefaultConfig()
	}
}

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
	if slices.Contains(cfg.GetDisabledRules(), ruleName) {
		return false
	}

	if cfg.IsEnabled() {
		return slices.Contains(cfg.GetEnabledRules(), ruleName)
	}

	return true
}
