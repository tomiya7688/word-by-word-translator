# Product specification

この文書を、単語単位翻訳アプリのMVPにおける**プロダクト挙動のSource of Truth**とする。

実装構造は `docs/architecture.md`、言語・多辞書拡張は `docs/language-dictionary-extension.md`、辞書形式は `docs/dictionary-package.md`、辞書ライセンス受入基準は `docs/license-policy.md`、配布区分は `docs/releases.md` を参照する。

## 1. Purpose

本アプリは、文章全体を自然な訳文へ変換する全文翻訳器ではない。

入力文を辞書引き可能な語単位へ分解し、各語を独立に辞書検索し、**原文の順序を保ったまま訳語候補を提示する**。

主な価値は、原文構造を残しながら、利用者の辞書引き作業を減らすことにある。

## 2. Non-goals

MVPでは次を行わない。

- 文脈に合わせた自然な全文翻訳
- 原文語順の積極的な並べ替え
- 文脈からの文章生成
- 文全体を解釈して訳語を1つへ強制的に決定すること
- Analyzer内部のsubword / model tokenをそのまま利用者へ表示すること

複数辞書から異なる訳語が得られる場合は、候補と出典を保持する。

## 3. MVP languages

初期対応言語は次の2言語。

- English (`en`)
- Japanese (`ja`)

対応方向:

- English -> Japanese
- Japanese -> English

言語追加はAnalyzerとDictionary Providerの追加を基本とし、コア処理を言語ごとに作り直さない。

## 4. Translation pipeline

概念上の処理順:

```text
input text
  -> Analyzer / Tokenizer
  -> lookup units
  -> selected Dictionary Providers
  -> dictionary result merge
  -> UI presentation
```

Analyzerは、辞書検索に適した単位への分割と、必要に応じてlemma / lookup unitの取得を行う。

Dictionary Providerは、指定されたlookup unitに対する辞書結果を返す。

Result Mergerは、複数辞書の候補を統合しつつ、出典辞書と個別結果を失わない。

UIはProcessから返された結果を表示し、辞書検索や統合判断を行わない。

## 5. Token model

各tokenは少なくとも次の情報を持てる。

| Field | Meaning |
| --- | --- |
| `surface` | 原文上の表記 |
| `lemma` | 原形。Analyzerが取得できる場合に使用 |
| `lookup_unit` | 辞書検索に使用する単位 |
| `subtokens` | Analyzer内部の分割情報。UI表示単位とはしない |
| `kind` | word / symbol |
| `status` | 解析・検索状態 |

翻訳後は、tokenに対して次も保持できる。

- merged translation candidates
- candidateごとのsource dictionary IDs
- raw per-dictionary lookup results

内部解析単位と、利用者へ見せる「語」単位は同一である必要はない。

## 6. Status model

MVPでは少なくとも次を区別する。

### `success`

解析と辞書検索が成功し、翻訳結果または記号表示を返せる状態。

### `segmentation_failed`

対象部分を辞書検索可能な語へ正常に解析できなかった状態。

表示時は原文surfaceを保持する。

### `dictionary_not_found`

解析は成功したが、選択された辞書群から訳語が見つからなかった状態。

表示時は原文surfaceを保持する。

### `lookup_error`

辞書アクセス中にエラーが発生し、その語について利用可能な訳語を得られなかった状態。

表示時は原文surfaceを保持する。

複数辞書利用時は、1辞書が失敗しても他辞書から有効な候補が得られれば、その候補を利用して処理を継続する。

## 7. Partial failure

部分失敗を許容し、原則として1語の失敗を文全体のエラーへ昇格させない。

例:

```text
Input:
I saw qwertymonster yesterday.

Display:
私 見た qwertymonster 昨日 .
```

この場合:

- `I`: success
- `saw`: success
- `qwertymonster`: dictionary_not_found
- `yesterday`: success
- `.`: success / symbol

`qwertymonster` は原文のまま残し、UIで未解決表示にする。

terminal rendererでは未解決tokenを赤色表示する。将来GUIへ移行しても、「未解決」という意味的styleをUI層で色やwidgetへ変換する。

## 8. Ordering and presentation

表示順は入力token順を保持する。

Processは自然文になるように語順を組み替えない。

success tokenでは辞書優先順位に従う候補を保持する。短縮表示では最優先候補を使えるが、1画面UIでは複数候補と出典辞書を確認できる。

例:

```text
{りんご [dict-a] / 林檎 [dict-b]}
```

同一訳語が複数辞書で一致した場合は重複表示せず、出典辞書IDをまとめて保持する。

## 9. Analyzer behavior

Analyzerの目的は全文翻訳ではなく、**辞書検索可能な単位を安定して作ること**。

現在のMVP:

- English: project tokenizer
- Japanese: Kagome adapter + IPADIC for segmentation and base-form extraction

日本語では活用形からlemmaを取得し、辞書形をlookup unitとして使用できる。

例:

```text
食べました
  surface: 食べ
  lemma / lookup_unit: 食べる
```

Analyzer固有実装はProcessのTokenizer interfaceの内側へ閉じ込める。

## 10. Dictionary behavior

辞書は交換可能なProviderとして扱う。

MVPでは:

- 複数辞書を登録可能
- UIから1つまたは複数を選択可能
- 選択された辞書だけを検索対象にできる
- 登録順を優先順位として扱う
- 同一訳語を重複除去できる
- 訳語ごとの出典辞書を保持する
- raw per-dictionary resultを保持する
- 1辞書のlookup errorで他辞書検索を停止しない

現在の実辞書:

- EJDict-hand: English -> Japanese
- TKG Japanese-English Learner's Dictionary: Japanese -> English

ライセンスとrelease tierは `docs/dictionaries.md` と `docs/releases.md` を参照する。

## 11. UI scope

MVP UIは1画面で次を扱える。

- input text
- translation direction
- one or multiple dictionary selections
- translate action
- word-by-word result
- merged candidates and dictionary provenance
- unresolved token indication
- validation / execution error

GUI frameworkはプロダクト仕様に含めない。UI状態と表示意味をframework固有実装から分離する。

## 12. MVP acceptance baseline

MVPの主要回帰条件:

1. English -> Japaneseを実行できる。
2. Japanese -> Englishを実行できる。
3. 日本語活用形から辞書形へ接続できる。
4. 複数辞書結果を統合できる。
5. 選択されていない辞書は検索しない。
6. 1辞書の失敗で、他辞書の有効結果を失わない。
7. 未知語は原文surfaceを保持する。
8. 未解決語だけをUIで識別できる。
9. 入力token順を保持する。
10. MIT / Full Releaseの辞書・asset条件を分離できる。

自動テストの対応表は `docs/test-coverage.md` を参照する。

## 13. Related specifications

- `docs/language-dictionary-extension.md` — language / analyzer / multi-dictionary extension policy
- `docs/architecture.md` — UI / Process / Data boundary
- `docs/analyzers.md` — Analyzer implementation and distribution notes
- `docs/dictionary-package.md` — Dictionary Package v1
- `docs/license-policy.md` — dictionary license acceptance and MIT / Full / Unsupported classification
- `docs/dictionaries.md` — adopted dictionaries
- `docs/ui.md` — screen state and presentation behavior
- `docs/releases.md` — MIT / Full Release composition
- `docs/test-coverage.md` — MVP regression coverage
