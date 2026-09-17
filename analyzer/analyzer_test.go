package analyzer_test

import (
	"testing"

	"github.com/Endi1/passthrough/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	tests := []struct {
		name         string
		maxExtraArgs int
		fixture      string
	}{
		{name: "default", maxExtraArgs: 0, fixture: "maxzero"},
		{name: "two extra arguments", maxExtraArgs: 2, fixture: "maxtwo"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			configured, err := analyzer.New(analyzer.Config{MaxExtraArgs: test.maxExtraArgs})
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			analysistest.Run(t, analysistest.TestData(), configured, test.fixture)
		})
	}
}

func TestNewRejectsNegativeMaxExtraArgs(t *testing.T) {
	if _, err := analyzer.New(analyzer.Config{MaxExtraArgs: -1}); err == nil {
		t.Fatal("New() accepted a negative MaxExtraArgs")
	}
}
