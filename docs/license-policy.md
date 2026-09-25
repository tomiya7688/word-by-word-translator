# Dictionary license acceptance policy

この文書を、**翻訳辞書のライセンス受入基準と MIT / Full Release 区分のSource of Truth**とする。

辞書package形式は `docs/dictionary-package.md`、採用済み辞書は `docs/dictionaries.md`、実際の配布構成は `docs/releases.md` を参照する。

## 1. Scope

この方針は、アプリへ同梱・再配布する**翻訳辞書データ**に適用する。

アプリ本体のコードライセンス、Analyzerライブラリ、Analyzer用辞書assetは別に扱う。

- application code: root `LICENSE`
- translation dictionaries: this policy
- analyzer dependencies / assets: `docs/analyzers.md` and `docs/releases.md`

コードと辞書データのライセンスを混同しない。

## 2. "MIT Release" is a project distribution class

本プロジェクトの **MIT Release** は、辞書のライセンス名ではない。

「導入者が通常利用する際に、辞書側の表示義務をほぼ意識せず、商用・非商用・改変・再配布できる辞書だけをまとめる」ための**プロジェクト独自の厳格な配布区分**である。

したがって、辞書ライセンスが文字列として `MIT` であっても、自動的にMIT Release対象にはならない。

### Standard MIT Licenseとの違い

一般的なMIT Licenseは、コピーまたは substantial portions にcopyright noticeとpermission noticeを含める条件を持つ。

そのため、標準MIT Licenseの辞書データを再配布する場合にそのnotice保持が必要なら、manifestは通常:

```json
{
  "modification": true,
  "redistribution": true,
  "attribution_required": false,
  "license_notice_required": true
}
```

となり、プロジェクトの `MIT Release` ではなく `Full Release` に分類される。

一方、アプリ本体は標準MIT Licenseであり、root `LICENSE` の条件に従う。

**アプリ本体がMIT Licenseであることと、辞書がMIT Release tierへ入ることは別問題。**

## 3. Manifest facts

各辞書は最低限、次のライセンス関連factをmanifestへ持つ。

| Field | Meaning |
| --- | --- |
| `source` | 条件確認に使った一次配布元 |
| `license` | upstreamが示すライセンス名・識別子 |
| `commercial_use` | 商用利用可能か |
| `noncommercial_use` | 非商用利用可能か |
| `modification` | 改変可能か |
| `redistribution` | 辞書データを再配布可能か |
| `attribution_required` | 利用・配布に帰属表示が必要か |
| `license_notice_required` | ライセンス文・copyright notice等の保持が必要か |
| `share_alike` | 派生物等に同一・互換ライセンス条件が課されるか |

`release_tier` はmanifestへ自己申告しない。

コード側の `DictionaryManifest.ReleaseTier()` が上記factから計算する。

## 4. Evidence rule

manifestのbooleanは推測で埋めない。

辞書追加時は、原則として**一次配布元のライセンス文・README・公式配布ページ**で条件を確認する。

最低限確認する事項:

1. commercial use
2. non-commercial use
3. modification
4. redistribution
5. attribution requirement
6. license/copyright notice requirement
7. share-alike requirement

可能なら固定revision / release / versionを使い、後から同じ条件を再確認できるようにする。

確認結果は:

- `dictionary.json`
- `docs/dictionaries.md`
- 必要なnotice原文

へ反映する。

## 5. Unknown or ambiguous terms

ライセンス条件が不明・曖昧な辞書を、推測でMIT / Fullへ入れない。

例:

- license名だけあり本文が確認できない
- redistribution可否が不明
- commercial useの範囲が不明
- upstream READMEとlicense本文が矛盾している
- 派生データへの条件が判断できない

この場合は**同梱候補のまま保留**する。

現在のmanifest schemaはunknownを三値で表現しないため、確認できていないassetを正式な配布manifestへ追加しない。

`ReleaseTier()` は法的審査器ではなく、**確認済みfactの機械分類器**である。

## 6. MIT Release acceptance

MIT Releaseへ入る辞書は、次をすべて満たす。

1. commercial use = allowed
2. non-commercial use = allowed
3. modification = allowed
4. redistribution = allowed
5. ordinary use / redistributionでattribution表示を要求しない
6. license / copyright noticeの同梱を要求しない
7. share-alikeを要求しない

manifest上は:

```text
commercial_use          = true
noncommercial_use       = true
modification            = true
redistribution          = true
attribution_required    = false
license_notice_required = false
share_alike             = false
```

この条件はIssue初期定義の5条件を、実際の配布判定で曖昧さが出ないよう:

- attribution
- license notice
- share-alike

へ分解したもの。

### Typical candidates

- CC0
- clearly dedicated Public Domain data
- equivalent licenses/terms with no attribution or notice burden for ordinary redistribution

ライセンス名だけでは判定しない。

## 7. Full Release acceptance

Full Releaseの**機械分類**に入る最低条件:

```text
commercial_use    = true
noncommercial_use = true
redistribution    = true
```

この3つを満たし、MIT Release条件を満たさない辞書は `ReleaseTierFull` になる。

例:

- attribution required
- license notice required
- share-alike required
- modification not allowed, but commercial/non-commercial use and redistribution are allowed

### Mechanical classification is not automatic adoption approval

Full tierに分類できることと、実際に採用することは別。

特に次は個別レビューする。

- `modification=false`
- share-alike
- 用途・表示方法に特殊条件がある
- データベース権等の追加条件がある
- アプリ配布形態との適合が不明

プロジェクト方針として、通常の改変・再配布・商用利用を大きく制約する辞書は可能な限り避ける。

ただし、現在の機械分類は「再配布可能な未改変asset」をFullで扱える余地を残すため、`modification=false` だけではUnsupportedへ落とさない。

## 8. Unsupported

次のいずれかがfalseなら、現在の機械判定では `ReleaseTierUnsupported`。

- commercial use
- non-commercial use
- redistribution

つまり:

```go
if !commercial || !noncommercial || !redistribution {
    return unsupported
}
```

代表例:

- Non-Commercial only
- redistribution prohibited
- commercial use prohibited
- personal-use-only data

UnsupportedはMIT / Fullのどちらにも同梱しない。

## 9. Exact release-tier algorithm

現在の実装と同じ判定:

```text
if commercial_use == false
   OR noncommercial_use == false
   OR redistribution == false
    => unsupported

else if modification == true
     AND attribution_required == false
     AND license_notice_required == false
     AND share_alike == false
    => mit

else
    => full
```

Source code:

`applications/main/data/processing/dictionary_manifest.go`

このアルゴリズムを変更する場合は、この文書・unit tests・release builderの期待値を同時に更新する。

## 10. Full Release notice files

Full Releaseで辞書manifestが追加義務を示す場合、release builderは次を要求する。

### Attribution

```text
attribution_required = true
```

dictionary package直下:

```text
ATTRIBUTION.txt
```

### License notice

```text
license_notice_required = true
```

dictionary package直下:

```text
LICENSE.txt
```

### Share-alike

```text
share_alike = true
```

少なくとも適用license原文を:

```text
LICENSE.txt
```

として要求する。

必要ファイルが欠ける場合、Full Release生成を失敗させる。

このチェックは「ファイルが存在すること」を保証するものであり、内容が法的に十分かを自動判定するものではない。

## 11. Application license separation

配布物ではアプリ本体と辞書ライセンスを分離する。

```text
release/
├─ LICENSE                     # application code license
├─ THIRD_PARTY_LICENSES.md     # included third-party summary
├─ release.json                # machine-readable composition
└─ dictionaries/
   └─ <dictionary-id>/
      ├─ dictionary.json
      ├─ entries.tsv
      ├─ ATTRIBUTION.txt       # when required
      └─ LICENSE.txt           # when required
```

root `LICENSE` はアプリ本体のMIT License。

辞書のlicense識別子や追加義務は、各dictionary packageとthird-party metadataで管理する。

## 12. Analyzer assets are independent

Analyzer用assetの配布条件は、翻訳辞書のrelease tierと独立して扱う。

現在:

- Kagome code: MIT License
- kagome-dict code: MIT License
- IPADIC-derived analyzer dictionary: redistribution notice required

そのため、IPADIC assetはプロジェクトのnotice-free MIT Release asset条件には入れず、Full Releaseで必要noticeを同梱する。

「翻訳辞書がMIT tierだから、そのrelease全体の全assetもnotice-free」という意味ではない。

release builderがcomponentごとに扱う。

## 13. Dictionary onboarding checklist

新辞書を追加する前に:

1. primary sourceを特定
2. version / revisionを固定
3. license原文を確認
4. 7つのmanifest factを判定
5. 不明条件があれば採用を保留
6. Dictionary Package v1またはadapterを用意
7. `ReleaseTier()` の結果を確認
8. Fullなら必要な `ATTRIBUTION.txt` / `LICENSE.txt` を同梱
9. `docs/dictionaries.md` に根拠とrevisionを記録
10. release builder testで意図したtierに入ることを確認

## 14. Current supported dictionaries

### EJDict-hand

- upstream statement: Public Domain / CC0
- project tier: `mit`
- ordinary attribution / license notice: not required

### TKG Japanese-English Learner's Dictionary

- upstream license: CC0-1.0
- project tier: `mit`
- ordinary attribution / license notice: not required

現時点の翻訳辞書は両方MIT tier。

Full Releaseとの差分は主にAnalyzer asset noticeだが、将来Full-tier翻訳辞書を追加できる構造になっている。

## 15. Policy invariants

- license名だけでrelease tierを決めない
- manifestは確認済みfactを記録する
- release tierはmanifest factから計算する
- ambiguous / unknownな条件を推測でtrueにしない
- MIT Releaseへnotice-required assetを混入させない
- Unsupported assetを配布物へ混入させない
- Full Releaseで必要noticeを欠落させない
- application code licenseとdictionary licenseを分離する
- automated classificationをlegal reviewの代替とみなさない
