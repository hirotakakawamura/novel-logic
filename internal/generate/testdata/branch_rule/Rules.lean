import Core
import Project

namespace test

open NovelLogic

def projectRules_branch_a : Rules ThingId PredId := {
  forbiddenStates := [
    (ThingId.hero, PredId.bad),
  ],
  forbiddenTransitions := [
    (PredId.mid, PredId.start),
  ]
}

def projectRules_main : Rules ThingId PredId := {
  forbiddenStates := [
  ],
  forbiddenTransitions := [
    (PredId.mid, PredId.start),
  ]
}

abbrev projectRules : Rules ThingId PredId := projectRules_main

end test
