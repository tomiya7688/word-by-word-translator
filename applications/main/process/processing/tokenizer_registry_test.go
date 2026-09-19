package processing

import (
	"errors"
	"testing"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
)

func TestTokenizerRegistryRoutesByLanguage(t *testing.T) {
	registry := NewTokenizerRegistry()
	if err := registry.Register(contracts.LanguageEnglish, NewEnglishTokenizer()); err != nil {
		t.Fatal(err)
	}

	tokens, err := registry.Tokenize(contracts.LanguageEnglish, "Hello.")
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 2 || tokens[0].LookupUnit != "hello" {
		t.Fatalf("tokens = %#v", tokens)
	}
}

func TestTokenizerRegistryRejectsDuplicateLanguage(t *testing.T) {
	registry := NewTokenizerRegistry()
	if err := registry.Register(contracts.LanguageEnglish, NewEnglishTokenizer()); err != nil {
		t.Fatal(err)
	}
	if err := registry.Register(contracts.LanguageEnglish, NewEnglishTokenizer()); err == nil {
		t.Fatal("duplicate tokenizer registration was accepted")
	}
}

func TestTokenizerRegistryReportsUnsupportedLanguage(t *testing.T) {
	registry := NewTokenizerRegistry()
	_, err := registry.Tokenize(contracts.LanguageJapanese, "日本語")
	if !errors.Is(err, ErrTokenizerNotRegistered) {
		t.Fatalf("error = %v", err)
	}
}
