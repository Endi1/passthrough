package passthrough

import (
	"strings"
	"testing"

	"github.com/golangci/plugin-module-register/register"
)

func TestRegisteredPlugin(t *testing.T) {
	constructor, err := register.GetPlugin(pluginName)
	if err != nil {
		t.Fatalf("GetPlugin() error = %v", err)
	}

	configured, err := constructor(map[string]any{"max-extra-args": 2})
	if err != nil {
		t.Fatalf("constructor() error = %v", err)
	}
	if mode := configured.GetLoadMode(); mode != register.LoadModeTypesInfo {
		t.Fatalf("GetLoadMode() = %q, want %q", mode, register.LoadModeTypesInfo)
	}

	analyzers, err := configured.BuildAnalyzers()
	if err != nil {
		t.Fatalf("BuildAnalyzers() error = %v", err)
	}
	if len(analyzers) != 1 || analyzers[0].Name != pluginName {
		t.Fatalf("BuildAnalyzers() = %#v, want one passthrough analyzer", analyzers)
	}
}

func TestSettings(t *testing.T) {
	tests := []struct {
		name string
		raw  any
		want int
	}{
		{name: "default", raw: nil, want: 0},
		{name: "zero", raw: map[string]any{"max-extra-args": 0}, want: 0},
		{name: "positive", raw: map[string]any{"max-extra-args": 3}, want: 3},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			created, err := New(test.raw)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			if got := created.(*plugin).config.MaxExtraArgs; got != test.want {
				t.Fatalf("MaxExtraArgs = %d, want %d", got, test.want)
			}
		})
	}
}

func TestInvalidSettings(t *testing.T) {
	tests := []struct {
		name string
		raw  any
	}{
		{name: "negative", raw: map[string]any{"max-extra-args": -1}},
		{name: "fractional", raw: map[string]any{"max-extra-args": 1.5}},
		{name: "string", raw: map[string]any{"max-extra-args": "one"}},
		{name: "null", raw: map[string]any{"max-extra-args": nil}},
		{name: "unknown", raw: map[string]any{"other": 1}},
		{name: "malformed root", raw: []any{1}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := New(test.raw)
			if err == nil {
				t.Fatal("New() unexpectedly succeeded")
			}
			if !strings.Contains(err.Error(), "max-extra-args") && test.name != "unknown" && test.name != "malformed root" {
				t.Fatalf("error %q does not identify max-extra-args", err)
			}
		})
	}
}

func TestPluginInstancesDoNotShareSettings(t *testing.T) {
	first, err := New(map[string]any{"max-extra-args": 1})
	if err != nil {
		t.Fatal(err)
	}
	second, err := New(map[string]any{"max-extra-args": 4})
	if err != nil {
		t.Fatal(err)
	}

	if got := first.(*plugin).config.MaxExtraArgs; got != 1 {
		t.Fatalf("first MaxExtraArgs = %d, want 1", got)
	}
	if got := second.(*plugin).config.MaxExtraArgs; got != 4 {
		t.Fatalf("second MaxExtraArgs = %d, want 4", got)
	}
}
