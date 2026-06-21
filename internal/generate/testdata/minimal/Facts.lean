import Core
import Project

namespace test

open NovelLogic

def fixedFacts_main : List (FixedFact ThingId PredId Scope) := [
]

def stateDecls_main : List (StateDecl ThingId PredId Scope) := [
  ⟨ThingId.hero, PredId.start, Scope.plot⟩,
]

abbrev allFixedFacts := fixedFacts_main

abbrev allStateDecls := stateDecls_main

def activeActions_main : List (ActionDecl ThingId PredId TimeId Scope) := [
  ⟨ThingId.hero, some PredId.start, PredId.mid, TimeId.t2, Scope.plot⟩,
]

def evolveBranch_main (t : TimeId) (thing : ThingId) : List PredId :=
  predsAt fixedFacts_main stateDecls_main activeActions_main timeOrder t thing

end test
