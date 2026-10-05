# word-by-word-translator

文章を**単語単位**に分解し、それぞれを辞書で引いて原文順に表示する翻訳アプリです。

一般的な機械翻訳のように文全体を自然な訳文へ作り直すのではなく、原文の構造を残したまま辞書引きの手間を減らすことを目的としています。

例:

```text
I saw qwertymonster yesterday.
↓
私 見た qwertymonster 昨日 .
```

辞書で見つからない語は原文のまま残し、未解決語として表示します。

## 現在の状態

現在は次の機能が実装されています。

- 英語 → 日本語
- 日本語 → 英語
- 単語単位の翻訳
- 複数辞書の利用と結果統合
- 未知語・未解決語の原文保持
- CLIからの実行
- MIT Release / Full Release 用の辞書・ライセンス構成生成

デスクトップGUIは **Wails v2 + Svelte + TypeScript** で実装予定です。

## ソースコードはどこにあるか

アプリ本体のソースコードは主に `applications/main/` にあります。

```text
applications/main/
├─ contracts/     共通データ型・インターフェース
├─ process/       単語分割、翻訳処理、辞書結果の統合
├─ data/          辞書の読み込み・検索・管理
└─ ui/            入力状態・表示用処理

cmd/word-by-word-translator/
└─ CLIの起動コード

tools/
├─ ejdict-import/       EJDict-handの取り込み
├─ tkg-ja-en-import/    TKG日英辞書の取り込み
└─ release-builder/     配布用データの生成
```

翻訳処理を追いたい場合は、まず以下を見ると全体を把握しやすいです。

- `applications/main/process/processing/translation_processor.go`
- `applications/main/process/processing/english_tokenizer.go`
- `applications/main/process/processing/kagome/japanese_tokenizer.go`
- `applications/main/data/processing/dictionary_store.go`
- `cmd/word-by-word-translator/main.go`

## 使用している辞書

現在の初期辞書は次の2つです。

| 翻訳方向 | 辞書 |
| --- | --- |
| 英語 → 日本語 | EJDict-hand |
| 日本語 → 英語 | TKG Japanese-English Learner's Dictionary |

どちらもプロジェクト定義の `MIT Release` 条件を満たす辞書として採用しています。

辞書の詳細・固定revision・ライセンス確認は `docs/dictionaries.md` を参照してください。

## CLIで試す

### 必要なもの

- Go 1.22 以上

まず辞書データを生成します。

Linux / macOS:

```bash
bash tools/release-builder/script/run.sh
```

Windows PowerShell:

```powershell
powershell -ExecutionPolicy Bypass -File tools/release-builder/script/run.ps1
```

### 英語 → 日本語

```bash
go run ./cmd/word-by-word-translator \
  -from en \
  -to ja \
  -dictionary dist/dictionaries/mit/ejdict-hand/dictionary.json \
  -- "I saw qwertymonster yesterday."
```

### 日本語 → 英語

```bash
go run ./cmd/word-by-word-translator \
  -from ja \
  -to en \
  -dictionary dist/dictionaries/mit/tkg-ja-en/dictionary.json \
  -- "食べました"
```

標準入力にも対応しています。

```bash
printf '食べました\n' | go run ./cmd/word-by-word-translator \
  -from ja \
  -to en \
  -dictionary dist/dictionaries/mit/tkg-ja-en/dictionary.json
```

CLIの全オプションは `docs/cli.md` を参照してください。

## このアプリがしないこと

このアプリは、少なくともMVPでは次のことを目的としていません。

- 文脈を考慮した自然な全文翻訳
- 原文語順の並べ替え
- 文脈から訳語を1つに決めること
- 学習者の代わりに文全体を解釈すること

Analyzerも全文を理解するためではなく、**辞書を引ける単位へ分解し、必要なら辞書形へ戻すため**に使います。

## ライセンスと配布

アプリ本体は MIT License です。

辞書や解析用データには、それぞれ別のライセンスが適用されます。

本プロジェクトでは配布物を次の2系統に分けます。

- **MIT Release**: 商用利用・非商用利用・改変・再配布が可能で、通常利用時の明示義務がない辞書だけを含む
- **Full Release**: noticeや帰属表示などが必要な辞書・解析データも条件を守って含められる

詳しくは以下を参照してください。

- `docs/license-policy.md`
- `docs/releases.md`

## ドキュメント

利用・仕様を知りたい場合:

- `docs/specification.md` — アプリの動作仕様
- `docs/cli.md` — CLIの使い方
- `docs/dictionaries.md` — 採用辞書
- `docs/ui.md` — UIの表示仕様

実装を知りたい場合:

- `docs/architecture.md` — UI / Process / Data の構成
- `docs/language-dictionary-extension.md` — 言語・辞書の追加方法
- `docs/dictionary-package.md` — 辞書package形式
- `docs/analyzers.md` — 単語分割・日本語解析
- `docs/testing.md` — テスト方針

## 開発者向け

主な実装言語は **Go** です。

GUIは **Wails v2 + Svelte + TypeScript** を採用します。

内部アーキテクチャは UI / Process / Data の責務を分離し、辞書固有処理はadapterとしてData側に閉じ込めます。

開発時の詳細な規約やAI向け入口は以下を参照してください。

- `AI_CONTEXT.md`
- `docs/architecture.md`
- GitHub Issues

READMEでは利用者向け情報を優先し、詳細な内部設計や開発運用ルールは各ドキュメントへ分離します。
