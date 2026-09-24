package processing

import (
	"reflect"
	"testing"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
)

func TestMergeDeduplicatesTranslationsAndKeepsDictionarySources(t *testing.T) {
	processor := NewTranslationProcessor(NewTokenizerRegistry())
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
	if len(result[0].DictionaryResults) != 2 ||
		result[0].DictionaryResults[0].Metadata.ID != "a" ||
		result[0].DictionaryResults[1].Metadata.ID != "b" {
		t.Fatalf("dictionary results = %#v", result[0].DictionaryResults)
	}
}

func TestMergeKeepsSuccessfulCandidatesWhenAnotherDictionaryFails(t *testing.T) {
	processor := NewTranslationProcessor(NewTokenizerRegistry())
	tokens := []contracts.Token{{Surface: "saw", LookupUnit: "saw", Kind: contracts.TokenKindWord}}
	lookups := contracts.DictionaryBatchLookupResponse{Tokens: []contracts.TokenDictionaryResult{{
		TokenIndex: 0,
		Dictionaries: []contracts.DictionaryLookupResult{
			{Metadata: contracts.DictionaryMetadata{ID: "broken"}, Error: "lookup failed"},
			{Metadata: contracts.DictionaryMetadata{ID: "working"}, Entries: []contracts.DictionaryEntry{{Translation: "見た"}}},
		},
	}}}

	result := processor.Merge(tokens, lookups)
	if result[0].Status != contracts.TokenStatusSuccess {
		t.Fatalf("status = %q, want success", result[0].Status)
	}
	if len(result[0].Candidates) != 1 || result[0].Candidates[0].Translation != "見た" {
		t.Fatalf("candidates = %#v", result[0].Candidates)
	}
	if !reflect.DeepEqual(result[0].Candidates[0].DictionaryIDs, []string{"working"}) {
		t.Fatalf("sources = %#v, want [working]", result[0].Candidates[0].DictionaryIDs)
	}
	if len(result[0].DictionaryResults) != 2 || result[0].DictionaryResults[0].Error != "lookup failed" {
		t.Fatalf("dictionary results = %#v", result[0].DictionaryResults)
	}
}

func TestMergePreservesDictionaryPriorityOrder(t *testing.T) {
	processor := NewTranslationProcessor(NewTokenizerRegistry())
	tokens := []contracts.Token{{Surface: "word", LookupUnit: "word", Kind: contracts.TokenKindWord}}
	lookups := contracts.DictionaryBatchLookupResponse{Tokens: []contracts.TokenDictionaryResult{{
		TokenIndex: 0,
		Dictionaries: []contracts.DictionaryLookupResult{
			{
				Metadata: contracts.DictionaryMetadata{ID: "high-priority"},
				Entries: []contracts.DictionaryEntry{
					{Translation: "first"},
					{Translation: "shared"},
				},
			},
			{
				Metadata: contracts.DictionaryMetadata{ID: "low-priority"},
				Entries: []contracts.DictionaryEntry{
					{Translation: "second"},
					{Translation: "shared"},
				},
			},
		},
	}}}

	result := processor.Merge(tokens, lookups)
	got := []string{
		result[0].Candidates[0].Translation,
		result[0].Candidates[1].Translation,
		result[0].Candidates[2].Translation,
	}
	if !reflect.DeepEqual(got, []string{"first", "shared", "second"}) {
		t.Fatalf("candidate order = %#v", got)
	}
	if !reflect.DeepEqual(result[0].Candidates[1].DictionaryIDs, []string{"high-priority", "low-priority"}) {
		t.Fatalf("shared sources = %#v", result[0].Candidates[1].DictionaryIDs)
	}
}

func TestMergeMarksUnknownWord(t *testing.T) {
	processor := NewTranslationProcessor(NewTokenizerRegistry())
	tokens := []contracts.Token{{Surface: "qwertymonster", LookupUnit: "qwertymonster", Kind: contracts.TokenKindWord}}

	result := processor.Merge(tokens, contracts.DictionaryBatchLookupResponse{})
	if result[0].Status != contracts.TokenStatusDictionaryNotFound {
		t.Fatalf("status = %q, want dictionary_not_found", result[0].Status)
	}
	if result[0].Token.Surface != "qwertymonster" {
		t.Fatalf("surface changed to %q", result[0].Token.Surface)
	}
}


func TestMergeMarksLookupErrorWhenEveryDictionaryFails(t *testing.T) {
	processor := NewTranslationProcessor(NewTokenizerRegistry())
	tokens := []contracts.Token{{Surface: "brokenword", LookupUnit: "brokenword", Kind: contracts.TokenKindWord}}
	lookups := contracts.DictionaryBatchLookupResponse{Tokens: []contracts.TokenDictionaryResult{{
		TokenIndex: 0,
		Dictionaries: []contracts.DictionaryLookupResult{
			{Metadata: contracts.DictionaryMetadata{ID: "broken-a"}, Error: "lookup failed"},
			{Metadata: contracts.DictionaryMetadata{ID: "broken-b"}, Error: "lookup failed"},
		},
	}}}

	result := processor.Merge(tokens, lookups)
	if len(result) != 1 {
		t.Fatalf("result count = %d, want 1", len(result))
	}
	if result[0].Status != contracts.TokenStatusLookupError {
		t.Fatalf("status = %q, want lookup_error", result[0].Status)
	}
	if result[0].Token.Surface != "brokenword" {
		t.Fatalf("surface = %q, want original", result[0].Token.Surface)
	}
	if len(result[0].Candidates) != 0 {
		t.Fatalf("candidates = %#v, want none", result[0].Candidates)
	}
}
