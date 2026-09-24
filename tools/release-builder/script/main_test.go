package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	dataprocessing "github.com/tomiya7688/word-by-word-translator/applications/main/data/processing"
)

func TestRunBuildsMITAndFullReleaseComposition(t *testing.T) {
	root := t.TempDir()
	dictionaries := filepath.Join(root, "dictionaries")
	output := filepath.Join(root, "releases")
	appLicense := filepath.Join(root, "APP_LICENSE")
	analyzerNotices := filepath.Join(root, "analyzers")

	writeFile(t, appLicense, "application license")
	writeAnalyzerNoticeFixture(t, analyzerNotices)

	writeDictionaryFixture(t, dictionaries, dataprocessing.DictionaryManifest{
		SchemaVersion:    1,
		ID:               "mit-dict",
		Name:             "MIT Dictionary",
		Source:           "https://example.invalid/mit",
		License:          "CC0-1.0",
		SourceLanguage:   "en",
		TargetLanguage:   "ja",
		CommercialUse:    true,
		NoncommercialUse: true,
		Modification:     true,
		Redistribution:   true,
		DataFile:         "entries.tsv",
		Format:           "tsv-v1",
	}, nil)
	writeDictionaryFixture(t, dictionaries, dataprocessing.DictionaryManifest{
		SchemaVersion:         1,
		ID:                    "full-dict",
		Name:                  "Full Dictionary",
		Source:                "https://example.invalid/full",
		License:               "Example-Attribution",
		SourceLanguage:        "en",
		TargetLanguage:        "ja",
		CommercialUse:         true,
		NoncommercialUse:      true,
		Modification:          true,
		Redistribution:        true,
		AttributionRequired:   true,
		LicenseNoticeRequired: true,
		DataFile:              "entries.tsv",
		Format:                "tsv-v1",
	}, map[string]string{
		attributionFile: "attribution",
		licenseFile:     "dictionary license",
	})
	writeDictionaryFixture(t, dictionaries, dataprocessing.DictionaryManifest{
		SchemaVersion:    1,
		ID:               "unsupported-dict",
		Name:             "Unsupported Dictionary",
		Source:           "https://example.invalid/unsupported",
		License:          "NC",
		SourceLanguage:   "en",
		TargetLanguage:   "ja",
		CommercialUse:    false,
		NoncommercialUse: true,
		Modification:     true,
		Redistribution:   true,
		DataFile:         "entries.tsv",
		Format:           "tsv-v1",
	}, nil)

	if err := run(dictionaries, output, appLicense, analyzerNotices, "all"); err != nil {
		t.Fatal(err)
	}

	mustExist(t, filepath.Join(output, "mit", "LICENSE"))
	mustExist(t, filepath.Join(output, "mit", "dictionaries", "mit-dict", "dictionary.json"))
	mustNotExist(t, filepath.Join(output, "mit", "dictionaries", "full-dict"))
	mustNotExist(t, filepath.Join(output, "mit", "dictionaries", "unsupported-dict"))
	mustNotExist(t, filepath.Join(output, "mit", "licenses", "analyzers"))

	mustExist(t, filepath.Join(output, "full", "dictionaries", "mit-dict", "dictionary.json"))
	mustExist(t, filepath.Join(output, "full", "dictionaries", "full-dict", "dictionary.json"))
	mustExist(t, filepath.Join(output, "full", "dictionaries", "full-dict", attributionFile))
	mustExist(t, filepath.Join(output, "full", "dictionaries", "full-dict", licenseFile))
	mustNotExist(t, filepath.Join(output, "full", "dictionaries", "unsupported-dict"))
	mustExist(t, filepath.Join(output, "full", "licenses", "analyzers", "kagome-dict-ipa", "NOTICE.txt"))

	mitManifest := loadReleaseManifest(t, filepath.Join(output, "mit", "release.json"))
	if mitManifest.ReleaseTier != releaseMIT || len(mitManifest.Dictionaries) != 1 || len(mitManifest.Analyzers) != 0 {
		t.Fatalf("MIT release manifest = %#v", mitManifest)
	}

	fullManifest := loadReleaseManifest(t, filepath.Join(output, "full", "release.json"))
	if fullManifest.ReleaseTier != releaseFull || len(fullManifest.Dictionaries) != 2 || len(fullManifest.Analyzers) != 3 {
		t.Fatalf("Full release manifest = %#v", fullManifest)
	}

	mitThirdParty := readFile(t, filepath.Join(output, "mit", "THIRD_PARTY_LICENSES.md"))
	if strings.Contains(mitThirdParty, "Full Dictionary") || strings.Contains(mitThirdParty, "kagome") {
		t.Fatalf("MIT third-party file contains Full-only components:\n%s", mitThirdParty)
	}
	fullThirdParty := readFile(t, filepath.Join(output, "full", "THIRD_PARTY_LICENSES.md"))
	if !strings.Contains(fullThirdParty, "Full Dictionary") || !strings.Contains(fullThirdParty, "kagome-dict-ipa") {
		t.Fatalf("Full third-party file is incomplete:\n%s", fullThirdParty)
	}
}

func TestFullReleaseFailsWhenRequiredDictionaryNoticeIsMissing(t *testing.T) {
	root := t.TempDir()
	dictionaries := filepath.Join(root, "dictionaries")
	output := filepath.Join(root, "releases")
	appLicense := filepath.Join(root, "APP_LICENSE")
	analyzerNotices := filepath.Join(root, "analyzers")

	writeFile(t, appLicense, "application license")
	writeAnalyzerNoticeFixture(t, analyzerNotices)
	writeDictionaryFixture(t, dictionaries, dataprocessing.DictionaryManifest{
		SchemaVersion:         1,
		ID:                    "full-dict",
		Name:                  "Full Dictionary",
		Source:                "https://example.invalid/full",
		License:               "Example-License",
		SourceLanguage:        "en",
		TargetLanguage:        "ja",
		CommercialUse:         true,
		NoncommercialUse:      true,
		Modification:          true,
		Redistribution:        true,
		LicenseNoticeRequired: true,
		DataFile:              "entries.tsv",
		Format:                "tsv-v1",
	}, nil)

	err := run(dictionaries, output, appLicense, analyzerNotices, "full")
	if err == nil || !strings.Contains(err.Error(), licenseFile) {
		t.Fatalf("error = %v, want missing %s", err, licenseFile)
	}
}

func TestReleaseBuilderRejectsUnsafeDictionaryID(t *testing.T) {
	root := t.TempDir()
	dictionaries := filepath.Join(root, "dictionaries")
	packageDir := filepath.Join(dictionaries, "unsafe")
	if err := os.MkdirAll(packageDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(packageDir, "entries.tsv"), "word\ttranslation\n")
	manifest := dataprocessing.DictionaryManifest{
		SchemaVersion:    1,
		ID:               "../escape",
		Name:             "Unsafe",
		Source:           "https://example.invalid",
		License:          "CC0-1.0",
		SourceLanguage:   "en",
		TargetLanguage:   "ja",
		CommercialUse:    true,
		NoncommercialUse: true,
		Modification:     true,
		Redistribution:   true,
		DataFile:         "entries.tsv",
		Format:           "tsv-v1",
	}
	writeManifest(t, filepath.Join(packageDir, dictionaryManifest), manifest)

	_, err := discoverDictionaries(dictionaries)
	if err == nil || !strings.Contains(err.Error(), "not safe") {
		t.Fatalf("error = %v", err)
	}
}

func writeDictionaryFixture(t *testing.T, root string, manifest dataprocessing.DictionaryManifest, extraFiles map[string]string) {
	t.Helper()
	packageDir := filepath.Join(root, manifest.ID)
	if err := os.MkdirAll(packageDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(packageDir, manifest.DataFile), "word\ttranslation\n")
	writeManifest(t, filepath.Join(packageDir, dictionaryManifest), manifest)
	for name, content := range extraFiles {
		writeFile(t, filepath.Join(packageDir, name), content)
	}
}

func writeAnalyzerNoticeFixture(t *testing.T, root string) {
	t.Helper()
	for _, analyzer := range fullReleaseAnalyzers {
		for _, releasePath := range analyzer.NoticeFiles {
			relative := strings.TrimPrefix(releasePath, "licenses/analyzers/")
			writeFile(t, filepath.Join(root, filepath.FromSlash(relative)), "notice")
		}
	}
}

func writeManifest(t *testing.T, destination string, manifest dataprocessing.DictionaryManifest) {
	t.Helper()
	file, err := os.Create(destination)
	if err != nil {
		t.Fatal(err)
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(manifest); err != nil {
		file.Close()
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}

func writeFile(t *testing.T, destination, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(destination, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustExist(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
}

func mustNotExist(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("%s unexpectedly exists (err=%v)", path, err)
	}
}

func loadReleaseManifest(t *testing.T, path string) releaseManifest {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	var manifest releaseManifest
	if err := json.NewDecoder(file).Decode(&manifest); err != nil {
		t.Fatal(err)
	}
	return manifest
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}
