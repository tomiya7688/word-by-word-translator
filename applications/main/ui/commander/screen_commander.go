package commander

import (
	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
	uiprocessing "github.com/tomiya7688/word-by-word-translator/applications/main/ui/processing"
)

type ScreenProcessing interface {
	SetText(state *uiprocessing.ScreenState, text string)
	SetDirection(state *uiprocessing.ScreenState, source, target contracts.Language) error
	ToggleDictionary(state *uiprocessing.ScreenState, dictionaryID string) error
	BuildRequest(state uiprocessing.ScreenState) (contracts.TranslationRequest, error)
	ApplyResponse(state *uiprocessing.ScreenState, response contracts.TranslationResponse)
	ApplyError(state *uiprocessing.ScreenState, err error)
}

type TranslationMessenger interface {
	Translate(request contracts.TranslationRequest) (contracts.TranslationResponse, error)
}

type ScreenCommander struct {
	processing  ScreenProcessing
	translation TranslationMessenger
}

func NewScreenCommander(processing ScreenProcessing, translation TranslationMessenger) *ScreenCommander {
	return &ScreenCommander{processing: processing, translation: translation}
}

func (c *ScreenCommander) SetText(state *uiprocessing.ScreenState, text string) {
	c.processing.SetText(state, text)
}

func (c *ScreenCommander) SetDirection(state *uiprocessing.ScreenState, source, target contracts.Language) error {
	err := c.processing.SetDirection(state, source, target)
	if err != nil {
		c.processing.ApplyError(state, err)
	}
	return err
}

func (c *ScreenCommander) ToggleDictionary(state *uiprocessing.ScreenState, dictionaryID string) error {
	err := c.processing.ToggleDictionary(state, dictionaryID)
	if err != nil {
		c.processing.ApplyError(state, err)
	}
	return err
}

func (c *ScreenCommander) Execute(state *uiprocessing.ScreenState) error {
	request, err := c.processing.BuildRequest(*state)
	if err != nil {
		c.processing.ApplyError(state, err)
		return err
	}
	response, err := c.translation.Translate(request)
	if err != nil {
		c.processing.ApplyError(state, err)
		return err
	}
	c.processing.ApplyResponse(state, response)
	return nil
}
