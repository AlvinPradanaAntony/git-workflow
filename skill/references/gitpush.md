# gitpush

Read `references/common.md` before acting. Follow the user's explicit scope and this command's mutation boundary.

Read `references/push.md`. Accept `--remote <name>` and `--branch <name>` with a complete unambiguous target. No --force or hidden amend. Resolve current branch/upstream and inspect the outgoing committed range, publication target and relevant project gates. Uncommitted changes remain local and should be reported, not staged. Use a targeted ordinary push only. If local history was amended/reset and remote is divergent, explain why normal push cannot proceed; require an explicitly authorized targeted lease push in follow-up, and do NOT call gitpushamend again when amendment already happened. Verify remote result and report target/hash.
