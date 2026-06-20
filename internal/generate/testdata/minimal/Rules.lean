import Core
import Project

namespace test

open NovelLogic

def projectRules_main : Rules ThingId PredId := {
  forbiddenStates := [
  ],
  forbiddenTransitions := [
  ]
}

abbrev projectRules : Rules ThingId PredId := projectRules_main

end test
