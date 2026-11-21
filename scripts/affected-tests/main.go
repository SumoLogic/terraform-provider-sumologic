// Command affected-tests picks the acceptance tests a PR needs to run.
//
// It maps changed files to the resource types they implement by parsing the
// ResourcesMap / DataSourcesMap in sumologic/provider.go, then selects every
// TestAcc* function in a test file whose HCL configs use one of those types
// (so a collector change also runs the source tests that create a collector).
// Anything it cannot map with confidence falls back to the full suite.
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const pkgDir = "sumologic"

var hclUse = regexp.MustCompile(`(?:resource|data)\s+"(sumologic_\w+)"`)

// Result is the selection: Full means run everything, otherwise run Tests.
type Result struct {
	Full   bool
	Reason string
	Tests  []string // TestAcc* function names
}

// RunRegex returns an anchored -run pattern for the selected tests.
func (r Result) RunRegex() string {
	return "^(" + strings.Join(r.Tests, "|") + ")$"
}

func full(reason string) Result { return Result{Full: true, Reason: reason} }

// registered returns resource type -> constructor function names for the
// entries of ResourcesMap and DataSourcesMap in provider.go (a type can be
// both a resource and a data source).
func registered(root string) (map[string][]string, error) {
	file, err := parser.ParseFile(token.NewFileSet(), filepath.Join(root, pkgDir, "provider.go"), nil, 0)
	if err != nil {
		return nil, err
	}
	out := map[string][]string{}
	ast.Inspect(file, func(n ast.Node) bool {
		kv, ok := n.(*ast.KeyValueExpr)
		if !ok {
			return true
		}
		if k, ok := kv.Key.(*ast.Ident); !ok || (k.Name != "ResourcesMap" && k.Name != "DataSourcesMap") {
			return true
		}
		lit, ok := kv.Value.(*ast.CompositeLit)
		if !ok {
			return true
		}
		for _, e := range lit.Elts {
			entry, ok := e.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			name, err := strconv.Unquote(entry.Key.(*ast.BasicLit).Value)
			call, isCall := entry.Value.(*ast.CallExpr)
			if err != nil || !isCall {
				continue
			}
			if fn, ok := call.Fun.(*ast.Ident); ok {
				out[name] = append(out[name], fn.Name)
			}
		}
		return false
	})
	return out, nil
}

// definingFiles maps top-level function name -> file (relative to root) for non-test files.
func definingFiles(root string) (map[string]string, error) {
	files, err := filepath.Glob(filepath.Join(root, pkgDir, "*.go"))
	if err != nil {
		return nil, err
	}
	out := map[string]string{}
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		parsed, err := parser.ParseFile(token.NewFileSet(), f, nil, 0)
		if err != nil {
			return nil, err
		}
		rel, _ := filepath.Rel(root, f)
		for _, d := range parsed.Decls {
			if fd, ok := d.(*ast.FuncDecl); ok && fd.Recv == nil {
				out[fd.Name.Name] = filepath.ToSlash(rel)
			}
		}
	}
	return out, nil
}

type testFile struct {
	path  string // relative to root
	funcs []string
	uses  map[string]bool
}

func testFiles(root string) ([]testFile, error) {
	files, err := filepath.Glob(filepath.Join(root, pkgDir, "*_test.go"))
	if err != nil {
		return nil, err
	}
	var out []testFile
	for _, f := range files {
		src, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		parsed, err := parser.ParseFile(token.NewFileSet(), f, src, 0)
		if err != nil {
			return nil, err
		}
		rel, _ := filepath.Rel(root, f)
		tf := testFile{path: filepath.ToSlash(rel), uses: map[string]bool{}}
		for _, d := range parsed.Decls {
			if fd, ok := d.(*ast.FuncDecl); ok && fd.Recv == nil && strings.HasPrefix(fd.Name.Name, "TestAcc") {
				tf.funcs = append(tf.funcs, fd.Name.Name)
			}
		}
		for _, m := range hclUse.FindAllStringSubmatch(string(src), -1) {
			tf.uses[m[1]] = true
		}
		out = append(out, tf)
	}
	return out, nil
}

// Analyze selects acceptance tests for the changed files (paths relative to root).
func Analyze(root string, changed []string) (Result, error) {
	if len(changed) == 0 {
		return full("no changed files"), nil
	}
	reg, err := registered(root)
	if err != nil {
		return Result{}, err
	}
	defs, err := definingFiles(root)
	if err != nil {
		return Result{}, err
	}
	tests, err := testFiles(root)
	if err != nil {
		return Result{}, err
	}

	// constructor file -> resource types it registers
	fileTypes := map[string][]string{}
	for typ, fns := range reg {
		for _, fn := range fns {
			if f, ok := defs[fn]; ok {
				fileTypes[f] = append(fileTypes[f], typ)
			}
		}
	}

	wantTypes := map[string]bool{}
	wantFiles := map[string]bool{}
	for _, f := range changed {
		f = filepath.ToSlash(f)
		switch {
		case !strings.HasSuffix(f, ".go") && f != "go.mod" && f != "go.sum":
			continue // docs, workflows, scripts: should_test already filters these
		case strings.HasPrefix(f, "scripts/"):
			continue
		case strings.HasPrefix(f, pkgDir+"/") && strings.HasSuffix(f, "_test.go"):
			wantFiles[f] = true
		case fileTypes[f] != nil:
			for _, t := range fileTypes[f] {
				wantTypes[t] = true
			}
		default:
			return full("shared or unmapped file changed: " + f), nil
		}
	}

	picked := map[string]bool{}
	for _, tf := range tests {
		hit := wantFiles[tf.path]
		for t := range wantTypes {
			hit = hit || tf.uses[t]
		}
		if hit {
			for _, fn := range tf.funcs {
				picked[fn] = true
			}
		}
	}
	if len(picked) == 0 {
		return full("no acceptance tests matched the changes"), nil
	}
	res := Result{Reason: fmt.Sprintf("%d resource types, %d tests", len(wantTypes), len(picked))}
	for fn := range picked {
		res.Tests = append(res.Tests, fn)
	}
	sort.Strings(res.Tests)
	return res, nil
}

// Untested lists registered resource types that no TestAcc test file uses.
func Untested(root string) ([]string, error) {
	reg, err := registered(root)
	if err != nil {
		return nil, err
	}
	tests, err := testFiles(root)
	if err != nil {
		return nil, err
	}
	used := map[string]bool{}
	for _, tf := range tests {
		if len(tf.funcs) == 0 {
			continue
		}
		for t := range tf.uses {
			used[t] = true
		}
	}
	var out []string
	for t := range reg {
		if !used[t] {
			out = append(out, t)
		}
	}
	sort.Strings(out)
	return out, nil
}

func main() {
	root := flag.String("root", ".", "repository root")
	check := flag.Bool("check", false, "list registered resources without an acceptance test and exit 1 if any")
	flag.Parse()

	if *check {
		missing, err := Untested(*root)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		for _, m := range missing {
			fmt.Println(m)
		}
		if len(missing) > 0 {
			os.Exit(1)
		}
		return
	}

	res, err := Analyze(*root, flag.Args())
	if err != nil {
		// never skip tests because the tool failed
		res = full("analysis failed: " + err.Error())
	}
	fmt.Fprintln(os.Stderr, "affected-tests:", map[bool]string{true: "full suite", false: "selective"}[res.Full], "-", res.Reason)
	if res.Full {
		fmt.Println("full=true")
		return
	}
	fmt.Println("full=false")
	fmt.Println("run=" + res.RunRegex())
}
