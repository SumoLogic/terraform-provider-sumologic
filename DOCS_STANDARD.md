# Docs standard

The standard for provider docs pages. [`AGENTS.md`](AGENTS.md) links here; read it before you add or change a resource or data source.

Every new resource or data source, and every change to an existing attribute, includes docs written to this standard. Use an existing page as a model, for example `website/docs/r/data_archiving_destination.html.markdown` (one example per variant, allowed values and limits, `ForceNew` / `Sensitive` markers, import).

## 1. Examples: one per common use case

- Do not stop at one minimal snippet. Add a separate, titled example for each common use case (`### Basic`, `### With <feature>`).
- Each example is valid HCL and uses only attributes that exist in the schema. Copy it from the acceptance test config where possible.
- Use realistic values, not `foo` / `bar`. Show required fields first.
- For resources with several `type`s or nested blocks, give one example per variant.
- Examples are self-contained: every `resource.x.y`, `data.x.y` and `var.x` an example references is declared in the same block of HCL. `terraform validate` fails on an undeclared reference, and so does a customer who copies the snippet.
- Cross-resource wiring (a role that references a user, a source that references a collector) gets its own example using resource references, not hard-coded IDs.

## 2. Field descriptions: complete

Every argument and exported attribute states, in this order:

1. What it does, in user terms. Use the UI label so people can find it.
2. Required or optional, and what happens when omitted (the default, or "the current value in Sumo Logic is kept").
3. Allowed values (list every one) and limits (min, max, length, pattern). Mirror the schema validation.
4. Side effects and gotchas: writes that reset something, things not reverted on destroy, `ForceNew`, eventual consistency.
5. Units (days, seconds, bytes) and format (ISO 8601, `7d`).

Check every statement against the [Sumo Logic API documentation](https://help.sumologic.com/docs/api/) or the code. Do not paraphrase from memory. If you cannot verify a claim, leave it out.

## 3. Guidance between fields

- Deprecated field: `Deprecated: use `x` instead.` with the replacement and when to migrate.
- Preferred field: say which one to use and why.
- Mutually exclusive or dependent fields: name the partner (`conflicts_with`, `required_with`, `ExactlyOneOf`) in both fields' descriptions.
- Nested blocks: heading ``### Schema for `x` `` (no space before the closing backtick), linked as `[x schema](#schema-for-x)`. A bare `#x` anchor is dead.

## 4. Other expected sections

- Import: the exact id format and a working `terraform import` line.
- Singleton or global resources: say only one may be defined, and list exactly what `destroy` resets (attribute and value) and what it leaves alone.
- A link to the help.sumologic.com page for the feature when there is one.
- Source resources (`sumologic_*_source`): the common source properties are documented once in `website/docs/index.html.markdown` under `## Common Source Properties`. A source page links to that section and documents only its own arguments. A new common property goes in the index section, not on each page.

## 5. Validate the examples

Do not trust examples by eye. Run this on every docs page you touched (needs `go` and `terraform` on your PATH):

```bash
python3 -I scripts/validate_docs_examples.py website/docs/r/<x>.html.markdown
```

It builds the provider from this checkout, runs `terraform validate` on each `resource` / `data` example, and runs a negative control (the first example plus a made-up argument must fail). Exit 0 means all examples are valid. A page with no examples fails. Examples that declare another provider's resources (for example `aws_*`) are skipped, so review those by eye. `validate` checks the schema only: it makes no API call and does no plan. Its unit tests: `python3 -I scripts/validate_docs_examples_test.py`.

CI runs the same script on every docs page a PR adds or changes (`Validate docs examples` workflow), and the PR fails if any example does not validate. Pages the PR does not touch are not checked.

## Checklist before hand-off

- [ ] At least one example per common use case; each passes `terraform validate`.
- [ ] Every new or changed attribute has what / default / allowed values / limits / side effects.
- [ ] Deprecations and preferred-field pointers are in both the docs and the schema.
- [ ] Import section is correct.
- [ ] Singleton resources: destroy behaviour is listed by value.
- [ ] Nested-block links use `#schema-for-<x>`.
- [ ] Every schema key is documented and every documented key exists in the schema (a source page counts keys listed under Common Source Properties; ignore `id`).
- [ ] `CHANGELOG.md` updated.
