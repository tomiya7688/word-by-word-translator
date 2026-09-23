package processing

import (
	"fmt"
	"strings"
)

const (
	ansiRed   = "\x1b[31m"
	ansiReset = "\x1b[0m"
)

func RenderTerminal(tokens []PresentedToken) string {
	rendered := make([]string, 0, len(tokens))
	for _, token := range tokens {
		rendered = append(rendered, renderPrimaryToken(token))
	}
	return strings.Join(rendered, " ")
}

func RenderScreen(state ScreenState) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Text: %s\n", state.Text)
	fmt.Fprintf(&builder, "Direction: %s -> %s\n", state.SourceLanguage, state.TargetLanguage)
	builder.WriteString("Dictionaries:\n")
	for _, option := range state.Dictionaries {
		if !dictionaryMatchesDirection(option, state.SourceLanguage, state.TargetLanguage) {
			continue
		}
		marker := " "
		if option.Selected {
			marker = "x"
		}
		name := option.Name
		if name == "" {
			name = option.ID
		}
		fmt.Fprintf(&builder, "[%s] %s (%s)\n", marker, name, option.ID)
	}
	builder.WriteString("Result: ")
	if len(state.Result) == 0 {
		builder.WriteString("<empty>")
	} else {
		rendered := make([]string, 0, len(state.Result))
		for _, token := range state.Result {
			rendered = append(rendered, renderDetailedToken(token))
		}
		builder.WriteString(strings.Join(rendered, " "))
	}
	if state.Error != "" {
		fmt.Fprintf(&builder, "\nError: %s", state.Error)
	}
	return builder.String()
}

func renderPrimaryToken(token PresentedToken) string {
	text := token.Text
	if token.Tone == TokenToneUnresolved {
		return ansiRed + text + ansiReset
	}
	return text
}

func renderDetailedToken(token PresentedToken) string {
	if token.Tone == TokenToneUnresolved || len(token.Candidates) == 0 {
		return renderPrimaryToken(token)
	}
	if len(token.Candidates) == 1 {
		return renderCandidate(token.Candidates[0])
	}

	choices := make([]string, 0, len(token.Candidates))
	for _, candidate := range token.Candidates {
		choices = append(choices, renderCandidate(candidate))
	}
	return "{" + strings.Join(choices, " / ") + "}"
}

func renderCandidate(candidate PresentedCandidate) string {
	if len(candidate.DictionaryIDs) == 0 {
		return candidate.Translation
	}
	return fmt.Sprintf("%s [%s]", candidate.Translation, strings.Join(candidate.DictionaryIDs, ","))
}
