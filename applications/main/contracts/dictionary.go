package contracts

type ReleaseTier string

const (
	ReleaseTierMIT         ReleaseTier = "mit"
	ReleaseTierFull        ReleaseTier = "full"
	ReleaseTierUnsupported ReleaseTier = "unsupported"
)

type DictionaryMetadata struct {
	ID                    string
	Name                  string
	Source                string
	License               string
	SourceLanguage        Language
	TargetLanguage        Language
	CommercialUse         bool
	NoncommercialUse      bool
	Modification          bool
	Redistribution        bool
	AttributionRequired   bool
	LicenseNoticeRequired bool
	ShareAlike            bool
	ReleaseTier           ReleaseTier
}

type DictionaryEntry struct {
	Headword    string
	Translation string
}

type DictionaryLookupResult struct {
	Metadata DictionaryMetadata
	Entries  []DictionaryEntry
	Error    string
}

type TokenDictionaryResult struct {
	TokenIndex   int
	Dictionaries []DictionaryLookupResult
}

type DictionaryBatchLookupRequest struct {
	Tokens         []Token
	SourceLanguage Language
	TargetLanguage Language
	DictionaryIDs  []string
}

type DictionaryBatchLookupResponse struct {
	Tokens []TokenDictionaryResult
}
