Novel body files for this work live here as plain text (.txt).
Edit them in your editor and commit with git; novel-logic does not overwrite prose.

Register a scene body:
  novel-logic novel add <scene_id>

After committing prose in git, pin the revision for check/CI:
  novel-logic novel revision pin <scene_id> [--note "PR #42"]

Default path: novels/<scene_id>.txt