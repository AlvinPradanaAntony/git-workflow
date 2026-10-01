# gitpushamend

Read `references/common.md` before acting. Follow the user's explicit scope and this command's mutation boundary.

Read `references/push.md`. Accept `--remote <name> --branch <name>` together when no verified upstream exists; allow overriding either only when the complete destination is unambiguous. Optional --no-edit is redundant and accepted. Reject message/language/format changes and --all staging.
Require normal operation state, an existing attached branch and nonempty staged content. Do NOT amend if nothing is staged; suggest gitpush for an existing amended commit awaiting publication. Inspect staged changes for relevance and preserve partial staging. Perform the push reference's entire remote/protection/reachability preflight BEFORE changing the local commit. Create the recovery ref. Amend exactly once via commit --amend --no-edit, verify message, authorship, parents and intended diff, then perform normal push or the pinned-lease rewrite specified in push.md. If any local verification fails, do not publish. If push fails, keep and report the successful local amend; retry only push after reinspection, never regenerate/amend again.
