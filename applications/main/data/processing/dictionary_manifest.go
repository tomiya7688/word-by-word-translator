package processing

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
)

const DictionaryManifestSchemaVersion = 1
const DictionaryDataFormatTSVV1 = "tsv-v1"

type DictionaryManifest struct {
	SchemaVersion         int                `json:"schema_version"`
	ID                    string             `json:"id"`
	Name                  string             `json:"name"`
	Source                string             `json:"source"`
	License               string             `json:"license"`
	SourceLanguage        contracts.Language `json:"source_language"`
	TargetLanguage        contracts.Language `json:"target_language"`
	CommercialUse         bool               `json:"commercial_use"`
	NoncommercialUse      bool               `json:"noncommercial_use"`
	Modification          bool               `json:"modification"`
	Redistribution        bool               `json:"redistribution"`
	AttributionRequired   bool               `json:"attribution_required"`
	LicenseNoticeRequired bool               `json:"license_notice_required"`
	ShareAlike            bool               `json:"share_alike"`
	DataFile              string             `json:"data_file"`
	Format                string             `json:"format"`
	CaseSensitive         bool               `json:"case_sensitive"`
}

func LoadDictionaryManifest(manifestPath string) (DictionaryManifest, error) {
	file, err := os.Open(manifestPath)
	if err != nil {
		return DictionaryManifest{}, fmt.Errorf("open dictionary manifest: %w", err)
	}
	defer file.Close()

	var manifest DictionaryManifest
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&manifest); err != nil {
		return DictionaryManifest{}, fmt.Errorf("parse dictionary manifest: %w", err)
	}
	if err := ensureDictionaryManifestEOF(decoder); err != nil {
		return DictionaryManifest{}, err
	}
	if err := manifest.Validate(); err != nil {
		return DictionaryManifest{}, err
	}
	return manifest, nil
}

func ensureDictionaryManifestEOF(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return fmt.Errorf("dictionary manifest contains multiple JSON values")
		}
		return fmt.Errorf("parse dictionary manifest trailing data: %w", err)
	}
	return nil
}

func (m DictionaryManifest) Validate() error {
	if m.SchemaVersion != DictionaryManifestSchemaVersion {
		return fmt.Errorf("unsupported dictionary manifest schema_version %d", m.SchemaVersion)
	}
	if strings.TrimSpace(m.ID) == "" || strings.TrimSpace(m.Name) == "" {
		return fmt.Errorf("dictionary manifest id and name are required")
	}
	if strings.TrimSpace(m.Source) == "" || strings.TrimSpace(m.License) == "" {
		return fmt.Errorf("dictionary manifest source and license are required")
	}
	if m.SourceLanguage == "" || m.TargetLanguage == "" {
		return fmt.Errorf("dictionary manifest source_language and target_language are required")
	}
	if err := validateDictionaryDataFile(m.DataFile); err != nil {
		return err
	}
	if m.Format != DictionaryDataFormatTSVV1 {
		return fmt.Errorf("unsupported dictionary format %q", m.Format)
	}
	return nil
}

func (m DictionaryManifest) Metadata() contracts.DictionaryMetadata {
	return contracts.DictionaryMetadata{
		ID:                    m.ID,
		Name:                  m.Name,
		Source:                m.Source,
		License:               m.License,
		SourceLanguage:        m.SourceLanguage,
		TargetLanguage:        m.TargetLanguage,
		CommercialUse:         m.CommercialUse,
		NoncommercialUse:      m.NoncommercialUse,
		Modification:          m.Modification,
		Redistribution:        m.Redistribution,
		AttributionRequired:   m.AttributionRequired,
		LicenseNoticeRequired: m.LicenseNoticeRequired,
		ShareAlike:            m.ShareAlike,
		ReleaseTier:           m.ReleaseTier(),
	}
}

func (m DictionaryManifest) ReleaseTier() contracts.ReleaseTier {
	if !m.CommercialUse || !m.NoncommercialUse || !m.Redistribution {
		return contracts.ReleaseTierUnsupported
	}
	if m.Modification && !m.AttributionRequired && !m.LicenseNoticeRequired && !m.ShareAlike {
		return contracts.ReleaseTierMIT
	}
	return contracts.ReleaseTierFull
}

func (m DictionaryManifest) ResolveDataPath(manifestPath string) string {
	normalized := strings.ReplaceAll(m.DataFile, "\\", "/")
	return filepath.Join(filepath.Dir(manifestPath), filepath.FromSlash(path.Clean(normalized)))
}

func validateDictionaryDataFile(dataFile string) error {
	normalized := strings.ReplaceAll(strings.TrimSpace(dataFile), "\\", "/")
	if normalized == "" {
		return fmt.Errorf("dictionary manifest data_file is required")
	}
	if strings.HasPrefix(normalized, "/") || strings.Contains(normalized, ":") {
		return fmt.Errorf("dictionary manifest data_file must be relative")
	}
	clean := path.Clean(normalized)
	if clean == ".." || strings.HasPrefix(clean, "../") {
		return fmt.Errorf("dictionary manifest data_file must stay inside the manifest directory")
	}
	return nil
}
