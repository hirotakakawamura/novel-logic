import Core
import Project

namespace Untitled

open NovelLogic

def projectRules_branch_dog : Rules ThingId PredId := {
  forbiddenStates := [
    (ThingId.momotaro, PredId.doubutsu),
  ],
  forbiddenTransitions := [
    (PredId.seinen, PredId.akachan),
    (PredId.ningen, PredId.doubutsu),
  ]
}

def projectRules_main : Rules ThingId PredId := {
  forbiddenStates := [
    (ThingId.momotaro, PredId.doubutsu),
  ],
  forbiddenTransitions := [
    (PredId.seinen, PredId.akachan),
    (PredId.ningen, PredId.doubutsu),
  ]
}

abbrev projectRules : Rules ThingId PredId := projectRules_main

end Untitled
