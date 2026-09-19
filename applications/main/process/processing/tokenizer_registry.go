package processing

import (
	"errors"
	"fmt"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
)

var ErrTokenizerNotRegistered = errors.New("tokenizer not registered")

type TokenizerRegistry struct {
	tokenizers map[contracts.Language]Tokenizer
}

func NewTokenizerRegistry() *TokenizerRegistry {
	return &TokenizerRegistry{tokenizers: make(map[contracts.Language]Tokenizer)}
}

func (r *TokenizerRegistry) Register(language contracts.Language, tokenizer Tokenizer) error {
	if language == "" {
		return errors.New("tokenizer language is required")
	}
	if tokenizer == nil {
		return errors.New("tokenizer is required")
	}
	if _, exists := r.tokenizers[language]; exists {
		return fmt.Errorf("tokenizer already registered for language %q", language)
	}
	r.tokenizers[language] = tokenizer
	return nil
}

func (r *TokenizerRegistry) Tokenize(language contracts.Language, text string) ([]contracts.Token, error) {
	tokenizer, exists := r.tokenizers[language]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrTokenizerNotRegistered, language)
	}
	return tokenizer.Tokenize(text)
}
