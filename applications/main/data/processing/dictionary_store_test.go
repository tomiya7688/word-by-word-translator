package processing

import (
	"errors"
	"testing"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
)

type failingDictionary struct {
	metadata contracts.DictionaryMetadata
	err      error
}

func (d *failingDictionary) Metadata() contracts.DictionaryMetadata {
	return d.metadata
}

func (d *failingDictionary) Lookup(string) ([]contracts.DictionaryEntry, error) {
	return nil, d.err
}

func TestDictionaryStoreKeepsOtherResultsWhenOneDictionaryFails(t *testing.T) {
	sourceLanguage := contracts.LanguageEnglish
	targetLanguage := contracts.LanguageJapanese

	broken := &failingDictionary{
		metadata: contracts.DictionaryMetadata{
			ID:             "broken",
			SourceLanguage: sourceLanguage,
			TargetLanguage: targetLanguage,
		},
		err: errors.New("lookup failed"),
	}
	working := NewMemoryDictionary(
		contracts.DictionaryMetadata{
			ID:             "working",
			SourceLanguage: sourceLanguage,
			TargetLanguage: targetLanguage,
		},
		map[string][]string{"apple": {"りんご"}},
	)

	store := NewDictionaryStore(broken, working)
	response := store.LookupBatch(contracts.DictionaryBatchLookupRequest{
		Tokens: []contracts.Token{{
			Surface:    "apple",
			LookupUnit: "apple",
			Kind:       contracts.TokenKindWord,
		}},
		SourceLanguage: sourceLanguage,
		TargetLanguage: targetLanguage,
	})

	if len(response.Tokens) != 1 {
		t.Fatalf("token results = %d, want 1", len(response.Tokens))
	}
	dictionaries := response.Tokens[0].Dictionaries
	if len(dictionaries) != 2 {
		t.Fatalf("dictionary results = %d, want 2", len(dictionaries))
	}
	if dictionaries[0].Metadata.ID != "broken" || dictionaries[0].Error == "" {
		t.Fatalf("broken result = %#v", dictionaries[0])
	}
	if dictionaries[1].Metadata.ID != "working" ||
		len(dictionaries[1].Entries) != 1 ||
		dictionaries[1].Entries[0].Translation != "りんご" {
		t.Fatalf("working result = %#v", dictionaries[1])
	}
}
