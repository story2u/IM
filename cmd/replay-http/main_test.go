package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveFixturePath(t *testing.T) {
	got := resolveFixturePath("fixtures/replay/suite.json", "events.json")
	want := filepath.Join("fixtures/replay", "events.json")
	if got != want {
		t.Fatalf("resolveFixturePath() = %q, want %q", got, want)
	}

	got = resolveFixturePath("/tmp/suite/root/suite.json", "/tmp/absolute.json")
	want = "/tmp/absolute.json"
	if got != want {
		t.Fatalf("resolveFixturePath absolute = %q, want %q", got, want)
	}
}

func TestValidationReportLoadsFixtureFiles(t *testing.T) {
	dir := t.TempDir()
	mustWrite := func(name, content string) string {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatalf("write fixture %s: %v", name, err)
		}
		return path
	}
	mustWrite("reference.json", `[{"event":"conversation.message","cursor":"1"}]`)
	mustWrite("go.json", `[{"event":"conversation.message","cursor":"1"}]`)
	suitePath := mustWrite("suite.json", `{"name":"fixture-load","cases":[{"name":"matching events","reference":"reference.json","go":"go.json"}]}`)

	suite, err := loadSuite(suitePath)
	if err != nil {
		t.Fatalf("loadSuite() = %v", err)
	}

	report, err := validationReport(suitePath, suite)
	if err != nil {
		t.Fatalf("validationReport() = %v", err)
	}
	if !report.Match {
		t.Fatalf("validationReport should match for committed fixtures, got %#v", report)
	}
	if report.CaseCount != 1 {
		t.Fatalf("case_count=%d, want 1", report.CaseCount)
	}
	if report.Cases[0].ReferenceCount != 1 || report.Cases[0].GoCount != 1 {
		t.Fatalf("event counts = reference %d go %d, want 1/1", report.Cases[0].ReferenceCount, report.Cases[0].GoCount)
	}
}

func TestLoadSuiteUsesReferenceFixturePath(t *testing.T) {
	dir := t.TempDir()
	referencePath := filepath.Join(dir, "reference.json")
	goPath := filepath.Join(dir, "go.json")
	if err := os.WriteFile(referencePath, []byte(`[{"event":"conversation.message","cursor":"1"}]`), 0o600); err != nil {
		t.Fatalf("write reference fixture: %v", err)
	}
	if err := os.WriteFile(goPath, []byte(`[{"event":"conversation.message","cursor":"1"}]`), 0o600); err != nil {
		t.Fatalf("write go fixture: %v", err)
	}
	suitePath := filepath.Join(dir, "suite.json")
	suiteJSON := `{"name":"reference-suite","cases":[{"name":"reference case","reference":"reference.json","go":"go.json"}]}`
	if err := os.WriteFile(suitePath, []byte(suiteJSON), 0o600); err != nil {
		t.Fatalf("write suite: %v", err)
	}

	suite, err := loadSuite(suitePath)
	if err != nil {
		t.Fatalf("loadSuite() = %v", err)
	}
	if got := suite.Cases[0].referenceEventsPath(); got != "reference.json" {
		t.Fatalf("referenceEventsPath() = %q, want reference.json", got)
	}
}

func TestCompareReportSummarizesMismatch(t *testing.T) {
	dir := t.TempDir()
	mustWrite := func(name, content string) string {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatalf("write fixture %s: %v", name, err)
		}
		return path
	}
	referencePath := mustWrite("reference.json", `[{"event":"conversation.message","cursor":"1","channel":"conversations","payload":{"id":"x"}}]`)
	goPath := mustWrite("go.json", `[{"event":"conversation.message","cursor":"2","channel":"conversations","payload":{"id":"x"}}]`)

	suite := replaySuite{
		Name: "mismatch-check",
		Cases: []suiteCase{
			{
				Name:                "cursor mismatch",
				ReferenceEventsPath: referencePath,
				GoEventsPath:        goPath,
			},
		},
	}

	report, err := compareReport("unused", suite)
	if err != nil {
		t.Fatalf("compareReport() = %v", err)
	}
	if report.Match {
		t.Fatal("report.Match = true, want false")
	}
	if len(report.Cases) != 1 {
		t.Fatalf("len(cases) = %d, want 1", len(report.Cases))
	}
	if report.Cases[0].Match {
		t.Fatal("case should mismatch")
	}
}
