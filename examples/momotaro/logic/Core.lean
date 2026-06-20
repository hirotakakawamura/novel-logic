/-!
# novel-logic Core — generic finite enumeration runtime (all works share this file).
-/
namespace NovelLogic

structure FixedFact (τ π σ : Type) where
  subject : τ
  pred : π
  scope : σ
  deriving DecidableEq, Repr

structure StateDecl (τ π σ : Type) where
  subject : τ
  pred : π
  scope : σ
  deriving DecidableEq, Repr

structure ActionDecl (τ π ι σ : Type) where
  subject : τ
  from? : Option π
  to : π
  time : ι
  scope : σ
  deriving DecidableEq, Repr

structure SceneWindow (σc ι : Type) where
  scene : σc
  start : ι
  stop : ι
  deriving DecidableEq, Repr

structure Rules (τ π : Type) where
  forbiddenStates : List (τ × π)
  forbiddenTransitions : List (π × π)

def timeIndex {ι : Type} [DecidableEq ι] (order : List ι) (t : ι) : Option Nat :=
  order.findIdx? (· == t)

def timeLe {ι : Type} [DecidableEq ι] (order : List ι) (a b : ι) : Bool :=
  match timeIndex order a, timeIndex order b with
  | some ia, some ib => decide (ia ≤ ib)
  | _, _ => false

def listContains {α : Type} [DecidableEq α] (xs : List α) (x : α) : Bool :=
  xs.elem x

def actionRespectsRules {τ π ι σ : Type} [DecidableEq τ] [DecidableEq π]
    (rules : Rules τ π) (a : ActionDecl τ π ι σ) : Bool :=
  let stateOk := ¬ listContains rules.forbiddenStates (a.subject, a.to)
  let transOk := match a.from? with
    | some f => ¬ listContains rules.forbiddenTransitions (f, a.to)
    | none => true
  decide (stateOk && transOk)

def applyAction {π : Type} [DecidableEq π] (active : List π) (from? : Option π) (to : π) : List π :=
  let without := match from? with
    | some f => active.erase f
    | none => active
  if listContains without to then without else to :: without

def predsAt {τ π ι σ : Type} [DecidableEq τ] [DecidableEq π] [DecidableEq ι]
    (fixed : List (FixedFact τ π σ)) (states : List (StateDecl τ π σ))
    (acts : List (ActionDecl τ π ι σ)) (order : List ι) (t : ι) (thing : τ) : List π :=
  let fixedPreds := (fixed.filter fun f => f.subject == thing).map (·.pred)
  let statePreds := (states.filter fun s => s.subject == thing).map (·.pred)
  let base := (fixedPreds ++ statePreds).eraseDups
  let relevant := acts.filter fun a =>
    a.subject == thing && timeLe order a.time t
  let sorted := relevant.mergeSort fun a b =>
    match timeIndex order a.time, timeIndex order b.time with
    | some ia, some ib => decide (ia ≤ ib)
    | some _, none => true
    | none, some _ => false
    | none, none => true
  sorted.foldl (fun acc a => applyAction acc a.from? a.to) base

def allActionsRespectRules {τ π ι σ : Type} [DecidableEq τ] [DecidableEq π]
    (rules : Rules τ π) (acts : List (ActionDecl τ π ι σ)) : Bool :=
  acts.all fun a => actionRespectsRules rules a

def noForbiddenStatesRegistered {τ π σ : Type} [DecidableEq τ] [DecidableEq π]
    (rules : Rules τ π) (states : List (StateDecl τ π σ)) : Bool :=
  states.all fun s => ¬ listContains rules.forbiddenStates (s.subject, s.pred)

def fixedFactsStable {τ π ι σ : Type} [DecidableEq τ] [DecidableEq π] [DecidableEq ι]
    (fixed : List (FixedFact τ π σ)) (states : List (StateDecl τ π σ))
    (acts : List (ActionDecl τ π ι σ)) (order : List ι) : Bool :=
  fixed.all fun f =>
    order.all fun t => listContains (predsAt fixed states acts order t f.subject) f.pred

def actionInSceneWindow {τ π ι σ σc : Type} [DecidableEq ι] [DecidableEq σc]
    (windows : List (SceneWindow σc ι)) (order : List ι)
    (a : ActionDecl τ π ι σ) (scopeToScene : σ → Option σc) : Bool :=
  match scopeToScene a.scope with
  | none => true
  | some sid =>
    match windows.find? (·.scene == sid) with
    | none => false
    | some w => timeLe order w.start a.time && timeLe order a.time w.stop

def allActionsInSceneWindows {τ π ι σ σc : Type}
    [DecidableEq τ] [DecidableEq π] [DecidableEq ι] [DecidableEq σ] [DecidableEq σc]
    (windows : List (SceneWindow σc ι)) (order : List ι)
    (acts : List (ActionDecl τ π ι σ)) (scopeToScene : σ → Option σc) : Bool :=
  acts.all fun a => actionInSceneWindow windows order a scopeToScene

end NovelLogic