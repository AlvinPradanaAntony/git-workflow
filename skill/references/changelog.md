# Changelog format

Use the structure observed in [TonzToon Komik CHANGELOG.md](https://github.com/AlvinPradanaAntony/tonztoon_komik-vibecode/blob/main/CHANGELOG.md) and [build-release.yml](https://github.com/AlvinPradanaAntony/tonztoon_komik-vibecode/blob/main/.github/workflows/build-release.yml), inspected 2026-10-01. This is a reusable structure; replace product name, description, versions, dates and release links with verified project values.

## File contract

- Reuse an existing `CHANGELOG.md` or case-equivalent `changelog.md`; do not create a second file on a case-sensitive system. If absent, create root `CHANGELOG.md` unless the chosen component has a documented changelog location.
- Start with `# Changelog`, a blank line, `Semua perubahan penting pada proyek **<Nama Proyek>** akan didokumentasikan di dalam file ini.`, then a short factual Indonesian description of the product. Include only verified technologies/capabilities; omit the description when unknown.
- Separate the introduction and release entries with `---`. Latest entry uses `## [X.Y.Z] - YYYY-MM-DD`, without a `v` prefix. Use the selected prerelease version when applicable. Date is the user's release date/timezone; do not invent a past date or call an unreleased draft published. If the repo has an Unreleased section, preserve it until the approved changes are assigned to a version.
- Categories use English headings, in order where present: `### Added`, `### Changed`, `### Fixed`, `### Removed`, `### Security`, `### Notes`. Use Indonesian `- ` bullets describing meaningful user/developer impact from the actual baseline diff. Omit empty categories. Initial Release/Technical headings from the reference may be used for an initial entry only when useful and supported.
- Consolidate repetitive commits into one factual item; exclude version bumps, release housekeeping and generated noise from the feature list. Do not copy raw commit subjects or list every touched file. Security items need evidence of a security improvement. Notes can cover real incompatibilities/migrations and limitations without inventing migration instructions.
- Prepend a new version entry, preserving all older entries, links, comments and existing `<details>` history. For an unpublished entry for the same version, merge deduplicated verified changes in place. For a published version, do not rewrite its content; choose a new version. A different existing layout needs a reviewed migration choice to adopt this requested format, with a preserve/defer option.

## New-file example

Illustrative content only; do not reuse these capabilities without diff evidence:

```markdown
# Changelog

Semua perubahan penting pada proyek **Nama Proyek** akan didokumentasikan di dalam file ini.

**Nama Proyek** adalah aplikasi desktop untuk mengelola kegiatan tim.

---

## [1.1.0] - 2026-10-01

### Added
- Filter status tugas pada daftar pekerjaan untuk memudahkan pemantauan tim.

### Fixed
- Memperbaiki pengiriman formulir berulang saat permintaan masih berlangsung.

---
```

Keep newest entries expanded. Preserve or optionally collapse older history using the exact `<details>` / `<summary><strong>Riwayat versi sebelumnya</strong></summary>` form below, with `###` version headings and `####` category headings. Never delete history to make it fit the release body. Do not manufacture a previous-version section when no verified predecessor is available.

## GitHub release notes

Extract the **exact approved version**, not whatever is currently first in the file. Include its categories and optionally **one** verified preceding released version in a collapsed history block. The reference file contains more history than its writing comment suggests; keep that history in the file while limiting the release body to one predecessor. Use `scripts/release_notes.py` to extract entries safely; it does not edit the changelog or publish anything. Verify the previous version/tag and URL before passing `--previous` and `--previous-url`.

```markdown
<details>
<summary><strong>Riwayat versi sebelumnya</strong></summary>

### [1.0.0](https://github.com/OWNER/REPO/releases/tag/v1.0.0) - 2026-09-20

#### Added
- Fitur awal proyek yang telah diverifikasi.

</details>
```

Release body follows the reference's `## 📋 Apa yang Baru di <tag>?`, extracted notes, `---`, then a required `## 📥 Download` section with `Platform | File | Keterangan` table. Include only the selected product outputs and exact verified filenames, package/container references or tagged-source links; update the proposed names from actual workflow outputs before reporting success. State real installation, signing and distribution limits. Missing promised outputs block publication. A source/script release may link its verified exact-tag archive instead of an app binary; a package release may link its verified package distribution. Do not fabricate native targets or empty download rows. Add actual build identifier and commit SHA when known. Keep changelog in Indonesian independently of commit-message language.
