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

theorem actions_in_scene_window_main :
    allActionsInSceneWindows sceneWindows timeOrder activeActions_main scopeToScene := by
  native_decide

theorem no_forbidden_states_main :
    noForbiddenStatesRegistered projectRules_main stateDecls_main := by
  native_decide

theorem no_forbidden_transitions_main :
    allActionsRespectRules projectRules_main activeActions_main := by
  native_decide

theorem fixed_facts_stable_main :
    fixedFactsStable fixedFacts_main stateDecls_main activeActions_main timeOrder := by
  native_decide

end test
