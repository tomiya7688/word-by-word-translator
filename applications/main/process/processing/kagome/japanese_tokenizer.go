package kagome

import (
	"strings"
	"unicode"

	"github.com/ikawaha/kagome-dict/ipa"
	kagometokenizer "github.com/ikawaha/kagome/v2/tokenizer"
	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
)

// JapaneseTokenizer adapts Kagome + IPADIC to the application's Tokenizer contract.
// Keep this adapter in its own package so builds that do not select it do not
// automatically pull the analyzer dictionary into the final binary.
type JapaneseTokenizer struct {
	tokenizer *kagometokenizer.Tokenizer
}

func NewJapaneseTokenizer() (*JapaneseTokenizer, error) {
	tokenizer, err := kagometokenizer.New(ipa.Dict(), kagometokenizer.OmitBosEos())
	if err != nil {
		return nil, err
	}
	return &JapaneseTokenizer{tokenizer: tokenizer}, nil
}

func (t *JapaneseTokenizer) Tokenize(text string) ([]contracts.Token, error) {
	morphemes := t.tokenizer.Tokenize(text)
	tokens := make([]contracts.Token, 0, len(morphemes))

	for _, morpheme := range morphemes {
		surface := morpheme.Surface
		if strings.TrimSpace(surface) == "" {
			continue
		}

		if isSymbol(morpheme) {
			tokens = append(tokens, contracts.Token{
				Surface: surface,
				Kind:    contracts.TokenKindSymbol,
				Status:  contracts.TokenStatusSuccess,
			})
			continue
		}

		lemma := surface
		if baseForm, ok := morpheme.BaseForm(); ok && baseForm != "" && baseForm != "*" {
			lemma = baseForm
		}

		tokens = append(tokens, contracts.Token{
			Surface:    surface,
			Lemma:      lemma,
			LookupUnit: lemma,
			Kind:       contracts.TokenKindWord,
			Status:     contracts.TokenStatusSuccess,
		})
	}

	return tokens, nil
}

func isSymbol(token kagometokenizer.Token) bool {
	pos := token.POS()
	if len(pos) > 0 && pos[0] == "記号" {
		return true
	}

	seen := false
	for _, value := range token.Surface {
		seen = true
		if !unicode.IsPunct(value) && !unicode.IsSymbol(value) {
			return false
		}
	}
	return seen
}
