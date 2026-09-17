package processing

import "github.com/tomiya7688/word-by-word-translator/applications/main/contracts"

type Dictionary interface {
	Metadata() contracts.DictionaryMetadata
	Lookup(lookupUnit string) ([]contracts.DictionaryEntry, error)
}
