package processing

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
)

type FileDictionary struct {
	metadata      contracts.DictionaryMetadata
	caseSensitive bool
	entries       map[string][]contracts.DictionaryEntry
}

func LoadFileDictionary(manifestPath string) (*FileDictionary, error) {
	manifest, err := LoadDictionaryManifest(manifestPath)
	if err != nil {
		return nil, err
	}
	entries, err := loadTSVV1(manifest.ResolveDataPath(manifestPath), manifest.CaseSensitive)
	if err != nil {
		return nil, err
	}
	return &FileDictionary{
		metadata:      manifest.Metadata(),
		caseSensitive: manifest.CaseSensitive,
		entries:       entries,
	}, nil
}

func (d *FileDictionary) Metadata() contracts.DictionaryMetadata {
	return d.metadata
}

func (d *FileDictionary) Lookup(lookupUnit string) ([]contracts.DictionaryEntry, error) {
	key := normalizeDictionaryKey(lookupUnit, d.caseSensitive)
	entries := d.entries[key]
	copied := make([]contracts.DictionaryEntry, len(entries))
	copy(copied, entries)
	return copied, nil
}

func loadTSVV1(dataPath string, caseSensitive bool) (map[string][]contracts.DictionaryEntry, error) {
	file, err := os.Open(dataPath)
	if err != nil {
		return nil, fmt.Errorf("open dictionary data: %w", err)
	}
	defer file.Close()

	entries := make(map[string][]contracts.DictionaryEntry)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.SplitN(line, "\t", 2)
		if len(fields) != 2 {
			return nil, fmt.Errorf("parse dictionary data line %d: expected headword<TAB>translation", lineNumber)
		}
		headword := strings.TrimSpace(fields[0])
		translation := strings.TrimSpace(fields[1])
		if headword == "" || translation == "" {
			return nil, fmt.Errorf("parse dictionary data line %d: headword and translation are required", lineNumber)
		}
		key := normalizeDictionaryKey(headword, caseSensitive)
		entries[key] = append(entries[key], contracts.DictionaryEntry{Headword: headword, Translation: translation})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read dictionary data: %w", err)
	}
	return entries, nil
}

func normalizeDictionaryKey(value string, caseSensitive bool) string {
	value = strings.TrimSpace(value)
	if caseSensitive {
		return value
	}
	return strings.ToLower(value)
}
