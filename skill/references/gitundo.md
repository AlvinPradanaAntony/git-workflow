# gitundo

Read `references/common.md` before acting. Follow the user's explicit scope and this command's mutation boundary.

Accept `--count <positive integer>` (1), `--to staged|changes` (staged), and `--parent 1|2` only for undoing a single merge commit. Read and follow common recovery rules.
Resolve the exact ancestor target. For ordinary linear history use first-parent ancestry for count, but stop if the requested span crosses a merge until the intended ancestry is agreed. For a tip merge commit, display its parents and require a concrete parent selection unless already explicit; reject invalid parent indexes. `--parent` with count other than 1 is invalid. If there is no parent (initial commit) or count exceeds history, stop with an explanation; never delete refs to emulate an unborn branch.
Reject an active merge/rebase/sequencer state and direct the user to the matching operation workflow. Check publication; if published or uncertain, explain the local/remote divergence and obtain concrete rewrite intent before proceeding unless already provided. Do not silently substitute revert.
Record backup ref and explain existing staged/unstaged changes will remain/merge into the resulting diff. `staged` maps to soft reset to the verified target; `changes` maps to mixed reset. Soft preserves current index and worktree, not magically stages preexisting unstaged hunks; mixed unstages everything and may expose added files as untracked. Verify old/new HEAD, index and worktree. Never push.
