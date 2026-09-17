# Changelog

## Unreleased

- Prepare passthrough for consideration as a public golangci-lint linter.
- Lower the minimum Go version from 1.26 to 1.22.
- Add a standalone `passthrough` diagnostic command and release workflow.
- Remove module-plugin registration and its `init()` function, as required by golangci-lint's public-linter checklist.

The module-plugin integration remains available in the immutable `v0.1.2` release. Existing custom golangci-lint builds should remain pinned to that version.

## v0.1.2

- Publish passthrough as a golangci-lint module plugin.
