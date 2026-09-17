package commander

import (
	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
	dataprocessing "github.com/tomiya7688/word-by-word-translator/applications/main/data/processing"
)

type DictionaryCommander struct {
	store *dataprocessing.DictionaryStore
}

func NewDictionaryCommander(store *dataprocessing.DictionaryStore) *DictionaryCommander {
	return &DictionaryCommander{store: store}
}

func (c *DictionaryCommander) LookupBatch(request contracts.DictionaryBatchLookupRequest) contracts.DictionaryBatchLookupResponse {
	return c.store.LookupBatch(request)
}
