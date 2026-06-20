import Core
import Project
import Facts
import Rules
import Timeline

namespace Untitled

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

theorem actions_in_scene_window_branch_dog :
    allActionsInSceneWindows sceneWindows timeOrder activeActions_branch_dog scopeToScene := by
  native_decide

theorem no_forbidden_transitions_branch_dog :
    allActionsRespectRules projectRules activeActions_branch_dog := by
  native_decide

theorem fixed_facts_stable_branch_dog :
    fixedFactsStable allFixedFacts allStateDecls activeActions_branch_dog timeOrder := by
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

theorem forbid_state_momotaro_doubutsu_at_end :
    ¬ listContains (predsAt allFixedFacts allStateDecls allActions timeOrder TimeId.t13 ThingId.momotaro) PredId.doubutsu := by
  native_decide

end Untitled
