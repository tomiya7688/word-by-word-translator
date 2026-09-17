package messenger

import (
	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
	datamessenger "github.com/tomiya7688/word-by-word-translator/applications/main/data/messenger"
)

type DictionaryMessenger struct {
	data *datamessenger.DictionaryMessenger
}

func NewDictionaryMessenger(data *datamessenger.DictionaryMessenger) *DictionaryMessenger {
	return &DictionaryMessenger{data: data}
}

func (m *DictionaryMessenger) LookupBatch(request contracts.DictionaryBatchLookupRequest) contracts.DictionaryBatchLookupResponse {
	return m.data.LookupBatch(request)
}
