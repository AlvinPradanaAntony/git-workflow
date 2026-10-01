# Release product, build and distribution discovery

Use for every gitrelease action after validating its arguments. This contract replaces the former mobile/desktop-only restriction. Support every non-web project type; absence of a GUI, executable, native folders or compiler is not an exclusion.

## Product boundary

- Exclude a website or browser-only web application as the selected released product, including a component released solely as that website's deployment. Explain the evidence and stop that component without edits/tag/publication. Do not invent a desktop wrapper to evade the exclusion.
- Support CLI/tools/installers, scripts/automation, libraries/packages/SDKs, plugins/extensions, independently distributed services/backends, firmware/embedded programs, mobile/desktop apps, and source-distributed projects. This list is illustrative, not a whitelist that excludes other non-web products.
- Classify by what is distributed and how consumers use it: entrypoints, exports/bin, manifests, README, build scripts, existing releases and CI. A language, package.json, HTML assets, HTTP endpoints or a web framework dependency alone does not establish a web product. Electron/Tauri are desktop; a reusable library or independently packaged server is not automatically a website.
- In a mixed repo select the requested non-web component and its version boundary. A shared product version may cover related backend/frontend only when established by evidence. Ask if the request could mean the web deployment or a separately distributed service/package; do not guess either an exclusion or inclusion.
- A bare gitrelease shows usage; help is documentation-only. An explicit preview remains read-only. Classification does not authorize unrelated builds, deployments or registry publication.

## Inspect the codebase before selecting build or assets

Read applicable AGENTS.md, manifests/lockfiles, VERSION/version code, actual entrypoints/APIs, build/package scripts, toolchain pins, existing workflows and historical output names. Derive one distribution contract for the selected release SHA: mode, required outputs/links, applicable targets, build/test/package commands, runners, version injection and applicable signing/runtime constraints. Existing consistent conventions or explicit arguments resolve routine choices without another questionnaire.

| Evidence/product | Appropriate outputs and checks |
| --- | --- |
| Go CLI/tool | Executable or portable archive for selected GOOS/GOARCH; inspect main package, build tags, version injection and CGO requirements before using cross-build. No GUI/DMG requirement for a macOS CLI. |
| Python script/CLI | Versioned .py/source or package where Python is a documented prerequisite; bundled executable only when the chosen product policy calls for it. Respect an existing bundler; do not choose PyInstaller automatically. |
| JS/TS CLI/library | Package tarball/source and exports/bin, declared Node/runtime and tested entrypoint. Browser framework presence does not exclude a reusable package. |
| Library/package/SDK | Build the project's wheel/sdist, crate, JAR, nupkg, package archive or verified tagged source as applicable. No OS matrix or native executable unless the library actually needs it. |
| Service/backend | Existing binary, container image reference or source/config archive policy; verify runtime/entrypoint/version. GitHub release publication does not deploy the service or push a registry image implicitly. |
| Mobile/desktop app | Selected actual platform packages and signing/distribution checks from release-publish.md; a missing promised APK/IPA/installer cannot be replaced by source. |
| Plugin/firmware/other | Existing consumer-specific package/image/source and platform/version checks. Ask about genuinely unknown format or target requirements. |
| Source/scripts/config/data distribution | Exact tagged distributable files and applicable validation. No compile step is necessary if the product has none. A git archive or GitHub automatic tag archive is valid when complete and selected by evidence/user choice. |

The table gives examples, not mandatory universal formats. Never force Windows/Linux/macOS, amd64/arm64, portable/installer pairs, APK/IPA, signing, Go or Python onto unrelated codebases. Missing build setup goes to init or a scoped authorized fix, not an out-of-scope rejection. Do not call source distribution a fallback for a failed required binary build.

## Ask only questions that change the release

Ask when the selected component/product, format, supported targets, signing method, compilation versus source/runtime distribution, version authority or competing publication route cannot be determined. Give actual evidence, viable alternatives, benefit/consequence/exact scope and one recommendation under gitrelease.md. Keep dependent work pending, and continue independent authorized work. Examples: Python-required script versus standalone executable when neither is established; tag-source distribution versus wheel/sdist for an unconfigured library; or API service package versus a web deployment in a mixed repository.

Do not ask merely because the project is CLI/library, no compiler exists, all frameworks support several OSes, or the same scope was already authorized. If version/changelog/CI consistently describe Go binaries for six targets plus a Python installer, use those exact outputs. Git Workflow's distribution repo is such a CLI/tool product and is eligible; verify its current files instead of rejecting it for lacking a GUI.

## Action and build boundaries

- init prepares the concrete local metadata/workflow/build/package helpers according to that contract, with no commit/tag/push/dispatch. Source mode gets applicable validation/package or tag-source verification steps rather than fake executable jobs.
- prepare inspects changes, determines version and updates metadata/changelog; check readiness without rebuilding/reconfiguring the product on every invocation. Missing external credentials or not-yet-built files do not block clear local metadata updates.
- An explicit publish authorizes the necessary existing scoped non-publishing tests/build/package commands for its approved SHA and output contract. Build in a clean isolated checkout if needed. Use existing pinned dependencies and project commands; do not install global toolchains, add packaging frameworks, change platforms, create keys or write secrets implicitly. Ask about unresolved material changes before dependent builds.
- Workflow route verifies the exact required validation/build/package gates before triggering, then verifies all actual outputs. Direct route builds/packages when required, verifies the results before remote writes and uses the selected source policy when there is no build. An unsupported environment or missing auth is a concrete blocker, not a reason to reject the product or improvise another publisher.
- Registry/container/store publication and service/web deployment are distinct side effects. GitHub Release authorization covers GitHub notes/tag/assets only unless the existing reviewed workflow and user's explicit scope also cover those operations.
- Verify complete required outputs, checksums/provenance/versions where applicable, and honest notes/download links. Source-only mode can have no custom attached binary, but must have verified deliverables in the exact tagged source/package. Automatic GitHub source archives may omit submodule content; inspect completeness and package any required missing content rather than claiming it is present.
- On missing promised assets, failed validators/build, ambiguous targets or auth failure, stop only the dependent stage and report what is complete and what remains. Never weaken branch/tag protection, force-push a branch, replace tags or publish incomplete targets automatically.
