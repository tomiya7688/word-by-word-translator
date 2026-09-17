package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	dictionaryID      = "ejdict-hand"
	dictionaryName    = "EJDict-hand"
	dictionaryLicense = "CC0-1.0"
)

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
	sourceDir := flag.String("source", "", "path to the EJDict src directory")
	outputDir := flag.String("output", "", "output dictionary package directory")
	revision := flag.String("revision", "", "pinned EJDict git revision")
	flag.Parse()

	if err := run(*sourceDir, *outputDir, *revision); err != nil {
		fmt.Fprintln(os.Stderr, "ejdict-import:", err)
		os.Exit(1)
	}
}

func run(sourceDir, outputDir, revision string) error {
	if strings.TrimSpace(sourceDir) == "" {
		return errors.New("source is required")
	}
	if strings.TrimSpace(outputDir) == "" {
		return errors.New("output is required")
	}
	if strings.TrimSpace(revision) == "" {
		return errors.New("revision is required")
	}

	files, err := filepath.Glob(filepath.Join(sourceDir, "*.txt"))
	if err != nil {
		return fmt.Errorf("find EJDict source files: %w", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no EJDict source files found in %s", sourceDir)
	}
	sort.Strings(files)

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	entriesPath := filepath.Join(outputDir, "entries.tsv")
	entriesFile, err := os.Create(entriesPath)
	if err != nil {
		return fmt.Errorf("create entries.tsv: %w", err)
	}

	count, importErr := importFiles(entriesFile, files)
	closeErr := entriesFile.Close()
	if importErr != nil {
		return importErr
	}
	if closeErr != nil {
		return fmt.Errorf("close entries.tsv: %w", closeErr)
	}
	if count == 0 {
		return errors.New("EJDict import produced no entries")
	}

	manifest := dictionaryManifest{
		SchemaVersion:         1,
		ID:                    dictionaryID,
		Name:                  dictionaryName,
		Source:                "https://github.com/kujirahand/EJDict/tree/" + revision,
		License:               dictionaryLicense,
		SourceLanguage:        "en",
		TargetLanguage:        "ja",
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

	fmt.Printf("ejdict-import: wrote %d entries to %s\n", count, outputDir)
	return nil
}

func importFiles(output *os.File, files []string) (int, error) {
	writer := bufio.NewWriter(output)
	seen := make(map[string]struct{})
	count := 0
	for _, sourcePath := range files {
		file, err := os.Open(sourcePath)
		if err != nil {
			return 0, fmt.Errorf("open %s: %w", sourcePath, err)
		}
		fileCount, importErr := importFile(writer, file, sourcePath, seen)
		closeErr := file.Close()
		if importErr != nil {
			return 0, importErr
		}
		if closeErr != nil {
			return 0, fmt.Errorf("close %s: %w", sourcePath, closeErr)
		}
		count += fileCount
	}
	if err := writer.Flush(); err != nil {
		return 0, fmt.Errorf("write entries.tsv: %w", err)
	}
	return count, nil
}

func importFile(writer *bufio.Writer, file *os.File, sourcePath string, seen map[string]struct{}) (int, error) {
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	lineNumber := 0
	count := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.SplitN(line, "\t", 2)
		if len(fields) != 2 {
			return 0, fmt.Errorf("parse %s line %d: expected headword<TAB>translation", sourcePath, lineNumber)
		}

		translation := strings.TrimSpace(fields[1])
		if translation == "" {
			return 0, fmt.Errorf("parse %s line %d: translation is required", sourcePath, lineNumber)
		}

		aliases := strings.Split(fields[0], ",")
		for _, alias := range aliases {
			headword := strings.TrimSpace(alias)
			if headword == "" {
				return 0, fmt.Errorf("parse %s line %d: headword is required", sourcePath, lineNumber)
			}
			key := strings.ToLower(headword) + "\x00" + translation
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = struct{}{}
			if _, err := fmt.Fprintf(writer, "%s\t%s\n", headword, translation); err != nil {
				return 0, fmt.Errorf("write entries.tsv: %w", err)
			}
			count++
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("read %s: %w", sourcePath, err)
	}
	return count, nil
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
