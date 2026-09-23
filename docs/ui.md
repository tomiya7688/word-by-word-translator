# UI presentation

UIの表示責務は `applications/main/ui/` に閉じ込める。

## One-screen MVP

MVPのUI状態は `ScreenState` に集約する。

- 翻訳する文章
- 翻訳方向
- 利用可能な辞書
- 複数辞書の選択状態
- 翻訳結果
- UI上に表示するエラー

UI Commanderは入力イベントをUI Processingへ渡し、実行時だけUI Messengerを介してProcess層の翻訳Commanderへ要求する。UIからData層を直接参照しない。

`TranslationRequest.DictionaryIDs` にはUIで選択された辞書IDだけを載せる。Data層の `DictionaryStore` は指定がある場合、その辞書だけを検索する。既存API互換のため、辞書ID指定が空の場合は従来どおり対応言語ペアの全辞書を検索する。

Terminalの1画面rendererは次を同時表示する。

```text
Text: apple
Direction: en -> ja
Dictionaries:
[x] Dictionary A (dict-a)
[x] Dictionary B (dict-b)
Result: {りんごA [dict-a] / りんごB [dict-b]}
```

複数辞書から異なる候補が返った場合は候補を併記し、各候補の出典辞書IDも表示する。

## Translation result presentation

Process層の `TranslationResponse` をUI表示用tokenへ変換する。

- `success`: 最優先の翻訳候補を短縮表示用textにする
- 全ての統合候補と出典辞書IDを表示モデルへ保持する
- `segmentation_failed`: 原文surfaceを表示し、未解決styleにする
- `dictionary_not_found`: 原文surfaceを表示し、未解決styleにする
- `lookup_error`: 原文surfaceを表示し、未解決styleにする
- 原文token順は変更しない

表示用tokenは `normal` / `unresolved` の意味的なtoneだけを持つ。色そのものはrenderer側で決めるため、将来のGUI frameworkへProcessの表示知識を漏らさない。

## Terminal renderer

短縮terminal rendererでは各語の最優先候補だけを表示し、`unresolved` をANSI redで表示する。

例:

```text
私 見た qwertymonster 昨日 .
```

このとき `qwertymonster` のみ赤色となる。

1画面rendererでは、入力・方向・辞書選択・全候補・エラーを同時に表示する。

GUIを追加するときも `ScreenState` / `PresentedToken` をGUI widget/styleへ写像し、翻訳・辞書処理側へ色やwidget依存を追加しない。
