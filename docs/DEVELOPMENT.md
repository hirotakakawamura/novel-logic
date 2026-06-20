# novel-logic — 開発者向けガイド

> 最終更新: 2026-06-21  
> 対象: 本リポジトリにコードを貢献する開発者（複数人での並行開発を想定）

ユーザー向けの入門・コマンド仕様は別ドキュメントを参照してください。

| ドキュメント | 内容 |
|--------------|------|
| [README.md](../README.md) | プロジェクト概要・クイックスタート |
| [REQUIREMENTS.md](REQUIREMENTS.md) | 要件定義・ドメインモデル（正本） |
| [COMMANDS.md](COMMANDS.md) | CLI コマンド仕様 |
| [examples/momotaro-walkthrough/README.md](../examples/momotaro-walkthrough/README.md) | 手順付きチュートリアル |

---

## 1. プロジェクトの位置づけ

novel-logic は **小説の論理構造を管理し、矛盾を検出する CUI ツール**です。

- **やること**: plot / scene / thing / fact / action / rule / branch の登録、整合性検証、Lean 4 へのコード生成
- **やらないこと（v0）**: 散文の執筆・リッチ編集、本文からの自動事実抽出、GUI

開発上の基本方針:

| 方針 | 内容 |
|------|------|
| **正本は YAML** | 作品データはユーザーが `init` で作るディレクトリ内の YAML + `novels/` テキスト。Lean は生成物 |
| **登録時拒否** | 矛盾がある登録は CLI が拒否する（ロールバック）。事後検証だけに頼らない |
| **二段階検証** | Stage 1（Go・高速・常時）→ Stage 2（Lean 生成 + `lake build`） |
| **サブコマンド正** | Cobra ベースの明示的 CLI。対話ウィザードは Phase 1 |
| **Phase 0 完了** | MVP（branch/fork/merge、テスト、CI）は実装済み。以降は Phase 1 機能を段階追加 |

---

## 2. 技術スタック

| 層 | 技術 | バージョン目安 |
|----|------|----------------|
| CLI / ドメイン | Go | 1.22+（`go.mod` 準拠） |
| CLI フレームワーク | [cobra](https://github.com/spf13/cobra) | v1.8 |
| 永続化 | YAML（[yaml.v3](https://github.com/go-yaml/yaml)） | — |
| 論理検証 | Lean 4 + lake | ローカル `elan` 推奨 |
| CI | GitHub Actions | `go test ./...` のみ（Lean 不要） |

---

## 3. リポジトリ構成

```
novel-logic/
├── cmd/novel-logic/          # main エントリ（cli.Execute のみ）
├── internal/
│   ├── cli/                  # Cobra コマンド・入出力・終了コード
│   ├── project/              # ドメインモデル・Load/Save・mutation・branch ロジック
│   ├── validate/             # Stage 1 検証・登録前チェック
│   ├── generate/             # Lean ファイル生成（embed Core.lean）
│   ├── lean/                 # toolchain 検出・lake build ラッパ
│   ├── template/             # init 用テンプレート（embed）
│   └── testfixture/          # テスト用 YAML フィクスチャ（cli/generate/validate 向け）
├── docs/                     # 要件・コマンド・本ドキュメント
├── examples/
│   ├── momotaro/             # テンプレート由来の完成サンプル
│   └── momotaro-walkthrough/ # チュートリアル用プロジェクト（branch デモ含む）
├── .github/workflows/        # CI
└── go.mod
```

**作品データはリポジトリ外でもよい。** `novel-logic init <path>` で作るディレクトリが 1 作品分のプロジェクトルートです。

### 作品ディレクトリ（ユーザーが管理）

```
my-work/
├── project.yaml              # title, time_order, last_check
├── plot.yaml
├── things.yaml
├── scenes.yaml
├── times.yaml
├── branches.yaml
├── forks.yaml                # 空でも [] を置く（テンプレート同梱）
├── merges.yaml
├── facts.yaml
├── actions.yaml
├── rules.yaml
├── novels.yaml
├── novels/
│   ├── main/<scene_id>.txt
│   └── <branch>/<scene_id>.txt
└── logic/                    # check 時に自動生成（手編集非推奨）
    ├── Core.lean             # リポジトリ同梱テンプレのコピー
    ├── Project.lean
    ├── Facts.lean
    ├── Rules.lean
    ├── Timeline.lean
    ├── Theorems.lean
    └── lakefile.toml
```

---

## 4. アーキテクチャ

### 4.1 レイヤーと依存関係

```
┌─────────────────────────────────────────┐
│  cmd/novel-logic                        │
└──────────────────┬──────────────────────┘
                   │
┌──────────────────▼──────────────────────┐
│  internal/cli                           │
│  register / update / remove / show / …  │
└─┬────────────┬────────────┬─────────────┘
  │            │            │
  ▼            ▼            ▼
project    validate     generate ──► lean
  │
  ▼
YAML ファイル群（作品ディレクトリ）
```

**依存の向き**: `cli` → `project` / `validate` / `generate` / `lean`。`project` は他 internal パッケージに依存しない（循環回避の要）。

### 4.2 パッケージ責務

| パッケージ | 責務 |
|------------|------|
| `internal/project` | 型定義、`Load`/`Save`、CRUD（`mutate.go` / `update.go` / `remove.go`）、branch 進化（`ActiveActions`, `EffectiveFactsOnBranch`, `EffectiveRulesOnBranch`）、novel/git 補助 |
| `internal/validate` | Stage 1 全体検証（`Run` / `RunForBranch`）、登録前チェック（`CheckForbidState`, `CheckActionRules`, `CheckPredNotThingID`）、ヒント（`Hints`） |
| `internal/generate` | 作品データから `logic/*.lean` を生成。`Core.lean` は embed |
| `internal/lean` | `elan`/`lake` の PATH 検出、`lake build` 実行 |
| `internal/cli` | ユーザー操作の入口。`saveValidated` で mutation → validate → save を共通化 |
| `internal/template` | `init --template` 用ファイルの materialize |
| `internal/testfixture` | 複数パッケージで共有するテスト用最小プロジェクト |

### 4.3 典型的な処理フロー

**登録系コマンド**（`fact add` 等）:

1. `loadProject()` — `project.Load(projectPath)`
2. 登録前チェック — `validate.Check*`（rule 違反などを即拒否）
3. `saveValidated(d, mutate)` — ドメイン mutation → `validate.Run` → `project.Save`

**`validate` / `check`**:

1. `validate.Run` または `RunForBranch`
2. （`check` のみ）`generate.Run` → `lean.LakeBuild`

**重要**: 登録前チェックと事後 `Run` は役割が異なる。登録前は **提案中の 1 件**に対する rule チェック、事後は **プロジェクト全体**の参照整合・branch 整合など。

---

## 5. ドメイン上の注意点（実装時に必ず押さえる）

詳細は [REQUIREMENTS.md](REQUIREMENTS.md) が正本。開発で頻出する点のみ抜粋。

### 5.1 thing ID

- 作品内で ID は **一意**。plot / novel で同名 ID は **同一 thing**
- 既存 ID への `thing add` はエラー。スコープ追加は `thing scope add`

### 5.2 scope

| 値 | 意味 |
|----|------|
| `plot` または空 | プロット層 |
| `novel:<scene_id>` | 本文（scene）層 |

### 5.3 branch / fork / merge

- `fact` / `action` / `rule` / `novel` に `branch` フィールド（省略時 `main`）
- **time 軸は全 branch で共有**（`project.yaml` の `time_order` は 1 本）
- 子 branch の有効データは **lineage**（親 → 子の鎖）で決まる
  - `EffectiveFactsOnBranch` / `EffectiveRulesOnBranch` — lineage 上の branch に登録された fact/rule
  - `ActiveActions` — fork/merge 境界を考慮した action 列
- merge 後の子 branch への新規 action 登録は拒否
- 登録時 rule チェック・Lean `projectRules_<branch>` は **branch スコープ**（他 branch の rule は継承しない）

### 5.4 novel 本文

- CLI は散文を **書き込まない**。`novels/<branch>/<scene_id>.txt` をエディタ + git で管理
- `novel revision pin` で git commit を `novels.yaml` に記録

---

## 6. 検証（Stage 1 / Stage 2）

### Stage 1（Go）

常時実行可能。主な issue コード例:

| コード | 例 |
|--------|-----|
| `rule.violation` | forbid-state / forbid-transition 違反 |
| `ref.thing` / `ref.time` | 存在しない参照 |
| `branch.unknown` / `fork.invalid` / `merge.action_mismatch` | branch 整合 |
| `branch.isolated_state` | 他 branch 専用 pred を from に参照 |
| `duplicate` | ID 重複 |
| `pred.thing_collision` | pred 文字列が thing ID と衝突 |

実装: `internal/validate/validate.go`（全体）、`internal/validate/register.go`（登録前）。

branch 関連の構造化 issue: `project.BranchIssue`（`branch_validate.go`）。文字列パースは使わない。

### Stage 2（Lean）

`generate.Run` 後に `lake build`。生成される主な定理（Tier 0）:

- `no_forbidden_states` / `no_forbidden_transitions`
- `actions_in_scene_window`
- `fixed_facts_stable`
- branch 別: `*_main`, `*_branch_dog` 等（`activeActions_<branch>`, `projectRules_<branch>`）

Lean 未インストール時は `check` が exit 5 で終了（Stage 1 は成功していても Stage 2 スキップは警告）。

---

## 7. 開発環境セットアップ

### 7.1 必須

```bash
# Go（1.22+）
export PATH=$HOME/.local/go/bin:$PATH   # 環境に応じて

git clone https://github.com/hirotakakawamura/novel-logic.git
cd novel-logic
go build -o bin/novel-logic ./cmd/novel-logic
```

### 7.2 推奨（Stage 2 ローカル確認用）

```bash
# elan + Lean 4（未導入の場合）
curl https://raw.githubusercontent.com/leanprover/elan/master/elan-init.sh -sSf | sh
```

### 7.3 動作確認

```bash
go test ./...

# サンプルプロジェクトで end-to-end
./bin/novel-logic init ./tmp-momo --template momotaro
./bin/novel-logic -C ./tmp-momo check --quick    # Stage 1 のみ（CI 同等）
./bin/novel-logic -C ./tmp-momo check            # Stage 1 + Lean（要 toolchain）

# ウォークスルー（branch デモ含む）
./bin/novel-logic -C ./examples/momotaro-walkthrough check --quick
```

---

## 8. テスト

### 8.1 方針

| 種別 | 場所 | 内容 |
|------|------|------|
| ドメイン単体 | `internal/project/*_test.go` | Load/Save、branch、novel、重複検出 |
| 検証 | `internal/validate/*_test.go` | issue 検出、登録前 rule、ウォークスルー全体 |
| CLI 結合 | `internal/cli/cli_test.go` | `runCLI` で終了コード・出力を検証 |
| Lean 生成 | `internal/generate/generate_test.go` | golden file スナップショット |

```bash
go test ./...
go test -v ./internal/project -run TestBranch
```

### 8.2 テストフィクスチャ

`internal/testfixture` に最小 YAML セットを集約（`cli` / `generate` / `validate` が利用）。

**注意**: `internal/project` の `_test.go` は **同パッケージ**のため `testfixture` を import すると循環になる。`project` 内テストは `helper_test.go` の `newTestProject` を使う。

新しいテストで作品データが必要な場合:

- `project` パッケージ内 → `newTestProject(t)` または `writeTestFile`
- それ以外 → `testfixture.LoadMinimal(t)` / `LoadValidate(t)` / `WriteMinimalDir(t)`

### 8.3 Lean 生成物の golden 更新

`internal/generate` は `testdata/minimal/` と出力を byte 比較する。

```bash
UPDATE_GOLDEN=1 go test ./internal/generate
git diff internal/generate/testdata/   # 差分を確認してからコミット
```

`generate.go` の出力形式を変えたら **必ず** golden を更新し、PR に含める。

### 8.4 ウォークスルー固定テスト

`validate.TestWalkthroughProjectValidates` が `examples/momotaro-walkthrough` を Load して全体検証する。チュートリアル用 YAML を壊すと CI が落ちる。

---

## 9. CI

[`.github/workflows/test.yml`](../.github/workflows/test.yml):

- トリガー: `push` / `pull_request` → `main`
- `go test ./...` のみ（Lean 不要、数分以内）
- **Stage 2（`check` フル）は CI に含めない**（Lean セットアップコスト・作品データ依存のため）

PR を出す前にローカルで `go test ./...` を通すこと。Lean 変更時はローカルで `check` も確認する。

---

## 10. コーディング規約・実装パターン

### 10.1 一般

- 既存ファイルのスタイル・命名に合わせる（過剰なリファクタは避ける）
- ドメインロジックは `internal/project` に置き、CLI は薄く保つ
- ユーザー向けエラーは `fmt.Errorf` で具体的に。CLI では `exitErr` / `exitErrf` で終了コードを付与
- 新規 YAML ファイルをテンプレートに追加する場合は `internal/template/data/default/` と `momotaro/` **両方**に反映

### 10.2 CLI 追加の型

| 操作種別 | パターン |
|----------|----------|
| 登録 | `validate.Check*` → `saveValidated(d, func() { return d.Add* })` |
| 更新 | 既存レコード取得 → 登録前チェック（branch 込み）→ `saveValidated` + `d.Update*` |
| 削除 | `saveValidated` + `d.Remove*`（参照チェックは project 側） |
| 参照のみ | `loadProject` → 表示（`show.go` / `misc.go`） |

`saveValidated`（`internal/cli/util.go`）は mutation 後に **全体** `validate.Run` を走らせてから `Save` する。登録前チェックだけでは足りないケースがある。

### 10.3 branch を扱う変更のチェックリスト

新しい fact/action/rule まわりの機能を足すとき:

- [ ] 登録前: `EffectiveRulesOnBranch(normalize(branch))` を使っているか
- [ ] 事後検証: `RunForBranch` / `branchScopedIssues` が整合しているか
- [ ] Lean 生成: `projectRules_<branch>` / `activeActions_<branch>` が対応しているか
- [ ] `timeline --branch` 等の表示が `ActiveActions` / `EffectiveFactsOnBranch` と一致しているか
- [ ] テスト: 他 branch の rule が影響しないケースがあるか

### 10.4 終了コード（CLI）

| コード | 意味 |
|--------|------|
| 0 | 成功 |
| 1 | Stage 1 検証失敗 / 登録前 rule 違反 |
| 2 | Lean 生成失敗 |
| 3 | Stage 2（lake build）失敗 |
| 4 | プロジェクト読込・保存・入力エラー |
| 5 | Lean toolchain 未検出 |

実装: `internal/cli/errors.go`（`ExitError`）、`cmd/novel-logic/main.go` で `os.Exit`。

---

## 11. 機能追加の流れ（例）

### 例: 新しい Stage 1 チェックを追加

1. 要件を [REQUIREMENTS.md](REQUIREMENTS.md) に追記（issue コード名を決める）
2. `internal/validate/validate.go` または `internal/project/branch_validate.go` に検出ロジック
3. `internal/validate/*_test.go` に再現ケース
4. 必要なら [COMMANDS.md](COMMANDS.md) にユーザー向け説明
5. `go test ./...`

### 例: 新しい CLI サブコマンド

1. `internal/cli/` にコマンドファイル追加、`init()` で `rootCmd` に登録
2. ドメイン操作は `internal/project` にメソッド追加
3. `internal/cli/cli_test.go` に `runCLI` テスト
4. [COMMANDS.md](COMMANDS.md) 更新

### 例: Lean 定理の追加（Tier 0/1）

1. `internal/generate/templates/Core.lean` に判定関数が必要なら追加
2. `generate.go` の `genTheorems` に定理を追加
3. golden 更新 + ローカル `check`
4. [REQUIREMENTS.md §6.5](REQUIREMENTS.md) の Tier 表を更新

---

## 12. Phase ロードマップ（開発者向け）

### Phase 0（完了）

- Go CLI 一式、YAML 永続化、branch/fork/merge
- Stage 1 + Lean 生成 + Stage 2
- 桃太郎 end-to-end、単体テスト、GitHub Actions

### Phase 1（予定・着手時は REQUIREMENTS と相談）

| 項目 | 概要 |
|------|------|
| 対話ウィザード | `novel-logic wizard plot/novel`（サブコマンドと同じ mutation 経路） |
| plot ↔ novel 厳格化 | Phase 0 はヒント中心。横断矛盾を Stage 1 で強化 |
| Tier 1 定理 | `novel_extends_plot` 等 |
| scene ↔ novel 1:N | `novel_id` 導入の検討 |
| thing 参照付き fact | 「A の主人は B」形式 |
| `--json` / `--dry-run` | COMMANDS.md に予約記載あり |

**未決定事項**は [REQUIREMENTS.md §10.2](REQUIREMENTS.md) を参照。実装前に仕様を固め、ドキュメントを先に更新する。

---

## 13. Git ブランチ管理・協業

### 13.1 基本方針（確定）

**Trunk-Based Development** をベースに、**GitHub Flow 的な PR ルール**を採用する。

| 項目 | 方針 |
|------|------|
| 長期ブランチ | **なし**（`develop` 等は作らない） |
| 正（trunk） | `main` のみ。常に `go test ./...` が通る状態を維持 |
| 作業ブランチ | `main` から短命の feature ブランチを切り、PR 経由で merge |
| merge 方式 | **Squash merge**（`main` 上は 1 PR = 1 コミット） |
| レビュー | **1 承認必須** + CI green |
| `main` への直接 push | **メンテナーのみ可**（通常は PR 経由。緊急時の例外） |

```
main ─────●─────●─────●─────●─────●─────►
           \   /       \   /
            feat-A      fix-B
         （PR + 1 approve + CI）
```

### 13.2 ブランチ命名

```
feat/<短い説明>     # 機能追加（例: feat/wizard-plot）
fix/<短い説明>      # バグ修正
docs/<短い説明>     # ドキュメントのみ
chore/<短い説明>    # CI・依存更新・リファクタ
test/<短い説明>     # テスト追加のみ（任意）
```

- 英語小文字 + ハイフン区切り
- Issue 番号がある場合は `feat/12-wizard-plot` も可
- 分岐元は常に **最新の `main`**

### 13.3 作業フロー

1. `git checkout main && git pull`
2. `git checkout -b feat/<name>`
3. 実装・`go test ./...`（Lean 変更時はローカル `check` も）
4. push して **Pull Request** を作成（Draft 可）
5. CI green + **1 名以上の Approve**
6. **Squash merge** で `main` に取り込み
7. 作業ブランチを削除

| 種別 | 分岐元 | merge 先 | 寿命の目安 |
|------|--------|----------|------------|
| 通常作業 | `main` | `main` | 1〜5 日 |
| hotfix | `main` | `main` | 当日 |
| 大きい Phase 1 機能 | `main` | `main`（**複数 PR に分割**） | 設計 doc PR を先に merge 推奨 |

**非推奨**

- 1 週間以上更新のないブランチ（rebase して早期 merge、または close）
- 1 PR に無関係な refactor や複数機能を混在

### 13.4 PR ルール

| ルール | 内容 |
|--------|------|
| CI | `go test ./...` が green であること（必須チェック） |
| 承認 | **1 承認必須**（承認者は PR 作者と異なることが望ましい） |
| スコープ | 1 PR = 1 目的 |
| docs | 仕様・CLI 変更は **同一 PR** に `docs/` 更新を含める |
| golden | `generate` 変更時は `testdata/minimal/` 差分を PR に含める |
| merge | **Squash merge** のみ使用 |

メンテナーが `main` に直接 push した場合も、**事後に PR 相当の説明**（コミットメッセージまたは Issue 参照）を残す。

### 13.5 GitHub `main` 保護設定（推奨）

リポジトリ Settings → Branches → Branch protection rules で `main` に以下を設定する。

| 設定 | 値 |
|------|-----|
| Require a pull request before merging | on（メンテナー bypass は許可） |
| Required approvals | **1** |
| Require status checks to pass | `test` job |
| Require branches to be up to date | 任意（初期は off でも可） |
| Allow force pushes | off |
| Allow deletions | off |

### 13.6 リリース・タグ（将来）

ユーザー向け配布を始める段階で `main` 上のコミットにセマンティックバージョンのタグを打つ。

```
v0.1.0  … Phase 0 相当
v0.2.0  … Phase 1 のマイルストーン
```

- タグは **`main` のみ**に付与
- CHANGELOG は Phase 1 着手時に `CHANGELOG.md` を新設予定

### 13.7 コンフリクトしやすい箇所

並行開発時は以下を同時に触らない、または早めに merge する。

| 箇所 | 理由 |
|------|------|
| `internal/generate/generate.go` + `testdata/minimal/` | 出力形式変更は golden 一式が変わる |
| `docs/REQUIREMENTS.md` | 仕様の正本。複数人で同時編集しやすい |
| `examples/momotaro-walkthrough/` | 固定テストが参照 |
| `internal/validate/validate.go` | issue コードの追加・変更 |

大きい変更は **設計 PR（docs のみ）→ 実装 PR** の 2 段階を推奨。

### 13.8 コミットメッセージ

Squash merge 後の `main` 上のコミットメッセージ（= PR タイトル）の書き方:

- 1 行目に変更の要旨（英語・日本語どちらでも可）
- PR 本文に「なぜ」「テスト方法」を書く

例:

```
fix: use EffectiveRulesOnBranch in fact registration
```

PR 本文:

```
Branch-scoped forbid-state rules were incorrectly applied on main.

Test: go test ./...
```

---

## 14. よくある落とし穴

| 問題 | 対処 |
|------|------|
| `go test ./internal/project` で import cycle | `project` テストから `testfixture` を import しない |
| 登録は通るが `validate` で落ちる | 登録前チェックと `Run` の両方を確認。`saveValidated` は全体検証する |
| branch A の rule が main に効く | `CheckForbidState(d, branch, ...)` の branch 引数を確認 |
| `check` だけ失敗 | ローカルに Lean があるか。`logic/.lake` は gitignore 済み |
| 日本語 pred の Lean 名 | `generate.go` の `leanPred`（ローマ字表 or hash）。新しい定番語は map 追加を検討 |
| init 後に forks/merges が無い | テンプレートに空 `[]` を入れる（`internal/template/data/`） |

---

## 15. 改訂履歴

| 日付 | 内容 |
|------|------|
| 2026-06-21 | 初版。Phase 0 完了時点の構成・方針・協業ガイド |
| 2026-06-21 | §13 Git ブランチ管理方針を確定（Trunk-Based + Squash + 1 承認） |