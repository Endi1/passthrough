# Passthrough

`passthrough` is a type-aware Go analyzer and golangci-lint module plugin that reports structural forwarding wrappers. A finding is a prompt for review, **not proof that a function should be removed**: wrappers can be intentional API, compatibility, instrumentation, or architectural boundaries.

## Detection rule

A named function or method is reported when:

1. its body has exactly one statement;
2. that statement is a `return` with exactly one expression;
3. the expression is a function or method call (not a type conversion);
4. every declared parameter occurs unchanged as a direct call argument at least once; and
5. no more than `max-extra-args` call arguments are anything other than unchanged parameter references.

Parameter identity comes from Go type information, not spelling. Parameters may be reordered or repeated; repeated direct references still count as forwarded and never as extras. A transformed use is an extra argument, but the function may qualify if that parameter is also passed directly:

```go
func NormalizeAndSend(p string) string {
	return send(p, normalize(p)) // matches when max-extra-args >= 1
}
```

Receivers are not required parameters. Their use as a call receiver does not count as an argument, so `return s.inner.Get(id)` can match. If a receiver is explicitly passed as an argument, it is an extra unless that expression independently refers to a declared parameter. A variadic expansion such as `target(args...)` is one syntactic argument, and `args` is forwarded unchanged. Zero-parameter wrappers match when their call has no more than the configured number of extras. Generic declarations and generic calls are supported.

Redundant parentheses around the returned call and its arguments are ignored. Function-valued variables and built-ins are calls and can match. Additional arguments can be arbitrary expressions—including literals, selectors, conversions, and calls with side effects.

### Examples

Reported with the default (`max-extra-args: 0`):

```go
func Get(id string) Item { return store.Get(id) }
func Swap(a, b string) string { return join(b, a) }
func Duplicate(p string) string { return join(p, p) }
func Count(p []string) int { return len(p) }
func Forward[T any](p T) T { return identity[T](p) }
```

Not reported:

```go
func Convert(p string) Name { return Name(p) }          // type conversion
func Partial(a, b string) string { return target(a) }   // b is omitted
func Changed(p string) string { return target(trim(p)) } // p is not direct
func Field(p Data) string { return target(p.Field) }     // p is not direct
func Extra(p string) string { return target(p, "x") }   // needs max-extra-args >= 1
```

### Boundaries

The initial implementation analyzes `*ast.FuncDecl` values only, not function literals. Declarations without bodies, blank (`_`) or unnamed parameters, and non-call return expressions do not qualify. It does not perform whole-program analysis, judge architectural intent, exempt exported/deprecated/recursive/interface methods, or provide deletion fixes. Recursive forwarding is therefore reported when it meets the structural rule.

## Configuration

The default and minimum value is zero:

```yaml
version: "2"
linters:
  enable:
    - passthrough
  settings:
    custom:
      passthrough:
        type: module
        settings:
          max-extra-args: 0
```

`max-extra-args` accepts nonnegative integers only. Negative, fractional, string, null, malformed, and unknown settings are rejected during plugin construction. Configuration belongs to each plugin/analyzer instance; instances do not share mutable state.

This is custom-linter configuration under `linters.settings.custom`, not a built-in golangci-lint setting.

## Build and run

Requirements are Go 1.26 or newer and golangci-lint 2.13.2. The versions and local plugin source are pinned in `go.mod` and `.custom-gcl.yml`.

```sh
go mod download
make build       # creates bin/golangci-lint-passthrough
make test
make test-integration
make lint
```

Run it directly with:

```sh
./bin/golangci-lint-passthrough run ./...
```

A stock golangci-lint binary cannot run the module plugin because module plugins are linked into a custom binary at build time. `golangci-lint custom` reads `.custom-gcl.yml`, imports this module, and produces that binary. The stock executable is only the builder.

Suppress an intentional wrapper through golangci-lint (the analyzer itself deliberately has no suppression logic):

```go
//nolint:passthrough // Intentional compatibility boundary.
func OldName(p string) string { return NewName(p) }
```

## Tests and CI

`go test ./...` runs `analysistest` fixtures and adapter/configuration tests. `make test-integration` builds the real custom binary and verifies default settings, a nondefault limit, plugin loading, and `//nolint:passthrough`. `make lint` runs the repository configuration with the custom binary.

GitHub Actions installs the pinned stock builder and invokes `make ci`; it never uses the stock executable for the lint run. For other CI systems, install `github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.13.2` and run the same target.
