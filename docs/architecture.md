# Architecture

## Goal

Goで単語単位翻訳を実装しつつ、辞書・UI・翻訳処理を交換可能に保つ。

このプロジェクトでは UPD Commander Base Design を、翻訳アプリ向けに次のように割り当てる。

```text
UI Layer
  Commander
  Messenger
  Processing
      ↕
Process Layer
  Commander
  Messenger
  Processing
      ↕
Data Layer
  Commander
  Messenger
  Processing
```

## UI Layer

責務:
- CLI / GUI からの入力受付
- 入力イベントのUI上の解釈
- Processから返された結果の表示形式決定
- 未知語などの視覚表現

持たない責務:
- 辞書検索
- tokenizationの本処理
- 辞書優先順位や統合判断
- 辞書ファイルI/O

GUIライブラリはUI層だけに閉じ込める。

## Process Layer

責務:
- 翻訳要求のオーケストレーション
- 言語ごとのtokenization処理
- 複数辞書への検索要求
- 辞書結果の統合・優先順位付け
- 未知語判定
- 原文tokenと訳語結果の対応付け

Processは表示方法や辞書ファイル形式を知らない。

## Data Layer

責務:
- 辞書adapter
- 辞書ファイルのload / parse
- index構築
- lookup実行
- cache
- 辞書metadata / license metadataの供給

Dataは「どの辞書結果を採用するか」という翻訳判断をしない。

## Commander

Commanderは一連の処理を指揮する。

例:

```text
Translate Commander
  -> tokenize processing
  -> dictionary messenger
  -> merge processing
  -> result messenger
```

Commander自身にtokenization、検索、mergeアルゴリズムを書かない。

## Messenger

Messengerは層を越える通信だけを担当する。

```text
UI <-> Process <-> Data
```

禁止:

```text
UI -> Data
UI Processing -> Process Processing
Process Processing -> Data Processing
```

層間通信にはcontract / DTOを使用し、相手層の内部型へ依存しない。

## Dictionary boundary

辞書は本体へ直接埋め込まず、adapterを通して共通contractへ変換する。

概念上の境界:

```go
type Dictionary interface {
    Lookup(token Token) ([]Entry, error)
    Metadata() DictionaryMetadata
}
```

実際のinterfaceは最初の辞書実装時に必要最小限で確定する。将来要件を推測して巨大なinterfaceを先に作らない。

## Tokenizer boundary

言語ごとの分割規則は交換可能にする。

```go
type Tokenizer interface {
    Tokenize(text string) ([]Token, error)
}
```

英語、日本語、中国語、韓国語などの実装差をProcess全体へ漏らさない。

## Multiple dictionaries

複数辞書を検索できるようにし、結果統合はProcess側で行う。

```text
Dictionary A ─┐
Dictionary B ─┼─> normalized entries -> merge/rank -> token result
Dictionary C ─┘
```

辞書adapterは原典情報を失わない。統合後も辞書名・由来・license等へ辿れる構造を維持する。

## Release boundary

MIT release と Full release は、可能な限り同一coreを使用し、搭載辞書セットとmetadataで差を表現する。

本体コードをreleaseごとにforkしない。

## Tools boundary

開発支援ツールは `tools/` 以下へ分離し、アプリ本体のpackageからimportしない。

- UPD checker
- context reduction support
- build / validation helpers

ツールのcacheやcloneは `.tools-cache/` 等の生成領域に置き、通常コンテキストと配布物から除外する。
