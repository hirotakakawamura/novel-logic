import Core
import Project

namespace Momotaro

open NovelLogic

def projectRules : Rules ThingId PredId := {
  forbiddenStates := [
    (ThingId.momotaro, PredId.doubutsu),
  ],
  forbiddenTransitions := [
    (PredId.seinen, PredId.akachan),
    (PredId.ningen, PredId.doubutsu),
  ]
}

end Momotaro
