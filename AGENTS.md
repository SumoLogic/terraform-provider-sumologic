# AGENTS.md

Guidance for AI coding agents (and humans) working in this repository. It complements [`CONTRIBUTING.md`](CONTRIBUTING.md), which covers adding resources and data sources and cutting releases.

This is the public, open-source Terraform provider for Sumo Logic. Everything you write here (code, docs, commits, PR text) is public: no internal ticket IDs, internal branch names, internal hostnames, internal class names or server internals.

## Layout

- `sumologic/` - provider code.
  - `resource_sumologic_<x>.go` / `_test.go`: resource schema, CRUD, acceptance tests.
  - `data_source_sumologic_<x>.go`: data sources.
  - `sumologic_<x>.go`: data types and API client calls.
  - `provider.go`: every resource and data source must be registered in `ResourcesMap` / `DataSourcesMap`.
- `website/docs/r/<x>.html.markdown` and `website/docs/d/<x>.html.markdown` - one docs page per resource / data source. `website/docs/index.html.markdown` documents provider arguments and the Common Source Properties.
- `scripts/` - helper scripts, including the docs example validator.
- `CHANGELOG.md` - release notes.

## Commands

```bash
make build            # gofmt check + go install
make fmt              # reformat Go code (CI fails on unformatted code)
make lint             # golangci-lint
make test             # unit tests (30s timeout)
make website-lint     # misspell check on website/
tfproviderlint ./...  # install: go install github.com/bflad/tfproviderlint/cmd/tfproviderlint@latest
```

Run `make fmt`, `make test` and `tfproviderlint ./...` before opening a PR. The Go version is in `go.mod`.

## Tests

- Unit tests run with `make test`. Acceptance tests (`make testacc`, sets `TF_ACC=1`) call the real Sumo Logic API and need credentials. Do not run them unless the user has provided credentials and asked for it.
- New or changed behaviour needs a test. For a bug fix, write the failing test first.
- PRs from forks do not get the CI secrets, so acceptance tests only run on branches in this repo.

## Code conventions

- Follow the neighbouring resource: copy its structure, naming and error handling rather than inventing a new pattern.
- Allowed values and limits belong in the schema (`validation.StringInSlice`, `IntBetween`, `StringMatch`, ...), and the docs must say the same thing.
- Mark `ForceNew` and `Sensitive` on the schema where they apply, and say so in the docs.
- Deprecating a field: set `Deprecated` in the schema with a message that names the replacement, and mirror it in the docs.
- Every schema `Description` / `Deprecated` string must stay in sync with the docs page.
- Do not make breaking changes to existing attributes without being asked.

## Pull requests

- Branch from the latest `master`. Do not push to `master` directly.
- Use the PR template in `.github/pull_request_template.md`. Keep the template and every checkbox in the PR description, and tick a box only when it is true. A CI check fails the PR if the template items are missing or unticked.
- Update `CHANGELOG.md`: keep the `## X.Y.Z (Unreleased)` header and the `* Add new change notes here` line, and add your entry under the right heading (`DOCS:`, `FEATURES:`, `ENHANCEMENTS:`, `BUG FIXES:`). Do not choose a version number.
- Do not create tags or releases (Sumo Logic team only, see `CONTRIBUTING.md`).
- Agents: do not push, open or update a PR, or post comments without the user's explicit confirmation.
- Code owners (`.github/CODEOWNERS`) are requested for review automatically.

## Docs

Every new resource or data source, and every change to an attribute, needs docs that follow [`DOCS_STANDARD.md`](DOCS_STANDARD.md). Read it before you write them. In short:

- One titled, valid example per common use case; examples are self-contained and copied from a working test config where possible.
- Every argument and exported attribute states what it does, required or optional and the default, allowed values and limits, side effects, and units.
- Deprecated and preferred fields point at each other; import and `destroy` behaviour are documented.
- Check examples with `python3 -I scripts/validate_docs_examples.py website/docs/r/<x>.html.markdown`. CI runs the same script on every docs page a PR adds or changes.
