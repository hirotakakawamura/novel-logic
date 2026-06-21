# novel-logic

小説の**論理構造・時系列・制約**を管理し、整合性を検証する CUI アプリケーション。

- 執筆エディタではない。事実・状態・ルール・遷移の管理と矛盾検出が目的
- 検証は **二段階**（Go 内蔵チェック → Lean 4 による厳密チェック）
- 実装: **Go**（CLI）+ **Lean 4**（`logic/` へ自動生成）

## はじめに読むもの

| ファイル | 内容 |
|----------|------|
| **[examples/momotaro-walkthrough/README.md](examples/momotaro-walkthrough/README.md)** | **入門チュートリアル**（桃太郎を空から登録する手順） |
| [docs/REQUIREMENTS.md](docs/REQUIREMENTS.md) | 要件定義・ドメインモデル |
| [docs/COMMANDS.md](docs/COMMANDS.md) | CUI コマンド仕様 |
| [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) | **開発者向け**（構成・テスト・協業ガイド） |

初心者はまず **桃太郎ウォークスルー** を上から順に試すのがおすすめです。

## クイックスタート

### ビルド

```bash
export PATH=$HOME/.local/go/bin:$PATH   # 必要に応じて
cd novel-logic
go build -o bin/novel-logic ./cmd/novel-logic
```

### サンプルで動作確認（最短）

```bash
./bin/novel-logic init ./examples/momotaro --template momotaro
./bin/novel-logic -C ./examples/momotaro check
# → OK: stage1 + stage2
```

### 手順を学びながら登録（推奨）

```bash
# チュートリアル付きサンプル（本リポジトリに同梱）
cat examples/momotaro-walkthrough/README.md

# 空から自分で積み上げる場合
./bin/novel-logic init ./my-work --template default
cd my-work
# README の Step 1 から順に実行
```

## 主なコマンド

```bash
novel-logic init <path> [--template default|momotaro]
novel-logic time add <id>
novel-logic scene add <id> --summary ... --time-start ... --time-end ...
novel-logic thing add <id> --name ... --tag ...          # 新規 ID のみ
novel-logic thing scope add <id> --scope novel:scene1   # 既存 thing にスコープ追加
novel-logic fact add --kind fixed|state --thing ... --pred ...
novel-logic rule add --kind forbid-state|forbid-transition ...
novel-logic action add --thing ... --from ... --to ... --at ...
novel-logic novel add <scene_id> [--branch main]        # メタ登録 + novels/<branch>/<scene>.txt
# 本文はエディタで novels/<branch>/<scene>.txt を編集（CLI は散文を書き込まない）
novel-logic branch list / fork add / merge add          # 物語分岐・合流
novel-logic timeline --branch <id>                      # branch 限定の action 列
novel-logic novel revision pin <scene_id> [--branch main]
novel-logic validate [--branch <id>]   # Stage 1 のみ（省略時は全 branch）
novel-logic check [--branch <id>]      # Stage 1 + Lean 生成 + lake build
novel-logic check --quick              # Stage 1 のみ（Lean 不要）

# 更新（同一内容の add は拒否 → update を使用）
novel-logic thing update <id> --name ... --tag ...
novel-logic fact update <id> --pred ...
novel-logic action update <id> --label ...

# 削除（参照がある場合は拒否）
novel-logic thing remove <id>
novel-logic thing scope remove <id> --scope novel:scene1
novel-logic scene remove <id>
novel-logic time remove <id>
novel-logic fact remove <id>
novel-logic action remove <id>
novel-logic rule remove <id>
novel-logic novel remove <scene_id>   # デフォルトで本文 .txt は残す（--keep-body）
```

### 状況確認（show / list）

登録したデータを読み取り専用で確認します。詳細は [桃太郎ウォークスルー Step 9](examples/momotaro-walkthrough/README.md#step-9-show--list-で状況確認読み取り専用) を参照。

```bash
PJ=./examples/momotaro-walkthrough

# 全体
./bin/novel-logic -C $PJ info
./bin/novel-logic -C $PJ plot show
./bin/novel-logic -C $PJ timeline

# 個別
./bin/novel-logic -C $PJ thing show momotaro
./bin/novel-logic -C $PJ scene show scene4
./bin/novel-logic -C $PJ action list
./bin/novel-logic -C $PJ novel show scene5
```

`bin/` から実行するときは `novel-logic` ではなく **`./novel-logic`**（または PATH へ `~/novel-logic/bin` を追加）。

## ステータス

- **Phase 0（MVP）**: 完了 — Go CLI、YAML 永続化、branch/fork/merge、Stage 1/2、桃太郎 end-to-end、単体テスト、GitHub Actions CI（Go test + examples Stage 2）
- **Phase 1 予定**: 対話ウィザード（`wizard`）、`--json` / `generate --dry-run`、plot ↔ novel 横断の厳格化、Tier 1 定理 — 詳細は [COMMANDS.md §8–§9](docs/COMMANDS.md)

## テスト・CI

```bash
go test ./...    # internal/cli, generate, project, validate（Lean 不要）

# Stage 2（Lean 要。ローカル確認）
./bin/novel-logic -C examples/momotaro-walkthrough check -q
```

[`.github/workflows/test.yml`](.github/workflows/test.yml) は 2 ジョブを実行します。

| ジョブ | 内容 |
|--------|------|
| `test` | `go test ./...` |
| `lean-check` | elan セットアップ後、`examples/momotaro-walkthrough` と `examples/momotaro` で `check`（Stage 2） |

## ディレクトリ構成

```
novel-logic/
  cmd/novel-logic/         # CLI エントリポイント
  internal/                # ドメイン・検証・Lean 生成
  .github/workflows/       # CI（go test + lean-check）
  docs/                    # 要件・コマンド仕様
  examples/
    momotaro/              # テンプレート由来の完成サンプル
    momotaro-walkthrough/  # 手順付きチュートリアル（branch_dog 分岐デモ含む）
```

作品データはユーザーが `init` で作成する別ディレクトリでも管理できます。本文は `novels/<branch>/<scene_id>.txt`（本線は `novels/main/`）。