package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func write(t *testing.T, root, rel, body string) {
	t.Helper()
	p := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func fixture(t *testing.T) string {
	root := t.TempDir()
	write(t, root, "sumologic/provider.go", `package sumologic
func Provider() {
	_ = Resources{
		ResourcesMap: map[string]*R{
			"sumologic_collector": resourceCollector(),
			"sumologic_source":    resourceSource(),
			"sumologic_lonely":    resourceLonely(),
		},
		DataSourcesMap: map[string]*R{
			"sumologic_collector": dataCollector(),
		},
	}
}`)
	write(t, root, "sumologic/resource_collector.go", "package sumologic\nfunc resourceCollector() {}\n")
	write(t, root, "sumologic/data_source_collector.go", "package sumologic\nfunc dataCollector() {}\n")
	write(t, root, "sumologic/resource_source.go", "package sumologic\nfunc resourceSource() {}\n")
	write(t, root, "sumologic/resource_lonely.go", "package sumologic\nfunc resourceLonely() {}\n")
	write(t, root, "sumologic/client.go", "package sumologic\nfunc helper() {}\n")
	write(t, root, "sumologic/resource_collector_test.go", "package sumologic\nfunc TestAccCollector_basic() {}\nvar _ = `resource \"sumologic_collector\" \"c\" {}`\n")
	write(t, root, "sumologic/resource_source_test.go", "package sumologic\nfunc TestAccSource_basic() {}\nfunc TestUnit() {}\nvar _ = `resource \"sumologic_collector\" \"c\" {}\nresource \"sumologic_source\" \"s\" {}`\n")
	return root
}

func TestChangedResourceRunsItsDependents(t *testing.T) {
	res, err := Analyze(fixture(t), []string{"sumologic/resource_collector.go"})
	if err != nil || res.Full {
		t.Fatalf("expected selective run, got %+v, %v", res, err)
	}
	if want := []string{"TestAccCollector_basic", "TestAccSource_basic"}; !reflect.DeepEqual(res.Tests, want) {
		t.Errorf("got %v, want %v", res.Tests, want)
	}
}

func TestLeafResourceRunsOnlyItself(t *testing.T) {
	res, _ := Analyze(fixture(t), []string{"sumologic/resource_source.go"})
	if res.Full || !reflect.DeepEqual(res.Tests, []string{"TestAccSource_basic"}) {
		t.Errorf("got %+v", res)
	}
}

func TestChangedTestFileRunsItsTests(t *testing.T) {
	res, _ := Analyze(fixture(t), []string{"sumologic/resource_source_test.go"})
	if res.Full || !reflect.DeepEqual(res.Tests, []string{"TestAccSource_basic"}) {
		t.Errorf("got %+v", res)
	}
}

func TestFallsBackToFull(t *testing.T) {
	root := fixture(t)
	for name, changed := range map[string][]string{
		"shared file":           {"sumologic/client.go"},
		"provider.go":           {"sumologic/provider.go"},
		"go.mod":                {"go.mod"},
		"root go file":          {"main.go"},
		"no files":              {},
		"resource with no test": {"sumologic/resource_lonely.go"},
		"mixed with shared":     {"sumologic/resource_source.go", "sumologic/client.go"},
	} {
		res, err := Analyze(root, changed)
		if err != nil || !res.Full {
			t.Errorf("%s: expected full suite, got %+v, %v", name, res, err)
		}
	}
}

func TestDocsOnlyChangeDoesNotForceFull(t *testing.T) {
	res, _ := Analyze(fixture(t), []string{"website/docs/r/x.markdown", "sumologic/resource_source.go"})
	if res.Full {
		t.Errorf("docs file should be ignored, got %+v", res)
	}
}

func TestUntested(t *testing.T) {
	got, err := Untested(fixture(t))
	if err != nil || !reflect.DeepEqual(got, []string{"sumologic_lonely"}) {
		t.Errorf("got %v, %v", got, err)
	}
}
