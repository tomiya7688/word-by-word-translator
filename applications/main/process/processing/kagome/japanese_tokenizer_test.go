package kagome

import (
	"testing"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
)

func TestJapaneseTokenizerUsesDictionaryBaseForm(t *testing.T) {
	tokenizer, err := NewJapaneseTokenizer()
	if err != nil {
		t.Fatal(err)
	}

	tokens, err := tokenizer.Tokenize("私はりんごを食べました。")
	if err != nil {
		t.Fatal(err)
	}

	var foundVerb bool
	var foundPeriod bool
	for _, token := range tokens {
		switch token.Surface {
		case "食べ":
			foundVerb = true
			if token.Lemma != "食べる" || token.LookupUnit != "食べる" {
				t.Fatalf("verb token = %#v", token)
			}
			if token.Kind != contracts.TokenKindWord {
				t.Fatalf("verb kind = %q", token.Kind)
			}
		case "。":
			foundPeriod = true
			if token.Kind != contracts.TokenKindSymbol || token.LookupUnit != "" {
				t.Fatalf("period token = %#v", token)
			}
		}
	}

	if !foundVerb {
		t.Fatal("inflected verb token was not found")
	}
	if !foundPeriod {
		t.Fatal("period token was not found")
	}
}
