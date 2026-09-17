# Dictionaries

この文書は、本プロジェクトで実際に採用・検証した辞書とrelease区分を記録します。

`MIT Release` / `Full Release` の判定ロジックと共通package形式は `docs/dictionary-package.md` を参照してください。

## Supported

### EJDict-hand

- Status: supported import source
- Direction: English -> Japanese
- Project: https://github.com/kujirahand/EJDict
- Upstream revision: `5e1a630bfabb2791a78d14e4e356d85bb6437e34`
- Upstream license statement: Public Domain / CC0-1.0
- Release tier: `mit`
- Importer: `tools/ejdict-import/`

MIT Release受入条件:

| Condition | Result |
| --- | --- |
| Commercial use | allowed |
| Non-commercial use | allowed |
| Modification | allowed |
| Redistribution | allowed |
| Attribution / notice required for ordinary use | no |

上流READMEには、辞書データがPublic Domain / CC0であり、個人・商用を問わず自由に利用・再配布できる旨が明記されています。

EJDict固有の `word1, word2<TAB>meaning` 形式はimport時に個別headwordへ展開し、アプリ実行時は通常の `Dictionary Package v1` として扱います。

### TKG Japanese-English Learner's Dictionary

- Status: supported import source
- Direction: Japanese -> English
- Project: https://github.com/tkgally/je-dict-1
- Upstream revision: `c954d75d3bd44cab1a8e4a6045fa16b88f539722`
- Upstream license: CC0-1.0
- Release tier: `mit`
- Importer: `tools/tkg-ja-en-import/`
- Imported source: `entries_index.json`
- Source entries at pinned revision: 30,743

MIT Release受入条件:

| Condition | Result |
| --- | --- |
| Commercial use | allowed |
| Non-commercial use | allowed |
| Modification | allowed |
| Redistribution | allowed |
| Attribution / notice required for ordinary use | no |

上流READMEはCC0-1.0であることと、商用を含む任意目的でデータとコードをコピーできることを明記しています。LICENSEもCC0-1.0 Universalの法文です。

全repositoryは大きいため、importerは固定revisionの `entries_index.json` だけを取得します。このindexに含まれる `headword`, `reading`, `gloss` を使い、`headword -> gloss` と `reading -> gloss` を標準TSVへ変換します。同じ読みを持つ別語は異なるglossとして保持し、複数候補を失いません。

## Full Release candidates

JMdict / Jitendex系は語彙量と品質の面で有力ですが、attribution / share-alike等の条件があるため、採用する場合はFull Release側で個別に扱います。

候補を追加するときは、ライセンス名だけでなく一次配布元の利用・改変・再配布・表示条件を確認してからこの一覧へ追加します。
