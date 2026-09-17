package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunExpandsAliasesAndWritesMITManifest(t *testing.T) {
	sourceDir := t.TempDir()
	outputDir := t.TempDir()
	input := "apple\tリンゴ、林檎\ncenter, centre\t中央、中心\ncenter\t中央、中心\n"
	if err := os.WriteFile(filepath.Join(sourceDir, "a.txt"), []byte(input), 0o600); err != nil {
		t.Fatal(err)
	}

	const revision = "0123456789abcdef"
	if err := run(sourceDir, outputDir, revision); err != nil {
		t.Fatal(err)
	}

	entries, err := os.ReadFile(filepath.Join(outputDir, "entries.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	gotEntries := string(entries)
	for _, want := range []string{
		"apple\tリンゴ、林檎\n",
		"center\t中央、中心\n",
		"centre\t中央、中心\n",
	} {
		if !strings.Contains(gotEntries, want) {
			t.Fatalf("entries.tsv does not contain %q:\n%s", want, gotEntries)
		}
	}
	if strings.Count(gotEntries, "center\t中央、中心\n") != 1 {
		t.Fatalf("duplicate center entry was not removed:\n%s", gotEntries)
	}

	rawManifest, err := os.ReadFile(filepath.Join(outputDir, "dictionary.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest dictionaryManifest
	if err := json.Unmarshal(rawManifest, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.ID != dictionaryID || manifest.License != dictionaryLicense {
		t.Fatalf("manifest identity = %#v", manifest)
	}
	if !manifest.CommercialUse || !manifest.NoncommercialUse || !manifest.Modification || !manifest.Redistribution {
		t.Fatalf("manifest usage flags = %#v", manifest)
	}
	if manifest.AttributionRequired || manifest.LicenseNoticeRequired || manifest.ShareAlike {
		t.Fatalf("manifest obligations = %#v", manifest)
	}
	if !strings.HasSuffix(manifest.Source, revision) {
		t.Fatalf("manifest source = %q", manifest.Source)
	}
}

func TestRunRequiresSourceFiles(t *testing.T) {
	err := run(t.TempDir(), t.TempDir(), "revision")
	if err == nil || !strings.Contains(err.Error(), "no EJDict source files") {
		t.Fatalf("error = %v", err)
	}
}
