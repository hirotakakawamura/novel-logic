import Core
import Project
import Facts
import Rules
import Timeline

namespace Momotaro

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

theorem forbid_state_main_momotaro_doubutsu_at_end :
    ¬ listContains (predsAt fixedFacts_main stateDecls_main activeActions_main timeOrder TimeId.t13 ThingId.momotaro) PredId.doubutsu := by
  native_decide

end Momotaro
