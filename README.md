# word-by-word-translator

単語単位で原文と訳語を対応させる翻訳アプリです。

## Implementation

- Main language: **Go**
- Architecture: **UPD Commander Base Design** をプロジェクト向けに適用
- AI-assisted development: **ai-context-reducer** の最小コアを標準採用
- Dictionary design: 複数辞書を独立した adapter として扱い、検索結果を統合可能にする
- Release policy: MIT release と Full release の辞書構成を分離できる設計にする

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

- `docs/architecture.md` — UI / Process / Data と辞書adapter境界
- `docs/dictionary-package.md` — 共通辞書package、manifest、TSV、release tier
- `AI_CONTEXT.md` — AI開発時の最小入口
