# AI Context

> AIが最初に読む小さい索引。詳細仕様をここへ複製しない。

## Project
- Name: word-by-word-translator
- Purpose: 単語単位で原文と訳語を対応させる多辞書対応翻訳アプリ
- Main language / runtime: Go

## Source of Truth
- Architecture: `docs/architecture.md`
- Project overview: `README.md`
- Issues / tasks: GitHub Issues
- Source / tests: 対象packageと対応する `_test.go`
- External design references:
  - ai-context-reducer: https://github.com/tomiya7688/ai-context-reducer
  - UPD Commander Base Design: https://github.com/tomiya7688/upd-commander-base-design

## Read First
1. current task / issue
2. この `AI_CONTEXT.md`
3. 対象sourceと対応test
4. 必要な場合だけ `docs/architecture.md`

## Ignore Normally
- build outputs / cache
- generated files
- large logs / dictionary datasets
- `.tools-cache/`
- unrelated Issues / docs / history

## Current Task Rules
- Goal / Required / Acceptance が十分なら追加探索を止める。
- Search first, read second。
- target source -> matching tests -> detailed docs の順を優先する。
- unrelated refactor を現在タスクへ混ぜない。
- 要約で不十分な場合だけ原典へ戻る。
- 複数AI/チャットから変更され得るため、実装前に remote delta を確認する。
- remote確認は commit summary / changed files / diff stat を先に使い、full diff は必要時だけ読む。
- 自動更新する場合は fast-forward only。dirty / diverged なら停止する。

## Architecture Invariants
- UI / Process / Data の責務を分離する。
- Commander は実処理をせず、処理呼び出しを指揮する。
- Messenger は層間通信だけを担当する。
- 実処理は各層の Processing に置く。
- UI -> Data の直接依存は禁止する。
- Process はGUI/CLI固有表示を知らない。
- Data は表示や翻訳判断を知らない。
- 辞書固有実装をProcessへ埋め込まず、Data側adapter境界で扱う。
- 複数辞書の統合判断はProcess側の責務とする。
- toolsの実装はアプリ本体から分離する。

## Validation
変更に必要な範囲だけ実行する。

- Format: `gofmt`
- Static: `go vet ./...`
- Tests: 対象package、必要なら `go test ./...`
- Architecture: Go UPD Commander Checker
- Final report: 実施した検証と `Unverified` を短く明示する

## Context Priority
- P0: current task / acceptance / architecture invariants
- P1: target source / matching tests
- P2: direct dependencies and interfaces
- P3: architecture/reference docs
- P4: history / unrelated information

## Adoption Profile
現在は小規模repoなので ai-context-reducer を過剰導入しない。

Adopted:
- Core AI index
- Search-first / exploration stop
- Source of Truth
- targeted validation
- Remote Delta First
- compact policy checks via UPD checker

Add later only when justified:
- Responsibility Map
- Task / Change Routing
- Source Structure Index
- generated context packs
- test impact routing
