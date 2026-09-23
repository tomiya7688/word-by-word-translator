package processing

import (
	"errors"
	"reflect"
	"testing"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
)

func TestScreenBuildsRequestWithMultipleSelectedDictionaries(t *testing.T) {
	processor := NewScreenProcessor()
	state := NewScreenState(
		contracts.LanguageEnglish,
		contracts.LanguageJapanese,
		[]DictionaryOption{
			{
				ID:             "dict-a",
				Name:           "Dictionary A",
				SourceLanguage: contracts.LanguageEnglish,
				TargetLanguage: contracts.LanguageJapanese,
			},
			{
				ID:             "dict-b",
				Name:           "Dictionary B",
				SourceLanguage: contracts.LanguageEnglish,
				TargetLanguage: contracts.LanguageJapanese,
			},
			{
				ID:             "ja-en",
				Name:           "JA-EN",
				SourceLanguage: contracts.LanguageJapanese,
				TargetLanguage: contracts.LanguageEnglish,
			},
		},
	)

	processor.SetText(&state, "I saw it.")
	if err := processor.ToggleDictionary(&state, "dict-a"); err != nil {
		t.Fatal(err)
	}
	if err := processor.ToggleDictionary(&state, "dict-b"); err != nil {
		t.Fatal(err)
	}

	request, err := processor.BuildRequest(state)
	if err != nil {
		t.Fatal(err)
	}
	if request.Text != "I saw it." ||
		request.SourceLanguage != contracts.LanguageEnglish ||
		request.TargetLanguage != contracts.LanguageJapanese {
		t.Fatalf("request = %#v", request)
	}
	if !reflect.DeepEqual(request.DictionaryIDs, []string{"dict-a", "dict-b"}) {
		t.Fatalf("dictionary ids = %#v", request.DictionaryIDs)
	}
}

func TestScreenDirectionSwitchDropsIncompatibleSelections(t *testing.T) {
	processor := NewScreenProcessor()
	state := NewScreenState(
		contracts.LanguageEnglish,
		contracts.LanguageJapanese,
		[]DictionaryOption{
			{
				ID:             "en-ja",
				SourceLanguage: contracts.LanguageEnglish,
				TargetLanguage: contracts.LanguageJapanese,
				Selected:       true,
			},
			{
				ID:             "ja-en",
				SourceLanguage: contracts.LanguageJapanese,
				TargetLanguage: contracts.LanguageEnglish,
			},
		},
	)

	if err := processor.SetDirection(&state, contracts.LanguageJapanese, contracts.LanguageEnglish); err != nil {
		t.Fatal(err)
	}
	if state.Dictionaries[0].Selected {
		t.Fatal("incompatible dictionary remained selected")
	}
	if err := processor.ToggleDictionary(&state, "ja-en"); err != nil {
		t.Fatal(err)
	}
	processor.SetText(&state, "食べる")

	request, err := processor.BuildRequest(state)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(request.DictionaryIDs, []string{"ja-en"}) {
		t.Fatalf("dictionary ids = %#v", request.DictionaryIDs)
	}
}

func TestScreenRequiresTextAndDictionarySelection(t *testing.T) {
	processor := NewScreenProcessor()
	state := NewScreenState(
		contracts.LanguageEnglish,
		contracts.LanguageJapanese,
		[]DictionaryOption{{
			ID:             "dict-a",
			SourceLanguage: contracts.LanguageEnglish,
			TargetLanguage: contracts.LanguageJapanese,
		}},
	)

	if _, err := processor.BuildRequest(state); !errors.Is(err, ErrScreenTextRequired) {
		t.Fatalf("error = %v", err)
	}
	processor.SetText(&state, "apple")
	if _, err := processor.BuildRequest(state); !errors.Is(err, ErrScreenDictionaryRequired) {
		t.Fatalf("error = %v", err)
	}
}

func TestScreenRejectsUnsupportedDirection(t *testing.T) {
	processor := NewScreenProcessor()
	state := NewScreenState(
		contracts.LanguageEnglish,
		contracts.LanguageJapanese,
		nil,
	)

	err := processor.SetDirection(&state, contracts.LanguageEnglish, contracts.LanguageEnglish)
	if !errors.Is(err, ErrScreenDirectionInvalid) {
		t.Fatalf("error = %v", err)
	}
}
