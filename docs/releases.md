# Release composition

本プロジェクトはアプリ本体のライセンスと、辞書・解析assetの配布条件を分離する。

翻訳辞書の受入基準とtier判定は `docs/license-policy.md` をSource of Truthとする。

- アプリ本体: MIT License
- MIT Release: プロジェクト定義の厳格な配布tier
- Full Release: notice / attribution等が必要なassetも同梱可能な配布tier

「MIT Release」は辞書データのライセンス名ではない。辞書側は `DictionaryManifest.ReleaseTier()` の判定条件に従う。

## Build

完全生成:

```bash
bash tools/release-builder/script/run.sh
```

PowerShell:

```powershell
powershell -ExecutionPolicy Bypass -File tools/release-builder/script/run.ps1
```

このrunnerは固定revisionからEJDict-hand / TKG JA-ENを生成した後、release builderを実行する。

既存dictionary packageだけから生成する場合:

```bash
go run ./tools/release-builder/script \
  -dictionaries dist/dictionaries \
  -output dist/releases \
  -app-license LICENSE \
  -analyzer-notices licenses/analyzers \
  -tier all
```

## Output

```text
dist/releases/
├─ mit/
│  ├─ LICENSE
│  ├─ release.json
│  ├─ THIRD_PARTY_LICENSES.md
│  └─ dictionaries/
└─ full/
   ├─ LICENSE
   ├─ release.json
   ├─ THIRD_PARTY_LICENSES.md
   ├─ dictionaries/
   └─ licenses/
      └─ analyzers/
```

## MIT Release

同梱条件:

- `ReleaseTierMIT` の辞書だけ
- `ReleaseTierFull` は除外
- `ReleaseTierUnsupported` は除外
- Kagome / IPADIC analyzer notice assetは含めない

このtierの辞書は、商用・非商用・改変・再配布が可能で、通常利用時のattribution / license notice / share-alikeを要求しないものに限定する。

## Full Release

同梱条件:

- `ReleaseTierMIT`
- `ReleaseTierFull`
- `ReleaseTierUnsupported` は除外
- Kagome v2.9.9 / kagome-dict v1.1.0 / kagome-dict-ipa v1.2.0 のnotice群を同梱

Full辞書packageの追加義務ファイル:

- `attribution_required=true`: package直下に `ATTRIBUTION.txt`
- `license_notice_required=true`: package直下に `LICENSE.txt`
- `share_alike=true`: package直下に `LICENSE.txt`

必要なファイルが無い場合、Full Release生成は失敗する。release builderはライセンス内容の法的妥当性そのものを判定せず、「manifestで宣言された義務に対応する原文ファイルが存在すること」を機械的に保証する。

## Analyzer notices

Full Release用の固定notice原文は次に保持する。

```text
licenses/analyzers/
├─ kagome/LICENSE.txt
├─ kagome-dict/LICENSE.txt
└─ kagome-dict-ipa/
   ├─ LICENSE.txt
   └─ NOTICE.txt
```

これらは現在使用している固定versionのupstream原文から保持する。

MIT ReleaseではIPADIC由来assetを同梱しない。Full Releaseではnoticeを欠落させずに同梱する。

## release.json

各releaseに、実際に含めた辞書とanalyzer componentを記録する。

辞書について最低限次を保持する。

- ID
- name
- source
- license
- project release tier
- attribution / license notice / share-alike条件
- 同梱noticeファイル

これにより、release directoryだけから配布構成を機械的に検査できる。

## THIRD_PARTY_LICENSES.md

人向けの一覧として自動生成する。

- 同梱辞書
- 元source
- license識別子
- project release tier
- 個別noticeファイル
- Full Releaseのanalyzer dependencyとnotice path

個別ライセンス原文が必要なcomponentでは、一覧だけで代替せず個別ファイルも保持する。

## Current dictionaries

現時点の実辞書2件は両方ともMIT tier。

- EJDict-hand: English -> Japanese
- TKG Japanese-English Learner's Dictionary: Japanese -> English

したがって現時点ではMIT / Fullの翻訳辞書集合は同じで、Full側には追加で日本語解析用third-party noticeが入る。将来Full tier辞書を追加するとrelease builderが自動的に差分を作る。


## Runnable CLI boundary

`cmd/word-by-word-translator` currently supports both MVP language directions and imports the Kagome/IPADIC Japanese Analyzer path. Therefore it is classified as a Full runtime.

The current release builder produces dictionary/license asset bundles; it does not yet compile or package this command as cross-platform binaries.

A bidirectional executable for the project-defined MIT Release requires a Japanese Analyzer whose distributed runtime assets satisfy the MIT Release asset policy.
