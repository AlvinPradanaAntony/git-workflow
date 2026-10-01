# Commit message contract

## Options
`--lang en|id|both` (en); `--format short|standard|detailed` (accept `standar` as an alias); `--scope <name>`; `--type <type>`; `--issue <reference>`. `commitpln` defaults to standard. For `commitmsg` without `--format`, use short for exactly one changed file in the repository and standard otherwise, as detailed in `commitmsg.md`. An explicit format takes precedence.
Types: feat, fix, refactor, perf, docs, test, build, ci, chore, revert, plus documented repository types. Infer the most accurate type/scope from the actual diff. Explain a requested type if it conflicts with the change instead of falsely labeling it.
Keep English type tokens in Indonesian messages. Optional scope denotes a functional area, not a random file path. Use an imperative/action-oriented description; target a subject of at most 72 characters unless the repository imposes a different rule. No period at subject end by default. Never pad bullets to reach a quota.

Derive the subject from the selected diff as one meaningful summary of its purpose or outcome. When several related changes are included, choose an umbrella description covering all material bullets; do not list files or cram every detail into the subject. If the selection spans unrelated purposes, propose separate commits with `commitpln` instead of a vague all-purpose heading. Each body bullet must add a distinct, verifiable implementation or behavior detail; do not restate or translate the subject into a bullet, and do not paraphrase one bullet in another. Honor explicitly selected formats without padding with invented or repeated bullets. Do not invent intent beyond the diff or repository context.

## Short
Subject only in the selected language. No body unless a required trailer needs explanation.
```text
fix(auth): prevent duplicate login submissions
```

## Standard (approved exact visual style)
Subject, ONE blank line, then factual `- ` bullets. No paragraph introduction, `Changes:`, `Perubahan:`, Markdown headings, emoji, test checklist or artificial sections. No fixed number of bullets. Start English bullets with action verbs; preserve identifiers as written in code. Avoid line wrapping short bullets unnecessarily. Bullets explain how the summary was achieved, not repeat its wording.
```text
fix(auth): prevent duplicate login submissions

- Guard the submit handler while the request is pending
- Disable the login button and display loading feedback
- Restore the button state when authentication finishes
```
Only include those claims if the diff supports them. The example is not boilerplate to copy into other commits.

For a one-file `.gitignore` change, do not enumerate ignored filenames merely to fill the body. With no explicit format and no other changed file in the repository, use short:
```text
chore(git): reduce accidental tracking of local setup files
```
When the ignored patterns actually target sensitive local credentials, a security-focused subject can describe the reduced risk, such as `chore(git): reduce accidental commits of local credentials`. Do not claim a security benefit solely because `.gitignore` changed; check the patterns and whether the relevant files were already tracked. `.gitignore` does not stop tracking existing indexed files or guarantee secrets cannot be committed.
For multiple related changes, use a heading that encompasses the distinct items:
```text
feat(lessons): streamline publishing and progress tracking

- Add a draft preview before lessons are published
- Record progress when students complete a lesson
- Show completion totals on the student dashboard
```

## Detailed
Subject, blank line, concise evidence-supported motivation, blank line, `Changes:` and bullets, then an optional behavior/impact paragraph. Indonesian uses `Perubahan:`. Include validation/migration details only when relevant and known; omit unsupported sections. Detailed does not mean verbose for a tiny change.

## Bilingual
Produce one complete EN subject/body, ONE blank line, the literal separator `---`, ONE blank line, then the corresponding ID subject/body. Match the selected format in both halves; preserve semantic equivalence. The first EN line is Git's single subject; the second subject is body text. Do not create two commits, put `[EN]`/`[ID]` headers, or translate code identifiers.
```text
fix(auth): prevent duplicate login submissions

- Guard the submit handler while the request is pending
- Restore the button state when authentication finishes

---

fix(auth): cegah pengiriman login berulang

- Blokir handler pengiriman selama permintaan berlangsung
- Pulihkan status tombol saat autentikasi selesai
```

## Trailers and breaking changes
Place shared machine-readable trailers ONCE, after both language blocks and one blank line. Default `Refs: <issue>` when supplied. Use `Closes: <issue>` only with explicit issue-closing intent; never infer it merely from `--issue`.
Use `!` for verified breaking changes and add `BREAKING CHANGE: <actual incompatibility and migration guidance>` at the end for this package. For both languages, keep one token with EN / ID explanation if practical; continuation lines are indented. Do not invent an upgrade path when unknown; say what consumers must investigate. Other footer tokens remain untranslated.
Store the exact UTF-8 message in a local temporary file and invoke Git with `commit -F <path>` (or the host's structured equivalent). Do not build multiline messages with shell interpolation or literal backslash-n sequences. Verify the stored message and commit diff after success.
