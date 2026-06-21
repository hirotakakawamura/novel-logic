import Core
import Project
import Facts
import Rules
import Timeline

namespace test

open NovelLogic

theorem actions_in_scene_window :
    allActionsInSceneWindows sceneWindows timeOrder activeActions_main scopeToScene := by
  native_decide

theorem no_forbidden_states :
    noForbiddenStatesRegistered projectRules stateDecls_main := by
  native_decide

theorem no_forbidden_transitions :
    allActionsRespectRules projectRules activeActions_main := by
  native_decide

theorem fixed_facts_stable :
    fixedFactsStable fixedFacts_main stateDecls_main activeActions_main timeOrder := by
  native_decide

end test
