package integration

import (
	"strings"
	"testing"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
	datacommander "github.com/tomiya7688/word-by-word-translator/applications/main/data/commander"
	datamessenger "github.com/tomiya7688/word-by-word-translator/applications/main/data/messenger"
	dataprocessing "github.com/tomiya7688/word-by-word-translator/applications/main/data/processing"
	processcommander "github.com/tomiya7688/word-by-word-translator/applications/main/process/commander"
	processmessenger "github.com/tomiya7688/word-by-word-translator/applications/main/process/messenger"
	processprocessing "github.com/tomiya7688/word-by-word-translator/applications/main/process/processing"
	uicommander "github.com/tomiya7688/word-by-word-translator/applications/main/ui/commander"
	uimessenger "github.com/tomiya7688/word-by-word-translator/applications/main/ui/messenger"
	uiprocessing "github.com/tomiya7688/word-by-word-translator/applications/main/ui/processing"
)

type trackingDictionary struct {
	metadata contracts.DictionaryMetadata
	entries  map[string][]contracts.DictionaryEntry
	calls    int
}

func (d *trackingDictionary) Metadata() contracts.DictionaryMetadata {
	return d.metadata
}

func (d *trackingDictionary) Lookup(lookupUnit string) ([]contracts.DictionaryEntry, error) {
	d.calls++
	return append([]contracts.DictionaryEntry(nil), d.entries[lookupUnit]...), nil
}

func TestTranslationScreenUsesOnlySelectedDictionariesAndDisplaysMergedResults(t *testing.T) {
	dictA := &trackingDictionary{
		metadata: contracts.DictionaryMetadata{
			ID:             "dict-a",
			Name:           "Dictionary A",
			SourceLanguage: contracts.LanguageEnglish,
			TargetLanguage: contracts.LanguageJapanese,
		},
		entries: map[string][]contracts.DictionaryEntry{
			"apple": {{Headword: "apple", Translation: "りんごA"}},
		},
	}
	dictB := &trackingDictionary{
		metadata: contracts.DictionaryMetadata{
			ID:             "dict-b",
			Name:           "Dictionary B",
			SourceLanguage: contracts.LanguageEnglish,
			TargetLanguage: contracts.LanguageJapanese,
		},
		entries: map[string][]contracts.DictionaryEntry{
			"apple": {{Headword: "apple", Translation: "りんごB"}},
		},
	}

	store := dataprocessing.NewDictionaryStore(dictA, dictB)
	dataCommander := datacommander.NewDictionaryCommander(store)
	dataMessenger := datamessenger.NewDictionaryMessenger(dataCommander)
	processDictionaryMessenger := processmessenger.NewDictionaryMessenger(dataMessenger)

	tokenizers := processprocessing.NewTokenizerRegistry()
	if err := tokenizers.Register(contracts.LanguageEnglish, processprocessing.NewEnglishTokenizer()); err != nil {
		t.Fatal(err)
	}
	processProcessor := processprocessing.NewTranslationProcessor(tokenizers)
	processCommander := processcommander.NewTranslateCommander(processProcessor, processDictionaryMessenger)

	uiTranslationMessenger := uimessenger.NewTranslationMessenger(processCommander)
	screenProcessor := uiprocessing.NewScreenProcessor()
	screenCommander := uicommander.NewScreenCommander(screenProcessor, uiTranslationMessenger)

	options := []uiprocessing.DictionaryOption{
		{
			ID:             "dict-a",
			Name:           "Dictionary A",
			SourceLanguage: contracts.LanguageEnglish,
			TargetLanguage: contracts.LanguageJapanese,
		},
		{
			ID:             "dict-b",
			Name:           "Dictionary B",
			SourceLanguage: contracts.LanguageEnglish,
			TargetLanguage: contracts.LanguageJapanese,
		},
	}

	selectedB := uiprocessing.NewScreenState(
		contracts.LanguageEnglish,
		contracts.LanguageJapanese,
		options,
	)
	screenCommander.SetText(&selectedB, "apple")
	if err := screenCommander.ToggleDictionary(&selectedB, "dict-b"); err != nil {
		t.Fatal(err)
	}
	if err := screenCommander.Execute(&selectedB); err != nil {
		t.Fatal(err)
	}

	if dictA.calls != 0 || dictB.calls != 1 {
		t.Fatalf("lookup calls: dict-a=%d dict-b=%d", dictA.calls, dictB.calls)
	}
	if len(selectedB.Result) != 1 ||
		selectedB.Result[0].Text != "りんごB" ||
		len(selectedB.Result[0].Candidates) != 1 {
		t.Fatalf("selected result = %#v", selectedB.Result)
	}

	selectedBoth := uiprocessing.NewScreenState(
		contracts.LanguageEnglish,
		contracts.LanguageJapanese,
		options,
	)
	screenCommander.SetText(&selectedBoth, "apple")
	if err := screenCommander.ToggleDictionary(&selectedBoth, "dict-a"); err != nil {
		t.Fatal(err)
	}
	if err := screenCommander.ToggleDictionary(&selectedBoth, "dict-b"); err != nil {
		t.Fatal(err)
	}
	if err := screenCommander.Execute(&selectedBoth); err != nil {
		t.Fatal(err)
	}

	if dictA.calls != 1 || dictB.calls != 2 {
		t.Fatalf("lookup calls after multiple selection: dict-a=%d dict-b=%d", dictA.calls, dictB.calls)
	}
	if len(selectedBoth.Result) != 1 || len(selectedBoth.Result[0].Candidates) != 2 {
		t.Fatalf("merged result = %#v", selectedBoth.Result)
	}

	rendered := uiprocessing.RenderScreen(selectedBoth)
	if !strings.Contains(rendered, "{りんごA [dict-a] / りんごB [dict-b]}") {
		t.Fatalf("rendered screen did not show merged candidates:\n%s", rendered)
	}
}
