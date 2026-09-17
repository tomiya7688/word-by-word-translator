package processing

import "github.com/tomiya7688/word-by-word-translator/applications/main/contracts"

type Tokenizer interface {
	Tokenize(text string) ([]contracts.Token, error)
}
