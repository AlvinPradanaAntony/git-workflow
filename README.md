# Git Workflow

CLI mandiri untuk memasang satu skill berisi **12 command Git + help**, dengan executable Windows, Linux, dan macOS. Jalankan `git-workflow init` dari proyek mana pun; root Git ditemukan dari folder terminal, termasuk subfolder dan linked worktree. Pengguna CLI tidak memerlukan Python, Node, atau Go.

Versi distribusi: **2.13.0**. Source: [AlvinPradanaAntony/git-workflow](https://github.com/AlvinPradanaAntony/git-workflow).

## Pasang CLI sekali

Unduh aset dari [GitHub Releases](https://github.com/AlvinPradanaAntony/git-workflow/releases). Workflow akan menyediakan unduhan setelah tag versi pertama dipush. Sebelum itu, paket build tersedia sebagai artifact `git-workflow-downloads` pada [GitHub Actions](https://github.com/AlvinPradanaAntony/git-workflow/actions).

### Windows

Unduh `install.ps1`, lalu jalankan dari PowerShell:

```powershell
.\install.ps1
git-workflow --version
```

Bootstrap mengambil executable sesuai CPU dan `SHA256SUMS` dari repo distribusi, memverifikasinya, lalu memasang ke `%LOCALAPPDATA%\Programs\GitWorkflow\bin`. PATH pengguna dan sesi PowerShell aktif diperbarui, tanpa administrator.

Untuk versi tertentu: `.\install.ps1 -Version v2.13.0`. Untuk pemasangan offline, simpan executable yang cocok dan `SHA256SUMS` di folder yang sama dengan `install.ps1`, atau gunakan `-SourceDirectory PATH`. `-BinDirectory PATH`, `-NoPath`, dan `-Replace` tersedia.

### Linux/macOS

Unduh `install.sh`, lalu jalankan:

```bash
sh install.sh
```

Bootstrap memasang executable ke `~/.local/bin` dan menambahkan PATH pada profil shell pengguna. Buka terminal baru atau jalankan baris `export` yang ditampilkan, kemudian:

```bash
git-workflow --version
```

Untuk versi tertentu: `sh install.sh --version v2.13.0`. Pemasangan offline memakai executable yang cocok dan `SHA256SUMS` di sebelah script, atau `--source-dir PATH`. `--bin-dir PATH`, `--no-path`, dan `--replace` tersedia.

Kedua bootstrap memakai repo ini sebagai default. `--repo OWNER/REPO` / `-Repo OWNER/REPO` atau `GIT_WORKFLOW_REPO` dapat mengganti sumber unduhan. Bootstrap memasang **CLI**; lanjutkan dengan `init` untuk memasang **skill**.

## Pasang skill ke proyek

Dari root atau subfolder repo:

```text
git-workflow init
git-workflow doctor
```

`init` langsung memasang. `--dry-run` menampilkan rencana sebelum menulis. `--project PATH` memilih repo lain dari lokasi mana pun, tanpa memindahkan installer.

```text
git-workflow init --dry-run
git-workflow init --project D:\Projects\my-app
```

Pemasangan proyek membuat satu skill bersama di `.agents/skills/git-workflow`, rules di `.agents/rules/git-workflow.md`, dan pointer terkelola dalam root `AGENTS.md`. File lain serta aturan tim tetap dipertahankan.

## Pemasangan global

```text
git-workflow init -g --agent codex
git-workflow init -g --agent antigravity
git-workflow init -g --agent antigravity-cli
git-workflow init -g --agent all
```

Tanpa `--agent`, CLI memakai pilihan pemasangan sebelumnya atau satu konfigurasi agent yang terdeteksi jelas. Jika belum jelas, CLI meminta pilihan melalui terminal; pemakaian noninteraktif perlu `--agent`. Global tidak mengubah file repo atau global `AGENTS.md`.

## Command terminal

| Command | Fungsi |
| --- | --- |
| `init` | Pasang skill; pengulangan tidak membuat duplikat. |
| `update` | Perbarui pemasangan menggunakan skill yang tertanam dalam executable saat ini. |
| `uninstall` | Backup dan hapus hanya file yang tercatat sebagai milik paket; file pengguna dipertahankan. |
| `doctor` | Periksa versi, Git, PATH, manifest, checksum, dan pointer proyek; read-only. |
| `help` | Tampilkan dokumentasi pemasangan CLI. |
| `--version` | Tampilkan versi CLI dan skill yang tertanam. |

`-g/--global` memilih pengguna saat ini. `--agent` memilih integrasi. `--project` tidak dapat digabung dengan `-g`. `init/update/uninstall` mendukung `--dry-run` dan `--replace`.

`update` tidak mengunduh CLI terbaru. Jalankan bootstrap lagi atau ganti executable terlebih dahulu, lalu gunakan `git-workflow update` atau `git-workflow update -g` untuk cakupan skill yang dipilih.

## Installer Python dan migrasi

`install_git_workflow.py` juga tersedia sebagai aset unduhan mandiri, dengan skill yang sama seperti CLI. Python 3.10+ dan Git diperlukan untuk pemasangan proyek. Script dapat disimpan di folder mana pun:

```text
python /path/to/install_git_workflow.py
python /path/to/install_git_workflow.py --apply
python /path/to/install_git_workflow.py --global --apply
python /path/to/install_git_workflow.py --uninstall --apply
```

Python mempertahankan preview sebagai default; tambahkan `--apply` untuk menulis. `--global` memasang ketiga integrasi. `--project PATH`, `--replace`, `--help`, dan `--version` tersedia.

CLI mengenali manifest Python 2.x. Perubahan manual pada file terkelola menghentikan update/uninstall; `--replace` menyimpan backup sebelum penggantian. CLI menyimpan backup pada cache pengguna di `git-workflow/backups`; lokasinya ditampilkan. Tidak ada staging, commit, push, hook, atau instalasi dependensi saat memasang skill.

## Command agent

CLI terminal memasang skill; agent menjalankan command Git. Di Codex:

```text
$git-workflow help --cmd
$git-workflow gitrelease prepare
$git-workflow commitmsg
$git-workflow gitrelease publish
```

Di Antigravity gunakan `/git-workflow` sebagai awalan. `git-workflow init` di terminal memasang skill; `gitrelease init` melalui agent menyiapkan sistem rilis aplikasi.

Nomor command tetap: `1 gitstatus`, `2 commitpln`, `3 commitmsg`, `4 branchname`, `5 prdesc`, `6 gitundo`, `7 gitreset`, `8 gitmergecancel`, `9 gitamend`, `10 gitpushamend`, `11 gitpush`, `12 gitrelease`. Help ringkas, `--detailed`, `--cmd`, dan pemilihan nomor/nama tetap tersedia.

## Aset unduhan

| Platform | amd64 | arm64 |
| --- | --- | --- |
| Windows | `git-workflow-windows-amd64.exe` | `git-workflow-windows-arm64.exe` |
| Linux | `git-workflow-linux-amd64` | `git-workflow-linux-arm64` |
| macOS | `git-workflow-darwin-amd64` | `git-workflow-darwin-arm64` |

Release juga menyertakan `install_git_workflow.py`, `install.sh`, `install.ps1`, arsip source, dan `SHA256SUMS`. Executable merupakan CLI portable. `darwin` adalah nama target Go untuk macOS.

## Build dan publish melalui workflow

[`.github/workflows/release.yml`](.github/workflows/release.yml) menguji CLI, migrasi Python, dan bootstrap pada native Windows/Linux/macOS, masing-masing amd64 dan arm64. Build menghasilkan seluruh aset sebagai artifact pada push ke `main`, pull request, atau run manual.

Push tag `vX.Y.Z` yang sesuai dengan `VERSION` dan changelog memicu publikasi setelah semua pengujian dan pemeriksaan aset berhasil. Release dibuat sebagai draft, aset diunggah dan diperiksa, lalu release dipublikasikan. Release notes mengambil hanya entri changelog versi tersebut dan menambahkan tabel unduhan. Release/tag yang sudah ada tidak ditimpa otomatis.

Untuk publikasi pertama setelah kode berada di `main`:

```bash
git tag -a v2.13.0 -m "Release v2.13.0"
git push origin main
git push origin refs/tags/v2.13.0
```

Untuk versi berikutnya, perbarui `VERSION` dan `CHANGELOG.md`, regenerasi paket, lalu commit sebelum membuat tag. Run manual dari branch hanya membangun; run manual dari tag yang cocok juga mempublikasikan release.

## Pengembangan

Build pengembang memerlukan Go 1.24+ dan Python 3.10+. Tidak ada modul Go pihak ketiga.

```bash
python scripts/package.py
python scripts/package.py --check
python -m unittest discover -s tests -v
go test ./...
go vet ./...
python scripts/build.py
python scripts/release.py finalize
```

`skill/` memuat sumber skill. `scripts/installer_template.py` memuat installer Python. `scripts/package.py` menghasilkan `bundle.json` untuk `go:embed` dan `install_git_workflow.py` dari sumber yang sama. Jangan mengedit file hasil generasi secara terpisah.
