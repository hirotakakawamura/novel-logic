import Core
import Project

namespace Untitled

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
  ⟨ThingId.momotaro, PredId.murazaiju, Scope.plot⟩,
  ⟨ThingId.momotaro, PredId.tabidachi, Scope.plot⟩,
  ⟨ThingId.inu, PredId.nora, Scope.plot⟩,
  ⟨ThingId.inu, PredId.nakama, Scope.plot⟩,
  ⟨ThingId.saru, PredId.nora, Scope.plot⟩,
  ⟨ThingId.saru, PredId.nakama, Scope.plot⟩,
  ⟨ThingId.kiji, PredId.nora, Scope.plot⟩,
  ⟨ThingId.kiji, PredId.nakama, Scope.plot⟩,
  ⟨ThingId.oni, PredId.kenzen, Scope.plot⟩,
  ⟨ThingId.oni, PredId.taijizumi, Scope.plot⟩,
  ⟨ThingId.momotaro, PredId.onitaijizumi, Scope.plot⟩,
  ⟨ThingId.momotaro, PredId.akachan, Scope.novel_scene2⟩,
  ⟨ThingId.momotaro, PredId.seinen, Scope.novel_scene2⟩,
]

def allActions : List (ActionDecl ThingId PredId TimeId Scope) := [
  ⟨ThingId.momotaro, some PredId.akachan, PredId.seinen, TimeId.t4, Scope.plot⟩,
  ⟨ThingId.momotaro, some PredId.murazaiju, PredId.tabidachi, TimeId.t6, Scope.plot⟩,
  ⟨ThingId.inu, some PredId.nora, PredId.nakama, TimeId.t8, Scope.plot⟩,
  ⟨ThingId.saru, some PredId.nora, PredId.nakama, TimeId.t9, Scope.plot⟩,
  ⟨ThingId.kiji, some PredId.nora, PredId.nakama, TimeId.t10, Scope.plot⟩,
  ⟨ThingId.oni, some PredId.kenzen, PredId.taijizumi, TimeId.t12, Scope.plot⟩,
  ⟨ThingId.momotaro, some PredId.tabidachi, PredId.onitaijizumi, TimeId.t12, Scope.plot⟩,
  ⟨ThingId.momotaro, some PredId.akachan, PredId.seinen, TimeId.t4, Scope.novel_scene2⟩,
]

def activeActions_main : List (ActionDecl ThingId PredId TimeId Scope) := [
  ⟨ThingId.momotaro, some PredId.akachan, PredId.seinen, TimeId.t4, Scope.plot⟩,
  ⟨ThingId.momotaro, some PredId.akachan, PredId.seinen, TimeId.t4, Scope.novel_scene2⟩,
  ⟨ThingId.momotaro, some PredId.murazaiju, PredId.tabidachi, TimeId.t6, Scope.plot⟩,
  ⟨ThingId.inu, some PredId.nora, PredId.nakama, TimeId.t8, Scope.plot⟩,
  ⟨ThingId.saru, some PredId.nora, PredId.nakama, TimeId.t9, Scope.plot⟩,
  ⟨ThingId.kiji, some PredId.nora, PredId.nakama, TimeId.t10, Scope.plot⟩,
  ⟨ThingId.oni, some PredId.kenzen, PredId.taijizumi, TimeId.t12, Scope.plot⟩,
  ⟨ThingId.momotaro, some PredId.tabidachi, PredId.onitaijizumi, TimeId.t12, Scope.plot⟩,
]

def evolveBranch_main (t : TimeId) (thing : ThingId) : List PredId :=
  predsAt allFixedFacts allStateDecls activeActions_main timeOrder t thing

end Untitled
