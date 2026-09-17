# TKG Japanese-English Importer

`tkgally/je-dict-1` の `entries_index.json` を `Dictionary Package v1` へ変換する開発ツールです。

## Source

- Project: https://github.com/tkgally/je-dict-1
- Revision: `source.lock` で固定
- Input: `entries_index.json`
- License: CC0-1.0
- Direction: Japanese -> English
- Release tier: MIT Release

上流READMEとLICENSEはCC0-1.0を明示し、商用を含む任意目的でのコピー・改変・再利用を許可しています。固定revisionの `entries_index.json` は `headword`, `reading`, `gloss` を持つため、初期の単語単位検索には個別entry JSONを全件取得する必要がありません。

## Import rules

- `headword -> gloss` を1行の標準TSVへ変換する
- `reading` が存在する場合は `reading -> gloss` も検索aliasとして追加する
- 同じ `lookup unit + gloss` は重複除去する
- 同音異義語でglossが異なる場合は複数候補として保持する
- 詳細説明、例文、品詞等は現時点では取り込まない

出力例:

```text
愛\tlove
あい\tlove
藍\tindigo
あい\tindigo
```

## Run

Linux / macOS / Git Bash:

```bash
bash tools/tkg-ja-en-import/script/run.sh
```

Windows PowerShell:

```powershell
powershell -ExecutionPolicy Bypass -File tools/tkg-ja-en-import/script/run.ps1
```

既定出力:

```text
dist/dictionaries/mit/tkg-ja-en/
├─ dictionary.json
└─ entries.tsv
```

取得した固定revisionの `entries_index.json` は `.tools-cache/tkg-ja-en/<revision>/` にキャッシュします。生成物とcacheはアプリ本体のsource of truthではなく、`source.lock` とimporterが再現元です。
