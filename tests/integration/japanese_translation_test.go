package integration

import (
	"testing"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
	datacommander "github.com/tomiya7688/word-by-word-translator/applications/main/data/commander"
	datamessenger "github.com/tomiya7688/word-by-word-translator/applications/main/data/messenger"
	dataprocessing "github.com/tomiya7688/word-by-word-translator/applications/main/data/processing"
	processcommander "github.com/tomiya7688/word-by-word-translator/applications/main/process/commander"
	processmessenger "github.com/tomiya7688/word-by-word-translator/applications/main/process/messenger"
	"github.com/tomiya7688/word-by-word-translator/applications/main/process/processing"
	kagomeprocessing "github.com/tomiya7688/word-by-word-translator/applications/main/process/processing/kagome"
)

func TestJapaneseInflectionReachesDictionaryLookup(t *testing.T) {
	tokenizer, err := kagomeprocessing.NewJapaneseTokenizer()
	if err != nil {
		t.Fatal(err)
	}

	dictionary := dataprocessing.NewMemoryDictionary(
		contracts.DictionaryMetadata{
			ID:             "ja-en-test",
			Name:           "JA-EN test",
			SourceLanguage: contracts.LanguageJapanese,
			TargetLanguage: contracts.LanguageEnglish,
			ReleaseTier:    contracts.ReleaseTierMIT,
		},
		map[string][]string{
			"食べる": {"to eat"},
		},
	)

	store := dataprocessing.NewDictionaryStore(dictionary)
	dataCommander := datacommander.NewDictionaryCommander(store)
	dataMessenger := datamessenger.NewDictionaryMessenger(dataCommander)
	processMessenger := processmessenger.NewDictionaryMessenger(dataMessenger)

	tokenizers := processing.NewTokenizerRegistry()
	if err := tokenizers.Register(contracts.LanguageJapanese, tokenizer); err != nil {
		t.Fatal(err)
	}
	processor := processing.NewTranslationProcessor(tokenizers)
	commander := processcommander.NewTranslateCommander(processor, processMessenger)

	response, err := commander.Translate(contracts.TranslationRequest{
		Text:           "りんごを食べました。",
		SourceLanguage: contracts.LanguageJapanese,
		TargetLanguage: contracts.LanguageEnglish,
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, token := range response.Tokens {
		if token.Token.Surface != "食べ" {
			continue
		}
		if token.Token.LookupUnit != "食べる" {
			t.Fatalf("lookup unit = %q", token.Token.LookupUnit)
		}
		if token.Status != contracts.TokenStatusSuccess {
			t.Fatalf("status = %q", token.Status)
		}
		if len(token.Candidates) != 1 || token.Candidates[0].Translation != "to eat" {
			t.Fatalf("candidates = %#v", token.Candidates)
		}
		return
	}

	t.Fatal("食べ token was not found")
}
