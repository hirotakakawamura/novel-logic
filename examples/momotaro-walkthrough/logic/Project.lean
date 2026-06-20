import Core

namespace Untitled

inductive ThingId
  | inu
  | kibidango
  | kiji
  | momotaro
  | obaasan
  | ojiisan
  | oni
  | onigashima
  | saru
  | village
  deriving DecidableEq, Repr

inductive TimeId
  | t1
  | t2
  | t3
  | t4
  | t5
  | t6
  | t7
  | t8
  | t9
  | t10
  | t11
  | t12
  | t13
  deriving DecidableEq, Repr

inductive BranchId
  | main
  deriving DecidableEq, Repr

inductive SceneId
  | scene1
  | scene2
  | scene3
  | scene4
  | scene5
  deriving DecidableEq, Repr

inductive PredId
  | ningen
  | nakama
  | kenzen
  | doubutsu
  | tabidachi
  | murazaiju
  | akachan
  | taijizumi
  | nora
  | seinen
  | onitaijizumi
  deriving DecidableEq, Repr

inductive Scope
  | plot
  | novel_scene1
  | novel_scene2
  | novel_scene3
  | novel_scene4
  | novel_scene5
  deriving DecidableEq, Repr

abbrev PlotScope : Scope := Scope.plot

def timeOrder : List TimeId := [TimeId.t1, TimeId.t2, TimeId.t3, TimeId.t4, TimeId.t5, TimeId.t6, TimeId.t7, TimeId.t8, TimeId.t9, TimeId.t10, TimeId.t11, TimeId.t12, TimeId.t13]

def scopeToScene : Scope → Option SceneId
  | Scope.plot => none
  | Scope.novel_scene1 => some SceneId.scene1
  | Scope.novel_scene2 => some SceneId.scene2
  | Scope.novel_scene3 => some SceneId.scene3
  | Scope.novel_scene4 => some SceneId.scene4
  | Scope.novel_scene5 => some SceneId.scene5

end Untitled
