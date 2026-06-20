import Core
import Project

namespace Momotaro

open NovelLogic

def sceneWindows : List (SceneWindow SceneId TimeId) := [
  ⟨SceneId.scene1, TimeId.t1, TimeId.t3⟩,
  ⟨SceneId.scene2, TimeId.t3, TimeId.t5⟩,
  ⟨SceneId.scene3, TimeId.t5, TimeId.t7⟩,
  ⟨SceneId.scene4, TimeId.t7, TimeId.t11⟩,
  ⟨SceneId.scene5, TimeId.t11, TimeId.t13⟩,
]

end Momotaro
