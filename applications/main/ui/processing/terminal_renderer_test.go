package processing

import "testing"

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
