# prdesc

Read `references/common.md` before acting. Follow the user's explicit scope and this command's mutation boundary.

Accept `--base <ref>` and `--lang en|id|both`. Determine the intended base from explicit input or verified repository context; ask when ambiguous, never assume main. Verify the base object exists; do not fetch in this read-only command. Compare commits and the full `base...HEAD` diff after resolving refs. Mention uncommitted work separately and exclude it from the committed PR draft. Write title, purpose, important changes, actual validation results, and breaking/migration notes when relevant. State validation not run if applicable. For both, use two equivalent sections separated by `---`. No fabricated issue closures, PR URL or test results. Do not invoke gh pr create, push, email or messaging.
