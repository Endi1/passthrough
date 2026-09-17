# golangci-lint upstream integration plan

Passthrough is not built into stock golangci-lint yet. Inclusion requires a proposal and a separate pull request to [`golangci/golangci-lint`](https://github.com/golangci/golangci-lint), following its [new-linter documentation](https://golangci-lint.run/docs/contributing/new-linters/) and [review checklist](https://github.com/golangci/golangci-lint/blob/HEAD/.github/new-linter-checklist.md).

## Standalone repository readiness

Before opening the upstream pull request:

- [x] uses `golang.org/x/tools/go/analysis`;
- [x] declares Go 1.22.0;
- [x] contains no `init()`, `panic()`, `log.Fatal()`, or `os.Exit()` calls;
- [x] does not modify the AST;
- [x] has an MIT license with author and year;
- [x] has unit and functional `analysistest` coverage;
- [x] tests a standard-library call;
- [x] has CI, linting, a README, and a `.gitignore`;
- [x] uses `main` as the default branch;
- [x] has a standalone diagnostic command;
- [ ] publish a new immutable semantic-version tag containing these compatibility changes;
- [ ] verify release CI publishes standalone binaries;
- [ ] confirm with maintainers that the rule does not duplicate an existing linter;
- [ ] sign the golangci-lint CLA.

The historical `v0.1.2` module-plugin release contains `init()` registration and is not the version that should be proposed for upstream inclusion. Create a new tag only after this preparation is merged and CI passes.

## Proposal

Open a new-linter proposal using golangci-lint's discussion template before implementing the upstream adapter. Include:

- repository: `https://github.com/Endi1/passthrough`;
- description: “Reports named functions and methods that structurally forward all parameters unchanged to one call.”;
- rationale and examples of wrappers the diagnostic helps reviewers inspect;
- an explanation that type identity prevents name-based parameter matching;
- the default behavior and the `max-extra-args` option;
- known boundaries from the project README.

Wait for maintainer feedback before investing in the upstream pull request. The checklist explicitly states that compliance does not guarantee acceptance.

## Upstream pull request

After approval, fork the current `golangci/golangci-lint` default branch and make the following changes there. Follow the exact layout and next release number present at that time.

1. Add the newly tagged `github.com/Endi1/passthrough` version to `go.mod` and `go.sum`.
2. Add `PassthroughSettings` to `pkg/config/linters_settings.go`:

   ```go
   type PassthroughSettings struct {
       MaxExtraArgs int `mapstructure:"max-extra-args"`
   }
   ```

   Wire it into `LintersSettings` under `mapstructure:"passthrough"`. Reject negative values during configuration validation before constructing the analyzer.
3. Add `pkg/golinters/passthrough/passthrough.go`. It should translate golangci-lint settings into `analyzer.Config`, call `analyzer.New`, wrap the result with `goanalysis.NewLinterFromAnalyzer`, and use `goanalysis.LoadModeTypesInfo`.
4. Register it alphabetically in `pkg/lint/lintersdb/builder_linter.go` with:

   - `WithSince("v2.<next-minor>.0")`;
   - `WithLoadForGoAnalysis()` because the analyzer requires type information;
   - `WithURL("https://github.com/Endi1/passthrough")`.
5. Add `pkg/golinters/passthrough/passthrough_integration_test.go` using `integration.RunTestdata(t)`.
6. Add `pkg/golinters/passthrough/testdata/passthrough.go` with a standard-library import and default-configuration findings.
7. Add a configured test fixture that proves a wrapper with extra arguments is reported when `max-extra-args` is nonzero but not by default. Follow the current integration-test conventions for per-fixture arguments.
8. Update `.golangci.next.reference.yml`, not `.golangci.reference.yml`:

   - add `passthrough` alphabetically to both available-linter lists;
   - document `max-extra-args` with a nondefault example and a comment stating that the default is `0`.
9. Update `jsonschema/golangci.next.jsonschema.json`, not the released schema.
10. Run the focused command from the contribution guide and the repository's complete test and lint suites.

The PR description must link this repository and briefly describe the linter. Address review feedback as additional commits rather than rewriting reviewed commits, per the upstream checklist.
