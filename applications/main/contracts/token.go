package contracts

type TokenKind string

const (
	TokenKindWord   TokenKind = "word"
	TokenKindSymbol TokenKind = "symbol"
)

type TokenStatus string

const (
	TokenStatusSuccess            TokenStatus = "success"
	TokenStatusSegmentationFailed TokenStatus = "segmentation_failed"
	TokenStatusDictionaryNotFound TokenStatus = "dictionary_not_found"
	TokenStatusLookupError        TokenStatus = "lookup_error"
)

type Token struct {
	Surface    string
	Lemma      string
	LookupUnit string
	Subtokens  []string
	Kind       TokenKind
	Status     TokenStatus
}
