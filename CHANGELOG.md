# Changelog

## [2.13.0] - 2026-10-01

### Added
- CLI mandiri `git-workflow` dengan `init`, `update`, `uninstall`, `doctor`, `help`, dan `--version`.
- Executable Windows, Linux, dan macOS untuk amd64 serta arm64; pengguna tidak perlu Python, Node, atau Go untuk menjalankan CLI.
- Bootstrap PowerShell/shell dengan deteksi platform, verifikasi SHA-256, pemasangan ke PATH pengguna, dan backup executable sebelumnya.
- Workflow GitHub Actions untuk pengujian native enam target, build seluruh aset, dan publikasi release otomatis pada push tag versi yang cocok.
- Release notes dari entri versi pada changelog, tabel unduhan, installer Python mandiri, dan checksum seluruh aset.

### Changed
- Root proyek ditemukan dari folder terminal, termasuk subfolder dan linked worktree; installer tidak perlu dipindahkan ke repo.
- CLI dan installer Python menggunakan sumber skill yang sama. `update` memperbarui skill dari executable saat ini; pembaruan CLI dilakukan melalui bootstrap.
- Manifest pemasangan lama dapat dimigrasikan. File pengguna dan isi AGENTS.md di luar blok terkelola dipertahankan; penggantian file yang diedit memerlukan `--replace` dan backup.
- Pemasangan baru mencatat separator pointer AGENTS.md agar uninstall mengembalikan byte asli, termasuk file yang sebelumnya kosong, tanpa menghapus aturan pengguna yang ditambahkan kemudian.

## [2.12.3] - 2026-09-30

### Changed
- `help --cmd` tanpa nilai menampilkan 12 nama command dengan numbering.
- `help --cmd NUMBER_OR_NAME` tetap menampilkan dokumentasi command yang dipilih.
