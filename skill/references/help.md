# Help modes and selection contract

Treat help as documentation only. Never inspect a repo, run Git, tests or the installer, or perform mutations. Read this reference; read matching command references only for additional explanation, never to execute the selected command. Reply in Indonesian. Adapt prefixes to `$git-workflow` for Codex, `/git-workflow` for Antigravity, or command text after selecting Git Workflow in ChatGPT. Parameters here belong to the agent prompt, not the installer.

## Modes and validation

- `help`: render the compact catalog below, with separate **No.**, **Command**, **Deskripsi**, **Opsi / nilai** columns. Keep description short and place flags/values only in the options column. Include shared message parameters once, then the help usage examples. Do not expand the detailed catalog by default.
- `help --detailed` or `help detailed`: render all detailed command sections below, with stable numbered headings and separate description/actions/parameter explanations. Help itself remains unnumbered. Include shared message options once. Do not duplicate the compact table.
- `help --cmd` without a value: render only the 12 Git command names as a numbered list using the compact catalog's canonical 1–12 mapping. Do not include descriptions, parameters, shared message options or detailed sections. One short example such as `help --cmd 2` may follow. This is valid, not a missing-argument error. When combined with --detailed but no selector value, the command list takes precedence.
- `help --cmd NUMBER_OR_NAME`: render only that command's full detailed section. This implies detailed display even without --detailed. Include shared message parameters only when selecting commitpln/commitmsg. Retain the canonical number in the selected heading. Combining --cmd and --detailed still shows only the selected command.
- Numbers are the 12 Git command numbers in the compact catalog, starting with gitstatus=1 and commitpln=2, ending with gitrelease=12. Help has no numeric index, and never shifts the Git numbering. Accept `--cmd help` by name. Resolve names case-insensitively, allowing one leading slash; normalize commitplan to commitpln. Numeric selectors must be integers in 1–12, never row offsets or action names.
- Accept only one --cmd with an optional following target and the optional detailed flag/word. Treat a value as omitted at end of input or before a recognized --detailed flag; do not consume that flag as a command name. A word directly after --cmd is a target, so unknown names remain errors. Reject an unknown supplied target, invalid number, repeated selector, unexpected argument or unknown option before any repo action; give a brief error and usage such as `help --cmd 2` or `help --cmd commitpln`. Do not reinterpret a selector as an executable command or invent new help flags.
- A natural-language request for detail about one named command follows the same selected-detail mode. It is a documentation request, not a new Git action.

## Presentation

Do not print this contract, reference navigation or catalog source labels. Command-list output is only numbered command names and an optional short selection example. Compact output is a short introduction, the summary table, shared message options and examples. Detailed output uses named numbered headings, a standalone **Deskripsi** paragraph, separate **Aksi/mode** tables where applicable, and separate **Parameter** tables with values/defaults and brief explanations. Preserve blank lines around tables/headings. Do not enclose the whole response in a code block or recombine descriptions with flag chains. Do not reintroduce removed release actions. Explain required selectors/defaults and the actual effect without lengthy implementation rules.

## Reference navigation — not part of the response

- [Compact catalog](#compact-catalog--default-output-source), [help options](#help).
- [1. Status](#1-gitstatus), [2. Commit preview](#2-commitpln), [3. Commit](#3-commitmsg), [4. Branch name](#4-branchname), [5. PR draft](#5-prdesc).
- [6. Undo](#6-gitundo), [7. Reset](#7-gitreset), [8. Merge cancellation](#8-gitmergecancel).
- [9. Amend](#9-gitamend), [10. Amend and push](#10-gitpushamend), [11. Push](#11-gitpush), [12. Releases](#12-gitrelease).
- [Shared message parameters](#parameter-pesan-commit--commitpln-dan-commitmsg).

# Compact catalog — default output source

Git Workflow menyediakan 12 command Git. Help hanya menampilkan panduan; tidak menjalankan operasi Git.

| No. | Command | Deskripsi | Opsi / nilai |
| --- | --- | --- | --- |
| — | `help` | Tampilkan bantuan ringkas atau detail. | `--detailed` atau `detailed`; `--cmd` tanpa nilai: daftar bernomor; dengan nomor/nama: detail command. |
| 1 | `gitstatus` | Ringkasan status repo; read-only. | Tidak ada. |
| 2 | `commitpln` | Rencana pembagian commit dan draft pesan; read-only. | Opsi pesan bersama; `--source` auto/staged/unstaged/all/description; `--files` PATH...; `--all`. Alias: commitplan. |
| 3 | `commitmsg` | Susun pesan dan commit lokal; tanpa push. | Opsi pesan bersama; `--source` auto/staged/unstaged/all; `--files` PATH...; `--all`. |
| 4 | `branchname` | Sarankan nama branch; read-only. | `--issue` ID. |
| 5 | `prdesc` | Buat draft judul/deskripsi PR; tanpa membuat PR. | `--base` REF; `--lang` en/id/both. |
| 6 | `gitundo` | Tarik commit lokal menjadi staged/changes. | `--count` N (1); `--to` staged/changes (staged); `--parent` 1/2 untuk satu merge commit. |
| 7 | `gitreset` | Reset lokal ke revision yang dipilih. | Pilih satu: `--soft`, `--mixed`, `--hard`; wajib `--to` REV. |
| 8 | `gitmergecancel` | Hentikan merge aktif dengan keep atau abort. | Wajib pilih `--keep` atau `--abort`; `--to` staged/changes (changes), khusus keep. |
| 9 | `gitamend` | Amend staged changes; pesan tetap, tanpa push. | `--no-edit` (selalu default). |
| 10 | `gitpushamend` | Amend no-edit lalu push ke target terverifikasi. | `--remote` NAME; `--branch` NAME; `--no-edit` (default). |
| 11 | `gitpush` | Push commit yang sudah ada. | `--remote` NAME; `--branch` NAME. |
| 12 | `gitrelease` | Versioning/release semua proyek selain web; build/aset mengikuti codebase. | Aksi: init/prepare/publish/delete; `--version` VERSION atau `--bump` patch/minor/major; `--from` REF; `--component` PATH; `--build-number` N; `--remote` NAME; `--branch` NAME; `--tag` NAME; `--route` auto/workflow/direct (auto); `--workflow` PATH; `--draft`; `--assets` PATH... (direct); `--replace-existing` (publish); `--delete-scope` release/tag/both (both, khusus delete). |

**Opsi pesan bersama — commitpln dan commitmsg:**

| Opsi | Nilai / default |
| --- | --- |
| `--lang` | en / id / both; default en. Both dipisahkan `---`. |
| `--format` | short / standard / detailed; alias standard: standar. Commitpln default standard; commitmsg otomatis short untuk satu file koheren, selain itu standard. |
| `--scope` | NAME, area fungsional. |
| `--type` | TYPE, misalnya feat / fix / docs / chore sesuai diff. |
| `--issue` | REF, referensi issue terverifikasi. |

Source auto memakai staged jika tersedia, lalu perubahan kerja yang relevan. Tanpa target push, gunakan upstream terverifikasi. Delete hanya menerima --tag (wajib), --delete-scope, --remote dan --component. Rincian setiap opsi tersedia pada detail command.

**Contoh:** `help --detailed`, `help --cmd`, `help --cmd 2`, `help --cmd commitmsg`. Gunakan awalan yang sesuai host.

# Detailed catalog — render only when selected

Satu skill menyediakan 12 command Git. `help` menampilkan panduan tanpa menjalankan operasi Git. Parameter di bawah digunakan dalam pesan ke agent, bukan pada skrip installer.

## help

**Deskripsi:** Menampilkan command, aksi, parameter dan contoh penggunaan.

**Parameter:**

| Parameter | Nilai / default | Penjelasan singkat |
| --- | --- | --- |
| `--detailed` | Flag; default mati | Tampilkan bantuan lengkap dengan numbering. `help detailed` juga diterima. |
| `--cmd` | Nilai opsional: nomor 1–12 atau nama command | Tanpa nilai, tampilkan daftar nama command bernomor. Dengan nilai, tampilkan detail pilihan tersebut. Mendukung nama help dan alias commitplan; tidak menjalankan command tersebut. |

**Contoh:** `help`, `help --detailed`, `help --cmd`, `help --cmd 2`, `help --cmd commitmsg`.

## 1. gitstatus

**Deskripsi:** Menampilkan branch, perubahan staged/unstaged, file baru dan operasi Git yang sedang berlangsung. Hanya membaca.

**Parameter:** Tidak ada.

## 2. commitpln

**Deskripsi:** Menganalisis perubahan, menyarankan pembagian commit, lalu membuat draft pesan. Tidak melakukan staging atau commit. Alias teks: `commitplan`.

**Parameter:** Mendukung parameter pesan commit pada bagian bersama di bawah, dengan default format `standard`, serta:

| Parameter | Nilai / default | Penjelasan singkat |
| --- | --- | --- |
| `--source` | `auto`, `staged`, `unstaged`, `all`, `description`; default `auto` | Pilih sumber draft. Auto memakai staged jika tersedia, lalu perubahan kerja yang relevan. Description memakai uraian user sebagai draft. |
| `--files` | Satu atau beberapa path | Batasi draft pada file yang dipilih. |
| `--all` | Flag | Pilih seluruh perubahan nonignored; setara dengan source all. |

**Contoh:** `commitpln --lang both --source staged`

## 3. commitmsg

**Deskripsi:** Menyusun pesan Conventional Commit dari perubahan nyata dan membuat commit lokal. Tidak push. Jika staged tersedia, default memakai staged; jika belum ada, agent dapat men-stage satu pekerjaan yang jelas sebelum commit.

**Parameter:** Mendukung parameter pesan commit pada bagian bersama di bawah, serta:

| Parameter | Nilai / default | Penjelasan singkat |
| --- | --- | --- |
| `--source` | `auto`, `staged`, `unstaged`, `all`; default `auto` | Pilih sumber commit. Staged tidak memakai fallback; unstaged memerlukan pemisahan dari staged lain. Description tidak tersedia untuk commit. |
| `--files` | Satu atau beberapa path | Pilih file tertentu; staging parsial atau staged lain diperiksa sebelum commit. |
| `--all` | Flag | Pilih seluruh perubahan nonignored, termasuk hunks yang belum staged dan file baru yang layak. |

**Default format:** Tanpa `--format`, tepat satu file berubah di repo memilih `short` untuk satu perubahan yang koheren; selain itu `standard`. Format eksplisit selalu didahulukan. Selector yang bertentangan ditolak.

**Contoh:** `commitmsg --lang both --format standard`

## 4. branchname

**Deskripsi:** Menyarankan nama branch berdasarkan tugas atau perubahan dan pola repo. Tidak membuat atau berpindah branch.

**Parameter:**

| Parameter | Nilai / default | Penjelasan singkat |
| --- | --- | --- |
| `--issue` | ID issue; opsional | Sertakan ID issue yang diberikan atau terverifikasi pada nama branch. |

## 5. prdesc

**Deskripsi:** Membuat draft judul/deskripsi Pull Request dari commit dan diff terhadap base. Tidak membuat PR atau push.

**Parameter:**

| Parameter | Nilai / default | Penjelasan singkat |
| --- | --- | --- |
| `--base` | Ref branch/commit; opsional | Tentukan pembanding PR. Tanpa nilai, gunakan base yang terverifikasi; tanyakan jika ambigu. |
| `--lang` | `en`, `id`, `both`; default `en` | Pilih bahasa draft. Both menampilkan versi EN dan ID yang dipisahkan `---`. |

## 6. gitundo

**Deskripsi:** Mengembalikan commit lokal menjadi staged atau changes, sambil mempertahankan isi perubahan. Tidak push.

**Parameter:**

| Parameter | Nilai / default | Penjelasan singkat |
| --- | --- | --- |
| `--count` | Bilangan positif; default `1` | Jumlah commit lokal yang ingin ditarik kembali. |
| `--to` | `staged`, `changes`; default `staged` | Staged memakai soft reset; changes memakai mixed reset. Perubahan unstaged sebelumnya tidak otomatis menjadi staged. |
| `--parent` | `1` atau `2`; khusus satu merge commit | Pilih parent tujuan ketika membatalkan merge yang sudah di-commit. |

**Contoh:** `gitundo --count 1 --to changes`

## 7. gitreset

**Deskripsi:** Mengembalikan HEAD ke revision tertentu dengan mode eksplisit. Tidak push.

**Aksi / mode — wajib pilih satu:**

| Mode | Penjelasan singkat |
| --- | --- |
| `--soft` | Pindahkan HEAD, pertahankan index dan file kerja. |
| `--mixed` | Pindahkan HEAD dan reset index; isi file kerja tetap sebagai changes. |
| `--hard` | Reset HEAD, index dan file kerja sesuai target; perubahan lokal yang terdampak dapat hilang dan harus dipreservasi/ditinjau. |

**Parameter:**

| Parameter | Nilai / default | Penjelasan singkat |
| --- | --- | --- |
| `--to` | Revision Git; wajib | Tentukan commit tujuan yang diverifikasi, misalnya `HEAD~1`. Tidak ada mode atau target yang diasumsikan. |

**Contoh:** `gitreset --soft --to HEAD~1`

## 8. gitmergecancel

**Deskripsi:** Menghentikan merge yang masih berlangsung. Untuk merge yang sudah di-commit, gunakan gitundo dengan parent yang jelas.

**Aksi — wajib pilih satu:**

| Aksi | Penjelasan singkat |
| --- | --- |
| `--keep` | Hapus status merge sambil mempertahankan file hasil saat ini sebagai perubahan biasa. Konflik isi tetap perlu ditinjau. |
| `--abort` | Coba kembali ke keadaan sebelum merge; resolusi/edit selama merge dapat hilang. |

**Parameter:**

| Parameter | Nilai / default | Penjelasan singkat |
| --- | --- | --- |
| `--to` | `staged`, `changes`; default `changes` untuk keep | Khusus `--keep`. Staged memerlukan index tanpa konflik; changes mempertahankan file kerja. Tidak boleh dipakai dengan abort. |

**Contoh:** `gitmergecancel --keep --to changes`

## 9. gitamend

**Deskripsi:** Memasukkan perubahan staged ke commit terakhir tanpa mengubah pesan commit. Lokal, tidak push; harus ada staged changes.

**Parameter:**

| Parameter | Nilai / default | Penjelasan singkat |
| --- | --- | --- |
| `--no-edit` | Flag opsional; selalu default | Pertahankan pesan commit terakhir, termasuk pesan bilingual dan trailer. |

## 10. gitpushamend

**Deskripsi:** Amend commit terakhir menggunakan staged changes tanpa mengubah pesan, lalu push ke target terverifikasi. Staged changes wajib; jika amend sudah selesai sebelumnya, lanjutkan push tanpa amend ulang.

**Parameter:**

| Parameter | Nilai / default | Penjelasan singkat |
| --- | --- | --- |
| `--remote` | Nama remote; opsional | Pilih remote tujuan. Tanpa selector, gunakan upstream terverifikasi. |
| `--branch` | Nama branch; opsional | Pilih branch tujuan. Jika tidak ada upstream, remote dan branch harus jelas. |
| `--no-edit` | Flag opsional; selalu default | Pertahankan pesan saat amend. |

**Catatan:** Jika perlu menulis ulang commit yang sudah dipush, agent memeriksa cakupan dan menggunakan lease yang terikat keadaan remote; tidak otomatis melewati penolakan push.

## 11. gitpush

**Deskripsi:** Push commit yang sudah ada ke target terverifikasi. Tidak otomatis stage, commit atau amend.

**Parameter:**

| Parameter | Nilai / default | Penjelasan singkat |
| --- | --- | --- |
| `--remote` | Nama remote; opsional | Pilih remote tujuan; default mengikuti upstream yang terverifikasi. |
| `--branch` | Nama branch; opsional | Pilih branch tujuan; target lengkap harus jelas. |

**Contoh:** `gitpush --remote origin --branch main`

## 12. gitrelease

**Deskripsi:** Menyiapkan versi/changelog dan menerbitkan GitHub Release untuk semua jenis proyek selain website/aplikasi web berbasis browser, termasuk CLI, library, tooling, skrip, service, mobile dan desktop. Build, format paket/aset atau distribusi source mengikuti codebase; agent bertanya jika keputusan penting belum jelas. Tanpa aksi, hanya menampilkan penggunaan.

**Aksi:**

| Aksi | Penjelasan singkat |
| --- | --- |
| `init` | Validasi/inisialisasi versi dan changelog; buat/lengkapi packaging, workflow build/release dan helper lokal untuk target yang dipilih. |
| `prepare` | Periksa perubahan kode, tentukan versi dari bukti, cek kesiapan pipeline lalu perbarui metadata/changelog lokal. Tanyakan hanya keputusan penting yang belum jelas; tidak bump ulang persiapan yang sama. |
| `publish` | Publikasikan versi yang sudah di-commit melalui workflow atau langsung, dengan keluaran/distribusi proyek yang terverifikasi. |
| `delete` | Hapus release/tag yang disebut secara eksplisit, setelah preservasi dan persetujuan cakupan yang diperlukan. Tidak membuat release pengganti. |

**Parameter versi dan cakupan:**

| Parameter | Nilai / default | Penjelasan singkat |
| --- | --- | --- |
| `--version` | Versi tepat | Pilih versi, misalnya `1.3.0` atau `2.0.0-rc.1`. Tidak boleh digabung dengan bump. |
| `--bump` | `patch`, `minor`, `major` | Tentukan kenaikan versi; jika tidak dipilih, prepare menentukannya dari perubahan dan kebijakan repo. |
| `--from` | Tag atau SHA | Pilih baseline perubahan yang relevan dan terverifikasi. |
| `--component` | Path app | Pilih aplikasi pada monorepo; komponen lain tidak ikut diubah. |
| `--build-number` | Nomor build native | Pilih counter Android/iOS jika diperlukan; ikuti otoritas counter proyek/CI. |

**Parameter publikasi / pipeline:**

| Parameter | Nilai / default | Penjelasan singkat |
| --- | --- | --- |
| `--remote` | Nama remote; opsional | Pilih remote yang diverifikasi berdasarkan konteks repo. |
| `--branch` | Nama branch; opsional | Pilih branch publikasi yang terverifikasi. |
| `--tag` | Nama tag | Pilih tag yang konsisten dengan versi/prefix repo. Wajib untuk delete. |
| `--route` | `auto`, `workflow`, `direct`; default `auto` | Auto memilih publisher yang sesuai; workflow memakai CI; direct memverifikasi hasil build/paket/source yang sesuai; build yang diperlukan mengikuti perintah proyek. |
| `--workflow` | Path workflow | Pilih workflow yang sudah ada. Pada init boleh menunjuk file baru di `.github/workflows/`. |
| `--draft` | Flag | Minta draft GitHub Release; workflow harus mendukung mode ini. |
| `--assets` | Satu atau beberapa path file | Pilih berkas distribusi proyek terverifikasi untuk route direct; tidak berlaku pada route workflow. |
| `--replace-existing` | Flag; khusus publish | Ganti tag/release dari versi yang sudah disiapkan melalui prosedur preservasi, pemeriksaan dan persetujuan cakupan. Bukan penghapusan otomatis saat push ditolak. |

**Parameter penghapusan:**

| Parameter | Nilai / default | Penjelasan singkat |
| --- | --- | --- |
| `--delete-scope` | `release`, `tag`, `both`; default `both` | Release menghapus GitHub Release/aset; tag menghapus tag lokal/remote; both menghapus keduanya. |

Delete hanya menerima `--tag` (wajib), `--delete-scope`, `--remote`, dan `--component`. `--replace-existing` tidak berlaku untuk init/prepare/delete. Pada init/prepare, parameter pipeline hanya menyiapkan atau menjelaskan rute; tidak melakukan publikasi. Publikasi menunggu aset valid dari seluruh target wajib; tag yang dipush saja belum membuktikan release selesai.

**Contoh alur:** `gitrelease prepare` → `commitmsg` → `gitrelease publish`.

## Parameter pesan commit — commitpln dan commitmsg

| Parameter | Nilai / default | Penjelasan singkat |
| --- | --- | --- |
| `--lang` | `en`, `id`, `both`; default `en` | Bahasa pesan. Both menghasilkan pesan EN dan ID lengkap dalam satu commit/draft, dipisahkan `---`. |
| `--format` | `short`, `standard`, `detailed` | Short: subject; standard: subject dan poin `-`; detailed: rincian relevan tambahan. Alias standard: `standar`. Default mengikuti command. |
| `--scope` | Nama area | Pilih area fungsional seperti auth, sesuai perubahan. |
| `--type` | Jenis Conventional Commit | Misalnya feat, fix, refactor, docs, chore; harus cocok dengan diff dan aturan repo. |
| `--issue` | Referensi issue | Sertakan referensi issue yang diberikan/terverifikasi pada trailer yang sesuai. |

Jika user meminta draft/pemeriksaan saja, agent mengikuti batas read-only tersebut. Opsi yang bertentangan atau tidak didukung ditolak sebelum perubahan dilakukan.
