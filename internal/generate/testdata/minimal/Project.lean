import Core

namespace test

inductive ThingId
  | ally
  | hero
  deriving DecidableEq, Repr

inductive TimeId
  | t1
  | t2
  | t3
  | t4
  deriving DecidableEq, Repr

inductive BranchId
  | main
  deriving DecidableEq, Repr

inductive SceneId
  | scene1
  | scene2
  deriving DecidableEq, Repr

inductive PredId
  | mid
  | start
  deriving DecidableEq, Repr

inductive Scope
  | plot
  | novel_scene1
  | novel_scene2
  deriving DecidableEq, Repr

abbrev PlotScope : Scope := Scope.plot

def timeOrder : List TimeId := [TimeId.t1, TimeId.t2, TimeId.t3, TimeId.t4]

def scopeToScene : Scope → Option SceneId
  | Scope.plot => none
  | Scope.novel_scene1 => some SceneId.scene1
  | Scope.novel_scene2 => some SceneId.scene2

end test
