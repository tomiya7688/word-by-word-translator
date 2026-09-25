# Dictionary Package v1

ライセンス受入基準とtier判定のSource of Truthは `docs/license-policy.md`。この文書ではpackage schema上の要約だけを記載する。

辞書固有形式を翻訳コアへ直接持ち込まず、共通の辞書packageとして読み込むための初期仕様です。

## 構成

推奨構成:

```text
dictionary/
├─ dictionary.json
└─ entries.tsv
```

manifestのファイル名自体は固定しません。`LoadFileDictionary()` へmanifestのパスを渡します。

## Manifest

例:

```json
{
  "schema_version": 1,
  "id": "example-en-ja",
  "name": "Example EN-JA",
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

未知のmanifest項目は受理しません。schemaを変更するときは `schema_version` を更新します。

`data_file` はmanifestと同じpackage内の相対パスだけを許可し、`../`、絶対パス、Windows drive pathなどpackage外へ出られる指定は拒否します。

## Release tier

`release_tier` はmanifestから自己申告させず、ライセンス条件から自動算出します。

### MIT

以下をすべて満たす辞書です。

- 商用利用可能
- 非商用利用可能
- 改変可能
- 再配布可能
- attribution不要
- license notice不要
- share-alike不要

### Full

商用利用・非商用利用・再配布が可能で、MIT tierの条件を満たさない辞書です。

例:

- attributionが必要
- license noticeが必要
- share-alikeが必要
- 改変不可だが、辞書そのものの再配布と商用利用は可能

個別ライセンスをFull releaseへ実際に同梱するかは、`docs/license-policy.md` の採用レビューに従います。manifestと自動tier判定は法的判断の代替ではありません。

### Unsupported

次のいずれかを満たす場合です。

- 商用利用不可
- 非商用利用不可
- 再配布不可

## TSV v1

UTF-8のテキストファイルです。1行を1訳語として扱います。

```text
headword<TAB>translation
```

同じ見出し語に複数訳語がある場合は複数行にします。

```text
hello	こんにちは
hello	やあ
world	世界
```

仕様:

- 空行は無視
- `#` から始まる行はコメント
- 先頭のTABより左を見出し語、右を訳語として扱う
- 見出し語・訳語の前後空白は除去
- 空の見出し語・訳語はエラー
- 不正UTF-8はエラー
- 同一の正規化キー + 同一訳語は読み込み時に重複除去
- 1行の上限は1 MiB
- `case_sensitive=false` の場合、検索キーを小文字化する

## Adapter boundary

`FileDictionary` は既存の `Dictionary` interfaceを実装します。

```text
manifest + TSV
      ↓
FileDictionary
      ↓
Dictionary interface
      ↓
DictionaryStore
      ↓
Process translation pipeline
```

JMdict等の固有形式は、将来この共通境界へ変換する専用adapterを追加します。Process層へXML/JSON/DB固有処理を漏らしません。
