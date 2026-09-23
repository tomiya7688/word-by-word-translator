package processing

import "strings"

const (
	ansiRed   = "\x1b[31m"
	ansiReset = "\x1b[0m"
)

func RenderTerminal(tokens []PresentedToken) string {
	rendered := make([]string, 0, len(tokens))
	for _, token := range tokens {
		text := token.Text
		if token.Tone == TokenToneUnresolved {
			text = ansiRed + text + ansiReset
		}
		rendered = append(rendered, text)
	}
	return strings.Join(rendered, " ")
}
