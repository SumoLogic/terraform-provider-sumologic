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

## Docs standard

Every new resource or data source, and every change to an existing attribute, includes docs written to this standard. Use an existing page as a model, for example `website/docs/r/data_archiving_destination.html.markdown` (one example per variant, allowed values and limits, `ForceNew` / `Sensitive` markers, import).

### 1. Examples: one per common use case

- Do not stop at one minimal snippet. Add a separate, titled example for each common use case (`### Basic`, `### With <feature>`).
- Each example is valid HCL and uses only attributes that exist in the schema. Copy it from the acceptance test config where possible.
- Use realistic values, not `foo` / `bar`. Show required fields first.
- For resources with several `type`s or nested blocks, give one example per variant.
- Examples are self-contained: every `resource.x.y`, `data.x.y` and `var.x` an example references is declared in the same block of HCL. `terraform validate` fails on an undeclared reference, and so does a customer who copies the snippet.
- Cross-resource wiring (a role that references a user, a source that references a collector) gets its own example using resource references, not hard-coded IDs.

### 2. Field descriptions: complete

Every argument and exported attribute states, in this order:

1. What it does, in user terms. Use the UI label so people can find it.
2. Required or optional, and what happens when omitted (the default, or "the current value in Sumo Logic is kept").
3. Allowed values (list every one) and limits (min, max, length, pattern). Mirror the schema validation.
4. Side effects and gotchas: writes that reset something, things not reverted on destroy, `ForceNew`, eventual consistency.
5. Units (days, seconds, bytes) and format (ISO 8601, `7d`).

Check every statement against the [Sumo Logic API documentation](https://help.sumologic.com/docs/api/) or the code. Do not paraphrase from memory. If you cannot verify a claim, leave it out.

### 3. Guidance between fields

- Deprecated field: `Deprecated: use `x` instead.` with the replacement and when to migrate.
- Preferred field: say which one to use and why.
- Mutually exclusive or dependent fields: name the partner (`conflicts_with`, `required_with`, `ExactlyOneOf`) in both fields' descriptions.
- Nested blocks: heading ``### Schema for `x` `` (no space before the closing backtick), linked as `[x schema](#schema-for-x)`. A bare `#x` anchor is dead.

### 4. Other expected sections

- Import: the exact id format and a working `terraform import` line.
- Singleton or global resources: say only one may be defined, and list exactly what `destroy` resets (attribute and value) and what it leaves alone.
- A link to the help.sumologic.com page for the feature when there is one.
- Source resources (`sumologic_*_source`): the common source properties are documented once in `website/docs/index.html.markdown` under `## Common Source Properties`. A source page links to that section and documents only its own arguments. A new common property goes in the index section, not on each page.

### 5. Validate the examples

Do not trust examples by eye. Run this on every docs page you touched (needs `go` and `terraform` on your PATH):

```bash
python3 -I scripts/validate_docs_examples.py website/docs/r/<x>.html.markdown
```

It builds the provider from this checkout, runs `terraform validate` on each `resource` / `data` example, and runs a negative control (the first example plus a made-up argument must fail). Exit 0 means all examples are valid. A page with no examples fails. Examples that declare another provider's resources (for example `aws_*`) are skipped, so review those by eye. `validate` checks the schema only: it makes no API call and does no plan. Its unit tests: `python3 -I scripts/validate_docs_examples_test.py`.

CI runs the same script on every docs page a PR adds or changes (`Validate docs examples` workflow), and the PR fails if any example does not validate. Pages the PR does not touch are not checked.

### Checklist before hand-off

- [ ] At least one example per common use case; each passes `terraform validate`.
- [ ] Every new or changed attribute has what / default / allowed values / limits / side effects.
- [ ] Deprecations and preferred-field pointers are in both the docs and the schema.
- [ ] Import section is correct.
- [ ] Singleton resources: destroy behaviour is listed by value.
- [ ] Nested-block links use `#schema-for-<x>`.
- [ ] Every schema key is documented and every documented key exists in the schema (a source page counts keys listed under Common Source Properties; ignore `id`).
- [ ] `CHANGELOG.md` updated.
