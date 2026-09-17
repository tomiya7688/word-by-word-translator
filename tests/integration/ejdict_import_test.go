package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/tomiya7688/word-by-word-translator/applications/main/contracts"
	dataprocessing "github.com/tomiya7688/word-by-word-translator/applications/main/data/processing"
)

func TestEJDictImportOutputLoadsAsMITDictionary(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(filename), "../.."))
	sourceDir := t.TempDir()
	outputDir := t.TempDir()

	input := "apple\tリンゴ、林檎\ncenter, centre\t中央、中心\n"
	if err := os.WriteFile(filepath.Join(sourceDir, "a.txt"), []byte(input), 0o600); err != nil {
		t.Fatal(err)
	}

	command := exec.Command(
		"go", "run", "./tools/ejdict-import/script",
		"-source", sourceDir,
		"-output", outputDir,
		"-revision", "test-revision",
	)
	command.Dir = root
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("ejdict importer failed: %v\n%s", err, output)
	}

	dictionary, err := dataprocessing.LoadFileDictionary(filepath.Join(outputDir, "dictionary.json"))
	if err != nil {
		t.Fatal(err)
	}
	metadata := dictionary.Metadata()
	if metadata.ID != "ejdict-hand" || metadata.ReleaseTier != contracts.ReleaseTierMIT {
		t.Fatalf("metadata = %#v", metadata)
	}

	entries, err := dictionary.Lookup("CENTRE")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Translation != "中央、中心" {
		t.Fatalf("entries = %#v", entries)
	}
}
