package integration

import (
	"reflect"
	"testing"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
	datacommander "github.com/tomiya7688/word-by-word-translator/applications/main/data/commander"
	datamessenger "github.com/tomiya7688/word-by-word-translator/applications/main/data/messenger"
	dataprocessing "github.com/tomiya7688/word-by-word-translator/applications/main/data/processing"
	processcommander "github.com/tomiya7688/word-by-word-translator/applications/main/process/commander"
	processmessenger "github.com/tomiya7688/word-by-word-translator/applications/main/process/messenger"
	processprocessing "github.com/tomiya7688/word-by-word-translator/applications/main/process/processing"
)

func TestEnglishToJapaneseWordByWordPipeline(t *testing.T) {
	metadataA := contracts.DictionaryMetadata{
		ID: "dict-a", Name: "Dictionary A", SourceLanguage: contracts.LanguageEnglish, TargetLanguage: contracts.LanguageJapanese,
	}
	metadataB := contracts.DictionaryMetadata{
		ID: "dict-b", Name: "Dictionary B", SourceLanguage: contracts.LanguageEnglish, TargetLanguage: contracts.LanguageJapanese,
	}
	dictA := dataprocessing.NewMemoryDictionary(metadataA, map[string][]string{
		"i": {"私"}, "saw": {"見た"}, "yesterday": {"昨日"},
	})
	dictB := dataprocessing.NewMemoryDictionary(metadataB, map[string][]string{
		"saw": {"見た", "見ました"},
	})
	store := dataprocessing.NewDictionaryStore(dictA, dictB)
	dataCommander := datacommander.NewDictionaryCommander(store)
	dataMessenger := datamessenger.NewDictionaryMessenger(dataCommander)
	processMessenger := processmessenger.NewDictionaryMessenger(dataMessenger)
	processor := processprocessing.NewTranslationProcessor(processprocessing.NewEnglishTokenizer())
	translate := processcommander.NewTranslateCommander(processor, processMessenger)

	response, err := translate.Translate(contracts.TranslationRequest{
		Text: "I saw qwertymonster yesterday.", SourceLanguage: contracts.LanguageEnglish, TargetLanguage: contracts.LanguageJapanese,
	})
	if err != nil {
		t.Fatalf("Translate returned error: %v", err)
	}
	if len(response.Tokens) != 5 {
		t.Fatalf("token count = %d, want 5", len(response.Tokens))
	}
	if response.Tokens[2].Status != contracts.TokenStatusDictionaryNotFound || response.Tokens[2].Token.Surface != "qwertymonster" {
		t.Fatalf("unknown token = %#v", response.Tokens[2])
	}
	if !reflect.DeepEqual(response.Tokens[1].Candidates[0].DictionaryIDs, []string{"dict-a", "dict-b"}) {
		t.Fatalf("merged sources = %#v", response.Tokens[1].Candidates[0].DictionaryIDs)
	}
	if response.Tokens[4].Candidates[0].Translation != "." {
		t.Fatalf("symbol translation = %q, want .", response.Tokens[4].Candidates[0].Translation)
	}
}
