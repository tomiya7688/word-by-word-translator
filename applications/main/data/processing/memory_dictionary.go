package processing

import "github.com/tomiya7688/word-by-word-translator/applications/main/contracts"

type MemoryDictionary struct {
	metadata contracts.DictionaryMetadata
	entries  map[string][]contracts.DictionaryEntry
}

func NewMemoryDictionary(metadata contracts.DictionaryMetadata, values map[string][]string) *MemoryDictionary {
	entries := make(map[string][]contracts.DictionaryEntry, len(values))
	for headword, translations := range values {
		items := make([]contracts.DictionaryEntry, 0, len(translations))
		for _, translation := range translations {
			items = append(items, contracts.DictionaryEntry{Headword: headword, Translation: translation})
		}
		entries[headword] = items
	}
	return &MemoryDictionary{metadata: metadata, entries: entries}
}

func (d *MemoryDictionary) Metadata() contracts.DictionaryMetadata {
	return d.metadata
}

func (d *MemoryDictionary) Lookup(lookupUnit string) ([]contracts.DictionaryEntry, error) {
	entries := d.entries[lookupUnit]
	copied := make([]contracts.DictionaryEntry, len(entries))
	copy(copied, entries)
	return copied, nil
}
