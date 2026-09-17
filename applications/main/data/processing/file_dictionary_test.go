package processing

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
)

func TestLoadFileDictionaryAndLookup(t *testing.T) {
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "dictionary.json")
	data := "# sample\nHello\tこんにちは\nhello\tやあ\nhello\tこんにちは\nworld\t世界\n"
	if err := os.WriteFile(filepath.Join(dir, "entries.tsv"), []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
	manifest := `{"schema_version":1,"id":"sample","name":"Sample","source":"test","license":"CC0-1.0","source_language":"en","target_language":"ja","commercial_use":true,"noncommercial_use":true,"modification":true,"redistribution":true,"attribution_required":false,"license_notice_required":false,"share_alike":false,"data_file":"entries.tsv","format":"tsv-v1","case_sensitive":false}`
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o600); err != nil {
		t.Fatal(err)
	}

	dictionary, err := LoadFileDictionary(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if dictionary.Metadata().ReleaseTier != contracts.ReleaseTierMIT {
		t.Fatalf("tier=%q", dictionary.Metadata().ReleaseTier)
	}
	entries, err := dictionary.Lookup("HELLO")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("len=%d, want 2", len(entries))
	}
}

func TestManifestReleaseTier(t *testing.T) {
	base := DictionaryManifest{
		CommercialUse:    true,
		NoncommercialUse: true,
		Modification:     true,
		Redistribution:   true,
	}
	if got := base.ReleaseTier(); got != contracts.ReleaseTierMIT {
		t.Fatalf("MIT tier=%q", got)
	}

	base.AttributionRequired = true
	if got := base.ReleaseTier(); got != contracts.ReleaseTierFull {
		t.Fatalf("Full tier with attribution=%q", got)
	}

	base.AttributionRequired = false
	base.Modification = false
	if got := base.ReleaseTier(); got != contracts.ReleaseTierFull {
		t.Fatalf("Full tier with modification restriction=%q", got)
	}

	base.CommercialUse = false
	if got := base.ReleaseTier(); got != contracts.ReleaseTierUnsupported {
		t.Fatalf("Unsupported tier=%q", got)
	}
}

func TestManifestRejectsEscapingDataFile(t *testing.T) {
	base := DictionaryManifest{
		SchemaVersion:  DictionaryManifestSchemaVersion,
		ID:             "x",
		Name:           "x",
		Source:         "test",
		License:        "CC0-1.0",
		SourceLanguage: contracts.LanguageEnglish,
		TargetLanguage: contracts.LanguageJapanese,
		Format:         DictionaryDataFormatTSVV1,
	}
	for _, dataFile := range []string{"../outside.tsv", `..\\outside.tsv`, "/tmp/outside.tsv", `C:\\outside.tsv`} {
		manifest := base
		manifest.DataFile = dataFile
		if err := manifest.Validate(); err == nil {
			t.Fatalf("expected validation error for %q", dataFile)
		}
	}
}

func TestLoadDictionaryManifestRejectsUnknownField(t *testing.T) {
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "dictionary.json")
	manifest := `{"schema_version":1,"id":"sample","name":"Sample","source":"test","license":"CC0-1.0","source_language":"en","target_language":"ja","commercial_use":true,"noncommercial_use":true,"modification":true,"redistribution":true,"data_file":"entries.tsv","format":"tsv-v1","unknown":true}`
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadDictionaryManifest(manifestPath); err == nil {
		t.Fatal("manifest with unknown field was accepted")
	}
}

func TestLoadFileDictionaryRejectsInvalidUTF8(t *testing.T) {
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "dictionary.json")
	manifest := `{"schema_version":1,"id":"sample","name":"Sample","source":"test","license":"CC0-1.0","source_language":"en","target_language":"ja","commercial_use":true,"noncommercial_use":true,"modification":true,"redistribution":true,"data_file":"entries.tsv","format":"tsv-v1"}`
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "entries.tsv"), []byte{'h', 'i', '\t', 0xff, '\n'}, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadFileDictionary(manifestPath); err == nil {
		t.Fatal("invalid UTF-8 dictionary data was accepted")
	}
}

func TestLoadFileDictionaryRejectsMalformedTSV(t *testing.T) {
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "dictionary.json")
	if err := os.WriteFile(filepath.Join(dir, "entries.tsv"), []byte("missing-tab\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	manifest := `{"schema_version":1,"id":"sample","name":"Sample","source":"test","license":"CC0-1.0","source_language":"en","target_language":"ja","commercial_use":true,"noncommercial_use":true,"modification":true,"redistribution":true,"data_file":"entries.tsv","format":"tsv-v1"}`
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadFileDictionary(manifestPath); err == nil {
		t.Fatal("expected malformed TSV error")
	}
}
