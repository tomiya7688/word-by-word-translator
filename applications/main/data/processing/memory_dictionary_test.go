package processing

import (
	"testing"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
)

func TestDictionaryStoreFiltersLanguagePair(t *testing.T) {
	enja := NewMemoryDictionary(contracts.DictionaryMetadata{
		ID: "enja", SourceLanguage: contracts.LanguageEnglish, TargetLanguage: contracts.LanguageJapanese,
	}, map[string][]string{"hello": {"こんにちは"}})
	jaen := NewMemoryDictionary(contracts.DictionaryMetadata{
		ID: "jaen", SourceLanguage: contracts.LanguageJapanese, TargetLanguage: contracts.LanguageEnglish,
	}, map[string][]string{"hello": {"should-not-match"}})
	store := NewDictionaryStore(enja, jaen)

	response := store.LookupBatch(contracts.DictionaryBatchLookupRequest{
		Tokens:         []contracts.Token{{Surface: "hello", LookupUnit: "hello", Kind: contracts.TokenKindWord}},
		SourceLanguage: contracts.LanguageEnglish,
		TargetLanguage: contracts.LanguageJapanese,
	})
	if len(response.Tokens) != 1 || len(response.Tokens[0].Dictionaries) != 1 {
		t.Fatalf("response = %#v, want exactly one matching dictionary", response)
	}
	if response.Tokens[0].Dictionaries[0].Metadata.ID != "enja" {
		t.Fatalf("dictionary = %q, want enja", response.Tokens[0].Dictionaries[0].Metadata.ID)
	}
}
