# gitreset

Read `references/common.md` before acting. Follow the user's explicit scope and this command's mutation boundary.

Require exactly one of `--soft`, `--mixed`, `--hard` and require `--to <revision>`. No guessed hard mode or target. Resolve target to a commit SHA, show current/target branch and outgoing commit range. Default permitted targets are HEAD or its ancestors; clarify other histories/branches. Reject detached HEAD unless its intent is explicit and reject active operations; use a dedicated merge command for merge cancellation.
Apply common publication/recovery rules. Soft preserves index and worktree; mixed resets index to target while preserving worktree; hard replaces index/tracked worktree content and can overwrite obstructing untracked paths. Before hard, identify actual losses, verify an outside-worktree snapshot and ask for concrete confirmation if not already approved. A backup ref alone does not protect uncommitted files. Never git clean, recurse into submodules, or silently remove ignored/untracked files beyond the approved effect. Use the verified object ID, recheck state, execute the chosen reset and verify resulting status. Remote stays unchanged.
