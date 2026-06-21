import Core
import Project

namespace Momotaro

open NovelLogic

def fixedFacts_main : List (FixedFact ThingId PredId Scope) := [
  ⟨ThingId.momotaro, PredId.ningen, Scope.plot⟩,
  ⟨ThingId.ojiisan, PredId.ningen, Scope.plot⟩,
  ⟨ThingId.obaasan, PredId.ningen, Scope.plot⟩,
  ⟨ThingId.inu, PredId.doubutsu, Scope.plot⟩,
  ⟨ThingId.saru, PredId.doubutsu, Scope.plot⟩,
  ⟨ThingId.kiji, PredId.doubutsu, Scope.plot⟩,
]

def stateDecls_main : List (StateDecl ThingId PredId Scope) := [
  ⟨ThingId.momotaro, PredId.akachan, Scope.plot⟩,
  ⟨ThingId.momotaro, PredId.seinen, Scope.plot⟩,
]

abbrev allFixedFacts := fixedFacts_main

abbrev allStateDecls := stateDecls_main

def activeActions_main : List (ActionDecl ThingId PredId TimeId Scope) := [
  ⟨ThingId.momotaro, some PredId.akachan, PredId.seinen, TimeId.t4, Scope.plot⟩,
]

def evolveBranch_main (t : TimeId) (thing : ThingId) : List PredId :=
  predsAt fixedFacts_main stateDecls_main activeActions_main timeOrder t thing

end Momotaro
