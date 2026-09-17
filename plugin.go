// Package passthrough registers the passthrough analyzer as a golangci-lint
// module plugin.
package passthrough

import (
	"fmt"

	"github.com/Endi1/passthrough/analyzer"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

const pluginName = "passthrough"

func init() {
	register.Plugin(pluginName, New)
}

type settings struct {
	MaxExtraArgs int `json:"max-extra-args"`
}

type plugin struct {
	config analyzer.Config
}

func New(rawSettings any) (register.LinterPlugin, error) {
	if values, ok := rawSettings.(map[string]any); ok {
		if value, exists := values["max-extra-args"]; exists && value == nil {
			return nil, fmt.Errorf("max-extra-args must be a nonnegative integer, got null")
		}
	}

	decoded, err := register.DecodeSettings[settings](rawSettings)
	if err != nil {
		return nil, fmt.Errorf("invalid passthrough configuration: max-extra-args must be a nonnegative integer: %w", err)
	}

	config := analyzer.Config{MaxExtraArgs: decoded.MaxExtraArgs}
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid passthrough configuration: %w", err)
	}

	return &plugin{config: config}, nil
}

func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	configured, err := analyzer.New(p.config)
	if err != nil {
		return nil, err
	}

	return []*analysis.Analyzer{configured}, nil
}

func (*plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
