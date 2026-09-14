# Development Tools

開発支援ツールはアプリ本体から分離して `tools/` 以下に置く。

## Unified check

Linux / Ubuntu / Git Bash:

```bash
bash tools/dev-check/script/check.sh
```

Windows PowerShell:

```powershell
powershell -ExecutionPolicy Bypass -File tools/dev-check/script/check.ps1
```

Go sourceが存在する場合は次を順に確認する。

1. `gofmt`
2. `go vet ./...`
3. `go test ./...`
4. UPD Commander Checker

初期状態でGo sourceがまだ無い場合は1〜3だけskipし、UPD checkerは実行する。

## UPD Commander Checker

Reference:
- https://github.com/tomiya7688/upd-commander-base-design
- pinned commit: `41143698dda8bf2bd1f985539dce15d82477ff9d`

runnerは `.tools-cache/` に参照repoを取得し、pinされたGo checkerを実行する。アプリ本体や配布物には含めない。

Direct run:

```bash
bash tools/upd-commander-checker/script/run.sh
```

```powershell
powershell -ExecutionPolicy Bypass -File tools/upd-commander-checker/script/run.ps1
```

## Context Reducer

Reference:
- https://github.com/tomiya7688/ai-context-reducer
- adoption reference commit: `09a8baf481d99c23a9959ea7295ec6642c2e56ad`

詳細は `tools/context-reducer/README.md`。

Remote Delta First:

```bash
bash tools/context-reducer/script/remote-context.sh
```

```powershell
powershell -ExecutionPolicy Bypass -File tools/context-reducer/script/remote-context.ps1
```

## CI

`.github/workflows/ci.yml` から unified check を呼び出す。ローカルとCIでチェック経路を分けない。
