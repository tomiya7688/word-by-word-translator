package processing

import (
	"testing"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
)

func TestEnglishTokenizerPreservesWordsAndSymbols(t *testing.T) {
	tokenizer := NewEnglishTokenizer()
	tokens, err := tokenizer.Tokenize("I saw qwertymonster yesterday.")
	if err != nil {
		t.Fatalf("Tokenize returned error: %v", err)
	}

	wantSurface := []string{"I", "saw", "qwertymonster", "yesterday", "."}
	if len(tokens) != len(wantSurface) {
		t.Fatalf("token count = %d, want %d", len(tokens), len(wantSurface))
	}
	for index, want := range wantSurface {
		if tokens[index].Surface != want {
			t.Fatalf("token[%d].Surface = %q, want %q", index, tokens[index].Surface, want)
		}
	}
	if tokens[0].LookupUnit != "i" {
		t.Fatalf("first lookup unit = %q, want i", tokens[0].LookupUnit)
	}
	if tokens[4].Kind != contracts.TokenKindSymbol || tokens[4].LookupUnit != "" {
		t.Fatalf("symbol token = %#v, want symbol with no lookup unit", tokens[4])
	}
}
