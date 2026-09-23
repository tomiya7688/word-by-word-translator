package processing

import (
	"testing"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
)

func TestDictionaryStoreSearchesOnlySelectedDictionaries(t *testing.T) {
	dictA := NewMemoryDictionary(
		contracts.DictionaryMetadata{
			ID:             "dict-a",
			SourceLanguage: contracts.LanguageEnglish,
			TargetLanguage: contracts.LanguageJapanese,
		},
		map[string][]string{"apple": {"りんごA"}},
	)
	dictB := NewMemoryDictionary(
		contracts.DictionaryMetadata{
			ID:             "dict-b",
			SourceLanguage: contracts.LanguageEnglish,
			TargetLanguage: contracts.LanguageJapanese,
		},
		map[string][]string{"apple": {"りんごB"}},
	)

	store := NewDictionaryStore(dictA, dictB)
	response := store.LookupBatch(contracts.DictionaryBatchLookupRequest{
		Tokens: []contracts.Token{{
			Surface:    "apple",
			LookupUnit: "apple",
			Kind:       contracts.TokenKindWord,
		}},
		SourceLanguage: contracts.LanguageEnglish,
		TargetLanguage: contracts.LanguageJapanese,
		DictionaryIDs:  []string{"dict-b"},
	})

	if len(response.Tokens) != 1 || len(response.Tokens[0].Dictionaries) != 1 {
		t.Fatalf("response = %#v", response)
	}
	if response.Tokens[0].Dictionaries[0].Metadata.ID != "dict-b" {
		t.Fatalf("dictionary = %#v", response.Tokens[0].Dictionaries[0].Metadata)
	}
}

func TestDictionaryStoreEmptySelectionKeepsBackwardCompatibleAllDictionaries(t *testing.T) {
	dictA := NewMemoryDictionary(
		contracts.DictionaryMetadata{
			ID:             "dict-a",
			SourceLanguage: contracts.LanguageEnglish,
			TargetLanguage: contracts.LanguageJapanese,
		},
		map[string][]string{"apple": {"りんごA"}},
	)
	dictB := NewMemoryDictionary(
		contracts.DictionaryMetadata{
			ID:             "dict-b",
			SourceLanguage: contracts.LanguageEnglish,
			TargetLanguage: contracts.LanguageJapanese,
		},
		map[string][]string{"apple": {"りんごB"}},
	)

	store := NewDictionaryStore(dictA, dictB)
	response := store.LookupBatch(contracts.DictionaryBatchLookupRequest{
		Tokens: []contracts.Token{{
			Surface:    "apple",
			LookupUnit: "apple",
			Kind:       contracts.TokenKindWord,
		}},
		SourceLanguage: contracts.LanguageEnglish,
		TargetLanguage: contracts.LanguageJapanese,
	})

	if len(response.Tokens) != 1 || len(response.Tokens[0].Dictionaries) != 2 {
		t.Fatalf("response = %#v", response)
	}
}
