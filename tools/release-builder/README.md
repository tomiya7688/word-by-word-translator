# Release Builder

Dictionary Package v1 の manifest とthird-party noticeを基に、プロジェクト定義の MIT Release / Full Release 配布bundleを生成する。

## Complete build

固定revisionから現在のMIT辞書を生成してから両releaseを作る。

Linux / Ubuntu / Git Bash:

```bash
bash tools/release-builder/script/run.sh
```

Windows PowerShell:

```powershell
powershell -ExecutionPolicy Bypass -File tools/release-builder/script/run.ps1
```

出力:

```text
dist/releases/
├─ mit/
└─ full/
```

## Builder only

すでに `dist/dictionaries/` に辞書packageがある場合:

```bash
go run ./tools/release-builder/script \
  -dictionaries dist/dictionaries \
  -output dist/releases \
  -app-license LICENSE \
  -analyzer-notices licenses/analyzers \
  -tier all
```

`-tier` は `mit`, `full`, `all`。

## Rules

- MIT Release: manifest判定が `mit` の辞書のみ
- Full Release: `mit` + `full` の辞書
- `unsupported` はどちらにも含めない
- Full辞書で `attribution_required=true` の場合はpackage直下の `ATTRIBUTION.txt` が必須
- Full辞書で `license_notice_required=true` または `share_alike=true` の場合はpackage直下の `LICENSE.txt` が必須
- Full ReleaseはKagome/IPADIC notice一式を `licenses/analyzers/` から同梱
- 必須noticeが欠落している場合は生成を失敗させる
- 両releaseにアプリ本体の `LICENSE`, `release.json`, `THIRD_PARTY_LICENSES.md` を生成する

詳細は `docs/releases.md`。
