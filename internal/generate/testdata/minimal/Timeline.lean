import Core
import Project

namespace test

open NovelLogic

def sceneWindows : List (SceneWindow SceneId TimeId) := [
  ⟨SceneId.scene1, TimeId.t1, TimeId.t2⟩,
  ⟨SceneId.scene2, TimeId.t2, TimeId.t4⟩,
]

end test
