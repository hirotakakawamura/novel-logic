import Core
import Project

namespace Momotaro

open NovelLogic

def allFixedFacts : List (FixedFact ThingId PredId Scope) := [
  ⟨ThingId.momotaro, PredId.ningen, Scope.plot⟩,
  ⟨ThingId.ojiisan, PredId.ningen, Scope.plot⟩,
  ⟨ThingId.obaasan, PredId.ningen, Scope.plot⟩,
  ⟨ThingId.inu, PredId.doubutsu, Scope.plot⟩,
  ⟨ThingId.saru, PredId.doubutsu, Scope.plot⟩,
  ⟨ThingId.kiji, PredId.doubutsu, Scope.plot⟩,
]

def allStateDecls : List (StateDecl ThingId PredId Scope) := [
  ⟨ThingId.momotaro, PredId.akachan, Scope.plot⟩,
  ⟨ThingId.momotaro, PredId.seinen, Scope.plot⟩,
]

def allActions : List (ActionDecl ThingId PredId TimeId Scope) := [
  ⟨ThingId.momotaro, some PredId.akachan, PredId.seinen, TimeId.t4, Scope.plot⟩,
]

end Momotaro
