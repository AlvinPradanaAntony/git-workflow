# Project version discovery and selection

## Inspect before choosing a version

Read release-targets.md. Identify the selected non-web product and its actual release boundaries. Verify entrypoints or package APIs, version authority and distribution policy; OS/architectures and native packaging apply only when relevant. Read project instructions, manifests, native platform version fields, lockfiles, VERSION files, app-info displays, changelog, reachable tags, GitHub releases and relevant Actions. Read locally or through GitHub tools; never send repository code or diffs to Context7. Public documentation queries must be generic.

| Project | Candidate sources to verify, not a blanket edit list |
| --- | --- |
| Flutter mobile / desktop | App `pubspec.yaml` version and build suffix, verified native target folders/build configuration; determine whether Gradle/Xcode or desktop packaging consume this value or override it. |
| Android native | App Gradle versionName/versionCode, convention plugin or version catalog when it is the authority. |
| iOS / macOS native | Xcode MARKETING_VERSION/CURRENT_PROJECT_VERSION, xcconfig or Info.plist according to build settings. |
| Electron desktop | Selected package.json version, matching lockfile package version where required, main entrypoint, app-info source and installer/portable packaging config. |
| Rust desktop / Tauri | Cargo.toml package/workspace version, Cargo.lock, desktop entrypoint, Tauri config if independent, and bundle targets; follow actual ownership. |
| .NET desktop | Central props/project Version, AssemblyVersion/FileVersion policies and installer metadata. |
| Go CLI/tool/service | VERSION or the established embedded version authority, Go entrypoints, build flags and ldflags/version injection; go.mod go/toolchain directives are toolchain versions, not product versions. |
| Python CLI/script/library | pyproject.toml/project version, setup.cfg/setup.py or established __version__/VERSION; verified package/bundler entrypoints and lockfile relationships. |
| JS/TS CLI or library | Selected package.json version, package exports/bin and appropriate lockfile root metadata; package.json alone does not imply a web product. |
| Rust CLI/library/service | Cargo package/workspace version and dependent lockfile metadata; bin/lib entrypoints and actual target policy. |
| Service/firmware/plugin/other non-web product | Established manifest/VERSION/embedded authority and output configuration; inspect release boundaries and project conventions rather than imposing an app schema. |
| Script/source-only product | Existing VERSION/header/tag-based version policy; no native build number or compiled executable is universally required. |

In a monorepo, select the non-web release component and distinguish unified app versions from independently released products. Ask for a component when scope is ambiguous; update only that product and its proven derived version fields. Follow links between authoritative and derived fields. Do not update dependency constraints or arbitrary text matching the old version. Use repo-prescribed package tools for lockfiles; do not run commands that also publish packages or create tags as a hidden version update.

Compare the current product version, changelog entries and published versions. List mismatches with file/value evidence; offer reconciliation choices instead of choosing the largest value automatically. A manifest already ahead of the latest release may be a prepared version: recommend finishing that version when its intended changes match, rather than bumping it twice. No remote access means remote publication state remains unknown.

## Baseline and evidence

- For continuation, find the previous released tag relevant to this component and target branch. Verify its peeled commit and ancestry. Use the selected baseline SHA..target SHA full diff and commit history, including merge/squash changes. Use an explicit `--from` only after verifying its relevance. Do not assume the numerically highest or most recent tag belongs to this product/branch.
- For initial setup without a relevant released tag, use current verified codebase capability and history as the initial entry. Do not pretend all past commits are new changes or reconstruct unverifiable release dates. Missing/shallow objects block reliable baseline inference; a scoped fetch may be chosen during preparation/publication when needed and authorized, never silently for an explicit read-only request.
- Separate committed release content from pending staged/unstaged/untracked work. Pending work may inform a preparation proposal, but never describe it as published or tag an uncommitted state. If there is no meaningful change, recommend retaining the version or deferring.

## Version recommendation

Use SemVer when the repo follows it: incompatible public API/data/config requirements → major; compatible new capabilities → minor; compatible fixes → patch. Conventional Commit types are evidence, not the sole decision: inspect actual compatibility and behavior. For 0.y.z or non-SemVer schemes, follow or ask about the existing policy; do not automatically claim 1.0.0 stability. For an unversioned new product, explain 0.1.0 for early development versus 1.0.0 for an explicitly stable contract, plus retaining the current state.

Validate prerelease identifiers and numeric ordering. Classify every SemVer suffix, including alpha/beta/rc, as prerelease; do not check only beta/rc as the reference workflow does. Promotion from prerelease to stable needs a deliberate choice. Build metadata does not raise SemVer precedence. Preserve the repo's tag naming scheme; validate tags and ensure uniqueness locally and remotely.

When the product has native platform build counters, keep them separate from its marketing version. Do not invent Android/iOS counters for CLI, libraries, scripts or services. Inspect prior published versionCode/CFBundleVersion or CI allocation; choose a monotonically valid build number only with evidence. Flutter +N is a build suffix, not automatically the GitHub tag suffix. If CI overrides N with a run number, report that authority instead of incrementing both blindly. Do not reset native counters to 1 on a marketing-version bump. Platform constraints may forbid arbitrary SemVer suffixes; adapt verified marketing/native fields without claiming an invalid native version is publishable.

## Update and validation

Apply the chosen initial or continued version to the exact authoritative files and necessary derived files; preserve unrelated formatting and local changes. Generate the changelog using `changelog.md`. Check that version fields, tag core, changelog header, release notes and artifact versions agree. Run only required scoped tests/builds/lint and accurately label missing platform runners, signing, credentials or unrun validation. Never create signing keys, upload credentials or mutate store/deployment settings implicitly.

Sources: [SemVer](https://semver.org/), [Android versioning](https://developer.android.com/studio/publish/versioning), [Flutter app versioning](https://docs.flutter.dev/deployment/android#update-the-apps-version-number). Check official platform docs for the detected project when needed.
