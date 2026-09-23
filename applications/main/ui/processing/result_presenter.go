package processing

import "github.com/tomiya7688/word-by-word-translator/applications/main/contracts"

type TokenTone string

const (
	TokenToneNormal     TokenTone = "normal"
	TokenToneUnresolved TokenTone = "unresolved"
)

type PresentedToken struct {
	Text string
	Tone TokenTone
}

func PresentTranslation(response contracts.TranslationResponse) []PresentedToken {
	presented := make([]PresentedToken, 0, len(response.Tokens))
	for _, translated := range response.Tokens {
		presented = append(presented, presentToken(translated))
	}
	return presented
}

func presentToken(translated contracts.TranslatedToken) PresentedToken {
	if translated.Status != contracts.TokenStatusSuccess {
		return PresentedToken{
			Text: translated.Token.Surface,
			Tone: TokenToneUnresolved,
		}
	}

	for _, candidate := range translated.Candidates {
		if candidate.Translation != "" {
			return PresentedToken{
				Text: candidate.Translation,
				Tone: TokenToneNormal,
			}
		}
	}

	return PresentedToken{
		Text: translated.Token.Surface,
		Tone: TokenToneNormal,
	}
}
