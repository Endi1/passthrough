# Passthrough

[![CI](https://github.com/Endi1/passthrough/actions/workflows/ci.yml/badge.svg)](https://github.com/Endi1/passthrough/actions/workflows/ci.yml)

`passthrough` is a type-aware Go analyzer that reports structural forwarding wrappers. A finding is a prompt for review, **not proof that a function should be removed**: wrappers can be intentional API, compatibility, instrumentation, or architectural boundaries.

The analyzer uses `golang.org/x/tools/go/analysis` and is being prepared for inclusion as a public golangci-lint linter. It is not built into stock golangci-lint yet. See the [upstream integration plan](docs/golangci-lint.md) for the remaining proposal and contribution steps.

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
func Convert(p string) Name { return Name(p) }           // type conversion
func Partial(a, b string) string { return target(a) }    // b is omitted
func Changed(p string) string { return target(trim(p)) } // p is not direct
func Field(p Data) string { return target(p.Field) }     // p is not direct
func Extra(p string) string { return target(p, "x") }   // needs max-extra-args >= 1
```

### Boundaries

The analyzer inspects `*ast.FuncDecl` values only, not function literals. Declarations without bodies, blank (`_`) or unnamed parameters, and non-call return expressions do not qualify. It does not perform whole-program analysis, judge architectural intent, exempt exported/deprecated/recursive/interface methods, or provide deletion fixes. Recursive forwarding is therefore reported when it meets the structural rule.

## Standalone use

Passthrough requires Go 1.22 or newer.

Before the next tagged release, install it from a checkout:

```sh
go install ./cmd/passthrough
passthrough ./...
```

After a release containing the command is published, install `github.com/Endi1/passthrough/cmd/passthrough` at that immutable tag. The standalone command uses the default `max-extra-args` value of zero. Integrators can construct a configured analyzer directly:

```go
configured, err := analyzer.New(analyzer.Config{MaxExtraArgs: 2})
```

`max-extra-args` accepts nonnegative integers only.

## golangci-lint

Once passthrough is accepted into golangci-lint, the intended configuration is:

```yaml
version: "2"
linters:
  enable:
    - passthrough
  settings:
    passthrough:
      max-extra-args: 0
```

Intentional wrappers can then be suppressed through golangci-lint:

```go
//nolint:passthrough // Intentional compatibility boundary.
func OldName(p string) string { return NewName(p) }
```

The previously published `v0.1.2` release remains available as a golangci-lint module plugin. Consumers of that plugin should continue pinning `github.com/Endi1/passthrough@v0.1.2`; the plugin registration was removed from the development branch to satisfy golangci-lint's public-linter requirement that linter repositories not contain `init()`.

## Development

```sh
go mod download
make fmt-check
make test
make vet
make build       # creates bin/passthrough
make lint
```

Tests use `analysistest` fixtures for the default and configured behavior. CI runs tests, vet, and a binary build with Go 1.22, and runs golangci-lint separately with its supported Go toolchain.

## License

MIT; see [LICENSE](LICENSE).
