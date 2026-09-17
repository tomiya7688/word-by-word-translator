package processing

import (
	"reflect"
	"testing"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
)

func TestMergeDeduplicatesTranslationsAndKeepsDictionarySources(t *testing.T) {
	processor := NewTranslationProcessor(NewEnglishTokenizer())
	tokens := []contracts.Token{{Surface: "saw", LookupUnit: "saw", Kind: contracts.TokenKindWord}}
	lookups := contracts.DictionaryBatchLookupResponse{Tokens: []contracts.TokenDictionaryResult{{
		TokenIndex: 0,
		Dictionaries: []contracts.DictionaryLookupResult{
			{Metadata: contracts.DictionaryMetadata{ID: "a"}, Entries: []contracts.DictionaryEntry{{Translation: "見た"}}},
			{Metadata: contracts.DictionaryMetadata{ID: "b"}, Entries: []contracts.DictionaryEntry{{Translation: "見た"}, {Translation: "見ました"}}},
		},
	}}}

	result := processor.Merge(tokens, lookups)
	if result[0].Status != contracts.TokenStatusSuccess {
		t.Fatalf("status = %q, want success", result[0].Status)
	}
	if len(result[0].Candidates) != 2 {
		t.Fatalf("candidate count = %d, want 2", len(result[0].Candidates))
	}
	if !reflect.DeepEqual(result[0].Candidates[0].DictionaryIDs, []string{"a", "b"}) {
		t.Fatalf("sources = %#v, want [a b]", result[0].Candidates[0].DictionaryIDs)
	}
}

func TestMergeMarksUnknownWord(t *testing.T) {
	processor := NewTranslationProcessor(NewEnglishTokenizer())
	tokens := []contracts.Token{{Surface: "qwertymonster", LookupUnit: "qwertymonster", Kind: contracts.TokenKindWord}}

	result := processor.Merge(tokens, contracts.DictionaryBatchLookupResponse{})
	if result[0].Status != contracts.TokenStatusDictionaryNotFound {
		t.Fatalf("status = %q, want dictionary_not_found", result[0].Status)
	}
	if result[0].Token.Surface != "qwertymonster" {
		t.Fatalf("surface changed to %q", result[0].Token.Surface)
	}
}
