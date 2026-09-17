package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
	dataprocessing "github.com/tomiya7688/word-by-word-translator/applications/main/data/processing"
)

func TestTKGJapaneseEnglishImportOutputLoadsAsMITDictionary(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(filename), "../.."))
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "entries_index.json")
	outputDir := filepath.Join(dir, "output")

	source := `{
  "metadata": {"total_entries": 3},
  "entries": [
    {"headword": "愛", "reading": "あい", "gloss": "love"},
    {"headword": "藍", "reading": "あい", "gloss": "indigo"},
    {"headword": "食べる", "reading": "たべる", "gloss": "to eat"}
  ]
}`
	if err := os.WriteFile(sourcePath, []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}

	command := exec.Command(
		"go", "run", "./tools/tkg-ja-en-import/script",
		"-source", sourcePath,
		"-output", outputDir,
		"-revision", "test-revision",
	)
	command.Dir = root
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("TKG importer failed: %v\n%s", err, output)
	}

	dictionary, err := dataprocessing.LoadFileDictionary(filepath.Join(outputDir, "dictionary.json"))
	if err != nil {
		t.Fatal(err)
	}
	metadata := dictionary.Metadata()
	if metadata.ID != "tkg-ja-en" || metadata.ReleaseTier != contracts.ReleaseTierMIT {
		t.Fatalf("metadata = %#v", metadata)
	}
	if metadata.SourceLanguage != contracts.LanguageJapanese || metadata.TargetLanguage != contracts.LanguageEnglish {
		t.Fatalf("language pair = %s -> %s", metadata.SourceLanguage, metadata.TargetLanguage)
	}

	entries, err := dictionary.Lookup("食べる")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Translation != "to eat" {
		t.Fatalf("食べる entries = %#v", entries)
	}

	homophones, err := dictionary.Lookup("あい")
	if err != nil {
		t.Fatal(err)
	}
	if len(homophones) != 2 || homophones[0].Translation != "love" || homophones[1].Translation != "indigo" {
		t.Fatalf("あい entries = %#v", homophones)
	}
}
