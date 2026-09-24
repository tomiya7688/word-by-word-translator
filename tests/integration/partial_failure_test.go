package integration

import (
	"errors"
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

type integrationFailingDictionary struct {
	metadata contracts.DictionaryMetadata
}

func (d *integrationFailingDictionary) Metadata() contracts.DictionaryMetadata {
	return d.metadata
}

func (d *integrationFailingDictionary) Lookup(string) ([]contracts.DictionaryEntry, error) {
	return nil, errors.New("simulated lookup failure")
}

func TestPipelineKeepsSuccessfulWordsWhenOneDictionaryFails(t *testing.T) {
	broken := &integrationFailingDictionary{
		metadata: contracts.DictionaryMetadata{
			ID:             "broken",
			Name:           "Broken Dictionary",
			SourceLanguage: contracts.LanguageEnglish,
			TargetLanguage: contracts.LanguageJapanese,
		},
	}
	working := dataprocessing.NewMemoryDictionary(
		contracts.DictionaryMetadata{
			ID:             "working",
			Name:           "Working Dictionary",
			SourceLanguage: contracts.LanguageEnglish,
			TargetLanguage: contracts.LanguageJapanese,
		},
		map[string][]string{
			"apple": {"りんご"},
		},
	)

	store := dataprocessing.NewDictionaryStore(broken, working)
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
	state := uiprocessing.NewScreenState(
		contracts.LanguageEnglish,
		contracts.LanguageJapanese,
		[]uiprocessing.DictionaryOption{
			{
				ID:             "broken",
				Name:           "Broken Dictionary",
				SourceLanguage: contracts.LanguageEnglish,
				TargetLanguage: contracts.LanguageJapanese,
				Selected:       true,
			},
			{
				ID:             "working",
				Name:           "Working Dictionary",
				SourceLanguage: contracts.LanguageEnglish,
				TargetLanguage: contracts.LanguageJapanese,
				Selected:       true,
			},
		},
	)
	screenCommander.SetText(&state, "apple qwertymonster")

	if err := screenCommander.Execute(&state); err != nil {
		t.Fatalf("Execute returned error despite partial dictionary failure: %v", err)
	}
	if state.Error != "" {
		t.Fatalf("screen error = %q", state.Error)
	}
	if len(state.Result) != 2 {
		t.Fatalf("result count = %d, want 2: %#v", len(state.Result), state.Result)
	}

	if state.Result[0].Text != "りんご" || state.Result[0].Tone != uiprocessing.TokenToneNormal {
		t.Fatalf("successful token = %#v", state.Result[0])
	}
	if len(state.Result[0].Candidates) != 1 ||
		state.Result[0].Candidates[0].Translation != "りんご" ||
		len(state.Result[0].Candidates[0].DictionaryIDs) != 1 ||
		state.Result[0].Candidates[0].DictionaryIDs[0] != "working" {
		t.Fatalf("successful candidates = %#v", state.Result[0].Candidates)
	}

	if state.Result[1].Text != "qwertymonster" ||
		state.Result[1].Tone != uiprocessing.TokenToneUnresolved {
		t.Fatalf("unresolved token = %#v", state.Result[1])
	}

	rendered := uiprocessing.RenderScreen(state)
	if !strings.Contains(rendered, "りんご [working]") {
		t.Fatalf("rendered screen missing successful translation:\n%s", rendered)
	}
	if !strings.Contains(rendered, "\x1b[31mqwertymonster\x1b[0m") {
		t.Fatalf("rendered screen did not preserve unresolved word in red:\n%s", rendered)
	}
}
