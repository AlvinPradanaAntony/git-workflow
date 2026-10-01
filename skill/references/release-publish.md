# GitHub release publication

For --replace-existing, apply release-manage.md before ordinary publication. Default publish never deletes a Release/tag. Each invocation performs checked stages; never blindly execute a semicolon-separated shell sequence.

## Discover and choose one publisher

Read the project's actual workflow files, trigger filters, version/notes extraction, required checks, artifacts, permissions and environment approvals. Prefer an existing matching release workflow when it can fulfill the selected stable/prerelease/draft mode. Do not copy Flutter runners, secrets or APK/IPA filenames into other projects. init prepares local automation under release-init.md; before triggering it, require the concrete committed diff and established scope to have been reviewed. Ask only for unresolved material changes or missing authorization, not duplicate routine permission.

Every helper referenced by CI must exist in the committed project at the release SHA. A globally installed skill path is unavailable on an Actions runner. Reuse an existing project extraction helper, or propose a reviewed copy of this skill's standalone `scripts/release_notes.py` into the project's scripts directory; include that exact file in the approved workflow/preparation scope. Do not add a CI command that points to an untracked or ignored local skill file.

| Product | Verify artifacts and distribution scope |
| --- | --- |
| Android | Downloadable APK with verified release signing. AAB may be an additional store artifact, but is not a directly installable download and cannot be the only app distribution asset. |
| iOS | IPA with a verified supported signing/distribution method; explain device/distribution limits. An unsigned package does not satisfy the installable-asset requirement. |
| Windows | Verified installer or portable app, such as EXE/MSI or a packaged portable ZIP; actual architecture, signing requirements and updater metadata. |
| Linux | Verified AppImage/DEB/RPM or portable application archive for selected architectures. |
| macOS | Verified DMG/PKG or packaged app ZIP, with the signing/notarization requirements of the chosen distribution method. |

Select only OS/architectures the app actually supports and the user chose. At least one installable/portable app asset is required; every selected target must have its expected asset. Source-code archives, symbols and checksums alone do not meet this requirement. If a target fails, stop publication or obtain a reviewed decision to release a smaller supported set; never advertise an unbuilt platform.

Plan a **workflow** route when Actions owns building/publishing the selected app packages. Inspect its build outputs, upload/download steps and existence checks; require publication to fail before creating/publishing a release if required app assets are missing. Offer a reviewed pipeline fix when those checks are absent. Plan a **direct** route when no workflow owns this release and the selected installable/portable app files already exist with verified provenance. If both publishers exist, choose one; never race Actions with `gh release create` for the same tag. If neither route is ready, recommend preparing the app build or deferring.

Drafts still write to GitHub. A push-triggered workflow with `draft: false` publishes on tag push: `--draft` cannot make that safe automatically. Offer a reviewed workflow change, direct draft only when the automatic publisher is deliberately disabled/bypassed, or defer. Check all workflow triggers, not just a selected file. SemVer prerelease tags must actually be marked prerelease by the chosen publisher; adapt a workflow that only detects beta/rc before using alpha or another suffix.

## Preflight the reviewed release

1. Resolve attached branch and exact owner/repo from the selected remote or explicitly chosen GitHub repository. Verify permissions, local/remote target branch, release SHA, tag naming, published release/tag state and branch/tag protection. Do not assume origin/main. Never print credential-bearing URLs.
2. Require approved code, version files and changelog committed at the release SHA. Pending unrelated work may remain local but must not enter that commit/artifacts. Build from a clean isolated checkout of the approved SHA when necessary; do not package dirty working-tree files. If release preparation needs a commit, first present exact files/diff and obtain the release-commit choice unless that exact scope is already authorized. Existing unrelated staged content needs isolation, not accidental inclusion. Follow gitrelease.md for specific questions and continue only independent authorized work while a required decision is pending.
3. Verify version/changelog/tag agreement, nonempty selected notes, baseline ancestry, validation/check results and app artifact provenance/versions. For direct publication, require nonempty existing approved app files for every selected target before creating a release or pushing its tag. For workflow-owned builds, verify the expected artifact map and build-to-publication gates before tag push, then wait for the required jobs and validate actual attached files before reporting success. If the pipeline can publish without its required app assets, fix that reviewed gate before triggering it. Do not upload source env files, signing keys or arbitrary directory contents as assets.
4. Show the concrete branch push range if needed and follow `push.md` for a normal targeted branch push. Tag creation/push is separately scoped to this release; no --tags, --all, --mirror, branch force-push, retagging or automatic deletion. Explicit deletion/replacement uses only the narrowly reviewed tag operations in release-manage.md. Validate and create an annotated tag at the approved SHA only after the publication decision. Recheck remote tags before creating/pushing. A matching existing tag may be reused to resume; a different target is a blocker for normal publication.
5. A fully specified prior authorization for this release can cover its exact branch/tag push and chosen publication route. Otherwise obtain the missing concrete decision after the files/notes are ready. Do not stop at an abstract permission request while preparation is unfinished.

## Workflow route

- Push only the approved branch if required and then the single approved tag, or invoke the existing workflow_dispatch route at the exact ref with its verified inputs. Explain that tag push may immediately publish. Do not dispatch arbitrarily after a tag-triggered run already started.
- Check trigger matching, workflow ref, token origin and permissions. A tag pushed using a workflow's GITHUB_TOKEN normally does not cause a new push-triggered workflow; use a supported explicitly authorized trigger/auth route, not an unexplained retry or secret creation.
- Observe the run corresponding to the tag and SHA, its required jobs, artifacts and final release metadata. Report queued/running state when observation is limited. On a failed build, retain the published tag and fix forward or retry the relevant job after inspecting what already exists. Never move the tag to new code silently.

## Direct route

- Prefer available supported GitHub plugin tools for release operations. If they cannot create releases, use authenticated `gh` or the official API in the user's environment when available. Do not assume a connected ChatGPT plugin grants the local CLI credentials. If access is unavailable, leave the ready files/notes and report the concrete blocker.
- Require the verified app asset list, then verify the tag exists remotely at the reviewed SHA. Use structured argv such as `gh release create <tag> <approved-app-files...> --repo <owner/repo> --verify-tag --title <title> --notes-file <exact-file>`; add --draft or --prerelease as selected. --verify-tag prevents GitHub from creating an unreviewed tag at the default branch. Do not use generated notes in place of the required changelog layout. Drafts also require the selected app assets; an upload failure is partial completion, not success.
- On an existing matching release, verify draft/prerelease status, body and assets before choosing resume, publish draft, retain or defer. No overwriting notes/assets, --clobber or editing published content without an explicit concrete update request. A prerelease must not be marked Latest stable.
- Verify the resulting tag target, release URL/status, body and exact asset names. Report partial completion if upload or verification fails; retry only the unfinished authorized stage.

Sources: [gh release create](https://cli.github.com/manual/gh_release_create), [GitHub workflow triggers](https://docs.github.com/en/actions/how-tos/write-workflows/choose-when-workflows-run/trigger-a-workflow), [release events](https://docs.github.com/en/actions/reference/workflows-and-actions/events-that-trigger-workflows). Recheck official docs when implementing a new project pipeline.
