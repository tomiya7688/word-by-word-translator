package processing

import (
	"reflect"
	"testing"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
)

func TestPresentTranslationKeepsUnresolvedSurfaceInOriginalOrder(t *testing.T) {
	response := contracts.TranslationResponse{Tokens: []contracts.TranslatedToken{
		{
			Token:      contracts.Token{Surface: "I"},
			Candidates: []contracts.TranslationCandidate{{Translation: "私"}},
			Status:     contracts.TokenStatusSuccess,
		},
		{
			Token:      contracts.Token{Surface: "saw"},
			Candidates: []contracts.TranslationCandidate{{Translation: "見た"}},
			Status:     contracts.TokenStatusSuccess,
		},
		{
			Token:  contracts.Token{Surface: "qwertymonster"},
			Status: contracts.TokenStatusDictionaryNotFound,
		},
		{
			Token:      contracts.Token{Surface: "yesterday"},
			Candidates: []contracts.TranslationCandidate{{Translation: "昨日"}},
			Status:     contracts.TokenStatusSuccess,
		},
		{
			Token:      contracts.Token{Surface: "."},
			Candidates: []contracts.TranslationCandidate{{Translation: "."}},
			Status:     contracts.TokenStatusSuccess,
		},
	}}

	got := PresentTranslation(response)
	want := []PresentedToken{
		{Text: "私", Tone: TokenToneNormal, Candidates: []PresentedCandidate{{Translation: "私"}}},
		{Text: "見た", Tone: TokenToneNormal, Candidates: []PresentedCandidate{{Translation: "見た"}}},
		{Text: "qwertymonster", Tone: TokenToneUnresolved},
		{Text: "昨日", Tone: TokenToneNormal, Candidates: []PresentedCandidate{{Translation: "昨日"}}},
		{Text: ".", Tone: TokenToneNormal, Candidates: []PresentedCandidate{{Translation: "."}}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("presented = %#v, want %#v", got, want)
	}
}

func TestPresentTranslationKeepsAllMergedCandidatesAndSources(t *testing.T) {
	got := PresentTranslation(contracts.TranslationResponse{Tokens: []contracts.TranslatedToken{{
		Token: contracts.Token{Surface: "saw"},
		Candidates: []contracts.TranslationCandidate{
			{Translation: "見た", DictionaryIDs: []string{"a", "b"}},
			{Translation: "見ました", DictionaryIDs: []string{"b"}},
		},
		Status: contracts.TokenStatusSuccess,
	}}})

	if len(got) != 1 || got[0].Text != "見た" {
		t.Fatalf("presented = %#v", got)
	}
	want := []PresentedCandidate{
		{Translation: "見た", DictionaryIDs: []string{"a", "b"}},
		{Translation: "見ました", DictionaryIDs: []string{"b"}},
	}
	if !reflect.DeepEqual(got[0].Candidates, want) {
		t.Fatalf("candidates = %#v, want %#v", got[0].Candidates, want)
	}
}

func TestPresentTranslationMarksEveryFailureStatusUnresolved(t *testing.T) {
	statuses := []contracts.TokenStatus{
		contracts.TokenStatusSegmentationFailed,
		contracts.TokenStatusDictionaryNotFound,
		contracts.TokenStatusLookupError,
	}

	for _, status := range statuses {
		t.Run(string(status), func(t *testing.T) {
			got := PresentTranslation(contracts.TranslationResponse{Tokens: []contracts.TranslatedToken{{
				Token:  contracts.Token{Surface: "original"},
				Status: status,
			}}})
			if len(got) != 1 || got[0].Text != "original" || got[0].Tone != TokenToneUnresolved {
				t.Fatalf("presented = %#v", got)
			}
		})
	}
}

func TestPresentTranslationFallsBackToSurfaceForSuccessfulEmptyCandidates(t *testing.T) {
	got := PresentTranslation(contracts.TranslationResponse{Tokens: []contracts.TranslatedToken{{
		Token:  contracts.Token{Surface: "."},
		Status: contracts.TokenStatusSuccess,
	}}})

	if len(got) != 1 || got[0].Text != "." || got[0].Tone != TokenToneNormal {
		t.Fatalf("presented = %#v", got)
	}
}
