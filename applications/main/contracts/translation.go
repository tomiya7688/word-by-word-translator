package contracts

type TranslationRequest struct {
	Text           string
	SourceLanguage Language
	TargetLanguage Language
	DictionaryIDs  []string
}

type TranslationCandidate struct {
	Translation   string
	DictionaryIDs []string
}

type TranslatedToken struct {
	Token             Token
	Candidates        []TranslationCandidate
	DictionaryResults []DictionaryLookupResult
	Status            TokenStatus
}

type TranslationResponse struct {
	Tokens []TranslatedToken
}
