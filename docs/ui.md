# UI presentation

UIの表示責務は `applications/main/ui/` に閉じ込める。

## Translation result presentation

Process層の `TranslationResponse` をUI表示用tokenへ変換する。

- `success`: 最優先の翻訳候補を表示する
- `segmentation_failed`: 原文surfaceを表示し、未解決styleにする
- `dictionary_not_found`: 原文surfaceを表示し、未解決styleにする
- `lookup_error`: 原文surfaceを表示し、未解決styleにする
- 原文token順は変更しない

表示用tokenは `normal` / `unresolved` の意味的なtoneだけを持つ。色そのものはrenderer側で決めるため、将来のGUI frameworkへProcessの表示知識を漏らさない。

## Terminal renderer

MVPのterminal rendererでは `unresolved` をANSI redで表示する。

例:

```text
私 見た qwertymonster 昨日 .
```

このとき `qwertymonster` のみ赤色となる。

GUIを追加するときも `PresentedToken` のtoneをGUI widget/styleへ写像し、翻訳・辞書処理側へ色やwidget依存を追加しない。
