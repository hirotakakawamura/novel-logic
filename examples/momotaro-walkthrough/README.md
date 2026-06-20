# 桃太郎ウォークスルー — novel-logic 入門

このディレクトリは、**空のプロジェクトから桃太郎の物語を登録していく**手順を実際に辿ったサンプルです。  
コマンドをそのままコピーして試せます。

関連ドキュメント:

- [要件定義](../../docs/REQUIREMENTS.md)
- [コマンド仕様](../../docs/COMMANDS.md)

---

## novel-logic とは？（30 秒で理解）

novel-logic は **小説の「論理の骨格」** を管理する CUI ツールです。

| 管理するもの | 例 |
|-------------|-----|
| 登場人物・物・場所（thing） | 桃太郎、犬、鬼が島 |
| 不変の事実（fixed_fact） | 桃太郎は人間 |
| 変化しうる状態（state） | 桃太郎は赤ちゃん / 青年 |
| 状態の変化（action） | t4 に赤ちゃん → 青年 |
| 禁則（rule） | 青年 → 赤ちゃんには戻れない |
| 時系列（time） | t1, t2, … t13 |
| プロットの区切り（scene） | scene1 = t1〜t3 |
| 本文（novel） | scene ごとのテキスト（`novels/<branch>/<scene_id>.txt`、git 管理。本線は `novels/main/`） |

**本文そのものの文学品質は検証しません。** 登録した fact / action / rule が矛盾していないかを機械的にチェックします。

---

## 事前準備

### ビルド（初回のみ）

```bash
export PATH=$HOME/.local/go/bin:$PATH   # Go をローカル導入している場合
cd /path/to/novel-logic
go build -o bin/novel-logic ./cmd/novel-logic
```

以降、CLI を `novel-logic` と表記します。`bin/` から実行するときは `./novel-logic`、PATH に通したあとは `novel-logic` だけで動きます。

### Lean について

- `novel-logic validate` や `check --quick` → **Lean 不要**（Stage 1 のみ）
- `novel-logic check`（デフォルト）→ **Lean + lake が必要**（Stage 2 まで実行）

Lean が無い環境では Stage 1 までで骨格の確認ができます。

---

## 登録の全体像

物語登録は大きく **2 フェーズ** に分かれます。

```
Phase A（プロット）  骨格を固める
  time → scene → thing → fact → rule → action

Phase B（本文）      scene ごとに本文を書く
  novel add → エディタで .txt 編集 → （任意）revision pin
  fact / action を --scope novel:<scene> で登録

最後                 矛盾がないか検証
  check
```

**ポイント**: 各登録コマンドは保存前に矛盾チェック（Stage 1）を行い、ルール違反があれば **登録を拒否** します。

---

## Step 0: プロジェクトを作る

```bash
novel-logic init ./my-momotaro --template default
cd my-momotaro
```

`project.yaml`, `things.yaml`, `scenes.yaml` などの空ファイルと `novels/` フォルダができます。

> **ヒント**: 手順を省略してサンプルデータだけ欲しい場合は  
> `novel-logic init ./momotaro --template momotaro` で桃太郎一式が入った状態から始められます。  
> この README は **学習用に空から積み上げる** 流れを説明します。

---

## Step 1: time（時系列の目盛り）を登録する

time はカレンダーの日付ではなく、**順序だけが重要な ID** です。  
scene の区間や action の「いつ」を表す目盛りになります。

桃太郎では t1〜t13 の 13 点を使います。

```bash
novel-logic time add t1
novel-logic time add t2
# …
novel-logic time add t13
```

登録した順序は `project.yaml` の `time_order` に記録されます。

```yaml
time_order:
  - t1
  - t2
  # ...
  - t13
```

---

## Step 2: scene（プロットの区切り）を登録する

各 scene には **開始 time** と **終了 time** を付けます。

| scene | time 区間 | 概要 |
|-------|-----------|------|
| scene1 | t1 – t3 | 桃発見・桃太郎誕生 |
| scene2 | t3 – t5 | 桃太郎を育てる |
| scene3 | t5 – t7 | 鬼退治へ出発 |
| scene4 | t7 – t11 | 犬・猿・雉を仲間にする |
| scene5 | t11 – t13 | 鬼退治・ハッピーエンド |

```bash
novel-logic scene add scene1 \
  --summary "桃発見・桃太郎誕生" \
  --time-start t1 --time-end t3

novel-logic scene add scene2 \
  --summary "桃太郎を育てる" \
  --time-start t3 --time-end t5

# scene3〜scene5 も同様
```

確認:

```bash
novel-logic plot show
```

---

## Step 3: thing（登場するもの）を登録する

キャラクター・場所・物を **ID + tag** で登録します。

- **ID**: 英数字（`momotaro`, `inu` など）。論理データの参照に使う
- **tag**: 分類用（`character`, `location`, `object`）。検証には使わない
- **name**: 人間向けの表示名（日本語可）

```bash
novel-logic thing add momotaro --name 桃太郎 --tag character
novel-logic thing add ojiisan  --name おじいさん --tag character
novel-logic thing add inu      --name 犬 --tag character
novel-logic thing add onigashima --name 鬼が島 --tag location
novel-logic thing add kibidango  --name 黍団子 --tag object
# …
```

**複数スコープ**（plot と複数 scene に登場）を一度に指定できます。

```bash
# 新規: plot と scene1 の両方
novel-logic thing add ojiisan --name おじいさん --tag character \
  --scope plot --scope novel:scene1

# 既存 thing にスコープを追加
novel-logic thing scope add momotaro --scope novel:scene1 --scope novel:scene5
novel-logic thing scope add momotaro --scope novel:scene2

# 既存 ID への thing add はエラー（スコープ追加は thing scope add を使う）
```

`things.yaml` の `scopes` は配列です。同一 ID の重複行は作りません。

---

## Step 4: fact（事実と状態）を登録する

### fixed_fact — 変わらない属性

「A は B である」のうち、**後から変化しない** ものです。

```bash
novel-logic fact add --kind fixed --thing momotaro --pred 人間
novel-logic fact add --kind fixed --thing inu       --pred 動物
```

### state — 変化しうる状態

「A は B である」のうち、**action で遷移しうる** ものです。  
先に「ありうる状態」を宣言しておき、後から action で実際の遷移を記録します。

```bash
novel-logic fact add --kind state --thing momotaro --pred 赤ちゃん
novel-logic fact add --kind state --thing momotaro --pred 青年
```

### `--pred` の注意

`--pred` は述語（状態の名前）です。**既存の thing ID と同名にはできません**（例: thing `village` があるとき `--pred village` は拒否）。

---

## Step 5: rule（禁則）を登録する

物語のルールを先に宣言しておくと、矛盾する action の登録を防げます。

```bash
# 青年から赤ちゃんへの逆戻りは禁止
novel-logic rule add --kind forbid-transition --from 青年 --to 赤ちゃん

# 人間が動物になる遷移は禁止
novel-logic rule add --kind forbid-transition --from 人間 --to 動物

# 桃太郎が動物である状態は禁止
novel-logic rule add --kind forbid-state --thing momotaro --pred 動物
```

---

## Step 6: action（状態遷移）を登録する

**「誰が・いつ・何から何へ変わったか」** を記録します。  
遷移元・遷移先の述語は、先に state として宣言しておく必要があります。

### 桃太郎の成長（scene2, t4）

```bash
novel-logic fact add --kind state --thing momotaro --pred 赤ちゃん
novel-logic fact add --kind state --thing momotaro --pred 青年

novel-logic action add \
  --thing momotaro --from 赤ちゃん --to 青年 \
  --at t4 --label 育った
```

### 旅立ち（scene3, t6）

```bash
novel-logic fact add --kind state --thing momotaro --pred 村在住
novel-logic fact add --kind state --thing momotaro --pred 旅立ち

novel-logic action add \
  --thing momotaro --from 村在住 --to 旅立ち \
  --at t6 --label 鬼退治へ出発
```

### 仲間入り（scene4, t8–t10）

犬・猿・雉それぞれに state を宣言し、time をずらして順序を表現します。

```bash
# 犬
novel-logic fact add --kind state --thing inu --pred 野良
novel-logic fact add --kind state --thing inu --pred 仲間
novel-logic action add --thing inu --from 野良 --to 仲間 \
  --at t8 --label 黍団子でもらい仲間になる

# 猿（t9）、雉（t10）も同様
```

### 鬼退治（scene5, t12）

鬼側と桃太郎側、両方の状態を更新します。

```bash
# 鬼
novel-logic fact add --kind state --thing oni --pred 健在
novel-logic fact add --kind state --thing oni --pred 退治済み
novel-logic action add --thing oni --from 健在 --to 退治済み \
  --at t12 --label 桃太郎一行に退治される

# 桃太郎
novel-logic fact add --kind state --thing momotaro --pred 鬼退治済み
novel-logic action add --thing momotaro --from 旅立ち --to 鬼退治済み \
  --at t12 --label 鬼が島で鬼を退治する
```

### 登録拒否の例

禁止ルールに抵触する action は **保存されません**。

```bash
novel-logic action add \
  --thing momotaro --from 青年 --to 赤ちゃん --at t6 --label 若返り
# → Error: forbidden transition "青年" -> "赤ちゃん"
```

---

## Step 6b: 分岐デモ（犬ルート `branch_dog`）

本ウォークスルーには **t7 で犬だけを仲間にする分岐** が含まれます。

| ファイル | 内容 |
|----------|------|
| `forks.yaml` | `fork_t7`（`main` @ t7） |
| `branches.yaml` | `branch_dog`（`via_action: act_fork_dog`） |
| `merges.yaml` | `merge_t11`（`branch_dog` と `main` を `鬼退治準備` へ合流） |
| `novels/branch_dog/scene4.txt` | 犬ルート専用の scene4 本文 |

```bash
novel-logic branch list
novel-logic fork show fork_t7
novel-logic merge show merge_t11
novel-logic timeline --branch main
novel-logic timeline --branch branch_dog
novel-logic validate --branch branch_dog
```

本線（`main`）は従来どおり三匹仲間 → 合流 → 鬼退治。`branch_dog` は犬のみ仲間にして t11 で本線に戻ります。

---

## Step 7: novel（本文）を登録する

Phase 0 では **1 scene × 1 branch に 1 本文**（`novels/main/<scene_id>.txt`）です。分岐ルートは `--branch` で別ディレクトリに置きます。

**重要**: CLI は **散文本文を書き込みません**。メタ登録と空ファイル作成だけを行い、本文はエディタで書きます。

### 7a. メタ登録

```bash
novel-logic novel add scene1
novel-logic novel add scene2
# scene3〜scene5 も同様
```

`novels.yaml` にメタが入り、空の `novels/scene1.txt` などができます（既にファイルがある場合はそのまま参照）。

### 7b. 本文を書く

エディタで各ファイルに散文を書きます。例（scene2）:

```
おじいさんとおばあさんは桃太郎を大切に育てた。桃太郎はすくすくと成長し、村で一番強い青年となった。
```

### 7c. git で版管理（任意・推奨）

作品ディレクトリが git リポジトリの場合:

```bash
git add novels/
git commit -m "add novel bodies"
novel-logic novel revision pin scene2 --note "初稿"
```

`revision pin` は「この scene の本文が、どの git commit 版か」を `novels.yaml` に記録します。PR / CI で本文の変更を検知する用途です。

> **本文は Lean 検証の対象外** です。論理層（fact / action / rule）との整合は `--scope novel:<scene>` 登録と、今後のフェーズで強化予定です。

### 7d. Phase B の論理登録（scene ごと）

本文に対応する fact / action は **`--scope novel:<scene_id>`** で登録します（plot 層とは別レコード）。

scene2 の例（育った）:

```bash
novel-logic thing scope add ojiisan --scope novel:scene2
novel-logic thing scope add obaasan --scope novel:scene2
novel-logic fact add --kind state --thing momotaro --pred 赤ちゃん --scope novel:scene2
novel-logic fact add --kind state --thing momotaro --pred 青年   --scope novel:scene2
novel-logic action add --thing momotaro --from 赤ちゃん --to 青年 \
  --at t4 --scope novel:scene2 --label 育った
```

`novel show scene2` で `alignment: novel layer only` または `mixed` になれば Phase B 登録済みです。

---

## Step 8: check（最終検証）

```bash
novel-logic check
```

実行内容:

| 段階 | 内容 |
|------|------|
| Stage 1 | Go による高速チェック（参照・rule・time 窓など） |
| generate | `logic/` に Lean 4 コードを自動生成 |
| Stage 2 | `lake build` で定理を機械証明 |

成功時の出力例:

```
OK: stage1
generated logic/
OK: stage1 + stage2
```

Stage 1 だけ確認したい場合:

```bash
novel-logic validate
# または
novel-logic check --quick
```

---

## Step 9: show / list で状況確認（読み取り専用）

登録後や作業の途中で、**YAML を開かずに** プロジェクトの状態を確認できます。  
いずれも **読み取り専用** で、データは変更しません。

### コマンドの実行方法

作品ディレクトリの外（例: `~/novel-logic/bin`）から実行する場合:

```bash
# bin/ にいるときは ./ を付ける（PATH 未設定の場合）
./novel-logic -C ~/novel-logic/examples/momotaro-walkthrough plot show

# 作品ディレクトリにいるとき
cd ~/novel-logic/examples/momotaro-walkthrough
../../bin/novel-logic plot show
```

以降、`-C <作品パス>` は省略して `novel-logic` と書きます。  
実際の環境に合わせて `./novel-logic` やフルパスに読み替えてください。

### 全体をざっと見る

| コマンド | 何がわかるか |
|----------|-------------|
| `novel-logic info` | 件数サマリ・最終 check の成否 |
| `novel-logic status` | Lean 検出結果・健全性 |
| `novel-logic plot show` | プロット要約 + scene 一覧 |
| `novel-logic timeline` | time 軸に scene 境界と action を並べた全体図 |

```bash
novel-logic info
novel-logic status
novel-logic plot show
novel-logic timeline
```

`timeline` は **物語の骨格を一望する** のに便利です。`--verbose` を付けると fact 一覧も出ます。

```
time_order:
  t4  action:act1          ← 桃太郎が青年に
  t12 action:act6, action:act7  ← 鬼退治（鬼と桃太郎が同時刻）
```

### エンティティ別に詳しく見る

```bash
# scene
novel-logic scene list
novel-logic scene show scene4     # その区間の action も表示

# 登場要素
novel-logic thing list
novel-logic thing list --tag character
novel-logic thing show momotaro   # 紐づく fact / action も表示

# time の順序
novel-logic time list

# 事実・遷移・ルール
novel-logic fact list
novel-logic fact list --thing momotaro
novel-logic fact list --kind state
novel-logic fact show fact7

novel-logic action list
novel-logic action show act7

novel-logic rule list
novel-logic rule show rule1

# 本文
novel-logic novel list
novel-logic novel show scene5
novel-logic novel show scene5 --full   # 全文表示（省略なし）
```

### コマンドと出力の対応（桃太郎 PJ の例）

**`thing show momotaro`** — 桃太郎に関する論理データをまとめて確認:

```
facts:
  fact1 [fixed] 人間
  fact7 [state] 赤ちゃん
  ...
actions:
  act1: 赤ちゃん → 青年 @ t4
  act2: 村在住 → 旅立ち @ t6
  act7: 旅立ち → 鬼退治済み @ t12
```

**`scene show scene4`** — 仲間入りの scene と、その time 窓内の action:

```
time: t7 .. t11
actions_in_window: 3
  - inu 野良 → 仲間 @ t8
  - saru 野良 → 仲間 @ t9
  - kiji 野良 → 仲間 @ t10
```

**`novel show scene5`** — 本文プレビュー + **関連 thing 一覧** + time 区間:

```
time: t11 .. t13
related_things (2):
  momotaro name=桃太郎 tags=[character]  (@t12:act7)
  oni name=鬼 tags=[character]  (@t12:act6)
body:
一行は鬼が島に渡り、鬼を退治して宝物を持ち帰った。...
```

### スコープと time の二層モデル（なぜ一致しないことがあるか）

| 層 | scope | いつ使うか | scene との関係 |
|----|-------|-----------|----------------|
| **プロット層** | `plot` | Phase A（骨格設計） | action の `at` が scene 窓内なら **継承表示**される |
| **本文層** | `novel:<scene_id>` | Phase B（本文執筆） | scene に **直接所属**。time 窓内であることも必須 |

`novel show` は意図的に **2 層を分けて表示**します。

```
layers: novel(facts=0 actions=0) plot_inherited(actions=2)
related_things (2):
  [novel] (none)
  [plot] momotaro ...  (@t12:act7)
  [plot] oni ...      (@t12:act6)
alignment: plot-inherited only — re-register with --scope novel:<scene> for Phase B
```

**一致させるルール（Phase B）**

1. scene 固有の fact / action は `--scope novel:scene5` で登録する（`plot` ではない）
2. `novel add` 時の time は **scene の time と自動同期**（`scenes.yaml` が正）
3. `novel:` スコープで fact/action を登録すると、主語 thing に `novel:<scene>` が **自動追加**される
4. 同一キーの `fact add` / `action add` は拒否。変更は `fact update` / `action update`

```bash
# Phase B の例（scene5 の鬼退治を本文層として登録）
novel-logic action add --thing momotaro --from 旅立ち --to 鬼退治済み \
  --at t12 --scope novel:scene5 --label 鬼が島で鬼を退治する
```

`validate --verbose` で、plot スコープのまま scene 窓内にある action に **hint** が出ます（エラーにはなりません）。

### 一覧フィルタの早見表

| コマンド | 主なフラグ |
|----------|-----------|
| `thing list` | `--tag <tag>`（繰り返し可・**すべて**持つ thing のみ） / `--scope plot` |
| `fact list` | `--tag <tag>` / `--kind fixed\|state\|all` / `--thing <id>` / `--scope <scope>` |
| `action list` | `--tag <tag>` / `--thing <id>` |
| `timeline` | `--verbose` |
| `novel show` | `--full` |

`--tag` の例:

```bash
novel-logic thing list --tag character      # キャラクターのみ
novel-logic thing list --tag location       # 場所のみ
novel-logic fact list --tag character       # キャラクターに紐づく fact のみ
novel-logic action list --tag character     # キャラクターの action のみ
```

フィルタが有効なとき、先頭に `# filter: tag=[character]` のような行が表示されます。

### ヘルプの見方

```bash
novel-logic --help
novel-logic thing --help
novel-logic scene show --help
```

---

## 完成したタイムライン

```
t1   scene1 [novel: scene1.txt]  桃太郎誕生
t3   scene2 [novel: scene2.txt]  育てられる
t4         桃太郎: 赤ちゃん → 青年
t5   scene3 [novel: scene3.txt]  旅立ち決意
t6         桃太郎: 村在住 → 旅立ち
t7   scene4 [novel: scene4.txt]  仲間入り
t8         犬:   野良 → 仲間
t9         猿:   野良 → 仲間
t10        雉:   野良 → 仲間
t11  scene5 [novel: scene5.txt]  鬼退治
t12        鬼:     健在 → 退治済み
           桃太郎: 旅立ち → 鬼退治済み
```

---

## このディレクトリの中身

| ファイル / フォルダ | 内容 |
|--------------------|------|
| `project.yaml` | タイトル、time 順序、最終 check 結果 |
| `plot.yaml` | プロット要約 |
| `things.yaml` | 登場要素 |
| `scenes.yaml` | scene 定義 |
| `times.yaml` | time の存在リスト |
| `facts.yaml` | fixed_fact + state |
| `actions.yaml` | 状態遷移 |
| `rules.yaml` | 禁則 |
| `novels.yaml` | 本文メタ（scene ごと） |
| `novels/*.txt` | 本文テキスト（git 管理。CLI は書き込まない） |
| `logic/` | `check` 時に自動生成される Lean プロジェクト（手編集非推奨） |

### 現在の登録件数（参考）

```
things: 10 / scenes: 5 / facts: 19 / actions: 7 / rules: 3 / novels: 5
```

---

## 削除（remove）

各エンティティを ID（または scene_id）で削除できます。他から参照されている場合は **削除を拒否** し、参照元を表示します。

```bash
novel-logic action remove act1
novel-logic fact remove fact7
novel-logic rule remove rule1
novel-logic thing scope remove momotaro --scope novel:scene5   # スコープのみ外す
novel-logic novel remove scene1                               # デフォルトで .txt は残す
novel-logic novel remove scene1 --keep-body=false             # メタと .txt 両方削除
novel-logic thing remove mob_only                             # 参照がなければ削除可
novel-logic scene remove scene3                               # novel / fact / action 参照がなければ可
novel-logic time remove t99                                 # scene / action / novel 参照がなければ可
```

削除のあとも `validate` / `check` で整合を確認できます。

---

## よくある質問

### Q. fixed_fact と state の違いは？

| 種別 | 変化 | 例 |
|------|------|-----|
| fixed_fact | 変わらない | 桃太郎は人間 |
| state | action で変わりうる | 桃太郎は赤ちゃん → 青年 |

fixed_fact を state に **昇格** することは可能です（`fact promote <id>`）。逆方向は不可です。

### Q. なぜ rule を先に登録するの？

rule を先に置いておくと、矛盾する action を **登録時点で拒否** できます。  
全体を組み立てたあと `check` でも検出しますが、早い段階で防げる方が修正コストが低いです。

### Q. 日本語の述語（人間、赤ちゃんなど）はそのまま使える？

YAML 上は日本語のまま書けます。Lean 生成時は内部 ID（`ningen`, `akachan` など）に変換されます。  
ユーザーが意識する必要はありません。

### Q. `novel-logic: command not found` と出る

バイナリは `novel-logic/bin/novel-logic` にありますが、PATH に入っていないとその名前では起動しません。

```bash
# bin/ にいるとき
./novel-logic -C ../examples/momotaro-walkthrough plot show

# 毎回使うなら PATH に追加（~/.bashrc）
export PATH="$HOME/novel-logic/bin:$PATH"
```

### Q. CLI で本文（散文）を登録できないの？

**意図的にありません。** 本文は `novels/<scene_id>.txt` をエディタで書き、git で管理します。  
CLI ができるのは `novel add`（メタ + 空ファイル）、`novel revision pin`（git commit 記録）、`novel show`（表示）です。

### Q. 他の作品を始めたい

```bash
novel-logic init ~/novels/my-story --template default
cd ~/novels/my-story
# Step 1 から繰り返す
```

---

## 次のステップ

- コマンドの全一覧: [docs/COMMANDS.md](../../docs/COMMANDS.md)
- ドメインモデルの詳細: [docs/REQUIREMENTS.md](../../docs/REQUIREMENTS.md)
- 対話型ウィザード（`novel-logic wizard`）は Phase 1 予定