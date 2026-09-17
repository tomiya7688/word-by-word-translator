package processing

import (
	"strings"
	"unicode"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
)

type EnglishTokenizer struct{}

func NewEnglishTokenizer() *EnglishTokenizer {
	return &EnglishTokenizer{}
}

func (t *EnglishTokenizer) Tokenize(text string) ([]contracts.Token, error) {
	runes := []rune(text)
	tokens := make([]contracts.Token, 0)

	for index := 0; index < len(runes); {
		current := runes[index]
		if unicode.IsSpace(current) {
			index++
			continue
		}

		if isWordRune(current) {
			start := index
			index++
			for index < len(runes) {
				if isWordRune(runes[index]) {
					index++
					continue
				}
				if isConnectorRune(runes[index]) && index+1 < len(runes) && isWordRune(runes[index+1]) {
					index++
					continue
				}
				break
			}

			surface := string(runes[start:index])
			tokens = append(tokens, contracts.Token{
				Surface:    surface,
				LookupUnit: strings.ToLower(surface),
				Kind:       contracts.TokenKindWord,
				Status:     contracts.TokenStatusSuccess,
			})
			continue
		}

		tokens = append(tokens, contracts.Token{
			Surface: string(current),
			Kind:    contracts.TokenKindSymbol,
			Status:  contracts.TokenStatusSuccess,
		})
		index++
	}

	return tokens, nil
}

func isWordRune(value rune) bool {
	return unicode.IsLetter(value) || unicode.IsDigit(value) || unicode.IsMark(value)
}

func isConnectorRune(value rune) bool {
	return value == '\'' || value == '’' || value == '-'
}
