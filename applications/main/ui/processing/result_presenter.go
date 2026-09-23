package processing

import "github.com/tomiya7688/word-by-word-translator/applications/main/contracts"

type TokenTone string

const (
	TokenToneNormal     TokenTone = "normal"
	TokenToneUnresolved TokenTone = "unresolved"
)

type PresentedCandidate struct {
	Translation   string
	DictionaryIDs []string
}

type PresentedToken struct {
	Text       string
	Tone       TokenTone
	Candidates []PresentedCandidate
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

	candidates := make([]PresentedCandidate, 0, len(translated.Candidates))
	for _, candidate := range translated.Candidates {
		if candidate.Translation == "" {
			continue
		}
		candidates = append(candidates, PresentedCandidate{
			Translation:   candidate.Translation,
			DictionaryIDs: append([]string(nil), candidate.DictionaryIDs...),
		})
	}
	if len(candidates) > 0 {
		return PresentedToken{
			Text:       candidates[0].Translation,
			Tone:       TokenToneNormal,
			Candidates: candidates,
		}
	}

	return PresentedToken{
		Text: translated.Token.Surface,
		Tone: TokenToneNormal,
	}
}
