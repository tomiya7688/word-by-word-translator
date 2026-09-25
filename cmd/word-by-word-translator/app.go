package main

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
	datacommander "github.com/tomiya7688/word-by-word-translator/applications/main/data/commander"
	datamessenger "github.com/tomiya7688/word-by-word-translator/applications/main/data/messenger"
	dataprocessing "github.com/tomiya7688/word-by-word-translator/applications/main/data/processing"
	processcommander "github.com/tomiya7688/word-by-word-translator/applications/main/process/commander"
	processmessenger "github.com/tomiya7688/word-by-word-translator/applications/main/process/messenger"
	processprocessing "github.com/tomiya7688/word-by-word-translator/applications/main/process/processing"
	kagomeprocessing "github.com/tomiya7688/word-by-word-translator/applications/main/process/processing/kagome"
	uicommander "github.com/tomiya7688/word-by-word-translator/applications/main/ui/commander"
	uimessenger "github.com/tomiya7688/word-by-word-translator/applications/main/ui/messenger"
	uiprocessing "github.com/tomiya7688/word-by-word-translator/applications/main/ui/processing"
)

type stringListFlag []string

func (values *stringListFlag) String() string {
	return strings.Join(*values, ",")
}

func (values *stringListFlag) Set(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("value must not be empty")
	}
	*values = append(*values, value)
	return nil
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("word-by-word-translator", flag.ContinueOnError)
	flags.SetOutput(stderr)

	var dictionaryManifests stringListFlag
	var selectedDictionaryIDs stringListFlag
	from := flags.String("from", "en", "source language: en or ja")
	to := flags.String("to", "ja", "target language: ja or en")
	compact := flags.Bool("compact", false, "print only the translated token sequence")
	flags.Var(&dictionaryManifests, "dictionary", "Dictionary Package v1 manifest path; repeat to load multiple dictionaries")
	flags.Var(&selectedDictionaryIDs, "use", "dictionary ID to use; repeat to select multiple dictionaries")
	flags.Usage = func() {
		fmt.Fprintln(stderr, "Usage: word-by-word-translator [options] [--] text")
		fmt.Fprintln(stderr)
		flags.PrintDefaults()
		fmt.Fprintln(stderr)
		fmt.Fprintln(stderr, "If text is omitted, input is read from stdin.")
	}

	if err := flags.Parse(args); err != nil {
		return 2
	}
	if len(dictionaryManifests) == 0 {
		fmt.Fprintln(stderr, "word-by-word-translator: at least one -dictionary manifest is required")
		return 2
	}

	text, err := readInput(flags.Args(), stdin)
	if err != nil {
		fmt.Fprintf(stderr, "word-by-word-translator: %v\n", err)
		return 1
	}
	if text == "" {
		fmt.Fprintln(stderr, "word-by-word-translator: translation text is required")
		return 2
	}

	source := contracts.Language(strings.TrimSpace(*from))
	target := contracts.Language(strings.TrimSpace(*to))

	dictionaries, options, err := loadDictionaries(dictionaryManifests)
	if err != nil {
		fmt.Fprintf(stderr, "word-by-word-translator: %v\n", err)
		return 1
	}

	tokenizers := processprocessing.NewTokenizerRegistry()
	if err := tokenizers.Register(contracts.LanguageEnglish, processprocessing.NewEnglishTokenizer()); err != nil {
		fmt.Fprintf(stderr, "word-by-word-translator: register English tokenizer: %v\n", err)
		return 1
	}
	if source == contracts.LanguageJapanese {
		japaneseTokenizer, err := kagomeprocessing.NewJapaneseTokenizer()
		if err != nil {
			fmt.Fprintf(stderr, "word-by-word-translator: initialize Japanese tokenizer: %v\n", err)
			return 1
		}
		if err := tokenizers.Register(contracts.LanguageJapanese, japaneseTokenizer); err != nil {
			fmt.Fprintf(stderr, "word-by-word-translator: register Japanese tokenizer: %v\n", err)
			return 1
		}
	}

	store := dataprocessing.NewDictionaryStore(dictionaries...)
	dataCommander := datacommander.NewDictionaryCommander(store)
	dataMessenger := datamessenger.NewDictionaryMessenger(dataCommander)
	processDictionaryMessenger := processmessenger.NewDictionaryMessenger(dataMessenger)
	translationProcessor := processprocessing.NewTranslationProcessor(tokenizers)
	translationCommander := processcommander.NewTranslateCommander(translationProcessor, processDictionaryMessenger)
	uiTranslationMessenger := uimessenger.NewTranslationMessenger(translationCommander)
	screenProcessor := uiprocessing.NewScreenProcessor()
	screenCommander := uicommander.NewScreenCommander(screenProcessor, uiTranslationMessenger)

	state := uiprocessing.NewScreenState(source, target, options)
	screenCommander.SetText(&state, text)

	useIDs := append([]string(nil), selectedDictionaryIDs...)
	if len(useIDs) == 0 {
		for _, option := range options {
			if option.SourceLanguage == source && option.TargetLanguage == target {
				useIDs = append(useIDs, option.ID)
			}
		}
	}
	for _, dictionaryID := range useIDs {
		if err := screenCommander.ToggleDictionary(&state, dictionaryID); err != nil {
			fmt.Fprintf(stderr, "word-by-word-translator: %v\n", err)
			return 2
		}
	}

	if err := screenCommander.Execute(&state); err != nil {
		fmt.Fprintf(stderr, "word-by-word-translator: %v\n", err)
		return 1
	}

	if *compact {
		fmt.Fprintln(stdout, uiprocessing.RenderTerminal(state.Result))
	} else {
		fmt.Fprintln(stdout, uiprocessing.RenderScreen(state))
	}
	return 0
}

func readInput(args []string, stdin io.Reader) (string, error) {
	if len(args) > 0 {
		return strings.TrimSpace(strings.Join(args, " ")), nil
	}
	content, err := io.ReadAll(stdin)
	if err != nil {
		return "", fmt.Errorf("read stdin: %w", err)
	}
	return strings.TrimSpace(string(content)), nil
}

func loadDictionaries(manifests []string) ([]dataprocessing.Dictionary, []uiprocessing.DictionaryOption, error) {
	dictionaries := make([]dataprocessing.Dictionary, 0, len(manifests))
	options := make([]uiprocessing.DictionaryOption, 0, len(manifests))
	seen := make(map[string]string)

	for _, manifestPath := range manifests {
		dictionary, err := dataprocessing.LoadFileDictionary(manifestPath)
		if err != nil {
			return nil, nil, fmt.Errorf("load dictionary %q: %w", manifestPath, err)
		}
		metadata := dictionary.Metadata()
		if previous, exists := seen[metadata.ID]; exists {
			return nil, nil, fmt.Errorf(
				"duplicate dictionary id %q in %q and %q",
				metadata.ID,
				previous,
				manifestPath,
			)
		}
		seen[metadata.ID] = manifestPath
		dictionaries = append(dictionaries, dictionary)
		options = append(options, uiprocessing.DictionaryOption{
			ID:             metadata.ID,
			Name:           metadata.Name,
			SourceLanguage: metadata.SourceLanguage,
			TargetLanguage: metadata.TargetLanguage,
		})
	}

	return dictionaries, options, nil
}
