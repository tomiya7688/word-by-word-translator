package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
	dataprocessing "github.com/tomiya7688/word-by-word-translator/applications/main/data/processing"
)

const (
	releaseSchemaVersion = 1
	dictionaryManifest    = "dictionary.json"
	attributionFile       = "ATTRIBUTION.txt"
	licenseFile           = "LICENSE.txt"
)

type releaseKind string

const (
	releaseMIT  releaseKind = "mit"
	releaseFull releaseKind = "full"
)

type dictionarySource struct {
	Manifest     dataprocessing.DictionaryManifest
	ManifestPath string
	PackageDir   string
}

type releaseDictionary struct {
	ID                    string                `json:"id"`
	Name                  string                `json:"name"`
	Source                string                `json:"source"`
	License               string                `json:"license"`
	ReleaseTier           contracts.ReleaseTier `json:"release_tier"`
	AttributionRequired   bool                  `json:"attribution_required"`
	LicenseNoticeRequired bool                  `json:"license_notice_required"`
	ShareAlike            bool                  `json:"share_alike"`
	NoticeFiles           []string              `json:"notice_files,omitempty"`
}

type releaseAnalyzer struct {
	ID          string   `json:"id"`
	Version     string   `json:"version"`
	NoticeFiles []string `json:"notice_files"`
}

type releaseManifest struct {
	SchemaVersion int                 `json:"schema_version"`
	ReleaseTier   releaseKind         `json:"release_tier"`
	Dictionaries  []releaseDictionary `json:"dictionaries"`
	Analyzers     []releaseAnalyzer   `json:"analyzers,omitempty"`
}

var fullReleaseAnalyzers = []releaseAnalyzer{
	{
		ID:      "kagome",
		Version: "v2.9.9",
		NoticeFiles: []string{
			"licenses/analyzers/kagome/LICENSE.txt",
		},
	},
	{
		ID:      "kagome-dict",
		Version: "v1.1.0",
		NoticeFiles: []string{
			"licenses/analyzers/kagome-dict/LICENSE.txt",
		},
	},
	{
		ID:      "kagome-dict-ipa",
		Version: "v1.2.0",
		NoticeFiles: []string{
			"licenses/analyzers/kagome-dict-ipa/LICENSE.txt",
			"licenses/analyzers/kagome-dict-ipa/NOTICE.txt",
		},
	},
}

func main() {
	dictionaries := flag.String("dictionaries", "", "root directory containing Dictionary Package v1 directories")
	output := flag.String("output", "", "release output root")
	appLicense := flag.String("app-license", "", "application LICENSE path")
	analyzerNotices := flag.String("analyzer-notices", "", "Full Release analyzer notice root")
	tier := flag.String("tier", "all", "release tier to build: mit, full, or all")
	flag.Parse()

	if err := run(*dictionaries, *output, *appLicense, *analyzerNotices, *tier); err != nil {
		fmt.Fprintln(os.Stderr, "release-builder:", err)
		os.Exit(1)
	}
}

func run(dictionariesRoot, outputRoot, appLicensePath, analyzerNoticesRoot, tier string) error {
	if strings.TrimSpace(dictionariesRoot) == "" {
		return errors.New("dictionaries is required")
	}
	if strings.TrimSpace(outputRoot) == "" {
		return errors.New("output is required")
	}
	if strings.TrimSpace(appLicensePath) == "" {
		return errors.New("app-license is required")
	}

	sources, err := discoverDictionaries(dictionariesRoot)
	if err != nil {
		return err
	}

	switch strings.ToLower(strings.TrimSpace(tier)) {
	case "mit":
		return buildRelease(releaseMIT, sources, outputRoot, appLicensePath, analyzerNoticesRoot)
	case "full":
		return buildRelease(releaseFull, sources, outputRoot, appLicensePath, analyzerNoticesRoot)
	case "all":
		if err := buildRelease(releaseMIT, sources, outputRoot, appLicensePath, analyzerNoticesRoot); err != nil {
			return err
		}
		return buildRelease(releaseFull, sources, outputRoot, appLicensePath, analyzerNoticesRoot)
	default:
		return fmt.Errorf("unsupported release tier %q", tier)
	}
}

func discoverDictionaries(root string) ([]dictionarySource, error) {
	seen := make(map[string]string)
	sources := make([]dictionarySource, 0)

	err := filepath.WalkDir(root, func(current string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || entry.Name() != dictionaryManifest {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("dictionary manifest must not be a symlink: %s", current)
		}

		manifest, err := dataprocessing.LoadDictionaryManifest(current)
		if err != nil {
			return fmt.Errorf("load %s: %w", current, err)
		}
		if err := validateReleaseComponentID(manifest.ID); err != nil {
			return err
		}
		if previous, exists := seen[manifest.ID]; exists {
			return fmt.Errorf("duplicate dictionary id %q: %s and %s", manifest.ID, previous, current)
		}
		seen[manifest.ID] = current
		sources = append(sources, dictionarySource{
			Manifest:     manifest,
			ManifestPath: current,
			PackageDir:   filepath.Dir(current),
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("discover dictionaries: %w", err)
	}
	if len(sources) == 0 {
		return nil, fmt.Errorf("no %s files found under %s", dictionaryManifest, root)
	}

	sort.Slice(sources, func(i, j int) bool {
		return sources[i].Manifest.ID < sources[j].Manifest.ID
	})
	return sources, nil
}

func buildRelease(kind releaseKind, sources []dictionarySource, outputRoot, appLicensePath, analyzerNoticesRoot string) error {
	releaseDir := filepath.Join(outputRoot, string(kind))
	if err := os.RemoveAll(releaseDir); err != nil {
		return fmt.Errorf("clean %s release: %w", kind, err)
	}
	if err := os.MkdirAll(releaseDir, 0o755); err != nil {
		return fmt.Errorf("create %s release: %w", kind, err)
	}
	if err := copyRegularFile(appLicensePath, filepath.Join(releaseDir, "LICENSE")); err != nil {
		return fmt.Errorf("copy application license: %w", err)
	}

	included := make([]releaseDictionary, 0)
	for _, source := range sources {
		tier := source.Manifest.ReleaseTier()
		if !includeDictionary(kind, tier) {
			continue
		}
		releaseDictionary, err := copyDictionaryPackage(kind, source, releaseDir)
		if err != nil {
			return err
		}
		included = append(included, releaseDictionary)
	}

	var analyzers []releaseAnalyzer
	var analyzerFiles []string
	if kind == releaseFull {
		if strings.TrimSpace(analyzerNoticesRoot) == "" {
			return errors.New("analyzer-notices is required for Full Release")
		}
		var err error
		analyzerFiles, err = copyAnalyzerNotices(analyzerNoticesRoot, releaseDir)
		if err != nil {
			return err
		}
		analyzers = append([]releaseAnalyzer(nil), fullReleaseAnalyzers...)
	}

	manifest := releaseManifest{
		SchemaVersion: releaseSchemaVersion,
		ReleaseTier:   kind,
		Dictionaries:  included,
		Analyzers:     analyzers,
	}
	if err := writeJSON(filepath.Join(releaseDir, "release.json"), manifest); err != nil {
		return err
	}
	if err := writeThirdPartyLicenses(filepath.Join(releaseDir, "THIRD_PARTY_LICENSES.md"), manifest, analyzerFiles); err != nil {
		return err
	}
	return nil
}

func includeDictionary(kind releaseKind, tier contracts.ReleaseTier) bool {
	switch kind {
	case releaseMIT:
		return tier == contracts.ReleaseTierMIT
	case releaseFull:
		return tier == contracts.ReleaseTierMIT || tier == contracts.ReleaseTierFull
	default:
		return false
	}
}

func copyDictionaryPackage(kind releaseKind, source dictionarySource, releaseDir string) (releaseDictionary, error) {
	manifest := source.Manifest
	tier := manifest.ReleaseTier()
	destination := filepath.Join(releaseDir, "dictionaries", manifest.ID)

	if err := copyRegularFile(source.ManifestPath, filepath.Join(destination, dictionaryManifest)); err != nil {
		return releaseDictionary{}, fmt.Errorf("copy dictionary %s manifest: %w", manifest.ID, err)
	}

	dataRelative := filepath.FromSlash(path.Clean(strings.ReplaceAll(manifest.DataFile, "\\", "/")))
	if err := copyRegularFile(manifest.ResolveDataPath(source.ManifestPath), filepath.Join(destination, dataRelative)); err != nil {
		return releaseDictionary{}, fmt.Errorf("copy dictionary %s data: %w", manifest.ID, err)
	}

	noticeFiles := make([]string, 0, 2)
	if kind == releaseFull && tier == contracts.ReleaseTierFull {
		if manifest.AttributionRequired {
			if err := copyRequiredDictionaryNotice(source, destination, attributionFile); err != nil {
				return releaseDictionary{}, err
			}
			noticeFiles = append(noticeFiles, attributionFile)
		}
		if manifest.LicenseNoticeRequired || manifest.ShareAlike {
			if err := copyRequiredDictionaryNotice(source, destination, licenseFile); err != nil {
				return releaseDictionary{}, err
			}
			noticeFiles = append(noticeFiles, licenseFile)
		}
	}

	return releaseDictionary{
		ID:                    manifest.ID,
		Name:                  manifest.Name,
		Source:                manifest.Source,
		License:               manifest.License,
		ReleaseTier:           tier,
		AttributionRequired:   manifest.AttributionRequired,
		LicenseNoticeRequired: manifest.LicenseNoticeRequired,
		ShareAlike:            manifest.ShareAlike,
		NoticeFiles:           noticeFiles,
	}, nil
}

func copyRequiredDictionaryNotice(source dictionarySource, destination, name string) error {
	sourcePath := filepath.Join(source.PackageDir, name)
	if err := copyRegularFile(sourcePath, filepath.Join(destination, name)); err != nil {
		return fmt.Errorf("Full Release dictionary %s requires %s: %w", source.Manifest.ID, name, err)
	}
	return nil
}

func copyAnalyzerNotices(sourceRoot, releaseDir string) ([]string, error) {
	for _, analyzer := range fullReleaseAnalyzers {
		for _, releasePath := range analyzer.NoticeFiles {
			relative := strings.TrimPrefix(releasePath, "licenses/analyzers/")
			sourcePath := filepath.Join(sourceRoot, filepath.FromSlash(relative))
			info, err := os.Lstat(sourcePath)
			if err != nil {
				return nil, fmt.Errorf("Full Release analyzer notice missing: %s: %w", relative, err)
			}
			if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
				return nil, fmt.Errorf("Full Release analyzer notice must be a regular file: %s", relative)
			}
		}
	}

	destinationRoot := filepath.Join(releaseDir, "licenses", "analyzers")
	files, err := copyDirectory(sourceRoot, destinationRoot)
	if err != nil {
		return nil, fmt.Errorf("copy analyzer notices: %w", err)
	}
	for index := range files {
		files[index] = filepath.ToSlash(filepath.Join("licenses", "analyzers", files[index]))
	}
	sort.Strings(files)
	return files, nil
}

func copyDirectory(sourceRoot, destinationRoot string) ([]string, error) {
	copied := make([]string, 0)
	err := filepath.WalkDir(sourceRoot, func(current string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if current == sourceRoot {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlinks are not allowed: %s", current)
		}
		relative, err := filepath.Rel(sourceRoot, current)
		if err != nil {
			return err
		}
		destination := filepath.Join(destinationRoot, relative)
		if entry.IsDir() {
			return os.MkdirAll(destination, 0o755)
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("unsupported file type: %s", current)
		}
		if err := copyRegularFile(current, destination); err != nil {
			return err
		}
		copied = append(copied, filepath.ToSlash(relative))
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(copied) == 0 {
		return nil, fmt.Errorf("no notice files found under %s", sourceRoot)
	}
	return copied, nil
}

func copyRegularFile(source, destination string) error {
	info, err := os.Lstat(source)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", source)
	}

	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()

	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	output, err := os.Create(destination)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(output, input)
	closeErr := output.Close()
	if copyErr != nil {
		return copyErr
	}
	if closeErr != nil {
		return closeErr
	}
	return nil
}

func writeJSON(destination string, value any) error {
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encodeErr := encoder.Encode(value)
	closeErr := file.Close()
	if encodeErr != nil {
		return fmt.Errorf("write %s: %w", destination, encodeErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close %s: %w", destination, closeErr)
	}
	return nil
}

func writeThirdPartyLicenses(destination string, manifest releaseManifest, analyzerFiles []string) error {
	var builder strings.Builder
	builder.WriteString("# Third-party licenses\n\n")
	fmt.Fprintf(&builder, "Release tier: %s\n\n", manifest.ReleaseTier)

	builder.WriteString("## Dictionaries\n\n")
	if len(manifest.Dictionaries) == 0 {
		builder.WriteString("No dictionaries are included.\n")
	}
	for _, dictionary := range manifest.Dictionaries {
		fmt.Fprintf(&builder, "### %s\n\n", dictionary.Name)
		fmt.Fprintf(&builder, "- ID: %s\n", dictionary.ID)
		fmt.Fprintf(&builder, "- Source: %s\n", dictionary.Source)
		fmt.Fprintf(&builder, "- License: %s\n", dictionary.License)
		fmt.Fprintf(&builder, "- Project release tier: %s\n", dictionary.ReleaseTier)
		if len(dictionary.NoticeFiles) > 0 {
			builder.WriteString("- Included notice files:\n")
			for _, name := range dictionary.NoticeFiles {
				fmt.Fprintf(&builder, "  - dictionaries/%s/%s\n", dictionary.ID, name)
			}
		}
		builder.WriteString("\n")
	}

	if len(manifest.Analyzers) > 0 {
		builder.WriteString("## Analyzer dependencies\n\n")
		for _, analyzer := range manifest.Analyzers {
			fmt.Fprintf(&builder, "### %s %s\n\n", analyzer.ID, analyzer.Version)
			builder.WriteString("Included notice files:\n")
			for _, name := range analyzer.NoticeFiles {
				fmt.Fprintf(&builder, "- %s\n", name)
			}
			builder.WriteString("\n")
		}
		if len(analyzerFiles) > 0 {
			builder.WriteString("All copied analyzer notice files are under licenses/analyzers/.\n")
		}
	}

	if err := os.WriteFile(destination, []byte(builder.String()), 0o644); err != nil {
		return fmt.Errorf("write THIRD_PARTY_LICENSES.md: %w", err)
	}
	return nil
}

func validateReleaseComponentID(id string) error {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(id) != id {
		return fmt.Errorf("dictionary id is not safe for release packaging: %q", id)
	}
	if id == "." || id == ".." || strings.ContainsAny(id, "/\\:") {
		return fmt.Errorf("dictionary id is not safe for release packaging: %q", id)
	}
	return nil
}
