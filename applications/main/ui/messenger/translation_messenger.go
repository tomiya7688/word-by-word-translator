package messenger

import "github.com/tomiya7688/word-by-word-translator/applications/main/contracts"

type TranslationCommander interface {
	Translate(request contracts.TranslationRequest) (contracts.TranslationResponse, error)
}

type TranslationMessenger struct {
	commander TranslationCommander
}

func NewTranslationMessenger(commander TranslationCommander) *TranslationMessenger {
	return &TranslationMessenger{commander: commander}
}

func (m *TranslationMessenger) Translate(request contracts.TranslationRequest) (contracts.TranslationResponse, error) {
	return m.commander.Translate(request)
}
