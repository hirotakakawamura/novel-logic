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

theorem forbid_state_momotaro_doubutsu_at_end :
    ¬ listContains (predsAt allFixedFacts allStateDecls allActions timeOrder TimeId.t13 ThingId.momotaro) PredId.doubutsu := by
  native_decide

end Untitled
