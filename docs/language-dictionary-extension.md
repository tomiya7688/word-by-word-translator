# Language and multi-dictionary extension policy

この文書を、**対応言語・言語ペア・辞書追加に関する拡張方針のSource of Truth**とする。

プロダクト挙動は `docs/specification.md`、層構造は `docs/architecture.md`、辞書packageとライセンスは `docs/dictionary-package.md` / `docs/releases.md` を参照する。

## 1. MVP language scope

MVPで有効化する言語は2言語。

| Language | Code | Analyzer |
| --- | --- | --- |
| English | `en` | project English tokenizer |
| Japanese | `ja` | Kagome adapter + IPADIC |

有効な翻訳方向:

- English -> Japanese
- Japanese -> English

次期候補として Korean / Chinese を優先候補に置くが、MVPの対応宣言には含めない。

## 2. Core principle

拡張の中心は「言語数」ではなく、**AnalyzerとDictionary Providerを独立して追加・選択・統合できること**に置く。

新しい言語や辞書を追加するために、次のコア責務を書き直さないことを設計目標とする。

- translation orchestration
- dictionary batch lookup contract
- dictionary result merge
- token status model
- UI / Process / Data boundary

新規対応は、既存interfaceへの実装追加と設定・登録の拡張で行う。

## 3. Language boundary

言語は `contracts.Language` で表す。

現在:

```go
const (
    LanguageEnglish  Language = "en"
    LanguageJapanese Language = "ja"
)
```

新言語を正式対応するときは、まず安定したlanguage codeをcontractへ追加する。

Language値だけを追加しても「対応言語」にはならない。次の条件が揃った時点で対応扱いとする。

1. source language用Analyzerがある
2. 対象方向のDictionary Providerが1つ以上ある
3. UIがその翻訳方向を選択可能
4. end-to-end integration testがある
5. 配布対象assetのlicense/release tierが確定している

## 4. Analyzer extension

AnalyzerはProcess層の `Tokenizer` interfaceへ適合させる。

```go
type Tokenizer interface {
    Tokenize(text string) ([]contracts.Token, error)
}
```

`TokenizerRegistry` はsource languageをkeyとしてAnalyzerを選択する。

```text
source_language
      |
      v
TokenizerRegistry
  en -> EnglishTokenizer
  ja -> KagomeTokenizer
  ko -> future KoreanTokenizer
  zh -> future ChineseTokenizer
```

したがって、言語固有の分割・lemma取得・lookup unit生成をTranslationProcessorへ直接追加しない。

### New analyzer checklist

新しいsource languageを追加するとき:

1. `contracts.Language` を追加
2. `Tokenizer` 実装を追加
3. application compositionで `TokenizerRegistry.Register()`
4. surface / lemma / lookup_unit / symbol判定のunit testを追加
5. 実辞書まで到達するintegration testを追加

Analyzerが使用する外部モデル・辞書assetは、翻訳辞書とは別にrelease条件を確認する。

## 5. Dictionary Provider boundary

辞書はData層の `Dictionary` interfaceへ適合させる。

```go
type Dictionary interface {
    Metadata() contracts.DictionaryMetadata
    Lookup(lookupUnit string) ([]contracts.DictionaryEntry, error)
}
```

言語ペアはadapter実装の型ではなく、`DictionaryMetadata.SourceLanguage` / `TargetLanguage` で表す。

`DictionaryStore` はrequestの言語ペアとmetadataを比較し、一致する辞書だけを検索する。

このため、同一のStoreへ異なる言語ペアの辞書を登録できる。

## 6. Adding another dictionary to an existing pair

同じ言語ペアへ辞書を増やす場合、AnalyzerやTranslationProcessorは変更しない。

必要な作業:

1. Dictionary Package v1または専用adapterを作る
2. source / target languageをmetadataへ設定
3. license条件をmanifestへ設定
4. storeへ登録
5. UIのdictionary optionへ公開
6. lookup / merge / provenanceのtestを追加

Dictionary Package v1へ変換できる辞書は `FileDictionary` を再利用する。

固有形式を直接Processへ持ち込まない。

## 7. Selection

UIは1つまたは複数のdictionary IDを `TranslationRequest.DictionaryIDs` へ載せる。

Data層は指定IDに一致する辞書だけを検索する。

```text
UI selected IDs
   ["a", "c"]
       |
       v
TranslationRequest
       |
       v
DictionaryStore
   A -> lookup
   B -> skip
   C -> lookup
```

低レベルAPIでdictionary ID指定が空の場合は、現在は後方互換のため対象言語ペアの全辞書を検索する。

UIのMVP操作では少なくとも1辞書の明示選択を要求する。

## 8. Priority

MVPの辞書優先順位は**DictionaryStoreへの登録順**である。

```go
NewDictionaryStore(highPriority, secondPriority, fallback)
```

選択された辞書の中でも、この登録順を保持してlookup resultをProcessへ返す。

TranslationProcessorはその順番でcandidateを構築するため、先に登録された辞書の候補が先に現れる。

現在はrequestごとの任意priority数値は持たない。

将来ユーザー設定可能なpriorityを追加する場合も、Data adapterへranking判断を入れず、Processへ渡す前のordered provider listまたはrequest/config側で順序を決める。

## 9. Merge semantics

複数辞書の統合はProcess層の責務。

MVPでは:

- candidate順は辞書優先順位を保持
- **同一translation文字列**を重複除去
- 同じtranslationを返したdictionary IDsを1candidateへ集約
- 異なるtranslationは個別candidateとして保持
- raw per-dictionary resultも保持

例:

```text
Dictionary A: 見た
Dictionary B: 見た, 見ました

Merged:
- 見た      [A, B]
- 見ました  [B]
```

Issue初期案にあった「近似訳語」の意味的deduplicationはMVPでは行わない。現在は完全一致のみ。

将来normalize / similarity処理を追加する場合もDictionary ProviderではなくResult Merger側の拡張とする。

## 10. Fallback and partial failure

特定辞書で見つからない語があっても、他の選択辞書を継続して検索する。

1辞書がlookup errorになっても、他辞書にcandidateがあればtoken全体は `success` とできる。

```text
Dictionary A -> error
Dictionary B -> translation found
Result       -> success with B candidate
```

全辞書でcandidateがなく、少なくとも1辞書にerrorがある場合は `lookup_error`。

全辞書でcandidateがなく、lookup errorもない場合は `dictionary_not_found`。

これによりprovider単位の故障を翻訳文全体の失敗へ拡大しない。

## 11. Provenance

統合後も元辞書を追跡できることを必須とする。

保持対象:

- candidateごとの `DictionaryIDs`
- tokenごとのraw `DictionaryResults`
- dictionary metadata
  - ID
  - name
  - source
  - license
  - language pair
  - release tier

UIは必要に応じてcandidateとsource dictionary IDsを表示できる。

結果統合のためにprovenanceを捨てない。

## 12. Adding a new language pair

例として Korean -> Japanese を追加する場合:

```text
1. LanguageKorean ("ko") を追加
2. Korean Tokenizer adapterを実装・登録
3. ko -> ja dictionary package/providerを追加
4. UIのsupported directionへ ko -> ja を追加
5. dictionary optionをUIへ公開
6. analyzer + dictionary + end-to-end testを追加
7. release tierを確認
```

既存の以下は原則変更不要。

- `TranslationProcessor.Merge`
- `DictionaryStore.LookupBatch` の言語ペアfilter
- status model
- candidate provenance model
- UI presenterの未解決語処理

もし新言語でこれらの変更が必要になった場合は、言語固有条件を直接埋め込む前にcontract/interfaceの一般化が必要かを検討する。

## 13. Current UI limitation

MVPの `ScreenProcessor` は、操作可能な翻訳方向を意図的に次へ限定している。

- `en -> ja`
- `ja -> en`

これはCoreの制限ではなくMVP UIの許可リスト。

新言語ペアを正式対応するときはUI validationも更新する。

将来的に対応方向が増えた場合、hard-coded pair判定をconfiguration/catalogへ移すことを検討する。現時点では2方向だけなので先行抽象化しない。

## 14. Next-language candidates

優先候補:

1. Korean
2. Chinese

これは実装順の候補であり、辞書license、Analyzer品質、配布可能性を確認する前にsupported扱いしない。

選定時は最低限次を評価する。

- word segmentation / lemma取得方法
- 再配布可能なDictionary Provider
- commercial useを含むlicense条件
- analyzer assetの配布条件
- word-by-word表示との相性
- regression test用の代表的な活用・未知語ケース

## 15. Extension invariants

言語・辞書拡張時も次を維持する。

- Analyzer固有処理をDictionaryへ入れない
- Dictionary固有formatをProcessへ入れない
- Data層にcandidate ranking判断を入れない
- UIからDataへ直接アクセスしない
- provider failureで他providerの結果を失わない
- provenanceを保持する
- 全文自然翻訳へ設計を変えない
- release tierとruntime capabilityを混同しない
