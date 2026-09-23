package processing

import "github.com/tomiya7688/word-by-word-translator/applications/main/contracts"

type TranslationProcessor struct {
	tokenizers *TokenizerRegistry
}

func NewTranslationProcessor(tokenizers *TokenizerRegistry) *TranslationProcessor {
	return &TranslationProcessor{tokenizers: tokenizers}
}

func (p *TranslationProcessor) Tokenize(text string, language contracts.Language) ([]contracts.Token, error) {
	return p.tokenizers.Tokenize(language, text)
}

func (p *TranslationProcessor) Merge(tokens []contracts.Token, lookups contracts.DictionaryBatchLookupResponse) []contracts.TranslatedToken {
	byIndex := make(map[int]contracts.TokenDictionaryResult, len(lookups.Tokens))
	for _, lookup := range lookups.Tokens {
		byIndex[lookup.TokenIndex] = lookup
	}

	result := make([]contracts.TranslatedToken, 0, len(tokens))
	for index, token := range tokens {
		if token.Kind == contracts.TokenKindSymbol || token.LookupUnit == "" {
			result = append(result, contracts.TranslatedToken{
				Token: token,
				Candidates: []contracts.TranslationCandidate{{
					Translation: token.Surface,
				}},
				Status: contracts.TokenStatusSuccess,
			})
			continue
		}

		dictionaryResults := byIndex[index].Dictionaries
		merged, hadError := mergeDictionaryResults(dictionaryResults)
		status := contracts.TokenStatusDictionaryNotFound
		if len(merged) > 0 {
			status = contracts.TokenStatusSuccess
		} else if hadError {
			status = contracts.TokenStatusLookupError
		}
		token.Status = status
		result = append(result, contracts.TranslatedToken{
			Token:             token,
			Candidates:        merged,
			DictionaryResults: dictionaryResults,
			Status:            status,
		})
	}

	return result
}

func mergeDictionaryResults(results []contracts.DictionaryLookupResult) ([]contracts.TranslationCandidate, bool) {
	candidates := make([]contracts.TranslationCandidate, 0)
	candidateIndex := make(map[string]int)
	hadError := false

	for _, result := range results {
		if result.Error != "" {
			hadError = true
		}
		for _, entry := range result.Entries {
			index, exists := candidateIndex[entry.Translation]
			if !exists {
				candidateIndex[entry.Translation] = len(candidates)
				candidates = append(candidates, contracts.TranslationCandidate{
					Translation:   entry.Translation,
					DictionaryIDs: []string{result.Metadata.ID},
				})
				continue
			}
			if !contains(candidates[index].DictionaryIDs, result.Metadata.ID) {
				candidates[index].DictionaryIDs = append(candidates[index].DictionaryIDs, result.Metadata.ID)
			}
		}
	}

	return candidates, hadError
}

func contains(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}
