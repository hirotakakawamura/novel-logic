# novel-logic — CUI コマンド仕様

> 最終更新: 2026-06-21
> 関連: [REQUIREMENTS.md](REQUIREMENTS.md) · 正本: [DRAFT.txt](DRAFT.txt)

バイナリ名は **`novel-logic`**（確定）。以下は `novel-logic` で記載する。

ドメインモデル: **plot · scene · novel · thing · fixed_fact · state · action · rule · time**（thing に **tag** 付与）。  
登録コマンドはいずれも **Stage 1 矛盾チェック**を通過しないと永続化しない（[REQUIREMENTS §4.1](REQUIREMENTS.md)）。

**CLI 体系（確定）**: サブコマンド型が正。対話ウィザード（`novel-logic wizard`）は Phase 1 の補助（[REQUIREMENTS §8.1](REQUIREMENTS.md)）。

**用語（YAML / CLI）**: ドキュメント上の概念名 `fixed_fact` は、YAML と CLI では **`kind: fixed`** と書く（`state` も同様）。

---

## 0. 共通仕様

### 0.0 用語集（混同しやすい語）

| 語 | 意味 |
|----|------|
| **branch**（`branches.yaml`） | 物語のルート分岐 ID（`main`, `branch_a` 等） |
| **revision pin**（`novel revision pin`） | 本文ファイルの **git commit** を `novels.yaml` に記録する操作 |
| **git branch** | ツール外の VCS 概念（作品データの履歴管理） |

### 0.1 実行コンテキスト

- 作品ルート（`project.yaml` があるディレクトリ）で実行するのが基本。
- 他ディレクトリからは **`-C <path>`** / **`--project <path>`** で作品ルートを指定。

```bash
novel-logic -C ~/novels/momotaro/ check
```

### 0.2 スコープ

plot / novel 二層のデータ操作で使う。

| 値 | 意味 | 対応フロー |
|----|------|-----------|
| `plot` | プロット全体 | Phase A |
| `novel:<scene_id>` | 特定 scene に紐づく本文側 | Phase B |

フラグ **`--scope`** のデフォルトはコマンドごとに §3–§5 で定義。

### 0.3 終了コード

| コード | 意味 |
|--------|------|
| `0` | 成功 |
| `1` | Stage 1 検証エラー、登録拒否（重複 ID（thing / fact / action / time / novel 等）、rule 抵触、fact 昇格/降格の不正経路、merge 後の branch 登録拒否、`time.registry_mismatch` 等） |
| `2` | Lean 生成エラー |
| `3` | Stage 2（`lake build`）失敗 |
| `4` | ユーザー入力・引数エラー（未知 ID、remove 拒否、必須フラグ不足等） |
| `5` | 環境エラー（Lean / elan 未検出等） |

### 0.4 出力形式

| フラグ | 効果 |
|--------|------|
| `-q` / `--quiet` | エラー時以外は出力抑制 |
| `-v` / `--verbose` | 詳細ログ（矛盾チェックの内訳等） |
| `--json` | 機械可読出力（**Phase 1 予定**。現行実装では未対応） |

### 0.5 登録拒否時の応答

矛盾がある場合、CLI は次を出力して **終了コード 1** とする。

1. 拒否理由（どの rule / fact / time と矛盾したか）
2. 関連エンティティの ID（可能なら）

---

## 1. プロジェクト・ユーティリティ

### `novel-logic init <path>`

新規作品ディレクトリを作成する（[REQUIREMENTS §7.2–7.3](REQUIREMENTS.md)）。

| 引数・フラグ | 必須 | 説明 |
|-------------|------|------|
| `<path>` | ○ | 作成先ディレクトリ |
| `--template <name>` | — | `default` / `momotaro`（空 / 桃太郎サンプル） |
| `--force` | — | 空でないディレクトリでも上書き（危険） |

**生成物（例）**

```
<path>/
  project.yaml          # 作品メタ
  plot.yaml             # plot 定義
  things.yaml
  scenes.yaml
  times.yaml
  branches.yaml         # story branch 定義（`main` 含む）
  forks.yaml            # 分岐点（初回 fork 登録時に作成）
  merges.yaml           # 合流点（初回 merge 登録時に作成）
  facts.yaml            # fixed_fact + state（`branch` 省略時 `main`）
  actions.yaml
  rules.yaml
  novels.yaml           # novel メタ（scene × branch ごと1件）
  novels/               # 本文（`<branch>/<scene_id>.txt`、git 管理）
  logic/                # generate 後（手編集非推奨）
```

---

### `novel-logic info`

作品名・パス・登録件数サマリ・最終 `check` 日時を表示。

---

### `novel-logic doctor`

環境診断: 自バイナリ版、`elan` / `lean` / `lake` の PATH、作品必須ファイル、推奨アクション。

---

### `novel-logic template list`

利用可能テンプレ一覧。

---

### `novel-logic version`

CLI バージョンと Lean Core テンプレバージョン。

---

## 2. 参照・一覧（読み取り専用）

### `novel-logic status`

健全性サマリ。

- 最終 `check` 結果（成功 / 失敗 / 未実行）
- エンティティ件数（thing / scene / fact / action / rule / time）
- Stage 1 / Stage 2 の個別状態
- Lean ツールチェーン検出結果

---

### `novel-logic timeline`

**time 軸**上の scene 区間・action 配置を表形式で表示（旧「章・視聴者」タイムラインの代替）。

| フラグ | 説明 |
|--------|------|
| `--branch <id>` | 表示する active action を branch 限定（デフォルト: `main`） |
| `--verbose` | 各 time の fact / action 概要 |

---

### `novel-logic plot show`

plot 要約と scene 一覧。

---

### `novel-logic scene list` / `novel-logic scene show <id>`

scene 一覧、または1件の概要・time 区間・紐づく novel の有無。

---

### `novel-logic novel list` / `novel-logic novel show <scene_id>`

scene に紐づく novel 一覧、または本文プレビュー + time 区間 + 二層（novel / plot）の fact/action 件数 + git pin 状態。

`novel show` は `body_file:`（パス・編集方法）と `git:`（pin 状態・次の手順）を分けて表示する。

| フラグ | 説明 |
|--------|------|
| `--full` | 本文全文表示 |

---

### `novel-logic thing list` / `novel-logic thing show <id>`

thing 一覧・詳細（**tag**、紐づく fixed_fact / state 件数）。

| フラグ | 説明 |
|--------|------|
| `--tag <tag>` | tag でフィルタ（例: `character`, `location`） |
| `--scope plot\|novel:<scene>\|all` | スコープフィルタ |

---

### `novel-logic time list`

登録済み time ID と順序。**表示順は `project.yaml` の `time_order`**（正本）。`times.yaml` は ID の存在レジストリ。両者の乖離は Stage 1 の `time.registry_mismatch` で検出する。

---

### `novel-logic fact list` / `novel-logic fact show <id>`

fixed_fact / state の一覧・詳細。

| フラグ | 説明 |
|--------|------|
| `--kind fixed\|state\|all` | 種別フィルタ（デフォルト: `all`） |
| `--thing <id>` | thing でフィルタ |
| `--scope <scope>` | スコープフィルタ |

---

### `novel-logic action list` / `novel-logic action show <id>`

action 一覧・詳細（subject、遷移 from→to、time、scope）。

---

### `novel-logic rule list` / `novel-logic rule show <id>`

rule 一覧・詳細（静的禁止 / 遷移禁止の区別を表示）。

---

## 3. 登録 — Phase A（プロット構築）

> フロー対応: [REQUIREMENTS §4 Phase A](REQUIREMENTS.md)

### `novel-logic plot set`

plot を登録または更新する（A1）。

| フラグ | 必須 | 説明 |
|--------|------|------|
| `--title <text>` | △* | 作品タイトル（`project.yaml`） |
| `--summary <text>` | △* | プロット概要（`plot.yaml`） |
| `--file <path>` | △* | 概要テキストファイル（`--summary` の代替） |

\* **`--title` / `--summary` / `--file` のいずれか1つ以上**が必須。新規作品では `--title` を推奨。

---

### `novel-logic thing add <id>`

**未登録 ID** の thing を新規作成（A2 / B2 共通。[REQUIREMENTS §3.8](REQUIREMENTS.md)）。既存 ID への `add` は **エラー**。

| フラグ | 必須 | 説明 |
|--------|------|------|
| `--name <text>` | — | 表示名 |
| `--tag <tag>` | ○* | 分類タグ（複数回指定可）。*新規時は1つ以上必須 |
| `--scope <scope>` | — | `plot`（デフォルト）または `novel:<scene_id>` |

例:

```bash
novel-logic thing add momotaro --name 桃太郎 --tag character --scope plot
novel-logic thing add mob_crowd --tag character --scope novel:scene3  # novel 初出の新規 thing
```

### `novel-logic thing scope add <id>`

既存 thing にスコープを追加（冪等）。

| フラグ | 必須 | 説明 |
|--------|------|------|
| `--scope <scope>` | ○ | 追加するスコープ（繰り返し可） |

```bash
novel-logic thing scope add momotaro --scope novel:scene5
```

### `novel-logic thing update <id>`

既存 thing の `name` / `tags` を更新（`--name` または `--tag` のいずれか必須）。

### `novel-logic thing scope remove <id>`

既存 thing からスコープを削除。

---

### `novel-logic scene add <id>`

scene を登録（A3）。

| フラグ | 必須 | 説明 |
|--------|------|------|
| `--summary <text>` | ○ | scene 概要 |
| `--time-start <time_id>` | ○ | 開始 time |
| `--time-end <time_id>` | ○ | 終了 time |

---

### `novel-logic time add <id>`

time を順序付きで登録（scene / action の前提）。**重複 ID は拒否（exit 1）**。

| フラグ | 必須 | 説明 |
|--------|------|------|
| `--after <time_id>` | — | 指定 time の直後に挿入（省略時は末尾） |

---

### `novel-logic fact add`

fixed_fact または state を登録（A4）。

| フラグ | 必須 | 説明 |
|--------|------|------|
| `--kind fixed\|state` | ○ | 種別 |
| `--thing <id>` | ○ | 主語 thing |
| `--pred <text>` | ○ | 述語（`PredId` に写像。既存 **thing ID と同名は拒否** — [REQUIREMENTS §6.3.1](REQUIREMENTS.md)） |
| `--scope plot` | — | デフォルト `plot` |
| `--branch <id>` | — | story branch（デフォルト: `main`） |

例:

```bash
novel-logic fact add --kind fixed --thing momotaro --pred 人間
novel-logic fact add --kind state --thing momotaro --pred 赤ちゃん
novel-logic fact add --kind state --thing momotaro --pred 犬仲間あり --branch main
```

**重複**: 同一 `(kind, thing, pred, scope)` の `add` は拒否。変更は `fact update <id>`。

### `novel-logic fact update <id>`

| フラグ | 説明 |
|--------|------|
| `--kind` / `--thing` / `--pred` / `--scope` | 指定した項目のみ更新 |

**`kind` の変更**: `fixed` ↔ `state` の昇格・降格は **`fact update --kind` では不可**（拒否時 exit 1）。昇格は `fact promote`、降格は不可（[REQUIREMENTS §3.2](REQUIREMENTS.md)）。

---

### `novel-logic fact promote <id>`

`kind: fixed` の fact を `state` に昇格する正規経路（降格は不可 — REQUIREMENTS §3.2）。

---

### `novel-logic action add`

action を登録（A5）。

| フラグ | 必須 | 説明 |
|--------|------|------|
| `--thing <id>` | ○ | 主語 thing |
| `--from <pred>` | — | 遷移元述語（省略可＝初期状態への遷移。`forbid-transition` は **from 非空時のみ**照合） |
| `--to <pred>` | ○ | 遷移先述語（`forbid-state` は to で照合） |
| `--at <time_id>` | ○ | 発生 time |
| `--label <text>` | — | 人間向けラベル（例: `育った`） |
| `--scope plot` | — | デフォルト `plot` |
| `--branch <id>` | — | story branch（デフォルト: `main`） |

例:

```bash
novel-logic action add --thing momotaro --from 赤ちゃん --to 青年 --at t4 --label 育った
novel-logic action add --branch branch_dog --thing inu --from 野良 --to 仲間 --at t8 --label 犬のみ仲間
```

**time 窓**: `scope=novel:<scene>` の action は scene の `[time_start, time_end]` 内であること（Stage 1 エラー）。`scope=plot`（デフォルト）は窓チェック対象外 — `validate` の `action.plot_scene_hint` で Phase B 整合を促す。

**重複**: 同一 `(thing, from, to, at, scope)` の `add` は拒否（`label` は重複キーに含まない）。変更は `action update <id>`。

### `novel-logic action update <id>`

| フラグ | 説明 |
|--------|------|
| `--thing` / `--from` / `--to` / `--at` / `--scope` / `--label` | 指定した項目のみ更新 |

---

### `novel-logic rule add`

rule を登録（プロット設計段階で推奨）。

| フラグ | 必須 | 説明 |
|--------|------|------|
| `--kind forbid-state` | ○* | 静的禁止: `--thing` + `--pred` と併用 |
| `--kind forbid-transition` | ○* | 遷移禁止: `--from` + `--to` と併用 |
| `--thing <id>` | △ | `forbid-state` 時必須 |
| `--pred <text>` | △ | `forbid-state` 時必須 |
| `--from <pred>` | △ | `forbid-transition` 時必須（rule 定義側。action 側の from 省略は初期遷移として許容） |
| `--to <pred>` | △ | `forbid-transition` 時必須 |
| `--branch <id>` | — | story branch（デフォルト: `main`） |

例:

```bash
novel-logic rule add --kind forbid-transition --from 青年 --to 赤ちゃん
novel-logic rule add --kind forbid-state --thing momotaro --pred 動物
```

**重複**: 同一 rule キーの `add` は拒否。変更は `rule update <id>`。

### `novel-logic rule update <id>`

| フラグ | 説明 |
|--------|------|
| `--kind` / `--thing` / `--pred` / `--from` / `--to` | 指定した項目のみ更新 |

---

## 3.5 分岐・合流（branch / fork / merge）

| コマンド | 用途 |
|----------|------|
| `branch list` / `show` / `add` / `remove` | branch 定義 |
| `fork add` / `fork choice add` / `fork show` | 分岐点 |
| `merge add` / `merge show` | 合流点 |

登録系コマンドに `--branch <id>`（デフォルト `main`）を付与。

```bash
# 分岐例（桃太郎 t7 で犬ルート / 単独ルート）
novel-logic fork add fork_t7 --parent main --at t7
novel-logic action add --branch main --thing momotaro --from 旅立ち --to 犬仲間あり --at t7 --label 犬ルート
novel-logic fork choice add --fork fork_t7 --action act_dog_route --branch branch_dog

# 合流（各 branch で merge action、to は共通）
novel-logic action add --branch branch_dog --thing momotaro --from 犬仲間あり --to 鬼退治準備 --at t11 --label 合流
novel-logic merge add merge_t11 --at t11 --into main --choice branch_dog:act_reunite_dog

novel-logic timeline --branch branch_dog
novel-logic novel add scene4 --branch branch_dog   # → novels/branch_dog/scene4.txt
```

---

## 4. 登録 — Phase B（本文）

> フロー対応: [REQUIREMENTS §4 Phase B](REQUIREMENTS.md)

**本文の扱い（確定）**

- 散文テキストは **`novels/<branch>/<scene_id>.txt`** に置き、**エディタ + git** で管理する。
- CLI は **本文を書き込まない**（`novel set --text` は廃止）。
- `novels.yaml` にはメタ（パス・time・git revision）のみ。time は scene と自動同期。

### `novel-logic novel add <scene_id>`

scene に紐づく novel メタを登録（B1）。デフォルトで空の `novels/<branch>/<scene_id>.txt` を作成。

**cardinality（確定）**: **1 scene × 1 branch : 1 novel**。同一 branch で再登録は **登録拒否（exit 1）**（[REQUIREMENTS §3.6](REQUIREMENTS.md)）。

**重複メッセージの違い**: `novel add` の再登録は `novel for scene … already registered`（CLI 操作時）。`novels.yaml` に同一 `(scene, branch)` が二重に存在する場合は `validate` が `novel.duplicate`（`duplicate novel: scene … registered twice`）を報告する。後者は手編集 YAML の整合チェック用。

| フラグ | 必須 | 説明 |
|--------|------|------|
| `--branch <id>` | — | story branch（デフォルト: `main`） |
| `--file <path>` | — | 本文パス（デフォルト: `novels/<branch>/<scene_id>.txt`）。**参照のみ** |
| `--init` | — | 本文ファイルが無ければ空ファイル作成（デフォルト: `true`） |
| `--pin` | — | 登録直後に `novel revision pin` を実行 |
| `--note <text>` | — | pin 時のメモ（PR 番号等） |
| `--allow-dirty` | — | 未コミット変更があっても pin を許可 |

```bash
novel-logic novel add scene1
# エディタで novels/main/scene1.txt を編集 → git commit
novel-logic novel revision pin scene1 --note "初稿"
```

### `novel-logic novel update <scene_id>`

`novels.yaml` のメタのみ更新（`body_path` 変更、scene time との再同期）。**本文テキストは変更しない**。

| フラグ | 説明 |
|--------|------|
| `--file <path>` | 新しい `body_path`（ファイルは既に存在すること） |

### `novel-logic novel revision pin <scene_id>`

本文ファイルに対応する **git commit** を `novels.yaml` に記録。

| フラグ | 説明 |
|--------|------|
| `--branch <id>` | story branch（デフォルト: `main`） |
| `--revision <sha>` | 明示的な commit SHA（省略時: 当該ファイルの最新 commit） |
| `--note <text>` | メモ（PR 番号等） |
| `--allow-dirty` | 作業ツリーに未コミット変更があっても pin |

### `novel-logic novel revision list <scene_id>`

pin 履歴（`revisions[]`）を表示。

### `novel-logic novel remove <scene_id>`

| フラグ | デフォルト | 説明 |
|--------|-----------|------|
| `--branch <id>` | `main` | story branch |
| `--keep-body` | `true` | 本文 `.txt` をディスクに残す（git 履歴保持） |

---

### novel スコープでの thing 登録（B2）

§3 `thing add` と同一コマンド。`--scope novel:<scene_id>` を指定する（[REQUIREMENTS §3.8](REQUIREMENTS.md)）。

---

### `novel-logic fact add` / `novel-logic action add`（novel スコープ）

A4 / A5 と同型。`--scope novel:<scene_id>` を指定（B3 / B4）。

Phase 0: Stage 1 は `novel:<scene>` の **time 窓**（`time.action_window`）のみ厳格。plot との整合は `validate` の **hint**（`action.plot_scene_hint`）。plot ↔ novel 横断の厳格検証（`novel_extends_plot`）は **Phase 1 / Tier 1**（[REQUIREMENTS §6.5](REQUIREMENTS.md)）。

---

## 4.5 削除（remove）

各エンティティを ID（または `scene_id`）で削除。他から参照されている場合は **拒否** し参照元を表示。

```bash
novel-logic thing remove <id>
novel-logic thing scope remove <id> --scope novel:scene1
novel-logic scene remove <id>
novel-logic time remove <id>
novel-logic fact remove <id>
novel-logic action remove <id>
novel-logic rule remove <id>
novel-logic novel remove <scene_id>    # デフォルト --keep-body
```

---

## 5. 検証・生成（コア）

### `novel-logic validate`

**Stage 1 のみ**。作品データ全体の構造・制約を検証。Lean 不要。

- スキーマ・参照 ID
- time 順序・`time_order` と `times.yaml` の整合（`time.registry_mismatch` — 欠落・重複・空 ID）
- scene / novel 区間
- rule / fact / action の抵触（**branch ごとの有効 action**）
- fork / merge 整合、`branch.isolated_state`
- tag 形式
- **hints**（`action.plot_scene_hint` 等）— 存在時または `--verbose` で表示

| フラグ | 説明 |
|--------|------|
| `--branch <id>` | 単一 branch のみ検証（省略時: 全 branch） |

**`validate` と `check --quick` の違い**

| | `validate` | `check --quick` |
|---|------------|-----------------|
| Stage 1 | ○ | ○ |
| hints 表示 | ○ | なし |
| `last_check` 更新 | なし | ○ |
| Lean | 不要 | 不要 |

---

### `novel-logic generate`

作品データ → `logic/` 以下の Lean 4 プロジェクトを生成（上書き）。  
構成は [REQUIREMENTS §6.2](REQUIREMENTS.md)（Core, Project, Facts, Rules, Timeline, Theorems）。

| フラグ | 説明 |
|--------|------|
| `--dry-run` | 書き込まず差分プレビュー（**Phase 1 予定**。現行実装では未対応） |

---

### `novel-logic check`

**メインコマンド**。保存・登録後の本番確認。

**デフォルト（確定）**: Stage 2 まで実行（[REQUIREMENTS §6.6](REQUIREMENTS.md)）。

```
validate → generate → lake build
```

高速確認は `novel-logic validate` または `novel-logic check --quick` を使う。

| フラグ | 説明 |
|--------|------|
| `--quick` | `validate` のみ（Stage 1） |
| `--branch <id>` | 単一 branch のみ Stage 1 検証（省略時: 全 branch） |
| `--no-generate` | 既存 `logic/` のまま Stage 2 のみ |
| `-j` / `--jobs` | `lake build -j` に渡す |

**成功時**: `OK: stage1 + stage2`、証明された定理名一覧（Tier 0+）。

**失敗時**: Stage 別にエラー（矛盾内容 / Lean 定理名 / ファイル行）。

---

## 6. コマンド体系図

```
novel-logic
├── init | info | doctor | template list | version
│
├── status | timeline
│
├── plot set | show
├── scene list | show | add | remove
├── novel list | show | add | update | remove
│   └── revision pin | list
├── thing list | show | add | update | remove
│   └── scope add | remove
├── time list | add | remove
├── fact list | show | add | update | promote | remove
├── action list | show | add | update | remove
├── rule list | show | add | update | remove
│
├── branch list | show | add | remove
├── fork add | choice add | show
├── merge add | show
│
├── validate [--branch]
├── generate
├── check [--branch | --quick]   … メイン（Stage 1 + 生成 + Stage 2）
│
└── wizard plot | novel      … Phase 1 補助（対話型）
```

---

## 7. 典型ワークフロー

### 桃太郎（新規）

```bash
novel-logic init ./momotaro --template momotaro
cd momotaro

# Phase A（テンプレで済んでいなければ）
novel-logic plot set --title 桃太郎 --summary "鬼退治の昔話"
novel-logic thing add momotaro --tag character
novel-logic time add 1
novel-logic scene add scene1 --summary "桃から誕生" --time-start 1 --time-end 3
novel-logic fact add --kind fixed --thing momotaro --pred 人間
novel-logic rule add --kind forbid-transition --from 青年 --to 赤ちゃん
novel-logic action add --thing momotaro --from 赤ちゃん --to 青年 --at 4

# Phase B
novel-logic novel add scene1
# エディタで novels/main/scene1.txt を編集
git add novels/main/scene1.txt && git commit -m "scene1 prose"
novel-logic novel revision pin scene1
novel-logic fact add --kind state --thing momotaro --pred 青年 --scope novel:scene5
novel-logic action add --thing momotaro --from 旅立ち --to 鬼退治済み \
  --at t12 --scope novel:scene5 --label 鬼が島で鬼を退治する

# 検証
novel-logic validate
novel-logic check
novel-logic timeline
```

### 執筆中

```bash
# 本文を直したあと
git commit -am "revise scene5"
novel-logic novel revision pin scene5

novel-logic fact add --kind state --thing momotaro --pred 青年 --scope novel:scene5
# → 同一キーの再 add は拒否。変更は fact update <id>

novel-logic validate
novel-logic check --quick
novel-logic check
```

### CI

**リポジトリ本体**（`novel-logic` 開発用）:

```bash
go test ./...
```

GitHub Actions（[`.github/workflows/test.yml`](../.github/workflows/test.yml)）が `push` / `pull_request` で上記を実行します。

**作品データ**（ユーザーが `init` したディレクトリ）:

```bash
novel-logic -C ./momotaro validate              # Stage 1（全 branch）
novel-logic -C ./momotaro validate --branch main
novel-logic -C ./momotaro check                 # Stage 1 + Lean 生成 + lake build
```

Lean toolchain が CI ランナーに無い場合は `validate` または `check --quick` を使います。

---

## 8. ウィザード（Phase 1・補助）

サブコマンドと **同一の登録・矛盾拒否ロジック**を呼ぶ対話型 UI。Phase 0 では未実装。

### `novel-logic wizard plot`

Phase A（A1–A5）を順に質問。完了後に YAML へ一括反映（各ステップで Stage 1）。

### `novel-logic wizard novel <scene_id>`

Phase B（B1–B4）を順に質問。

| フラグ | 説明 |
|--------|------|
| `--from <step>` | 途中ステップから再開（例: `scene`） |

---

## 9. Phase 1 以降（その他予約）

| コマンド | 概要 |
|----------|------|
| `novel-logic diff` | 前回成功 `check` からの作品データ / 生成 Lean 差分 |
| `novel-logic export constraints` | AI 執筆用の rule + fact + action テキスト出力 |
| `novel-logic import extract` | 本文から fact / action 候補を提案（承認フロー） |
| `novel-logic thing tag` | 既存 thing への tag 追加・削除 |
| `novel-logic novel add`（1:N） | Phase 1: 1 scene 複数 novel |

---

## 10. 実装設計の決定・未決定

### 10.1 決定済み

- [x] **永続化形式**: YAML ファイル群（§1 `init` の分割。詳細は [REQUIREMENTS §7.3](REQUIREMENTS.md)）
- [x] **CLI 体系**: サブコマンド型を正 + `novel-logic wizard` は Phase 1 補助（§8）
- [x] **バイナリ名**: `novel-logic`（`nl` エイリアスは採用しない）
- [x] **scene ↔ novel**: branch 内 1:1。本文は `novels/<branch>/<scene>.txt` + git。1:N は Phase 1 で拡張検討
- [x] **branch / fork / merge**: `branches.yaml` / `forks.yaml` / `merges.yaml` + `--branch` フラグ
- [x] **単体テスト**: `internal/cli`, `generate`, `project`, `validate`
- [x] **GitHub Actions**: `go test ./...`（Lean 不要）
- [x] **thing ID**: 作品内一意。`thing add` は新規のみ。スコープ追加は `thing scope add`
- [x] **add / update 分離**: fact / action / rule / thing の同一キー `add` は拒否
- [x] **novel revision pin**: git commit を `novels.yaml` に記録
- [x] **`check` デフォルト**: Stage 2 込み（`--quick` で Stage 1 のみ）
- [x] **Lean テンプレ**: `PredId` / `Scope` は `inductive` 生成。Phase 0 は述語形式のみ（[REQUIREMENTS §6.3.1](REQUIREMENTS.md)）

### 10.2 未決定・フェーズ予定

（[REQUIREMENTS §10.2](REQUIREMENTS.md)）

---

## 11. 改訂履歴

| 日付 | 内容 |
|------|------|
| 2026-06-20 | 初版（旧ドメイン: chapter / first-time / timeline） |
| 2026-06-20 | [REQUIREMENTS.md](REQUIREMENTS.md) / [DRAFT.txt](DRAFT.txt) に合わせ全面改訂（plot/scene/novel/thing/fact/action/rule/time） |
| 2026-06-20 | 永続化形式を YAML ファイル群に確定 |
| 2026-06-20 | CLI 体系を確定（サブコマンド正 + wizard Phase 1） |
| 2026-06-20 | バイナリ名を `novel-logic` に確定 |
| 2026-06-20 | scene ↔ novel を Phase 0 は 1:1 に確定 |
| 2026-06-20 | thing ID マージ規則を確定（同一 ID = 同一 thing） |
| 2026-06-20 | `check` デフォルトを Stage 2 込みに確定 |
| 2026-06-20 | Lean テンプレ詳細を確定（§6.3.1） |
| 2026-06-21 | novel: git 管理 `.txt` + `novel add` / `revision pin`。`novel set` 廃止 |
| 2026-06-21 | add/update 分離、thing scope add、`update` コマンド追加 |
| 2026-06-21 | branch/fork/merge、`validate`/`check --branch`、`timeline --branch` を反映 |
| 2026-06-21 | CI（GitHub Actions）、単体テスト、未実装フラグ（`--json`/`--dry-run`）を明記 |
| 2026-06-21 | 設計判断 #14/#16/#20: plot time 窓は hint のみ、`novel_extends_plot` は Phase 1、空 `from` と forbid-transition |