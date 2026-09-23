package commander

import (
	"errors"
	"reflect"
	"testing"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
	uiprocessing "github.com/tomiya7688/word-by-word-translator/applications/main/ui/processing"
)

type fakeTranslationMessenger struct {
	request contracts.TranslationRequest
	calls   int
	result  contracts.TranslationResponse
	err     error
}

func (m *fakeTranslationMessenger) Translate(request contracts.TranslationRequest) (contracts.TranslationResponse, error) {
	m.calls++
	m.request = request
	return m.result, m.err
}

func TestScreenCommanderExecutesTranslationFromScreenState(t *testing.T) {
	screenProcessing := uiprocessing.NewScreenProcessor()
	translation := &fakeTranslationMessenger{
		result: contracts.TranslationResponse{Tokens: []contracts.TranslatedToken{{
			Token: contracts.Token{Surface: "apple"},
			Candidates: []contracts.TranslationCandidate{{
				Translation:   "りんご",
				DictionaryIDs: []string{"dict-a"},
			}},
			Status: contracts.TokenStatusSuccess,
		}}},
	}
	commander := NewScreenCommander(screenProcessing, translation)
	state := uiprocessing.NewScreenState(
		contracts.LanguageEnglish,
		contracts.LanguageJapanese,
		[]uiprocessing.DictionaryOption{{
			ID:             "dict-a",
			SourceLanguage: contracts.LanguageEnglish,
			TargetLanguage: contracts.LanguageJapanese,
			Selected:       true,
		}},
	)
	commander.SetText(&state, "apple")

	if err := commander.Execute(&state); err != nil {
		t.Fatal(err)
	}
	if translation.calls != 1 {
		t.Fatalf("translation calls = %d, want 1", translation.calls)
	}
	if !reflect.DeepEqual(translation.request.DictionaryIDs, []string{"dict-a"}) {
		t.Fatalf("request = %#v", translation.request)
	}
	if len(state.Result) != 1 || state.Result[0].Text != "りんご" || state.Error != "" {
		t.Fatalf("state = %#v", state)
	}
}

func TestScreenCommanderKeepsValidationErrorOnScreen(t *testing.T) {
	screenProcessing := uiprocessing.NewScreenProcessor()
	translation := &fakeTranslationMessenger{}
	commander := NewScreenCommander(screenProcessing, translation)
	state := uiprocessing.NewScreenState(
		contracts.LanguageEnglish,
		contracts.LanguageJapanese,
		nil,
	)
	commander.SetText(&state, "apple")

	err := commander.Execute(&state)
	if !errors.Is(err, uiprocessing.ErrScreenDictionaryRequired) {
		t.Fatalf("error = %v", err)
	}
	if translation.calls != 0 {
		t.Fatalf("translation calls = %d, want 0", translation.calls)
	}
	if state.Error == "" {
		t.Fatal("screen error was not populated")
	}
}
