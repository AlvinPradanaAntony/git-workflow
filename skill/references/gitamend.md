# gitamend

Read `references/common.md` before acting. Follow the user's explicit scope and this command's mutation boundary.

Accept optional `--no-edit` (always the default); reject message/language/format rewrite flags and automatic --all staging. Require an existing HEAD, normal operation state and nonempty staged diff. If no staged content, report nothing to amend and do not create a timestamp-only replacement. Preserve partial staging and unrelated work.
Inspect whether staged changes belong to the current commit; ask if clearly unrelated. Check publication and apply common rewrite/backup rules; tell the user a published commit's local identity will change and remote will remain unchanged. Preserve existing authorship, parents and full original message including bilingual content/trailers. Use `git commit --amend --no-edit`, obey hooks, verify the exact resulting message/parents/diff and report old/new hash, remaining changes and no push. If hooks alter the message, report that discrepancy rather than claiming no-edit equality.
