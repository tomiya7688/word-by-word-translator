package messenger

import (
	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
	datacommander "github.com/tomiya7688/word-by-word-translator/applications/main/data/commander"
)

type DictionaryMessenger struct {
	commander *datacommander.DictionaryCommander
}

func NewDictionaryMessenger(commander *datacommander.DictionaryCommander) *DictionaryMessenger {
	return &DictionaryMessenger{commander: commander}
}

func (m *DictionaryMessenger) LookupBatch(request contracts.DictionaryBatchLookupRequest) contracts.DictionaryBatchLookupResponse {
	return m.commander.LookupBatch(request)
}
