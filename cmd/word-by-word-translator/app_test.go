package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
	dataprocessing "github.com/tomiya7688/word-by-word-translator/applications/main/data/processing"
)

func TestCLIEnglishUnknownWordAndDictionarySelection(t *testing.T) {
	root := t.TempDir()
	dictA := writeDictionaryPackage(t, root, dataprocessing.DictionaryManifest{
		SchemaVersion:    dataprocessing.DictionaryManifestSchemaVersion,
		ID:               "dict-a",
		Name:             "Dictionary A",
		Source:           "https://example.invalid/a",
		License:          "CC0-1.0",
		SourceLanguage:   contracts.LanguageEnglish,
		TargetLanguage:   contracts.LanguageJapanese,
		CommercialUse:    true,
		NoncommercialUse: true,
		Modification:     true,
		Redistribution:   true,
		DataFile:         "entries.tsv",
		Format:           dataprocessing.DictionaryDataFormatTSVV1,
	}, "i\t私\nsaw\t見た\nyesterday\t昨日\n")
	dictB := writeDictionaryPackage(t, root, dataprocessing.DictionaryManifest{
		SchemaVersion:    dataprocessing.DictionaryManifestSchemaVersion,
		ID:               "dict-b",
		Name:             "Dictionary B",
		Source:           "https://example.invalid/b",
		License:          "CC0-1.0",
		SourceLanguage:   contracts.LanguageEnglish,
		TargetLanguage:   contracts.LanguageJapanese,
		CommercialUse:    true,
		NoncommercialUse: true,
		Modification:     true,
		Redistribution:   true,
		DataFile:         "entries.tsv",
		Format:           dataprocessing.DictionaryDataFormatTSVV1,
	}, "saw\t見ました\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{
		"-dictionary", dictA,
		"-dictionary", dictB,
		"-use", "dict-a",
		"-compact",
		"--",
		"I saw qwertymonster yesterday.",
	}, strings.NewReader(""), &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("exit code = %d, stderr = %q", exitCode, stderr.String())
	}
	want := "私 見た \x1b[31mqwertymonster\x1b[0m 昨日 .\n"
	if stdout.String() != want {
		t.Fatalf("stdout = %q, want %q", stdout.String(), want)
	}
	if strings.Contains(stdout.String(), "見ました") {
		t.Fatalf("unselected dictionary result leaked into output: %q", stdout.String())
	}
}

func TestCLIMultipleDictionariesShowsMergedProvenance(t *testing.T) {
	root := t.TempDir()
	dictA := writeDictionaryPackage(t, root, dataprocessing.DictionaryManifest{
		SchemaVersion:    dataprocessing.DictionaryManifestSchemaVersion,
		ID:               "dict-a",
		Name:             "Dictionary A",
		Source:           "https://example.invalid/a",
		License:          "CC0-1.0",
		SourceLanguage:   contracts.LanguageEnglish,
		TargetLanguage:   contracts.LanguageJapanese,
		CommercialUse:    true,
		NoncommercialUse: true,
		Modification:     true,
		Redistribution:   true,
		DataFile:         "entries.tsv",
		Format:           dataprocessing.DictionaryDataFormatTSVV1,
	}, "apple\tりんご\n")
	dictB := writeDictionaryPackage(t, root, dataprocessing.DictionaryManifest{
		SchemaVersion:    dataprocessing.DictionaryManifestSchemaVersion,
		ID:               "dict-b",
		Name:             "Dictionary B",
		Source:           "https://example.invalid/b",
		License:          "CC0-1.0",
		SourceLanguage:   contracts.LanguageEnglish,
		TargetLanguage:   contracts.LanguageJapanese,
		CommercialUse:    true,
		NoncommercialUse: true,
		Modification:     true,
		Redistribution:   true,
		DataFile:         "entries.tsv",
		Format:           dataprocessing.DictionaryDataFormatTSVV1,
	}, "apple\tりんご\napple\t林檎\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{
		"-dictionary", dictA,
		"-dictionary", dictB,
		"apple",
	}, strings.NewReader(""), &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("exit code = %d, stderr = %q", exitCode, stderr.String())
	}
	output := stdout.String()
	for _, expected := range []string{
		"[x] Dictionary A (dict-a)",
		"[x] Dictionary B (dict-b)",
		"{りんご [dict-a,dict-b] / 林檎 [dict-b]}",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("stdout missing %q:\n%s", expected, output)
		}
	}
}

func TestCLIJapaneseToEnglishReadsStdin(t *testing.T) {
	root := t.TempDir()
	dictionary := writeDictionaryPackage(t, root, dataprocessing.DictionaryManifest{
		SchemaVersion:    dataprocessing.DictionaryManifestSchemaVersion,
		ID:               "ja-en",
		Name:             "Japanese English",
		Source:           "https://example.invalid/ja-en",
		License:          "CC0-1.0",
		SourceLanguage:   contracts.LanguageJapanese,
		TargetLanguage:   contracts.LanguageEnglish,
		CommercialUse:    true,
		NoncommercialUse: true,
		Modification:     true,
		Redistribution:   true,
		DataFile:         "entries.tsv",
		Format:           dataprocessing.DictionaryDataFormatTSVV1,
	}, "食べる\tto eat\nます\tpolite\nた\tpast\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{
		"-from", "ja",
		"-to", "en",
		"-dictionary", dictionary,
	}, strings.NewReader("食べました\n"), &stdout, &stderr)

	if exitCode != 0 {
		t.Fatalf("exit code = %d, stderr = %q", exitCode, stderr.String())
	}
	output := stdout.String()
	for _, expected := range []string{
		"Direction: ja -> en",
		"[x] Japanese English (ja-en)",
		"to eat [ja-en]",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("stdout missing %q:\n%s", expected, output)
		}
	}
}

func TestCLIRejectsDuplicateDictionaryIDs(t *testing.T) {
	root := t.TempDir()
	first := writeDictionaryPackageAt(t, filepath.Join(root, "first"), dataprocessing.DictionaryManifest{
		SchemaVersion:    dataprocessing.DictionaryManifestSchemaVersion,
		ID:               "duplicate",
		Name:             "First",
		Source:           "https://example.invalid/first",
		License:          "CC0-1.0",
		SourceLanguage:   contracts.LanguageEnglish,
		TargetLanguage:   contracts.LanguageJapanese,
		CommercialUse:    true,
		NoncommercialUse: true,
		Modification:     true,
		Redistribution:   true,
		DataFile:         "entries.tsv",
		Format:           dataprocessing.DictionaryDataFormatTSVV1,
	}, "apple\tりんご\n")
	second := writeDictionaryPackageAt(t, filepath.Join(root, "second"), dataprocessing.DictionaryManifest{
		SchemaVersion:    dataprocessing.DictionaryManifestSchemaVersion,
		ID:               "duplicate",
		Name:             "Second",
		Source:           "https://example.invalid/second",
		License:          "CC0-1.0",
		SourceLanguage:   contracts.LanguageEnglish,
		TargetLanguage:   contracts.LanguageJapanese,
		CommercialUse:    true,
		NoncommercialUse: true,
		Modification:     true,
		Redistribution:   true,
		DataFile:         "entries.tsv",
		Format:           dataprocessing.DictionaryDataFormatTSVV1,
	}, "apple\t林檎\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{
		"-dictionary", first,
		"-dictionary", second,
		"apple",
	}, strings.NewReader(""), &stdout, &stderr)

	if exitCode == 0 {
		t.Fatalf("exit code = 0, stdout = %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "duplicate dictionary id") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func writeDictionaryPackage(t *testing.T, root string, manifest dataprocessing.DictionaryManifest, entries string) string {
	t.Helper()
	return writeDictionaryPackageAt(t, filepath.Join(root, manifest.ID), manifest, entries)
}

func writeDictionaryPackageAt(t *testing.T, packageDir string, manifest dataprocessing.DictionaryManifest, entries string) string {
	t.Helper()
	if err := os.MkdirAll(packageDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(packageDir, manifest.DataFile), []byte(entries), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(packageDir, "dictionary.json")
	file, err := os.Create(manifestPath)
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
	return manifestPath
}
