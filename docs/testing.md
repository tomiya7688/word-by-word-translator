# MVP test coverage

Issue #20 のMVP回帰テスト対象を、現在の自動テストへ対応付ける。

| 対象 | 主な自動テスト |
| --- | --- |
| 英語Analyzer | `applications/main/process/processing/english_tokenizer_test.go` |
| 日本語Analyzer | `applications/main/process/processing/kagome/japanese_tokenizer_test.go`, `tests/integration/japanese_translation_test.go` |
| Dictionary Provider | `applications/main/data/processing/memory_dictionary_test.go`, `file_dictionary_test.go` |
| 辞書未登録 | `tests/integration/translation_test.go` の `qwertymonster` |
| 複数辞書統合 | `translation_processor_test.go`, `tests/integration/translation_screen_test.go` |
| 部分失敗 | `dictionary_store_test.go`, `tests/integration/partial_failure_test.go` |
| 日本語→英語 | `tests/integration/japanese_translation_test.go`, `bidirectional_translation_test.go` |
| 英語→日本語 | `tests/integration/translation_test.go`, `bidirectional_translation_test.go` |
| 未解決語の原文保持/UI | `result_presenter_test.go`, `terminal_renderer_test.go`, `partial_failure_test.go` |
| Release Tier判定 | `file_dictionary_test.go`, `tools/release-builder/script/main_test.go` |

## Required regression cases

### Unknown word

`I saw qwertymonster yesterday.` を英語Analyzerから辞書統合まで通し、`qwertymonster` が `dictionary_not_found` のまま原文保持されることを検証する。

### Partial dictionary failure

壊れた辞書と正常辞書を同時に選択し、一方のlookup失敗によって正常辞書の結果や画面全体が失われないことを検証する。

- 正常辞書の候補は `success`
- 訳語が得られず辞書errorだけがある語は `lookup_error`
- UIでは未解決語の原文を保持する
- UI全体の実行はerrorにしない

## Development check

通常のMVP回帰確認は次で行う。

```bash
bash tools/dev-check/script/check.sh
```

この経路で `gofmt`, `go vet ./...`, `go test ./...`, UPD Commander Checker をまとめて実行する。
