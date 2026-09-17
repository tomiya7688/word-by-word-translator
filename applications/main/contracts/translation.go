package contracts

type TranslationRequest struct {
	Text           string
	SourceLanguage Language
	TargetLanguage Language
}

type TranslationCandidate struct {
	Translation   string
	DictionaryIDs []string
}

type TranslatedToken struct {
	Token      Token
	Candidates []TranslationCandidate
	Status     TokenStatus
}

type TranslationResponse struct {
	Tokens []TranslatedToken
}
