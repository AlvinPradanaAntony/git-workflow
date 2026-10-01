#!/usr/bin/env sh
# Install the Git Workflow executable once; skill installation remains `git-workflow init`.
set -eu

repo=${GIT_WORKFLOW_REPO:-}
version=latest
source_dir=
bin_dir=${GIT_WORKFLOW_BIN_DIR:-"$HOME/.local/bin"}
modify_path=1
replace=0

usage() {
    cat <<'EOF'
Install Git Workflow CLI for Linux/macOS (no Python, Node or Go required).
  sh install.sh --source-dir /path/to/release-files
  sh install.sh [--version v2.13.0]
Default repository: AlvinPradanaAntony/git-workflow
Options: --bin-dir PATH, --no-path, --replace, --help
The release files must include the platform binary and SHA256SUMS.
No repository is selected and no skill is installed by this bootstrap.
EOF
}
fail() { printf '%s\n' "Stopped: $*" >&2; exit 2; }
while [ "$#" -gt 0 ]; do
    case "$1" in
        --repo|--version|--source-dir|--bin-dir)
            [ "$#" -ge 2 ] || fail "Missing value for $1"
            case "$1" in
                --repo) repo=$2;;
                --version) version=$2;;
                --source-dir) source_dir=$2;;
                --bin-dir) bin_dir=$2;;
            esac
            shift 2;;
        --no-path) modify_path=0; shift;;
        --replace) replace=1; shift;;
        -h|--help) usage; exit 0;;
        *) fail "Unknown option: $1";;
    esac
done

case "$(uname -s)" in Linux) os=linux;; Darwin) os=darwin;; *) fail "Use install.ps1 on Windows";; esac
case "$(uname -m)" in x86_64|amd64) arch=amd64;; aarch64|arm64) arch=arm64;; *) fail "Unsupported CPU architecture";; esac
asset="git-workflow-$os-$arch"
if command -v sha256sum >/dev/null 2>&1; then
    hash_file() { sha256sum "$1" | awk '{print $1}'; }
elif command -v shasum >/dev/null 2>&1; then
    hash_file() { shasum -a 256 "$1" | awk '{print $1}'; }
else
    fail "sha256sum or shasum is required to verify the executable"
fi

# Adjacent release files work without any repository URL.
if [ -z "$source_dir" ] && [ -z "$repo" ]; then
    script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
    if [ -f "$script_dir/$asset" ] && [ -f "$script_dir/SHA256SUMS" ]; then
        source_dir=$script_dir
    else
        repo=AlvinPradanaAntony/git-workflow
    fi
fi
[ -z "$source_dir" ] || [ -z "$repo" ] || fail "Choose --source-dir or --repo"
[ -n "$bin_dir" ] || fail "Empty bin directory"
case "$bin_dir" in
    *'
'*) fail "Bin directory cannot contain newlines";;
    /*) :;; *) fail "Use an absolute --bin-dir path";;
esac
[ ! -L "$bin_dir" ] || fail "Bin directory is a symlink"
destination="$bin_dir/git-workflow"
receipt="$bin_dir/.git-workflow-cli.sha256"
[ ! -L "$destination" ] && [ ! -L "$receipt" ] || fail "Destination or receipt is a symlink"
if [ -n "$repo" ]; then
    case "$repo" in *[!A-Za-z0-9_./-]*|/*|*/|*/*/*) fail "Expected GitHub OWNER/REPOSITORY";; esac
    case "$repo" in */*) :;; *) fail "Expected GitHub OWNER/REPOSITORY";; esac
    case "$version" in latest) base="https://github.com/$repo/releases/latest/download";;
        v[0-9]*) case "$version" in *[!A-Za-z0-9_.-]*) fail "Invalid version tag";; esac
                 base="https://github.com/$repo/releases/download/$version";;
        *) fail "Use latest or a version tag such as v2.13.0";;
    esac
    command -v curl >/dev/null 2>&1 || fail "curl is required for a GitHub download"
fi

temp_dir=$(mktemp -d)
trap 'rm -rf "$temp_dir"' EXIT HUP INT TERM
if [ -n "$source_dir" ]; then
    [ -f "$source_dir/$asset" ] && [ -f "$source_dir/SHA256SUMS" ] || fail "Missing $asset or SHA256SUMS"
    cp "$source_dir/$asset" "$temp_dir/$asset"
    cp "$source_dir/SHA256SUMS" "$temp_dir/SHA256SUMS"
else
    curl --proto '=https' --tlsv1.2 -fsSL "$base/$asset" -o "$temp_dir/$asset"
    curl --proto '=https' --tlsv1.2 -fsSL "$base/SHA256SUMS" -o "$temp_dir/SHA256SUMS"
fi
expected=$(awk -v name="$asset" '$2 == name {print $1}' "$temp_dir/SHA256SUMS")
[ "${#expected}" -eq 64 ] || fail "Missing/duplicate SHA256SUMS entry for $asset"
case "$expected" in *[!0-9a-f]*) fail "Invalid SHA-256 entry";; esac
actual=$(hash_file "$temp_dir/$asset")
[ "$expected" = "$actual" ] || fail "Executable checksum mismatch; nothing installed"

if [ "$modify_path" -eq 1 ]; then
    case "${SHELL:-}" in */zsh) profile="$HOME/.zshrc";; */bash) profile="$HOME/.bashrc";; *) profile="$HOME/.profile";; esac
    [ ! -L "$profile" ] || fail "Shell profile is a symlink; use --no-path and manage PATH yourself"
    # Quote an arbitrary absolute path as a literal shell word, never execute it.
    quoted_dir=$(printf '%s' "$bin_dir" | sed "s/'/'\\\\''/g")
    path_line="export PATH='$quoted_dir':\"\$PATH\""
    if [ -f "$profile" ] && grep -Fq '# git-workflow:path:begin' "$profile"; then
        grep -Fxq "$path_line" "$profile" || fail "Existing Git Workflow PATH block differs; use --no-path or review it"
    fi
fi

mkdir -p "$bin_dir"
if [ -e "$destination" ]; then
    [ -f "$destination" ] || fail "Destination is not a regular file"
    old_hash=$(hash_file "$destination")
    if [ "$old_hash" != "$expected" ]; then
        recorded=
        if [ -f "$receipt" ]; then recorded=$(cat "$receipt"); fi
        [ "$recorded" = "$old_hash" ] || [ "$replace" -eq 1 ] || fail "Existing executable is unmanaged/edited; review it or use --replace"
        backup=$(mktemp "$bin_dir/git-workflow.backup.XXXXXX")
        cp -p "$destination" "$backup"
        printf '%s\n' "Previous executable backup: $backup"
    fi
fi
temp_binary=$(mktemp "$bin_dir/.git-workflow-install.XXXXXX")
cp "$temp_dir/$asset" "$temp_binary"
chmod 755 "$temp_binary"
mv -f "$temp_binary" "$destination"
printf '%s\n' "$expected" > "$receipt"

if [ "$modify_path" -eq 1 ]; then
    if [ ! -f "$profile" ] || ! grep -Fq '# git-workflow:path:begin' "$profile"; then
        printf '\n# git-workflow:path:begin\n%s\n# git-workflow:path:end\n' "$path_line" >> "$profile"
    fi
    printf '%s\n' "PATH configured in $profile. Open a new terminal, or run: $path_line"
fi
"$destination" --version
printf '%s\n' "CLI installed at $destination. In a repository run: git-workflow init"
