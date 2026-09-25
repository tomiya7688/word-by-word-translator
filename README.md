# word-by-word-translator

単語単位で原文と訳語を対応させる翻訳アプリです。

## Product specification

MVPのプロダクト挙動は `docs/specification.md` をSource of Truthとします。

## Implementation

- Main language: **Go**
- Architecture: **UPD Commander Base Design** をプロジェクト向けに適用
- AI-assisted development: **ai-context-reducer** の最小コアを標準採用
- Dictionary design: 複数辞書を独立した adapter として扱い、検索結果を統合可能にする
- Release policy: MIT release と Full release の辞書構成を分離できる設計にする

## Run the CLI

The current runnable MVP is a Full-runtime CLI because Japanese source analysis uses Kagome + IPADIC.

```bash
go run ./cmd/word-by-word-translator \
  -from en \
  -to ja \
  -dictionary /path/to/dictionary.json \
  -- "I saw qwertymonster yesterday."
```

See `docs/cli.md` for dictionary selection, stdin, compact output, and Japanese -> English usage.

## Dictionaries

MIT Release向けの最初の実辞書として、両方向を用意しています。

- English -> Japanese: **EJDict-hand** — Public Domain / CC0-1.0
- Japanese -> English: **TKG Japanese-English Learner's Dictionary** — CC0-1.0

完全な辞書packageは上流の固定revisionから再現生成します。採用辞書は `docs/dictionaries.md`、ライセンス受入基準は `docs/license-policy.md` を参照してください。

## Japanese analysis

日本語の形態素解析adapterとして Kagome v2.9.9 + IPADIC を追加しています。活用形から基本形を取得して翻訳辞書の見出し語へ接続します。

IPADICには再配布時のnotice条件があるため、このadapterは現時点ではFull Release向けとして独立packageに隔離しています。詳細は `docs/analyzers.md` を参照してください。

## UI

UI層には、文章入力・日本語↔英語の方向切替・複数辞書選択・翻訳実行・単語単位結果を1つの状態として扱うMVP screen modelがあります。

未解決語は原文を保持して赤色表示し、複数辞書の統合候補は出典辞書ID付きで表示できます。詳細は `docs/ui.md` を参照してください。

## Releases

配布物は `MIT Release` / `Full Release` を自動生成できます。MIT版は厳格な辞書条件を満たすものだけ、Full版はnotice付きassetも同梱します。

```bash
bash tools/release-builder/script/run.sh
```

アプリ本体はルートの `LICENSE` にあるMIT Licenseです。辞書・解析assetのライセンスは分離して `release.json` / `THIRD_PARTY_LICENSES.md` / 個別noticeへ記録します。詳細は `docs/releases.md`。

## Development standards

このリポジトリでは次を開発標準とします。

### ai-context-reducer

- `AI_CONTEXT.md` をAI向けの小さな入口にする
- Search first, read second
- Goal / Required / Acceptance が揃ったら探索を止める
- Source of Truth を明示する
- unrelated refactor を現在タスクへ混ぜない
- targeted validation を優先する
- 未確認領域は `Unverified` として明示する
- remote更新があり得るため、実装前は compact remote delta を優先する

Reference: https://github.com/tomiya7688/ai-context-reducer

### UPD Commander Base Design

- UI / Process / Data の責務を分離する
- Commander は呼び出しの指揮だけを担当し、実処理を持たない
- Messenger は層間通信だけを担当する
- 実処理は各層の Processing に置く
- UI から Data への直接依存を作らない
- Process はUI表示方式を知らない
- Data はUIや翻訳処理の判断を知らない
- 開発チェックに Go UPD Commander Checker を組み込む

Reference: https://github.com/tomiya7688/upd-commander-base-design

## Planned layer mapping

```text
UI
  CLI / GUI input and presentation
        ↕ Messenger
Process
  tokenization / lookup orchestration / result merge / language pipeline
        ↕ Messenger
Data
  dictionary loading / adapters / indexing / cache
```

## Documents

- `docs/specification.md` — MVPの目的・非目標・token/status・部分失敗・受入基準
- `docs/cli.md` — 起動可能なCLI、辞書選択、stdin、Full-runtime境界
- `docs/language-dictionary-extension.md` — 対応言語、Analyzer登録、多辞書選択・優先順位・統合・言語追加手順
- `docs/architecture.md` — UI / Process / Data と辞書adapter境界
- `docs/dictionary-package.md` — 共通辞書package、manifest、TSV、release tier
- `docs/license-policy.md` — 辞書ライセンス受入基準、MIT / Full / Unsupported判定、追加時レビュー
- `docs/dictionaries.md` — 採用辞書、ライセンス確認、固定revision
- `docs/analyzers.md` — 日本語形態素解析adapter、基本形、配布区分
- `docs/ui.md` — 翻訳画面状態、辞書選択、結果表示、未解決語の視覚表現
- `docs/releases.md` — MIT / Full Release の自動選別、notice、配布構成
- `AI_CONTEXT.md` — AI開発時の最小入口
