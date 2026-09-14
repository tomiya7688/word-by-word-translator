# Context Reducer Integration

このディレクトリは `ai-context-reducer` の原則を、このプロジェクトの開発フローへ適用するための補助領域です。

Reference repository:
- https://github.com/tomiya7688/ai-context-reducer
- Adoption reference commit: `09a8baf481d99c23a9959ea7295ec6642c2e56ad`

## Adopted now

- root `AI_CONTEXT.md`
- Search first, read second
- Goal / Required / Acceptance が揃ったら探索停止
- Source of Truth の明示
- targeted validation
- unrelated refactor の分離
- Unverified の明示
- Remote Delta First
- compact policy checks through UPD checker

## Remote context

Linux / macOS / Git Bash:

```bash
bash tools/context-reducer/script/remote-context.sh
```

Windows PowerShell:

```powershell
powershell -ExecutionPolicy Bypass -File tools/context-reducer/script/remote-context.ps1
```

このスクリプトは remote を勝手にmerge/rebaseしません。`fetch` 後に local/remote HEAD、ahead/behind、remote commit summary、changed files、diff stat、上限付きdiff抜粋だけを表示します。

## Not adopted yet

現在のrepo規模では維持コストが上回るため、次は必要になるまで導入しません。

- Source Structure Index
- generated Context Pack
- large Responsibility Map
- test impact analyzer
- large-scale routing tables

repoが成長し、AIが対象sourceへ到達しづらくなった時点で追加します。
