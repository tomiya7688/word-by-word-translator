package integration

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestReleaseBuilderCLIProducesMITAndFullBundles(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(filename), "../.."))
	fixtureRoot := t.TempDir()
	dictionaries := filepath.Join(fixtureRoot, "dictionaries")
	output := filepath.Join(fixtureRoot, "releases")
	appLicense := filepath.Join(fixtureRoot, "LICENSE")
	analyzerNotices := filepath.Join(fixtureRoot, "analyzers")

	if err := os.WriteFile(appLicense, []byte("application license"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, relative := range []string{
		"kagome/LICENSE.txt",
		"kagome-dict/LICENSE.txt",
		"kagome-dict-ipa/LICENSE.txt",
		"kagome-dict-ipa/NOTICE.txt",
	} {
		path := filepath.Join(analyzerNotices, filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("notice"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	packageDir := filepath.Join(dictionaries, "mit-dict")
	if err := os.MkdirAll(packageDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(packageDir, "entries.tsv"), []byte("apple\tりんご\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifest := map[string]any{
		"schema_version":          1,
		"id":                      "mit-dict",
		"name":                    "MIT Dictionary",
		"source":                  "https://example.invalid/mit",
		"license":                 "CC0-1.0",
		"source_language":         "en",
		"target_language":         "ja",
		"commercial_use":          true,
		"noncommercial_use":       true,
		"modification":            true,
		"redistribution":          true,
		"attribution_required":    false,
		"license_notice_required": false,
		"share_alike":             false,
		"data_file":               "entries.tsv",
		"format":                  "tsv-v1",
		"case_sensitive":          false,
	}
	file, err := os.Create(filepath.Join(packageDir, "dictionary.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(file).Encode(manifest); err != nil {
		file.Close()
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	command := exec.Command(
		"go", "run", "./tools/release-builder/script",
		"-dictionaries", dictionaries,
		"-output", output,
		"-app-license", appLicense,
		"-analyzer-notices", analyzerNotices,
		"-tier", "all",
	)
	command.Dir = root
	if combined, err := command.CombinedOutput(); err != nil {
		t.Fatalf("release builder failed: %v\n%s", err, combined)
	}

	for _, path := range []string{
		filepath.Join(output, "mit", "release.json"),
		filepath.Join(output, "mit", "dictionaries", "mit-dict", "entries.tsv"),
		filepath.Join(output, "full", "release.json"),
		filepath.Join(output, "full", "licenses", "analyzers", "kagome-dict-ipa", "NOTICE.txt"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("%s: %v", path, err)
		}
	}

	fullThirdParty, err := os.ReadFile(filepath.Join(output, "full", "THIRD_PARTY_LICENSES.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(fullThirdParty), "kagome-dict-ipa") {
		t.Fatalf("Full THIRD_PARTY_LICENSES is incomplete:\n%s", fullThirdParty)
	}
}
