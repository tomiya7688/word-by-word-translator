package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunBuildsDictionaryPackage(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "entries_index.json")
	outputDir := filepath.Join(dir, "output")
	source := `{
  "metadata": {"description": "test", "total_entries": 3},
  "entries": [
    {"id": "1", "headword": "愛", "reading": "あい", "gloss": "love", "ignored": true},
    {"id": "2", "headword": "藍", "reading": "あい", "gloss": "indigo"},
    {"id": "3", "headword": "食べる", "reading": "たべる", "gloss": "to eat"}
  ]
}`
	if err := os.WriteFile(sourcePath, []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := run(sourcePath, outputDir, "test-revision"); err != nil {
		t.Fatal(err)
	}

	entries, err := os.ReadFile(filepath.Join(outputDir, "entries.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	want := strings.Join([]string{
		"愛\tlove",
		"あい\tlove",
		"藍\tindigo",
		"あい\tindigo",
		"食べる\tto eat",
		"たべる\tto eat",
		"",
	}, "\n")
	if string(entries) != want {
		t.Fatalf("entries.tsv = %q, want %q", string(entries), want)
	}

	rawManifest, err := os.ReadFile(filepath.Join(outputDir, "dictionary.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest dictionaryManifest
	if err := json.Unmarshal(rawManifest, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.ID != dictionaryID || manifest.SourceLanguage != "ja" || manifest.TargetLanguage != "en" {
		t.Fatalf("manifest = %#v", manifest)
	}
	if !manifest.CommercialUse || !manifest.Modification || !manifest.Redistribution || manifest.AttributionRequired {
		t.Fatalf("license flags = %#v", manifest)
	}
}

func TestRunRejectsSourceCountMismatch(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "entries_index.json")
	source := `{"metadata":{"total_entries":2},"entries":[{"headword":"愛","reading":"あい","gloss":"love"}]}`
	if err := os.WriteFile(sourcePath, []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}

	err := run(sourcePath, filepath.Join(dir, "output"), "test-revision")
	if err == nil || !strings.Contains(err.Error(), "total_entries") {
		t.Fatalf("error = %v, want total_entries mismatch", err)
	}
}

func TestRunDeduplicatesEquivalentLookupRows(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "entries_index.json")
	outputDir := filepath.Join(dir, "output")
	source := `{"metadata":{"total_entries":2},"entries":[{"headword":"かな","reading":"かな","gloss":"kana"},{"headword":"かな","reading":"かな","gloss":"kana"}]}`
	if err := os.WriteFile(sourcePath, []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := run(sourcePath, outputDir, "test-revision"); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadFile(filepath.Join(outputDir, "entries.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	if string(entries) != "かな\tkana\n" {
		t.Fatalf("entries.tsv = %q", string(entries))
	}
}
