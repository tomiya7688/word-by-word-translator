package integration

import (
	"errors"
	"testing"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
	datacommander "github.com/tomiya7688/word-by-word-translator/applications/main/data/commander"
	datamessenger "github.com/tomiya7688/word-by-word-translator/applications/main/data/messenger"
	dataprocessing "github.com/tomiya7688/word-by-word-translator/applications/main/data/processing"
	processcommander "github.com/tomiya7688/word-by-word-translator/applications/main/process/commander"
	processmessenger "github.com/tomiya7688/word-by-word-translator/applications/main/process/messenger"
	processprocessing "github.com/tomiya7688/word-by-word-translator/applications/main/process/processing"
	kagomeprocessing "github.com/tomiya7688/word-by-word-translator/applications/main/process/processing/kagome"
)

func TestOnePipelineRoutesEnglishAndJapaneseTokenizers(t *testing.T) {
	japaneseTokenizer, err := kagomeprocessing.NewJapaneseTokenizer()
	if err != nil {
		t.Fatal(err)
	}

	tokenizers := processprocessing.NewTokenizerRegistry()
	if err := tokenizers.Register(contracts.LanguageEnglish, processprocessing.NewEnglishTokenizer()); err != nil {
		t.Fatal(err)
	}
	if err := tokenizers.Register(contracts.LanguageJapanese, japaneseTokenizer); err != nil {
		t.Fatal(err)
	}

	enJa := dataprocessing.NewMemoryDictionary(
		contracts.DictionaryMetadata{
			ID: "en-ja", Name: "EN-JA", SourceLanguage: contracts.LanguageEnglish, TargetLanguage: contracts.LanguageJapanese,
		},
		map[string][]string{"apple": {"りんご"}},
	)
	jaEn := dataprocessing.NewMemoryDictionary(
		contracts.DictionaryMetadata{
			ID: "ja-en", Name: "JA-EN", SourceLanguage: contracts.LanguageJapanese, TargetLanguage: contracts.LanguageEnglish,
		},
		map[string][]string{"食べる": {"to eat"}},
	)

	store := dataprocessing.NewDictionaryStore(enJa, jaEn)
	dataCommander := datacommander.NewDictionaryCommander(store)
	dataMessenger := datamessenger.NewDictionaryMessenger(dataCommander)
	processMessenger := processmessenger.NewDictionaryMessenger(dataMessenger)
	processor := processprocessing.NewTranslationProcessor(tokenizers)
	translate := processcommander.NewTranslateCommander(processor, processMessenger)

	english, err := translate.Translate(contracts.TranslationRequest{
		Text: "apple", SourceLanguage: contracts.LanguageEnglish, TargetLanguage: contracts.LanguageJapanese,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(english.Tokens) != 1 || len(english.Tokens[0].Candidates) != 1 || english.Tokens[0].Candidates[0].Translation != "りんご" {
		t.Fatalf("english response = %#v", english)
	}

	japanese, err := translate.Translate(contracts.TranslationRequest{
		Text: "食べました", SourceLanguage: contracts.LanguageJapanese, TargetLanguage: contracts.LanguageEnglish,
	})
	if err != nil {
		t.Fatal(err)
	}
	var found bool
	for _, token := range japanese.Tokens {
		if token.Token.LookupUnit == "食べる" {
			found = true
			if len(token.Candidates) != 1 || token.Candidates[0].Translation != "to eat" {
				t.Fatalf("japanese token = %#v", token)
			}
		}
	}
	if !found {
		t.Fatalf("japanese response = %#v", japanese)
	}

	_, err = translate.Translate(contracts.TranslationRequest{
		Text: "test", SourceLanguage: contracts.Language("xx"), TargetLanguage: contracts.LanguageEnglish,
	})
	if !errors.Is(err, processprocessing.ErrTokenizerNotRegistered) {
		t.Fatalf("unsupported language error = %v", err)
	}
}
