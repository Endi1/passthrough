// Command passthrough runs the passthrough analyzer.
package main

import (
	"github.com/Endi1/passthrough/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.Default())
}
