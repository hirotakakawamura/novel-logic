import Core
import Project
import Facts
import Rules
import Timeline

namespace test

open NovelLogic

theorem actions_in_scene_window :
    allActionsInSceneWindows sceneWindows timeOrder allActions scopeToScene := by
  native_decide

theorem no_forbidden_states :
    noForbiddenStatesRegistered projectRules allStateDecls := by
  native_decide

theorem no_forbidden_transitions :
    allActionsRespectRules projectRules allActions := by
  native_decide

theorem fixed_facts_stable :
    fixedFactsStable allFixedFacts allStateDecls allActions timeOrder := by
  native_decide

theorem actions_in_scene_window_main :
    allActionsInSceneWindows sceneWindows timeOrder activeActions_main scopeToScene := by
  native_decide

theorem no_forbidden_transitions_main :
    allActionsRespectRules projectRules activeActions_main := by
  native_decide

theorem fixed_facts_stable_main :
    fixedFactsStable allFixedFacts allStateDecls activeActions_main timeOrder := by
  native_decide

end test
