#!/usr/bin/env python3
"""Run `terraform validate` on every HCL example in a provider docs page.

usage: validate_docs_examples.py <docs-page.html.markdown> [more pages...]
env:   TF_REPO (default: the repo containing this script)

Builds the provider from TF_REPO, points Terraform at it with dev_overrides, then validates each
```hcl block that declares a `resource` or `data` (import commands and other snippets are skipped).
Per page it also runs a negative control: a copy of the first example with a made-up argument must FAIL
validation, otherwise the check proves nothing.

`validate` checks the schema only. It does not call the API or plan. Read-only for the repo; all work is in a temp dir.
Exit code: 0 all examples valid and the control failed as expected, 1 otherwise, 2 setup problem.
"""
import os
import re
import shutil
import subprocess
import sys
import tempfile

REPO_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
TF_REPO = os.path.expanduser(os.environ.get('TF_REPO', REPO_ROOT))
HEADER = '''terraform {
  required_providers {
    sumologic = { source = "SumoLogic/sumologic" }
  }
}
provider "sumologic" {
  access_id  = "x"
  access_key = "y"
}
'''
DECL = re.compile(r'^(resource|data)\s+"', re.M)
FIRST_DECL = re.compile(r'^((?:resource|data)\s+"[^"]+"\s+"[^"]+"\s*\{)[ \t]*\n', re.M)


TOP_LEVEL = re.compile(r'^(terraform|provider)\b[^{\n]*\{', re.M)
FILE_CALL = re.compile(r'(jsondecode\()?file\(\s*"\$\{path\.module\}/([^"]+)"\s*\)')


def strip_provider_blocks(body):
    """Drop top-level `terraform {}` / `provider {}` blocks: the header supplies its own, a second one fails validation."""
    out, pos = [], 0
    for m in TOP_LEVEL.finditer(body):
        if m.start() < pos:
            continue
        depth, i, quote = 1, m.end(), False
        while i < len(body) and depth:
            c = body[i]
            if c == '"' and body[i - 1] != '\\':
                quote = not quote
            elif not quote:
                depth += (c == '{') - (c == '}')
            i += 1
        out.append(body[pos:m.start()])
        pos = i
    out.append(body[pos:])
    return ''.join(out)


def uses_other_providers(body):
    """True if the example declares a resource/data of another provider (e.g. aws_*); only the sumologic provider is installed."""
    return any(not t.startswith('sumologic_') for t in re.findall(r'^(?:resource|data)\s+"([^"]+)"', body, re.M))


def examples_with_lines(text):
    """[(line, body)] for each hcl example that declares a resource/data. `line` is the 1-based line of its ```hcl fence."""
    found = []
    for m in re.finditer(r'```hcl\n(.*?)```', text, re.S):
        body = strip_provider_blocks(m.group(1))
        if DECL.search(body):
            found.append((text.count('\n', 0, m.start()) + 1, body))
    return found


def examples_from_text(text):
    return [body for _, body in examples_with_lines(text)]


def file_stubs(body):
    """(name, content) for each `file("${path.module}/x")` the example reads: `{}` if it is parsed as JSON, else a line of text."""
    return [(name, '{}' if json else 'x') for json, name in FILE_CALL.findall(body)]


def examples(path):
    return examples_with_lines(open(path, encoding='utf-8').read())


HEADER_LINES = HEADER.count('\n')


def example_relative(out):
    """Terraform reports lines of the generated main.tf; shift them so they count from the first line of the example."""
    out = re.sub(r'on main\.tf line (\d+)', lambda m: f'at line {int(m.group(1)) - HEADER_LINES} of the example', out)
    return re.sub(r'^(\s+)(\d+):', lambda m: f'{m.group(1)}{int(m.group(2)) - HEADER_LINES}:', out, flags=re.M)


def drop_warnings(out):
    """Remove Warning: blocks (the dev_overrides notice is always there and hides the real error)."""
    blocks = re.split(r'\n\s*\n(?=(?:Warning: |Error: |Success))', out.strip())
    kept = [b for b in blocks if not b.startswith('Warning: ')]
    return '\n\n'.join(kept) if kept else out.strip()


def first_error(detail):
    found = re.findall(r'Error:.*', detail)
    return found[0] if found else detail.strip().splitlines()[0] if detail.strip() else ''


def escape_annotation(text):
    return text.replace('%', '%25').replace('\r', '%0D').replace('\n', '%0A')


def validate(workdir, body, env):
    os.makedirs(workdir, exist_ok=True)
    with open(os.path.join(workdir, 'main.tf'), 'w') as f:
        f.write(HEADER + body)
    for name, content in file_stubs(body):
        with open(os.path.join(workdir, name), 'w') as f:
            f.write(content + '\n')
    r = subprocess.run(['terraform', 'validate', '-no-color'], cwd=workdir, env=env, capture_output=True, text=True)
    out = r.stdout + r.stderr
    return r.returncode == 0, example_relative(drop_warnings(out)).strip()


def main(pages):
    if not pages:
        print(__doc__)
        return 2
    for tool in ('terraform', 'go'):
        if not shutil.which(tool):
            print(f'{tool} not found on PATH')
            return 2
    if not os.path.isdir(TF_REPO):
        print(f'provider repo not found: {TF_REPO} (set TF_REPO)')
        return 2

    tmp = tempfile.mkdtemp(prefix='tfval-')
    binary_dir = os.path.join(tmp, 'bin')
    os.makedirs(binary_dir)
    print(f'building provider from {TF_REPO} ...')
    b = subprocess.run(['go', 'build', '-o', os.path.join(binary_dir, 'terraform-provider-sumologic'), '.'],
                       cwd=TF_REPO, capture_output=True, text=True)
    if b.returncode != 0:
        print(b.stdout + b.stderr)
        return 2
    rc = os.path.join(tmp, 'tfrc')
    with open(rc, 'w') as f:
        f.write('provider_installation {\n  dev_overrides { "SumoLogic/sumologic" = "%s" }\n  direct {}\n}\n' % binary_dir)
    env = dict(os.environ, TF_CLI_CONFIG_FILE=rc)

    failed = False
    problems = []  # (page, line, title, detail) for annotations and the job summary
    for page in pages:
        found = examples(page)
        exs = [body for _, body in found]
        checkable = [b for b in exs if not uses_other_providers(b)]
        print(f'\n{page}: {len(exs)} example(s)')
        if not exs:
            print('  FAIL: no resource/data examples found')
            problems.append((page, 1, 'No examples', 'The page has no ```hcl example that declares a resource or data source.'))
            failed = True
            continue
        for i, (line, body) in enumerate(found):
            ok, detail = validate(os.path.join(tmp, f'p{abs(hash(page))}', f'ex{i}'), body, env)
            first = body.strip().splitlines()[0][:70]
            if uses_other_providers(body):
                print(f'  skip example {i + 1}: {first} (uses a provider other than sumologic; not checked)')
                continue
            print(f'  {"ok  " if ok else "FAIL"} example {i + 1} (page line {line}): {first}')
            if not ok:
                print('\n'.join('       ' + l for l in detail.splitlines()))
                problems.append((page, line, f'Example {i + 1} does not validate', detail))
            failed |= not ok
        if not checkable:
            print('  FAIL: every example uses another provider, so nothing could be checked')
            problems.append((page, 1, 'Nothing could be checked', 'Every example uses a provider other than sumologic.'))
            failed = True
            continue
        control = FIRST_DECL.sub(r'\1\n  zz_not_a_real_argument = true\n', checkable[0], count=1)
        if control == checkable[0]:
            print('  FAIL: could not build the negative control (first example has no `resource "t" "n" {` line)')
            failed = True
            continue
        ok, detail = validate(os.path.join(tmp, f'p{abs(hash(page))}', 'control'), control, env)
        if ok:
            print('  FAIL: negative control PASSED validation, so this check cannot catch errors')
            problems.append((page, 1, 'Negative control passed', 'A made-up argument was accepted, so the check cannot catch errors.'))
            failed = True
        else:
            print(f'  ok   negative control failed as expected: {first_error(detail)}')
    shutil.rmtree(tmp, ignore_errors=True)
    print('\nNote: validate checks the schema only; no API call, no plan.')
    if problems:
        report(problems)
    return 1 if failed else 0


HOW_TO_FIX = (
    'To reproduce locally (needs go and terraform on your PATH):\n'
    '  python3 -I scripts/validate_docs_examples.py <page>\n'
    'Line numbers above count from the first line of the example. Each example must be valid HCL, use only '
    'attributes that exist in the schema, and declare everything it references. See "Docs standard" in AGENTS.md.'
)


def report(problems):
    """Print a how-to-fix footer; inside GitHub Actions also emit annotations and a job summary."""
    print('\n' + HOW_TO_FIX)
    if os.environ.get('GITHUB_ACTIONS') == 'true':
        for page, line, title, detail in problems:
            print(f'::error file={page},line={line},title={escape_annotation(title)}::{escape_annotation(detail)}')
    summary = os.environ.get('GITHUB_STEP_SUMMARY')
    if summary:
        with open(summary, 'a', encoding='utf-8') as f:
            f.write('### Docs examples that do not validate\n\n')
            for page, line, title, detail in problems:
                f.write(f'**{page}** (page line {line}): {title}\n\n```\n{detail}\n```\n\n')
            f.write('```\n' + HOW_TO_FIX + '\n```\n')


if __name__ == '__main__':
    sys.exit(main(sys.argv[1:]))
