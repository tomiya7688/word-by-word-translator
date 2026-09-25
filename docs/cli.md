# Command-line application

`cmd/word-by-word-translator` is the runnable composition root for the current MVP pipeline.

It wires the existing UI, Process, and Data layers without reimplementing translation behavior in the command package.

## Runtime classification

The command supports both:

- English -> Japanese
- Japanese -> English

Japanese source analysis uses Kagome + IPADIC. Because the analyzer dictionary carries redistribution notice requirements, this executable is currently treated as a **Full runtime**.

This does not change the project-defined MIT Release dictionary policy. A notice-free Japanese Analyzer is still required before a bidirectional executable can be classified as MIT Release-compatible.

## Run

Load one or more Dictionary Package v1 manifests with repeated `-dictionary` flags.

Example: English -> Japanese

```bash
go run ./cmd/word-by-word-translator \
  -from en \
  -to ja \
  -dictionary dist/dictionaries/mit/ejdict-hand/dictionary.json \
  -- "I saw qwertymonster yesterday."
```

Example: Japanese -> English

```bash
go run ./cmd/word-by-word-translator \
  -from ja \
  -to en \
  -dictionary dist/dictionaries/mit/tkg-ja-en/dictionary.json \
  -- "食べました"
```

Text can also be read from stdin:

```bash
printf '食べました\n' | go run ./cmd/word-by-word-translator \
  -from ja \
  -to en \
  -dictionary dist/dictionaries/mit/tkg-ja-en/dictionary.json
```

## Dictionary selection

Every `-dictionary` flag loads one Dictionary Package v1 package.

By default, all loaded dictionaries matching the requested language direction are selected.

Use repeated `-use` flags to select a subset by dictionary ID:

```bash
go run ./cmd/word-by-word-translator \
  -from en \
  -to ja \
  -dictionary /path/to/dict-a/dictionary.json \
  -dictionary /path/to/dict-b/dictionary.json \
  -use dict-b \
  -- "apple"
```

Multiple dictionaries can be selected:

```text
-use dict-a -use dict-b
```

The selected IDs are passed through the existing UI screen state into `TranslationRequest.DictionaryIDs`. The CLI does not bypass the normal dictionary selection path.

Duplicate dictionary IDs across loaded manifests are rejected because `-use <id>` would otherwise be ambiguous.

## Output

Default output uses the existing one-screen terminal renderer.

Example:

```text
Text: apple
Direction: en -> ja
Dictionaries:
[x] Dictionary A (dict-a)
[x] Dictionary B (dict-b)
Result: {りんご [dict-a,dict-b] / 林檎 [dict-b]}
```

Unresolved words retain their original surface and are rendered in ANSI red.

For a compact word-by-word sequence, use `-compact`:

```bash
go run ./cmd/word-by-word-translator \
  -compact \
  -dictionary /path/to/dictionary.json \
  -- "I saw qwertymonster yesterday."
```

The compact output is the same `RenderTerminal` representation already used by the UI layer.

## Options

```text
-from <language>        source language: en or ja
-to <language>          target language: ja or en
-dictionary <manifest> Dictionary Package v1 manifest; repeatable
-use <dictionary-id>   selected dictionary ID; repeatable
-compact               print only the translated token sequence
```

Defaults:

```text
-from en
-to ja
```

At least one `-dictionary` is required.

## Architecture

The command directory is a composition root.

```text
cmd/word-by-word-translator
        |
        +-> UI ScreenCommander
              |
              v
          UI Messenger
              |
              v
        Process TranslateCommander
              |
              v
        Process DictionaryMessenger
              |
              v
        Data DictionaryMessenger
              |
              v
        Data DictionaryStore
```

The command is allowed to instantiate objects from all layers for dependency wiring, but it must not contain:

- tokenization algorithms
- dictionary lookup algorithms
- result merging
- unresolved-token presentation rules

Those remain in their existing layer-specific Processing packages.

## Build

A local binary can be built with:

```bash
go build -o build/word-by-word-translator ./cmd/word-by-word-translator
```

The current release builder creates dictionary/license asset bundles and does not yet produce cross-platform application binaries. Binary packaging can be added as a separate release task.
