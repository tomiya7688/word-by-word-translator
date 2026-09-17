# Dictionary Package v1

辞書固有形式をアプリ本体へ直接持ち込まず、Data層のadapterが共通形式へ変換して扱うための初期パッケージ仕様です。

## Layout

```text
dictionary-name/
├─ dictionary.json
└─ entries.tsv
```

`dictionary.json` の `data_file` はマニフェストと同じディレクトリ以下の相対パスだけを許可します。

## Manifest

```json
{
  "schema_version": 1,
  "id": "example-en-ja",
  "name": "Example English-Japanese",
  "source": "https://example.invalid/dictionary",
  "license": "CC0-1.0",
  "source_language": "en",
  "target_language": "ja",
  "commercial_use": true,
  "noncommercial_use": true,
  "modification": true,
  "redistribution": true,
  "attribution_required": false,
  "license_notice_required": false,
  "share_alike": false,
  "data_file": "entries.tsv",
  "format": "tsv-v1",
  "case_sensitive": false
}
```

`release_tier` はマニフェストへ手入力しません。ライセンス条件からアプリ側で計算します。

- `mit`: commercial / noncommercial / modification / redistribution がすべて可能で、attribution / license notice / share-alike が不要
- `full`: 上記4利用条件を満たすが、attribution / license notice / share-alike のいずれかが必要
- `unsupported`: 上記4利用条件のいずれかを満たさない

ここでの `mit` はプロジェクトの MIT Release 受入区分であり、辞書自身のライセンス名が MIT License であることを意味しません。

## TSV v1

1行につき1訳語です。

```text
hello\tこんにちは
hello\tやあ
world\t世界
```

- 区切りは最初のTAB
- 空行は無視
- `#` で始まる行はコメント
- 同じheadwordを複数行記述可能
- 1行あたり最大1 MiB
- `case_sensitive: false` の場合、検索キーは小文字化して比較する

TSV v1 は外部辞書の原形式そのものを標準化するものではありません。JMdict等の外部形式は辞書固有adapterで読み込み、この共通Dictionary interfaceへ接続します。
