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

theorem actions_in_scene_window_branch_a :
    allActionsInSceneWindows sceneWindows timeOrder activeActions_branch_a scopeToScene := by
  native_decide

theorem no_forbidden_transitions_branch_a :
    allActionsRespectRules projectRules_branch_a activeActions_branch_a := by
  native_decide

theorem fixed_facts_stable_branch_a :
    fixedFactsStable allFixedFacts allStateDecls activeActions_branch_a timeOrder := by
  native_decide

theorem actions_in_scene_window_main :
    allActionsInSceneWindows sceneWindows timeOrder activeActions_main scopeToScene := by
  native_decide

theorem no_forbidden_transitions_main :
    allActionsRespectRules projectRules_main activeActions_main := by
  native_decide

theorem fixed_facts_stable_main :
    fixedFactsStable allFixedFacts allStateDecls activeActions_main timeOrder := by
  native_decide

theorem forbid_state_hero_bad_at_end :
    ¬ listContains (predsAt allFixedFacts allStateDecls allActions timeOrder TimeId.t4 ThingId.hero) PredId.bad := by
  native_decide

end test
