// Package analyzer implements the passthrough analysis.Analyzer.
package analyzer

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type Config struct {
	MaxExtraArgs int
}

func (c Config) Validate() error {
	if c.MaxExtraArgs < 0 {
		return fmt.Errorf("max-extra-args must be a nonnegative integer, got %d", c.MaxExtraArgs)
	}

	return nil
}

func New(cfg Config) (*analysis.Analyzer, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &analysis.Analyzer{
		Name: "passthrough",
		Doc:  "report functions and methods that structurally pass their parameters through to one call",
		Run: func(pass *analysis.Pass) (any, error) {
			run(pass, cfg)
			return nil, nil
		},
	}, nil
}

func run(pass *analysis.Pass, cfg Config) {
	for _, file := range pass.Files {
		for _, declaration := range file.Decls {
			function, ok := declaration.(*ast.FuncDecl)
			if !ok {
				continue
			}

			inspectFunction(pass, cfg, function)
		}
	}
}

func inspectFunction(pass *analysis.Pass, cfg Config, function *ast.FuncDecl) {
	if function.Body == nil || len(function.Body.List) != 1 {
		return
	}

	returnStatement, ok := function.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(returnStatement.Results) != 1 {
		return
	}

	returned := unparen(returnStatement.Results[0])
	call, ok := returned.(*ast.CallExpr)
	if !ok || !isCallable(pass, call.Fun) {
		return
	}

	object, ok := pass.TypesInfo.Defs[function.Name].(*types.Func)
	if !ok {
		return
	}

	signature, ok := object.Type().(*types.Signature)
	if !ok {
		return
	}

	parameters := make(map[*types.Var]struct{}, signature.Params().Len())
	for parameter := range signature.Params().Variables() {
		if parameter.Name() == "" || parameter.Name() == "_" {
			return
		}
		parameters[parameter] = struct{}{}
	}

	observed := make(map[*types.Var]struct{}, len(parameters))
	extraArguments := 0
	for _, argument := range call.Args {
		identifier, ok := unparen(argument).(*ast.Ident)
		if ok {
			if parameter, parameterOK := pass.TypesInfo.Uses[identifier].(*types.Var); parameterOK {
				if _, declared := parameters[parameter]; declared {
					observed[parameter] = struct{}{}
					continue
				}
			}
		}

		extraArguments++
	}

	if len(observed) != len(parameters) || extraArguments > cfg.MaxExtraArgs {
		return
	}

	var callee bytes.Buffer
	if err := printer.Fprint(&callee, pass.Fset, call.Fun); err != nil {
		return
	}

	pass.Reportf(function.Name.Pos(), "function %s is a passthrough to %s", function.Name.Name, callee.String())
}

func isCallable(pass *analysis.Pass, expression ast.Expr) bool {
	expression = unparen(expression)
	if typeAndValue, ok := pass.TypesInfo.Types[expression]; ok && typeAndValue.IsType() {
		return false
	}

	typ := pass.TypesInfo.TypeOf(expression)
	if typ == nil {
		return false
	}

	_, ok := typ.Underlying().(*types.Signature)
	return ok
}

func unparen(expression ast.Expr) ast.Expr {
	for {
		parenthesized, ok := expression.(*ast.ParenExpr)
		if !ok {
			return expression
		}
		expression = parenthesized.X
	}
}
