# novel-logic — 要件定義

> 最終更新: 2026-06-21  
> ステータス: **Phase 0（MVP）実装完了** — branch/fork/merge・単体テスト・CI 含む
> 正本: [DRAFT.txt](DRAFT.txt)

---

## 1. 背景・目的

### 1.1 きっかけ

短編『配信の部屋』において、プロット・各章本文を **Lean 4** で形式化し、時系列・状態・「初めて」描写・世界観ルールの整合を `lake build` で機械検証した。  
本プロジェクトはその知見を一般化し、**任意の小説**について論理構造を管理する CUI ツールを構築する。

### 1.2 本アプリの目的

**人間または AI が小説を執筆する際の、プロット・本文・設定・状態遷移の管理と矛盾検出**を行う CUI ツールを提供する。

- 小説本文そのもののリッチ編集機能は **スコープ外**（v0）
- 登録された事実・状態・行為・ルールが「正しい骨格」となり、執筆・レビュー・AI 生成の参照元になる
- 可視化（GUI）は **将来**。v0 は CUI のみ

### 1.3 非目的（v0）

- WYSIWYG 執筆エディタ
- 本文からの事実自動抽出（完全自動同期）
- 文学品質・文体の評価
- 配信・出版パイプライン

---

## 2. 確定した方針

| 項目 | 決定内容 |
|------|----------|
| **ドメインモデル** | [DRAFT.txt](DRAFT.txt) で定義（§3 要約） |
| **検証方式** | **A3. 二段階**（§6） |
| **実装言語（CLI）** | **Go** |
| **論理エンジン** | **Lean 4**（`elan` / `lake`）— Go 生成物として利用（§6.1–6.7） |
| **登録時の挙動** | 矛盾がある場合は **登録を拒否**し、矛盾内容をユーザーに応答 |
| **永続化形式** | **YAML ファイル群**（§7.3）。人間・AI が編集する正本 |
| **CLI 体系** | **サブコマンド型を正** + **対話ウィザードを補助**（§8.1）。ウィザードは Phase 1 |
| **バイナリ名** | **`novel-logic`**（`nl` エイリアスは採用しない） |
| **scene ↔ novel** | **Phase 0 は 1:1**（§3.6）。1:N は Phase 1 で拡張検討 |
| **thing ID** | **作品内で ID は一意**。plot / novel 同一 ID は **同一 thing**（§3.8） |
| **`check` デフォルト** | **Stage 1 + 生成 + Stage 2**（§6.6）。`--quick` で Stage 1 のみ |
| **Lean テンプレ** | `PredId` / `Scope` は作品ごと `inductive`（§6.3.1） |
| **作業場所** | `/home/zenak/novel-logic/`（本リポジトリ） |

---

## 3. 管理対象（ドメインモデル）

[DRAFT.txt](DRAFT.txt) に基づく。各エンティティは作品（プロジェクト）スコープで管理する。

| 概念 | 説明 |
|------|------|
| **plot** | 小説のプロット。scene の定義を含む |
| **scene** | plot 内のシーン。novel と紐づく。執筆の単位（cardinality は §3.6） |
| **novel** | 小説本文。特定の scene に紐づいて登録される |
| **thing** | もの・ことの総称。キャラクター・場所・概念等を同一型で扱う。**任意の tag** で種別を付与する（§3.4） |
| **fixed_fact** | 「thingA は thingB である」形式。**変わらない**事実 |
| **state** | 「thingA は thingB である」形式。**action によって遷移**する |
| **action** | 「thingA は thingB をした」「thingA は thingB になった」等の状態遷移。time と紐づく |
| **rule** | 禁則。「thingA は thingB ではない」「thingA は thingB にはならない」等 |
| **time** | 時系列上の位置。DATE 型でなくてよい。**順序性のある数値**で表現する |

### 3.1 エンティティ間の関係

```
plot
 └── scene ── novel
       │         │
       │         ├── thing（tag 付き）── fixed_fact / state
       │         │              └── action ── time
       │         └── time（開始・終了）
       │
       └── thing（plot 全体に共通、scene 非依存、tag 付き）
              ├── fixed_fact
              └── state ── action ── time
```

| 紐づけ | 内容 |
|--------|------|
| plot ↔ scene | plot は scene の集合として定義される |
| scene ↔ novel | scene に基づいて novel（本文）が執筆・登録される。**Phase 0 は 1:1**（§3.6） |
| scene ↔ time | 各 scene に **開始 time** と **終了 time** を設定する |
| novel ↔ time | novel 内にも **開始 time** と **終了 time** の設定が必要 |
| plot / novel ↔ thing | プロット全体または本文スコープの thing を登録できる |
| thing ↔ tag | thing 登録時に **任意の文字列タグ** を1つ以上付与できる（分類・検索用） |
| thing ↔ fixed_fact / state | thing に属性として紐づく。thing 間の関係もここで表現する |
| state ↔ action | action は state を変更する。action 登録時に time を指定する |
| action ↔ time | action がいつ行われたかを表す |

### 3.2 fixed_fact と state の扱い

| ルール | 内容 |
|--------|------|
| 構造の類似 | fixed_fact と state は同じ「A は B である」形式 |
| **昇格** | `novel-logic fact promote <id>` のみ（`fact update --kind state` は拒否） |
| **降格** | state を fixed_fact に変更することは **できない**（`fact update --kind fixed` も拒否） |

### 3.3 thing 間の関係の表現

**relation は管理対象に含めない。** thing 間の関係・属性・禁則は次で足りる。

| 表現したい内容 | 使う概念 |
|----------------|----------|
| 不変の属性・分類 | fixed_fact（例: 桃太郎は人間） |
| 変化する属性・状態 | state + action（例: 赤ちゃん ⇒ 青年） |
| 禁止される状態・遷移 | rule（例: 人間 ⇒ 動物にはならない） |

### 3.4 thing の tag

thing は概念が広いため、登場要素の**種別は tag で表現**する。tag は thing のメタデータであり、独立した論理エンティティではない。

| 項目 | 内容 |
|------|------|
| 付与数 | **1つ以上**（複数可） |
| 値 | ユーザーが任意に定義する文字列（例: `character`, `object`, `location`, `concept`, `item`） |
| 用途 | 一覧・フィルタ・執筆支援。論理矛盾の主判定は fixed_fact / state / rule が担う |
| 登録 | thing 登録（A2 / B2）時に指定。後から追加・変更可能とする |
| 推奨 | 作品内で tag 名を揃える（同一意味に `char` と `character` を混在させない等） |

tag 自体は作品ごとの語彙として自由。システムが `character` 等を**強制しない**（テンプレが候補を提示するのみ）。

### 3.6 scene ↔ novel の cardinality（確定）

| Phase | 方針 |
|-------|------|
| **0（MVP）** | **branch 内 1:1**。1 scene × 1 branch につき novel は1つ（`novels/<branch>/<scene_id>.txt`） |
| **1 以降** | **1:N 拡張を検討**（草稿・改訂版、`novel_id` 導入等） |

Phase 0 のルール:

- 本文は **`novels/<branch>/<scene_id>.txt` を git 管理**（CLI は散文を書き込まない。エディタで編集）。
- `novel-logic novel add <scene_id>` でメタ登録。同一 `scene_id` に2つ目の novel は作れない（更新は `novel update`）。
- `novel-logic novel revision pin` で本文ファイルの **git commit** を `novels.yaml` に記録（CI / PR 向け）。
- 改稿・履歴は **git** で管理する（ツール内バージョン管理は Phase 1）。
- Lean の `novel_extends_plot`（Tier 1）は **Phase 1 で実装**（§6.5）。Phase 0 は scene 1:1 前提のメタ整合と **hint**（`action.plot_scene_hint`）まで。

### 3.7 分岐・合流（branch / fork / merge）

物語のルート分岐を first-class エンティティとして管理する。

| 概念 | ファイル | 説明 |
|------|----------|------|
| branch | `branches.yaml` | ルート ID（`main` 含む）。子 branch は `parent` / `via_fork` / `via_action` |
| fork | `forks.yaml` | 分岐点（親 branch・`at`・choice 一覧） |
| merge | `merges.yaml` | 合流点（`into_branch`・各 from branch の merge action） |

**確定ルール**

- `fact` / `action` / `rule` / `novel` に `branch` フィールド（省略時 `main`）
- **time 軸は全 branch で共有**（`project.yaml` の `time_order` は 1 本）
- novel 本文は **`novels/<branch_id>/<scene_id>.txt`**（`main` も `novels/main/scene1.txt`）
- 合流は **merge action** を各 from branch に登録し、**同じ `to` pred** へ遷移させる
- merge 登録後、終了した子 branch への新規 action 登録は拒否

**検証（Stage 1）**: `branch.unknown`, `fork.invalid`, `fork.exclusive`, `merge.action_mismatch`, `merge.after_action`, `branch.isolated_state`, `novel.missing_body`, `time.registry_mismatch` 等。`validate` / `check` に `--branch <id>`（省略時は全 branch）。

### 3.8 time の扱い

- カレンダー日付である必要はない
- 整数等の **全順序可能な ID** でよい（例: 桃太郎例では `1` … `13`）
- scene の `[開始, 終了]` は time 軸上の区間として解釈する
- action は単一の time ポイントに紐づく

### 3.9 thing ID とスコープ（確定）

**同一 ID は常に同一 thing**（plot / novel で別レコードにしない）。

| 操作 | 挙動 |
|------|------|
| `thing add <id>` で **未登録 ID** | 新規 thing を作成し、`scopes` に指定スコープを追加 |
| `thing add <id>` で **登録済み ID** | **エラー**（`thing scope add` または `thing update` を案内） |
| `thing scope add <id>` | 既存 thing に `scopes` を追加（冪等） |
| `thing update <id>` | `name` / `tags` の変更 |

- `scopes` の値: `plot` / `novel:<scene_id>`
- 桃太郎を plot で登録後、novel で `momotaro` を指定 → **同一レコード**に `novel:scene5` が追加される
- novel のみ初出の ID（例: モブ）→ `scopes: [novel:scene3]` のみで新規作成
- fact / action の `--thing` は、この **グローバル thing ID** を参照する

---

## 4. 利用フロー（CUI）

登録は次の順序を基本とする。各ステップで **既存の rule / fixed_fact / state / action / time を勘案した矛盾チェック**を行い、矛盾時は登録を拒否する。

### Phase A — プロット構築

| # | 操作 | 内容 |
|---|------|------|
| A1 | plot 登録 | 小説のプロットを plot として登録する |
| A2 | thing 登録（plot スコープ） | scene に紐づかない、プロット全体共通の thing を登録（**tag 付き**） |
| A3 | scene 登録 | plot 内の scene を登録。開始 time・終了 time を設定 |
| A4 | fixed_fact / state 登録 | plot スコープの thing に紐づく fixed_fact・state を登録 |
| A5 | action 登録 | plot スコープの state に紐づく action を登録（time 付き） |

### Phase B — 本文登録

| # | 操作 | 内容 |
|---|------|------|
| B1 | novel 登録 | プロットに従い執筆した本文を、scene に紐づけて novel として登録 |
| B2 | thing 登録（novel スコープ） | novel 内の thing を登録（**tag 付き**） |
| B3 | fixed_fact / state 登録 | novel スコープの thing に紐づく fixed_fact・state を登録 |
| B4 | action 登録 | novel スコープの state に紐づく action を登録（time 付き） |

### 4.1 登録時の矛盾検出（必須）

**すべての登録操作**において、少なくとも次を検証する。

| 検証 | 例 |
|------|-----|
| rule 違反 | 「青年 ⇒ 赤ちゃん」action の登録拒否 |
| state 矛盾 | 「桃太郎は動物」state の登録拒否（rule と整合） |
| time 整合 | `novel:<scene>` スコープ action の `at` が当該 scene の `[time_start, time_end]` 外なら拒否。`plot` スコープは **窓チェック対象外**（Phase A はプロット全体の time 軸。scene 整合は hint のみ — §6 Stage 1） |
| 遷移不可能 | rule で禁止された state 遷移を伴う action |
| fact 矛盾 | 既存 fixed_fact・state・rule と両立しない新規登録 |

矛盾時の応答:

1. **何が矛盾しているか**を具体的に表示する
2. 当該登録を **実行しない**（ロールバック）

---

## 5. 参照例：桃太郎（DRAFT より）

[DRAFT.txt](DRAFT.txt) の具体例を要件の参照データとする。  
（DRAFT 内の scene 表記の誤記・time 区間の誤りは、要件例として解釈時に補正する。実装テストデータ作成時に正規化する。）

### 5.1 plot / scene

| scene | 概要 | time 区間（意図） |
|-------|------|-------------------|
| scene1 | 桃発見・桃太郎誕生 | 1–3 |
| scene2 | 桃太郎を育てる | 3–5 |
| scene3 | 鬼退治へ出発 | 5–7 |
| scene4 | 犬・猿・雉を仲間にする | 7–11 |
| scene5 | 鬼退治・ハッピーエンド | 11–13 |

### 5.2 thing（抜粋）

| thing | tag（例） |
|-------|-----------|
| 桃太郎、おじいさん、おばあさん、犬、猿、雉、鬼、男の子 | `character` |
| 山、川、鬼が島、村 | `location` |
| 桃、黍団子、宝物 | `object` |

### 5.3 論理データ（抜粋）

**fixed_fact**

- 桃太郎は人間 / おじいさんは人間 / おばあさんは人間
- 犬は動物 / 猿は動物 / 雉は動物

**state**

- 桃太郎は赤ちゃん
- 桃太郎は青年

**action**

- 桃太郎は育った（赤ちゃん ⇒ 青年）@ time4

**rule**

- 青年 ⇒ 赤ちゃんにはならない
- 人間 ⇒ 動物にはならない
- 桃太郎は動物ではない

---

## 6. 二段階検証（A3）

登録時矛盾チェック（§4.1）を二層で実装する。

### Stage 1: Go 内蔵チェック（常時・高速）

Lean 未インストールでも実行可能。

| カテゴリ | チェック例 |
|----------|-----------|
| スキーマ | 必須フィールド、参照 ID の存在、thing の tag 形式（空文字・重複 tag の扱い） |
| time | 順序性、`time_order` と `times.yaml` の整合（欠落・重複・空 ID）、`time add` 重複は登録拒否（exit 1）、scene/novel の区間包含、`novel:<scene>` action の time 窓（`time.action_window`）。`plot` スコープ action は窓エラーにしない |
| rule | 明示 rule との抵触（登録拒否条件）。`forbid-transition` は action の `from` が **非空のときのみ**照合（`from` 省略＝初期状態遷移。Lean `actionRespectsRules` と同型） |
| state | 禁止 state の登録、fixed_fact ↔ state 昇格規則 |
| 重複 | 同一スコープでの一意性違反 |

### Stage 2: Lean 厳密チェック（要 toolchain）

Stage 1 通過後、論理モデルを Lean に生成し `lake build` で検証。

| カテゴリ | チェック例 |
|----------|-----------|
| 遷移閉包 | 許可されない state 遷移の導出不能 |
| rule 定理化 | rule から導かれる禁止事項の一貫性 |
| plot ↔ novel | **Phase 1（Tier 1）**: `novel_extends_plot` 等。Phase 0 の Stage 2 は Tier 0 のみ |
| 時系列 | 全 action の time 順と state の単調整合 |

**Lean が無い環境**: Stage 1 のみ実行し、Stage 2 スキップを明示（警告表示）。

### 6.1 Lean の役割分担

| 層 | 担当 | 理由 |
|----|------|------|
| **Go（Stage 1）** | スキーマ、参照存在、登録拒否、tag 形式 | 高速・常時・Lean 不要 |
| **Lean（Stage 2）** | 時系列上の状態整合、rule の閉包、plot ↔ novel 横断（**Tier 1 は Phase 1**。Phase 0 は Tier 0 のみ） | 型と定理で整合性を保証 |
| **Lean 対象外** | tag 分類、本文の文学品質、自然言語パース | メタデータ／非形式データ |

Lean の価値は、**登録時に見逃した矛盾が全体を組み立てたときに定理として破綻する**ことを `lake build` で検出できる点にある（『配信の部屋』プロトタイプと同型）。

人間・AI が編集する正本は作品データ（YAML 等）。**Lean は Go が生成する検証用コード**とし、手編集は非推奨。

### 6.2 生成物の構成

作品ディレクトリ内の `logic/` に、次を生成する。

```
logic/
  Core.lean          # 汎用ランタイム（全作品共通・novel-logic リポジトリのテンプレ）
  Project.lean       # 当該作品の ThingId / TimeId / BranchId / SceneId / PredId（生成）
  Facts.lean         # fixed_fact / state / action / activeActions_<branch>（生成）
  Rules.lean         # rule の列（生成）
  Timeline.lean      # scene / novel の time 区間（生成）
  Theorems.lean      # 整合性定理（生成）
  lakefile.toml
```

- **Core.lean** は `novel-logic` 本体に1つ。mathlib なしの軽量構成を維持する。
- 作品ごとに変わるのは **ID 列とデータ** のみ。
- **novel（本文テキスト）** は Lean に入れない。登録済み fact / action の参照整合のみ検証する。

### 6.3 ドメイン → Lean の写像

#### 識別子

作品内の thing / time / scene は、生成される `inductive` または有限列挙 ID とする。

```lean
-- 生成例（桃太郎）
inductive ThingId | momotaro | ojiisan | inu | ...
inductive TimeId  | t1 | t2 | ... | t13
inductive SceneId | scene1 | scene2 | ...
```

**tag は Lean に送らない**（§3.4）。分類・フィルタは Go 側のみ。

#### fixed_fact / state / action

「thingA は thingB である」を **Predicate** に正規化してから生成する。

| 概念 | Lean 構造 |
|------|-----------|
| fixed_fact | `structure FixedFact`（`subject : ThingId`, `pred : PredId`, `scope : Scope`） |
| state | `structure StateDecl`（同上） |
| action | `structure Action`（`subject`, `from? : Option PredId`, `to : PredId`, `at : TimeId`, `scope`） |

#### 6.3.1 実装詳細（確定）

| 論点 | 決定 |
|------|------|
| **PredId** | 作品ごとに **`inductive PredId` を生成**（`facts` / `rules` / `actions` から出現述語を収集） |
| **Scope** | **`inductive Scope \| plot \| novel (sceneId : SceneId)`** — 文字列パースではなく型で scene 参照を保証 |
| **Phase 0 の文形式** | **「A は B（述語）」のみ**。`--pred` は常に `PredId` に写像する文字列 |
| **B が thing 名の関係** | Phase 0 **非対応**（例: 「犬の主人は桃太郎」）。Phase 1 で `object : ThingId` 付き fact を検討 |
| **述語と thing ID の衝突** | Stage 1 で検出: `--pred` が既存 `ThingId` と同名なら **登録拒否**（別名を促す） |

生成例（`Project.lean`）:

```lean
inductive PredId | 人間 | 青年 | 赤ちゃん | 動物 | ...
inductive Scope | plot | novel (sceneId : SceneId)
```

桃太郎例:

- `FixedFact ⟨momotaro, 人間, Plot⟩`
- `Action ⟨momotaro, some 赤ちゃん, 青年, t4, Plot⟩`

#### rule

| rule の種類 | Lean での表現 | 桃太郎例 |
|-------------|---------------|----------|
| **静的禁止** | `forbiddenStates : Set (ThingId × PredId)` | 桃太郎 × 動物 |
| **遷移禁止** | `forbiddenTransitions : Set (PredId × PredId)` | 青年→赤ちゃん、人間→動物 |

#### time / scene

```lean
structure SceneWindow where
  scene : SceneId
  start stop : TimeId

def timeOrder : List TimeId   -- [t1, t2, ..., t13]
def timeLe : TimeId → TimeId → Bool  -- order から導出
```

`time` が日付でない要件は、`TimeId` の全順序リストで満たす。

### 6.4 状態シミュレーション（evolve）

Lean の核は、**time 順に action を適用した全体の状態展開**である。Stage 1 は1件ずつの登録拒否、Stage 2 は履歴全体を見る。

```lean
/-- time 順に action を適用し、各 (thing, time) の述語集合を構築 -/
def evolve (facts : List FixedFact) (actions : List Action) (t : TimeId) :
    ThingId → Set PredId
```

| 入力 | evolve への扱い |
|------|----------------|
| fixed_fact | 全 time で常に真として初期注入 |
| action | 該当 `at` で subject の述語を差し替え |
| rule | `evolve` の結果が `forbiddenStates` / `forbiddenTransitions` を満たすことを定理化 |

同一 subject の排他 state は、生成時に正規化ルールを適用する。fixed_fact → state **昇格**後は、fixed_fact 述語が action で破られないことを別定理で検証する（§3.2 降格不可の補強）。

### 6.5 証明すべき定理（Tier）

#### Tier 0 — Phase 0（MVP・桃太郎 end-to-end）

| 定理 | 内容 |
|------|------|
| `actions_in_scene_window` | 各 action の `at` が所属 scope の time 区間内（`plot` は `scopeToScene => none` で常に OK。`novel:<scene>` のみ scene 窓を検証 — §10.1 #8） |
| `no_forbidden_states` | 登録 state が `forbiddenStates` に含まれない |
| `no_forbidden_transitions` | 各 action の遷移が `forbiddenTransitions` に含まれない |
| `fixed_facts_stable` | fixed_fact 述語が後続 action で破られない |
| `actions_in_scene_window_<branch>` | branch ごとの active action について上記と同型 |
| `no_forbidden_transitions_<branch>` | branch ごとの active action が rule を尊重 |
| `fixed_facts_stable_<branch>` | branch ごとの active action で fixed_fact が破られない |

`activeActions_<branch>` は fork / merge を反映した「その branch で有効な action 列」。証明スタイルは有限列挙の `native_decide` を基本とする（配信の部屋プロトタイプ踏襲）。

#### Tier 1 — Phase 1

| 定理 | 内容 |
|------|------|
| `novel_extends_plot` | novel スコープの fact / action が同一 scene の plot 側と矛盾しない |
| `state_at_time_consistent` | time 順 evolve の結果と登録 state が一致 |
| `scene_time_monotone` | scene の `[start, stop]` が time 軸と整合 |

#### Tier 2 — Phase 1 以降（配信の部屋再モデル化）

| 旧プロトタイプ | 新モデル + Lean |
|----------------|-----------------|
| 視聴者数単調増加 | thing + state（視聴者数レンジ）+ action + 単調性定理 |
| ディルド挿入章制限 | rule: 特定 state は scene 以降のみ |
| FirstTime 二重登録防止 | 特定 `(thing, pred)` への初回 action を1回に制限 |
| WorldConstraints | fixed_fact + rule の集合 |

### 6.6 登録フローとの連携

```
ユーザーが action 等を登録
    → Go Stage 1: 即時拒否（明示 rule 抵触等）or 受理・永続化
    → check 時: Go が logic/ を再生成 → lake build
        → 成功: Tier 0+ 定理がすべて証明済み
        → 失敗: 定理名・Lean エラー行を CLI が表示
```

| タイミング | 検証 |
|------------|------|
| **各登録時** | Stage 1 のみ（高速） |
| **`check`（デフォルト）** | Stage 1 → 生成 → Stage 2（全体整合） |
| **`check --quick`** | Stage 1 のみ（Lean 不要） |
| **CI** | `novel-logic check`（デフォルト＝Stage 2 込み） |

Lean 未検出時: `check` は Stage 1 完了後に **警告を出して Stage 2 をスキップ**（終了コード 5 または 1—実装で定義）。`--quick` なら Lean 不要で完結。

定理が増えても **人間は作品データだけを編集**する。Lean は常に再生成する。

### 6.7 桃太郎での Theorems.lean 骨格（生成例）

```lean
namespace Momotaro

theorem rules_consistent :
  ∀ a ∈ allActions, actionRespectsRules momotaroRules a := by native_decide

theorem momotaro_not_animal_at_end :
  ¬ (PredId.動物 ∈ evolve allFacts allActions t13 ThingId.momotaro) := by
  native_decide

theorem no_reverse_growth :
  ¬ ∃ a ∈ allActions, a.subject = ThingId.momotaro ∧
    a.from? = some PredId.青年 ∧ a.to = PredId.赤ちゃん := by native_decide

end Momotaro
```

`lake build` 成功 = 上記を含む Tier 0 定理がすべて証明済み、と運用する。

---

## 7. アーキテクチャ概要

```
┌─────────────────────────────────────────────────────────┐
│  人間 / AI                                               │
│  CUI で plot / scene / novel / thing / fact / action 登録 │
└──────────────────────────┬──────────────────────────────┘
                           │ 各登録時に矛盾チェック
                           ▼
┌─────────────────────────────────────────────────────────┐
│  novel-logic CLI (Go)                                    │
│  · エンティティ CRUD（登録フロー §4）                      │
│  · Stage 1 矛盾検出 → 拒否 or 受理                         │
│  · Lean コード生成                                        │
│  · lake build（Stage 2）                                  │
└──────────────────────────┬──────────────────────────────┘
                           │
              ┌────────────┴────────────┐
              ▼                         ▼
     ┌─────────────────┐      ┌─────────────────┐
     │  Stage 1        │      │  Stage 2        │
     │  Go 内蔵        │      │  Lean 4         │
     └─────────────────┘      └─────────────────┘
```

### 7.1 作品データの置き場所

- **CLI 本体**: `novel-logic/`（本リポジトリ）
- **作品データ**: ユーザーが作成する別ディレクトリ（init 相当の操作で生成）

ドメインエンティティ名（plot, scene, novel, thing 等）は DRAFT に従い固定する。  
作品データの正本は **YAML ファイル群**（§7.3）。Lean は `logic/` のみ（生成物）。

### 7.2 作品ディレクトリ

```
my-work/
  project.yaml          # 作品メタ・time 順序
  plot.yaml
  things.yaml
  scenes.yaml
  times.yaml
  branches.yaml         # story branch（`main` 含む）
  forks.yaml            # 分岐点
  merges.yaml           # 合流点
  facts.yaml            # fixed_fact + state（`branch` 省略時 `main`）
  actions.yaml
  rules.yaml
  novels.yaml           # novel メタ（scene × branch ごと1件）
  novels/               # 本文
    main/
      scene1.txt
    branch_dog/         # 分岐ルート例
      scene4.txt
  logic/                # §6.2 の Lean 生成物（手編集非推奨）
    Core.lean
    Project.lean
    Facts.lean
    Rules.lean
    Timeline.lean
    Theorems.lean
    lakefile.toml
```

### 7.3 YAML 永続化（確定）

| ファイル | 内容 |
|----------|------|
| `project.yaml` | タイトル、作成日、最終 check 結果、`time_order`（TimeId の全順序） |
| `plot.yaml` | plot 要約 |
| `things.yaml` | thing 定義（`id` **作品内一意**, `name`, `tags[]`, `scopes[]`。同一 ID の重複行は禁止） |
| `scenes.yaml` | scene 定義（`id`, `summary`, `time_start`, `time_end`） |
| `times.yaml` | time メタ（`id` の存在確認用。順序は `project.yaml` の `time_order` が正） |
| `branches.yaml` | branch 定義（`id`, `label`, `parent`, `via_fork`, `via_action`） |
| `forks.yaml` | fork 定義（`parent_branch`, `at`, `scope`, `choices[]`） |
| `merges.yaml` | merge 定義（`at`, `into_branch`, `scope`, `choices[]`） |
| `facts.yaml` | fixed_fact / state（`kind`, `thing`, `pred`, `scope`, `branch`） |
| `actions.yaml` | action（`thing`, `from`, `to`, `at`, `scope`, `label`, `branch`） |
| `rules.yaml` | rule（`kind`, `thing`/`pred`/`from`/`to`, `branch`） |
| `novels.yaml` | novel メタ（`scene_id` + `branch` で一意。`time_start`, `time_end`, `body_path`, `revision`, `revisions[]`） |
| `novels/<branch>/<scene_id>.txt` | 本文（branch 内 1:1。**git 管理**。論理検証の対象外。CLI は書き込まない） |

**編集ルール**

- CUI 登録コマンドは Stage 1 通過後に上記 YAML を更新する（失敗時は書き込まない）。
- 同一キーへの `add` は拒否し、変更は `update` サブコマンドで行う（fact / action / rule / thing）。
- 本文（`novels/<branch>/*.txt`）は **YAML に含めない**。エディタで編集し git で管理する。CLI はメタと revision pin のみ。
- レガシー `novels/<scene_id>.txt` は Load 時に `novels/main/<scene_id>.txt` へ正規化する。
- 人手での直接編集を許容する。`novel-logic validate` / `novel-logic check` で整合を確認する。
- `logic/` は `novel-logic generate` の出力のみ。正本ではない。
- SQLite 等への移行は Phase 1 以降に再検討可（現時点では採用しない）。

---

## 8. CLI コマンド（方針）

[DRAFT.txt](DRAFT.txt) は **ドメインと利用フロー**を定義する。  
CLI は §4 の登録フローに対応する操作を提供する。詳細は [COMMANDS.md](COMMANDS.md)。

### 8.1 CLI 体系（確定）

| 層 | 内容 | Phase |
|----|------|-------|
| **正（メイン）** | サブコマンド型（`novel-logic thing add`, `novel-logic fact add` 等） | 0 |
| **補助** | 対話型ウィザード（`novel-logic wizard`）— 同じ登録を質問形式で実行 | 1 |

- Phase 0 の MVP・CI・テストは **サブコマンドのみ**で完結させる。
- ウィザードはサブコマンドの薄いラッパとし、登録・矛盾拒否のロジックは共有する。
- 人手・AI はサブコマンドまたは YAML 直接編集を使う（ウィザード必須ではない）。

| 操作群 | 対応フロー | 備考 |
|--------|-----------|------|
| plot 操作 | A1 | 作成・表示・更新 |
| thing 操作 | A2, B2 | スコープ（plot / novel）指定、**tag 付与・一覧フィルタ** |
| scene 操作 | A3 | time 区間付き |
| fact 操作 | A4, B3 | fixed_fact / state の登録・昇格 |
| action 操作 | A5, B4 | time 付き。state 変更 |
| rule 操作 | — | rule の登録・一覧 |
| novel 操作 | B1 | scene 紐づけ本文登録 |
| wizard | A / B 全体 | Phase 1。`novel-logic wizard plot` / `novel-logic wizard novel` 等 |
| 検証 | §6 | 登録時 Stage 1 / `check` 時 Stage 1 + 2（§6.6） |

---

## 9. フェーズ計画

### Phase 0（MVP）

- [x] Go CLI 骨格
- [x] ドメインモデル（§3）の YAML スキーマ（§7.3）
- [x] 登録フロー §4 の CRUD + **登録時矛盾拒否**（Stage 1）
- [x] Stage 1 検証（rule / time / state / branch / fork / merge）
- [x] `novel-logic` 本体に **Core.lean** テンプレ（mathlib なし）
- [x] Go コード生成器（§6.2: Project / Facts / Rules / Timeline / Theorems）
- [x] Stage 2: Tier 0 定理 + `lake build`（桃太郎例で end-to-end）
- [x] branch / fork / merge + branch 別 Lean 定理（`activeActions_*`）
- [x] `examples/momotaro/` / `examples/momotaro-walkthrough/`（`branch_dog` 分岐デモ）
- [x] 単体テスト（`internal/cli`, `generate`, `project`, `validate`）
- [x] GitHub Actions CI（`go test ./...`）

### Phase 1

- [ ] **`novel-logic wizard`**（対話型。サブコマンドと登録ロジック共有）
- [ ] **evolve** 一般化 + Tier 1 定理（`novel_extends_plot` 等）
- [ ] plot ↔ novel 横断整合の強化
- [ ] 差分表示・履歴
- [ ] 配信の部屋型テンプレ + Tier 2 定理（旧プロトタイプ再モデル化）

### Phase 2

- [ ] AI 向け制約エクスポート
- [ ] 本文からの事実抽出 assist
- [ ] GUI

---

## 10. 実装設計の決定・未決定

### 10.1 決定済み

| # | 論点 | 決定内容 |
|---|------|----------|
| 1 | **永続化形式** | **YAML ファイル群**（§7.2–7.3）。SQLite / ハイブリッドは Phase 0 では不採用 |
| 2 | **CLI コマンド体系** | **サブコマンド型を正** + **ウィザード補助**（§8.1）。ウィザードは Phase 1 |
| 3 | **バイナリ名** | **`novel-logic`**（`nl` エイリアスは採用しない） |
| 4 | **scene ↔ novel** | **Phase 0 は 1:1**、1:N は Phase 1 で拡張検討（§3.6） |
| 5 | **plot / novel 同一 thing ID** | **同一 ID = 同一 thing**。スコープ追加のみ（§3.8） |
| 6 | **`check` デフォルト** | **Stage 2 込み**（`validate` → `generate` → `lake build`）。`--quick` で Stage 1 のみ |
| 7 | **Lean テンプレ詳細** | `PredId` / `Scope` は `inductive` 生成。Phase 0 は述語形式のみ（§6.3.1） |
| 8 | **plot スコープ action の time 窓** | Stage 1 は **`novel:<scene>` のみ厳格**。`plot` は hint（`action.plot_scene_hint`）のみ（§6 Stage 1）。Lean も `scopeToScene plot => none` で同型 |
| 9 | **`forbid-transition` と空 `from`** | action の `from` 省略時は遷移禁止 rule を **照合しない**（初期状態への遷移）。`forbid-state` は `to` で照合。Lean と一致 |
| 10 | **`novel_extends_plot` のフェーズ** | **Phase 1 / Tier 1**（§6.5）。Phase 0 では未実装。§3.6 の 1:1 前提はメタ・hint 用 |

### 10.2 未決定・フェーズ予定

| # | 論点 | 備考 |
|---|------|------|
| 11 | **旧モデル再モデル化** | Phase 1 予定（§6.5 Tier 2 / §11）。時期の追加調整のみ |
| 12 | **thing 参照付き fact** | Phase 1（§6.3.1「B が thing 名」） |

---

## 11. 旧プロトタイプとの対応（参考）

『配信の部屋』Lean 資産は Phase 1 で新ドメインへ再モデル化する。

| 旧概念 | 新モデルでの位置づけ（案） |
|--------|---------------------------|
| Chapter / 章 | scene |
| 本文 md | novel |
| ChapterState | thing + state の複合 |
| SceneRecord | action + time |
| FirstTime | 特定 state 遷移の action（初回性は rule またはタグで表現） |
| WorldConstraints | rule + fixed_fact |
| StoryDay | time（順序 ID） |

| 旧プロトタイプ（Lean） | 新モデルでの定理化（§6.5 Tier 2） |
|------------------------|----------------------------------|
| `ChapterState.allValid` | `no_forbidden_states` + `no_forbidden_transitions` |
| `FirstTime.uniqueKinds` | 初回 action 一意性 rule + 定理 |
| 視聴者数単調 | evolve + 単調性定理 |
| `WorldConstraints` | `Rules.lean` + fixed_fact 初期注入 |

参照パス: `/home/zenak/配信の部屋/Plotlogic/`

---

## 12. 用語集

| 用語 | 意味 |
|------|------|
| **plot** | プロット。scene の容器 |
| **scene** | プロット内の場面単位。novel と対応 |
| **novel** | 本文。scene に紐づく |
| **thing** | 登場要素の総称。tag で種別付け |
| **tag** | thing に付与する任意の分類ラベル（`character`, `object` 等） |
| **fixed_fact** | 不変の「A は B」 |
| **state** | 可変の「A は B」。action で遷移 |
| **action** | 状態遷移・行為。time 付き |
| **rule** | 禁止・不変制約 |
| **time** | 順序 ID（日付でなくてよい） |
| **branch** | 物語ルート ID（`main` が本線。fork で子 branch が生える） |
| **fork** | 分岐点。親 branch 上の action から子 branch へ |
| **merge** | 合流点。複数 branch の merge action を同一 `to` pred へ揃える |
| **PredId** | 作品内述語 ID（`人間`, `青年` 等）。Lean 生成時に列挙化 |
| **Scope** | fact / action の所属（`plot` \| `novel sceneId`。Lean では `inductive`） |
| **evolve** | time 順に action を適用して得る述語集合（§6.4） |

---

## 13. 改訂履歴

| 日付 | 内容 |
|------|------|
| 2026-06-20 | 初版。A3・Go+Lean 確定 |
| 2026-06-20 | [DRAFT.txt](DRAFT.txt) に基づき全面再定義。ドメインモデル・利用フロー・桃太郎例を反映 |
| 2026-06-20 | relation を削除。thing 間の関係は fixed_fact / state / rule で表現 |
| 2026-06-20 | thing に任意 tag（character, object 等）を付与可能に |
| 2026-06-20 | §6.1–6.7 Lean 4 活用方針を追加（生成構成・写像・定理 Tier・evolve） |
| 2026-06-20 | 永続化形式を YAML ファイル群に確定（§7.3） |
| 2026-06-20 | CLI 体系を確定（サブコマンド正 + ウィザード補助・Phase 1） |
| 2026-06-20 | バイナリ名を `novel-logic` に確定 |
| 2026-06-20 | scene ↔ novel を Phase 0 は 1:1 に確定（§3.6） |
| 2026-06-20 | thing ID を作品内一意に確定（§3.8）。スコープ追加は `thing scope add` |
| 2026-06-20 | `check` デフォルトを Stage 2 込みに確定（`--quick` で Stage 1 のみ） |
| 2026-06-20 | Lean テンプレ詳細を確定（§6.3.1: PredId / Scope inductive、Phase 0 は述語のみ） |
| 2026-06-21 | novel 本文を `novels/*.txt` + git 管理に確定。`novel add` / `revision pin`。`novel set` 廃止 |
| 2026-06-21 | add/update 分離（重複 `add` 拒否）、`thing scope add`、各 `update` コマンド |
| 2026-06-21 | Phase 0 完了を反映（branch/fork/merge、Lean branch 定理、テスト、CI） |
| 2026-06-21 | §7.2–7.3 を `novels/<branch>/` 構成に更新 |
| 2026-06-21 | 設計判断 #14/#16/#20 確定（plot time 窓は hint のみ、`novel_extends_plot` は Phase 1、空 `from` と forbid-transition） |