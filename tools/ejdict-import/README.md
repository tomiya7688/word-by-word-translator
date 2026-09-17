# EJDict Importer

EJDict-hand の上流データを、本プロジェクトの `Dictionary Package v1` へ変換する開発・release用ツールです。

## Upstream

- Repository: https://github.com/kujirahand/EJDict
- Pinned revision: `source.lock` に固定
- Direction: English -> Japanese
- License: Public Domain / CC0-1.0

上流READMEは、データを Public Domain / CC0 とし、個人・商用を問わず利用・再配布可能と説明しています。本プロジェクトではこれを MIT Release の辞書受入条件を満たすものとして扱います。

## Why import instead of loading the source format directly?

EJDictの元データはTSVに近い形式ですが、次のように複数の綴りを1つの見出し語欄へ持てます。

```text
center, centre<TAB>中央、中心
```

importerは各綴りを独立した標準entryへ展開します。

```text
center<TAB>中央、中心
centre<TAB>中央、中心
```

翻訳coreや汎用 `FileDictionary` にEJDict固有規則を持ち込まないため、この変換は `tools/` に閉じ込めます。

## Run

Linux / macOS / Git Bash:

```bash
bash tools/ejdict-import/script/run.sh
```

Windows PowerShell:

```powershell
powershell -ExecutionPolicy Bypass -File tools/ejdict-import/script/run.ps1
```

既定出力:

```text
dist/dictionaries/mit/ejdict-hand/
├─ dictionary.json
└─ entries.tsv
```

別の出力先を指定する場合:

```bash
bash tools/ejdict-import/script/run.sh /path/to/output
```

```powershell
powershell -ExecutionPolicy Bypass -File tools/ejdict-import/script/run.ps1 -Output C:\path\to\output
```

runnerは `.tools-cache/ejdict-hand/` に上流repoを取得し、`source.lock` のcommitへdetachしてから変換します。

## Update upstream revision

上流更新を取り込む場合は、ライセンス条件とデータ形式が変わっていないことを確認したうえで `source.lock` を更新し、以下を実行します。

```bash
go test ./...
bash tools/ejdict-import/script/run.sh
```

生成manifestの `source` には固定commit URLが記録されます。
