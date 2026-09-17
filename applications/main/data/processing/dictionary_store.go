package processing

import "github.com/tomiya7688/word-by-word-translator/applications/main/contracts"

type DictionaryStore struct {
	dictionaries []Dictionary
}

func NewDictionaryStore(dictionaries ...Dictionary) *DictionaryStore {
	copied := append([]Dictionary(nil), dictionaries...)
	return &DictionaryStore{dictionaries: copied}
}

func (s *DictionaryStore) LookupBatch(request contracts.DictionaryBatchLookupRequest) contracts.DictionaryBatchLookupResponse {
	results := make([]contracts.TokenDictionaryResult, 0, len(request.Tokens))
	for index, token := range request.Tokens {
		if token.Kind != contracts.TokenKindWord || token.LookupUnit == "" {
			continue
		}

		dictionaryResults := make([]contracts.DictionaryLookupResult, 0, len(s.dictionaries))
		for _, dictionary := range s.dictionaries {
			metadata := dictionary.Metadata()
			if metadata.SourceLanguage != request.SourceLanguage || metadata.TargetLanguage != request.TargetLanguage {
				continue
			}

			entries, err := dictionary.Lookup(token.LookupUnit)
			lookup := contracts.DictionaryLookupResult{Metadata: metadata, Entries: entries}
			if err != nil {
				lookup.Error = err.Error()
			}
			dictionaryResults = append(dictionaryResults, lookup)
		}

		results = append(results, contracts.TokenDictionaryResult{
			TokenIndex:   index,
			Dictionaries: dictionaryResults,
		})
	}
	return contracts.DictionaryBatchLookupResponse{Tokens: results}
}
