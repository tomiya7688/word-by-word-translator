package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	dictionaryID      = "tkg-ja-en"
	dictionaryName    = "TKG Japanese-English Learner's Dictionary"
	dictionaryLicense = "CC0-1.0"
)

type sourceIndex struct {
	Metadata sourceMetadata `json:"metadata"`
	Entries  []sourceEntry  `json:"entries"`
}

type sourceMetadata struct {
	TotalEntries int `json:"total_entries"`
}

type sourceEntry struct {
	Headword string `json:"headword"`
	Reading  string `json:"reading"`
	Gloss    string `json:"gloss"`
}

type dictionaryManifest struct {
	SchemaVersion         int    `json:"schema_version"`
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Source                string `json:"source"`
	License               string `json:"license"`
	SourceLanguage        string `json:"source_language"`
	TargetLanguage        string `json:"target_language"`
	CommercialUse         bool   `json:"commercial_use"`
	NoncommercialUse      bool   `json:"noncommercial_use"`
	Modification          bool   `json:"modification"`
	Redistribution        bool   `json:"redistribution"`
	AttributionRequired   bool   `json:"attribution_required"`
	LicenseNoticeRequired bool   `json:"license_notice_required"`
	ShareAlike            bool   `json:"share_alike"`
	DataFile              string `json:"data_file"`
	Format                string `json:"format"`
	CaseSensitive         bool   `json:"case_sensitive"`
}

func main() {
	sourcePath := flag.String("source", "", "path to TKG entries_index.json")
	outputDir := flag.String("output", "", "output dictionary package directory")
	revision := flag.String("revision", "", "pinned TKG dictionary git revision")
	flag.Parse()

	if err := run(*sourcePath, *outputDir, *revision); err != nil {
		fmt.Fprintln(os.Stderr, "tkg-ja-en-import:", err)
		os.Exit(1)
	}
}

func run(sourcePath, outputDir, revision string) error {
	if strings.TrimSpace(sourcePath) == "" {
		return errors.New("source is required")
	}
	if strings.TrimSpace(outputDir) == "" {
		return errors.New("output is required")
	}
	if strings.TrimSpace(revision) == "" {
		return errors.New("revision is required")
	}

	index, err := readSourceIndex(sourcePath)
	if err != nil {
		return err
	}
	if len(index.Entries) == 0 {
		return errors.New("source index contains no entries")
	}
	if index.Metadata.TotalEntries > 0 && index.Metadata.TotalEntries != len(index.Entries) {
		return fmt.Errorf("source index total_entries=%d but decoded %d entries", index.Metadata.TotalEntries, len(index.Entries))
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	entriesPath := filepath.Join(outputDir, "entries.tsv")
	entriesFile, err := os.Create(entriesPath)
	if err != nil {
		return fmt.Errorf("create entries.tsv: %w", err)
	}
	count, importErr := importIndex(entriesFile, index)
	closeErr := entriesFile.Close()
	if importErr != nil {
		return importErr
	}
	if closeErr != nil {
		return fmt.Errorf("close entries.tsv: %w", closeErr)
	}
	if count == 0 {
		return errors.New("TKG dictionary import produced no entries")
	}

	manifest := dictionaryManifest{
		SchemaVersion:         1,
		ID:                    dictionaryID,
		Name:                  dictionaryName,
		Source:                "https://github.com/tkgally/je-dict-1/tree/" + revision,
		License:               dictionaryLicense,
		SourceLanguage:        "ja",
		TargetLanguage:        "en",
		CommercialUse:         true,
		NoncommercialUse:      true,
		Modification:          true,
		Redistribution:        true,
		AttributionRequired:   false,
		LicenseNoticeRequired: false,
		ShareAlike:            false,
		DataFile:              "entries.tsv",
		Format:                "tsv-v1",
		CaseSensitive:         false,
	}
	if err := writeManifest(filepath.Join(outputDir, "dictionary.json"), manifest); err != nil {
		return err
	}

	fmt.Printf("tkg-ja-en-import: wrote %d lookup entries from %d source entries to %s\n", count, len(index.Entries), outputDir)
	return nil
}

func readSourceIndex(sourcePath string) (sourceIndex, error) {
	file, err := os.Open(sourcePath)
	if err != nil {
		return sourceIndex{}, fmt.Errorf("open source index: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var index sourceIndex
	if err := decoder.Decode(&index); err != nil {
		return sourceIndex{}, fmt.Errorf("parse source index: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return sourceIndex{}, errors.New("source index contains multiple JSON values")
		}
		return sourceIndex{}, fmt.Errorf("parse source index trailing data: %w", err)
	}
	return index, nil
}

func importIndex(output *os.File, index sourceIndex) (int, error) {
	writer := bufio.NewWriter(output)
	seen := make(map[string]struct{})
	count := 0
	for sourceIndex, entry := range index.Entries {
		headword := strings.TrimSpace(entry.Headword)
		reading := strings.TrimSpace(entry.Reading)
		gloss := strings.TrimSpace(entry.Gloss)
		if headword == "" || gloss == "" {
			return 0, fmt.Errorf("source entry %d: headword and gloss are required", sourceIndex+1)
		}
		if err := validateTSVField(headword); err != nil {
			return 0, fmt.Errorf("source entry %d headword: %w", sourceIndex+1, err)
		}
		if err := validateTSVField(gloss); err != nil {
			return 0, fmt.Errorf("source entry %d gloss: %w", sourceIndex+1, err)
		}

		written, err := writeLookupEntry(writer, seen, headword, gloss)
		if err != nil {
			return 0, err
		}
		if written {
			count++
		}

		if reading == "" {
			continue
		}
		if err := validateTSVField(reading); err != nil {
			return 0, fmt.Errorf("source entry %d reading: %w", sourceIndex+1, err)
		}
		written, err = writeLookupEntry(writer, seen, reading, gloss)
		if err != nil {
			return 0, err
		}
		if written {
			count++
		}
	}
	if err := writer.Flush(); err != nil {
		return 0, fmt.Errorf("write entries.tsv: %w", err)
	}
	return count, nil
}

func writeLookupEntry(writer *bufio.Writer, seen map[string]struct{}, lookupUnit, translation string) (bool, error) {
	key := strings.ToLower(lookupUnit) + "\x00" + translation
	if _, exists := seen[key]; exists {
		return false, nil
	}
	seen[key] = struct{}{}
	if _, err := fmt.Fprintf(writer, "%s\t%s\n", lookupUnit, translation); err != nil {
		return false, fmt.Errorf("write entries.tsv: %w", err)
	}
	return true, nil
}

func validateTSVField(value string) error {
	if strings.ContainsAny(value, "\t\r\n") {
		return errors.New("tabs and line breaks are not supported")
	}
	return nil
}

func writeManifest(path string, manifest dictionaryManifest) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create dictionary.json: %w", err)
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encodeErr := encoder.Encode(manifest)
	closeErr := file.Close()
	if encodeErr != nil {
		return fmt.Errorf("write dictionary.json: %w", encodeErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close dictionary.json: %w", closeErr)
	}
	return nil
}
