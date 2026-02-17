package loglinter

import (
	"github.com/golangci/plugin-module-register/register"
	"github.com/hel1th/loglinter/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

type plugin struct{}

func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		analyzer.Analyzer,
	}, nil
}

func (p *plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}

func New(settings any) (register.LinterPlugin, error) {
	return &plugin{}, nil
}

func init() {
	register.Plugin("loglinter", New)
}
