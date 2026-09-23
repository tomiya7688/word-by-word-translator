package processing

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
)

var (
	ErrScreenTextRequired       = errors.New("translation text is required")
	ErrScreenDictionaryRequired = errors.New("at least one dictionary must be selected")
	ErrScreenDirectionInvalid   = errors.New("unsupported translation direction")
)

type DictionaryOption struct {
	ID             string
	Name           string
	SourceLanguage contracts.Language
	TargetLanguage contracts.Language
	Selected       bool
}

type ScreenState struct {
	Text           string
	SourceLanguage contracts.Language
	TargetLanguage contracts.Language
	Dictionaries   []DictionaryOption
	Result         []PresentedToken
	Error          string
}

type ScreenProcessor struct{}

func NewScreenProcessor() *ScreenProcessor {
	return &ScreenProcessor{}
}

func NewScreenState(source, target contracts.Language, dictionaries []DictionaryOption) ScreenState {
	return ScreenState{
		SourceLanguage: source,
		TargetLanguage: target,
		Dictionaries:   append([]DictionaryOption(nil), dictionaries...),
	}
}

func (p *ScreenProcessor) SetText(state *ScreenState, text string) {
	state.Text = text
	state.Error = ""
}

func (p *ScreenProcessor) SetDirection(state *ScreenState, source, target contracts.Language) error {
	if !supportedDirection(source, target) {
		return fmt.Errorf("%w: %s -> %s", ErrScreenDirectionInvalid, source, target)
	}
	state.SourceLanguage = source
	state.TargetLanguage = target
	state.Error = ""
	for index := range state.Dictionaries {
		if !dictionaryMatchesDirection(state.Dictionaries[index], source, target) {
			state.Dictionaries[index].Selected = false
		}
	}
	return nil
}

func (p *ScreenProcessor) ToggleDictionary(state *ScreenState, dictionaryID string) error {
	for index := range state.Dictionaries {
		option := &state.Dictionaries[index]
		if option.ID != dictionaryID {
			continue
		}
		if !dictionaryMatchesDirection(*option, state.SourceLanguage, state.TargetLanguage) {
			return fmt.Errorf("dictionary %q is not available for %s -> %s", dictionaryID, state.SourceLanguage, state.TargetLanguage)
		}
		option.Selected = !option.Selected
		state.Error = ""
		return nil
	}
	return fmt.Errorf("dictionary %q is not available", dictionaryID)
}

func (p *ScreenProcessor) BuildRequest(state ScreenState) (contracts.TranslationRequest, error) {
	if strings.TrimSpace(state.Text) == "" {
		return contracts.TranslationRequest{}, ErrScreenTextRequired
	}
	if !supportedDirection(state.SourceLanguage, state.TargetLanguage) {
		return contracts.TranslationRequest{}, fmt.Errorf(
			"%w: %s -> %s",
			ErrScreenDirectionInvalid,
			state.SourceLanguage,
			state.TargetLanguage,
		)
	}

	dictionaryIDs := selectedDictionaryIDs(state)
	if len(dictionaryIDs) == 0 {
		return contracts.TranslationRequest{}, ErrScreenDictionaryRequired
	}
	return contracts.TranslationRequest{
		Text:           state.Text,
		SourceLanguage: state.SourceLanguage,
		TargetLanguage: state.TargetLanguage,
		DictionaryIDs:  dictionaryIDs,
	}, nil
}

func (p *ScreenProcessor) ApplyResponse(state *ScreenState, response contracts.TranslationResponse) {
	state.Result = PresentTranslation(response)
	state.Error = ""
}

func (p *ScreenProcessor) ApplyError(state *ScreenState, err error) {
	if err == nil {
		state.Error = ""
		return
	}
	state.Error = err.Error()
}

func selectedDictionaryIDs(state ScreenState) []string {
	selected := make([]string, 0, len(state.Dictionaries))
	for _, option := range state.Dictionaries {
		if option.Selected && dictionaryMatchesDirection(option, state.SourceLanguage, state.TargetLanguage) {
			selected = append(selected, option.ID)
		}
	}
	return selected
}

func dictionaryMatchesDirection(option DictionaryOption, source, target contracts.Language) bool {
	return option.SourceLanguage == source && option.TargetLanguage == target
}

func supportedDirection(source, target contracts.Language) bool {
	return (source == contracts.LanguageEnglish && target == contracts.LanguageJapanese) ||
		(source == contracts.LanguageJapanese && target == contracts.LanguageEnglish)
}
