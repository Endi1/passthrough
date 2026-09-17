package maxzero

import "strings"

var packageFunction = func(value string) string { return value + "!" }

type pair struct{ Field string }
type service struct{ inner store }
type store struct{}
type namedString string

func target(value string) string                       { return value + "!" }
func targetTwo(first, second string) string            { return first + second }
func targetZero() string                               { return "target" }
func pairTarget(first, second string) (string, string) { return first, second }
func variadicTarget(values ...string) []string         { return append([]string(nil), values...) }
func identity[T any](value T) T {
	var zero T
	if false {
		return zero
	}
	return value
}

func ordinary(value string) string { // want `function ordinary is a passthrough to target`
	return target(value)
}

func multipleResults(first, second string) (string, string) { // want `function multipleResults is a passthrough to pairTarget`
	return pairTarget(first, second)
}

func reordered(first, second string) string { // want `function reordered is a passthrough to targetTwo`
	return targetTwo(second, first)
}

func repeated(value string) string { // want `function repeated is a passthrough to targetTwo`
	return targetTwo(value, value)
}

func variadic(values ...string) []string { // want `function variadic is a passthrough to variadicTarget`
	return variadicTarget(values...)
}

func zeroParameters() string { // want `function zeroParameters is a passthrough to targetZero`
	return targetZero()
}

func genericDeclaration[T any](value T) T { // want `function genericDeclaration is a passthrough to identity`
	return identity(value)
}

func genericCall[T any](value T) T { // want `function genericCall is a passthrough to identity\[T\]`
	return identity[T](value)
}

func functionValue(value string) string { // want `function functionValue is a passthrough to packageFunction`
	return packageFunction(value)
}

func builtin(values []string) int { // want `function builtin is a passthrough to len`
	return len(values)
}

func standardLibrary(value string) string { // want `function standardLibrary is a passthrough to strings.TrimSpace`
	return strings.TrimSpace(value)
}

func parenthesized(value string) string { // want `function parenthesized is a passthrough to \(target\)`
	return (target)((value))
}

func parenthesizedReturn(value string) string { // want `function parenthesizedReturn is a passthrough to target`
	return (target((value)))
}

func grouped(first, second string) string { // want `function grouped is a passthrough to targetTwo`
	return targetTwo(first, second)
}

func recursive(value string) string { // want `function recursive is a passthrough to recursive`
	return recursive(value)
}

func (receiver service) method(value string) string { // want `function method is a passthrough to receiver.inner.Get`
	return receiver.inner.Get(value)
}

func (store) Get(value string) string { return value + "!" }

// Negative cases.
func twoStatements(value string) string { _ = value; return target(value) }
func empty()                            {}
func declarationOnly(value string) string
func bare(value string) (result string)                { return }
func multipleExpressions(value string) (string, error) { return target(value), nil }
func notCall(value string) string                      { return value }
func conversion(value string) namedString              { return namedString(value) }
func omitted(first, second string) string              { return target(first) }
func transformed(value string) string                  { return target(target(value)) }
func selectorOnly(value pair) string                   { return target(value.Field) }
func calleeOnly(function func() string) string         { return function() }
func receiverOnly(value service) string                { return value.inner.NoArgs() }
func (store) NoArgs() string                           { return "x" }
func extra(value string) string                        { return targetTwo(value, "extra") }
func blank(_ string) string                            { return target("x") }
func unnamed(string) string                            { return target("x") }

var literal = func(value string) string {
	return target(value)
}
