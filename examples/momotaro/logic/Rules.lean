import Core
import Project

namespace Momotaro

open NovelLogic

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

end Momotaro
