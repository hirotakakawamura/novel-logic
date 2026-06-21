import Core
import Project

namespace Untitled

open NovelLogic

def fixedFacts_branch_dog : List (FixedFact ThingId PredId Scope) := [
  ⟨ThingId.momotaro, PredId.ningen, Scope.plot⟩,
  ⟨ThingId.ojiisan, PredId.ningen, Scope.plot⟩,
  ⟨ThingId.obaasan, PredId.ningen, Scope.plot⟩,
  ⟨ThingId.inu, PredId.doubutsu, Scope.plot⟩,
  ⟨ThingId.saru, PredId.doubutsu, Scope.plot⟩,
  ⟨ThingId.kiji, PredId.doubutsu, Scope.plot⟩,
]

def stateDecls_branch_dog : List (StateDecl ThingId PredId Scope) := [
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
  ⟨ThingId.momotaro, PredId.pred_daf4895b, Scope.plot⟩,
  ⟨ThingId.momotaro, PredId.pred_fb795e08, Scope.plot⟩,
]

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
  ⟨ThingId.momotaro, PredId.pred_daf4895b, Scope.plot⟩,
  ⟨ThingId.momotaro, PredId.pred_fb795e08, Scope.plot⟩,
]

abbrev allFixedFacts := fixedFacts_main

abbrev allStateDecls := stateDecls_main

def activeActions_branch_dog : List (ActionDecl ThingId PredId TimeId Scope) := [
  ⟨ThingId.momotaro, some PredId.akachan, PredId.seinen, TimeId.t4, Scope.plot⟩,
  ⟨ThingId.momotaro, some PredId.akachan, PredId.seinen, TimeId.t4, Scope.novel_scene2⟩,
  ⟨ThingId.momotaro, some PredId.murazaiju, PredId.tabidachi, TimeId.t6, Scope.plot⟩,
  ⟨ThingId.momotaro, some PredId.tabidachi, PredId.pred_daf4895b, TimeId.t7, Scope.plot⟩,
  ⟨ThingId.inu, some PredId.nora, PredId.nakama, TimeId.t8, Scope.plot⟩,
  ⟨ThingId.momotaro, some PredId.tabidachi, PredId.pred_fb795e08, TimeId.t11, Scope.plot⟩,
  ⟨ThingId.oni, some PredId.kenzen, PredId.taijizumi, TimeId.t12, Scope.plot⟩,
  ⟨ThingId.momotaro, some PredId.pred_fb795e08, PredId.onitaijizumi, TimeId.t12, Scope.plot⟩,
]

def evolveBranch_branch_dog (t : TimeId) (thing : ThingId) : List PredId :=
  predsAt fixedFacts_branch_dog stateDecls_branch_dog activeActions_branch_dog timeOrder t thing

def activeActions_main : List (ActionDecl ThingId PredId TimeId Scope) := [
  ⟨ThingId.momotaro, some PredId.akachan, PredId.seinen, TimeId.t4, Scope.plot⟩,
  ⟨ThingId.momotaro, some PredId.akachan, PredId.seinen, TimeId.t4, Scope.novel_scene2⟩,
  ⟨ThingId.momotaro, some PredId.murazaiju, PredId.tabidachi, TimeId.t6, Scope.plot⟩,
  ⟨ThingId.momotaro, some PredId.tabidachi, PredId.pred_daf4895b, TimeId.t7, Scope.plot⟩,
  ⟨ThingId.inu, some PredId.nora, PredId.nakama, TimeId.t8, Scope.plot⟩,
  ⟨ThingId.saru, some PredId.nora, PredId.nakama, TimeId.t9, Scope.plot⟩,
  ⟨ThingId.kiji, some PredId.nora, PredId.nakama, TimeId.t10, Scope.plot⟩,
  ⟨ThingId.momotaro, some PredId.tabidachi, PredId.pred_fb795e08, TimeId.t11, Scope.plot⟩,
  ⟨ThingId.oni, some PredId.kenzen, PredId.taijizumi, TimeId.t12, Scope.plot⟩,
  ⟨ThingId.momotaro, some PredId.pred_fb795e08, PredId.onitaijizumi, TimeId.t12, Scope.plot⟩,
]

def evolveBranch_main (t : TimeId) (thing : ThingId) : List PredId :=
  predsAt fixedFacts_main stateDecls_main activeActions_main timeOrder t thing

end Untitled
