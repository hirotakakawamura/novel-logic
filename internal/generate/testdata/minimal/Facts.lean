import Core
import Project

namespace test

open NovelLogic

def allFixedFacts : List (FixedFact ThingId PredId Scope) := [
]

def allStateDecls : List (StateDecl ThingId PredId Scope) := [
  ⟨ThingId.hero, PredId.start, Scope.plot⟩,
]

def allActions : List (ActionDecl ThingId PredId TimeId Scope) := [
  ⟨ThingId.hero, some PredId.start, PredId.mid, TimeId.t2, Scope.plot⟩,
]

def activeActions_main : List (ActionDecl ThingId PredId TimeId Scope) := [
  ⟨ThingId.hero, some PredId.start, PredId.mid, TimeId.t2, Scope.plot⟩,
]

def evolveBranch_main (t : TimeId) (thing : ThingId) : List PredId :=
  predsAt allFixedFacts allStateDecls activeActions_main timeOrder t thing

end test
