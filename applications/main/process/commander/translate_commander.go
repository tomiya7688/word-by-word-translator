package commander

import "github.com/tomiya7688/word-by-word-translator/applications/main/contracts"

type TranslationProcessing interface {
	Tokenize(text string, language contracts.Language) ([]contracts.Token, error)
	Merge(tokens []contracts.Token, lookups contracts.DictionaryBatchLookupResponse) []contracts.TranslatedToken
}

type DictionaryMessenger interface {
	LookupBatch(request contracts.DictionaryBatchLookupRequest) contracts.DictionaryBatchLookupResponse
}

type TranslateCommander struct {
	processing TranslationProcessing
	dictionary DictionaryMessenger
}

func NewTranslateCommander(processing TranslationProcessing, dictionary DictionaryMessenger) *TranslateCommander {
	return &TranslateCommander{processing: processing, dictionary: dictionary}
}

func (c *TranslateCommander) Translate(request contracts.TranslationRequest) (contracts.TranslationResponse, error) {
	tokens, err := c.processing.Tokenize(request.Text, request.SourceLanguage)
	if err != nil {
		return contracts.TranslationResponse{}, err
	}
	lookups := c.dictionary.LookupBatch(contracts.DictionaryBatchLookupRequest{
		Tokens:         tokens,
		SourceLanguage: request.SourceLanguage,
		TargetLanguage: request.TargetLanguage,
		DictionaryIDs:  append([]string(nil), request.DictionaryIDs...),
	})
	return contracts.TranslationResponse{Tokens: c.processing.Merge(tokens, lookups)}, nil
}
