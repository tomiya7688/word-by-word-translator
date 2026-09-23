package processing

import (
	"strings"
	"testing"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
)

func TestRenderTerminalColorsOnlyUnresolvedTokensRed(t *testing.T) {
	tokens := []PresentedToken{
		{Text: "私", Tone: TokenToneNormal},
		{Text: "見た", Tone: TokenToneNormal},
		{Text: "qwertymonster", Tone: TokenToneUnresolved},
		{Text: "昨日", Tone: TokenToneNormal},
		{Text: ".", Tone: TokenToneNormal},
	}

	got := RenderTerminal(tokens)
	want := "私 見た \x1b[31mqwertymonster\x1b[0m 昨日 ."
	if got != want {
		t.Fatalf("rendered = %q, want %q", got, want)
	}
}

func TestRenderTerminalHandlesEmptyResult(t *testing.T) {
	if got := RenderTerminal(nil); got != "" {
		t.Fatalf("rendered = %q, want empty", got)
	}
}

func TestRenderScreenShowsInputDirectionSelectionAndMergedCandidates(t *testing.T) {
	state := ScreenState{
		Text:           "saw qwertymonster",
		SourceLanguage: contracts.LanguageEnglish,
		TargetLanguage: contracts.LanguageJapanese,
		Dictionaries: []DictionaryOption{
			{
				ID:             "a",
				Name:           "Dictionary A",
				SourceLanguage: contracts.LanguageEnglish,
				TargetLanguage: contracts.LanguageJapanese,
				Selected:       true,
			},
			{
				ID:             "b",
				Name:           "Dictionary B",
				SourceLanguage: contracts.LanguageEnglish,
				TargetLanguage: contracts.LanguageJapanese,
				Selected:       true,
			},
		},
		Result: []PresentedToken{
			{
				Text: "見た",
				Tone: TokenToneNormal,
				Candidates: []PresentedCandidate{
					{Translation: "見た", DictionaryIDs: []string{"a", "b"}},
					{Translation: "見ました", DictionaryIDs: []string{"b"}},
				},
			},
			{Text: "qwertymonster", Tone: TokenToneUnresolved},
		},
	}

	rendered := RenderScreen(state)
	for _, expected := range []string{
		"Text: saw qwertymonster",
		"Direction: en -> ja",
		"[x] Dictionary A (a)",
		"[x] Dictionary B (b)",
		"{見た [a,b] / 見ました [b]}",
		"\x1b[31mqwertymonster\x1b[0m",
	} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("rendered screen missing %q:\n%s", expected, rendered)
		}
	}
}
