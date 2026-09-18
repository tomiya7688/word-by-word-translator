# Japanese analyzers

日本語の単語分割・基本形取得を翻訳辞書そのものから分離して扱います。

## Kagome + IPADIC

最初の高精度な日本語形態素解析adapterとして、次を固定採用します。

- Kagome: `github.com/ikawaha/kagome/v2` v2.9.9
- Analyzer dictionary: `github.com/ikawaha/kagome-dict/ipa` v1.2.0
- Adapter: `applications/main/process/processing/kagome/`
- Purpose: segmentation, POS-based symbol detection, base-form/lemma extraction

Kagome v2.9.9 は Go 1.19 を要求するため、本プロジェクトの Go 1.22 方針を維持したまま利用できます。現行Kagome v2はより新しいGoを要求するため、互換性を意図して v2.9.9 を固定します。

### Lookup behavior

Kagomeのsurfaceと基本形を分けます。

例:

```text
食べました
  食べ  -> lemma / lookup_unit: 食べる
  まし  -> lemma / lookup_unit: ます
  た    -> lemma / lookup_unit: た
```

これにより、翻訳辞書に辞書形 `食べる` が登録されていれば、活用形を入力しても同じ見出し語へ検索できます。

句読点・記号は `TokenKindSymbol` とし、辞書検索対象にはしません。

## Release classification

Kagome本体のコードはMIT Licenseです。一方、採用するIPADICデータはMeCab IPADIC由来で、再配布時にcopyright / license noticeを保持する条件があります。

そのため、本プロジェクトの「表示・notice不要」を条件とする **MIT Releaseの辞書/asset受入条件には現時点では含めません**。Kagome adapterは独立packageに置き、明示的に選択したビルドだけがIPADICを取り込むようにします。

- Translation dictionary tierとは別に扱う
- Full Releaseでは必要なthird-party noticeを同梱して利用可能
- MIT Release向けには、将来notice不要の解析データまたは別tokenizerを用意する

この分類はプロジェクトの配布ポリシー上の判定であり、個別ライセンスの法的判断を置き換えるものではありません。
