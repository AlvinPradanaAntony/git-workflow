#!/usr/bin/env python3
"""Validate release assets, render changelog notes and publish one GitHub Release."""
import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import subprocess

from build import TARGETS
from release_notes import render

ROOT = Path(__file__).resolve().parents[1]
SUPPORT = ("install.sh", "install.ps1", "install_git_workflow.py")


def version():
    value = (ROOT / "VERSION").read_text(encoding="ascii").strip()
    if not re.fullmatch(r"(?:0|[1-9][0-9]*)\.(?:0|[1-9][0-9]*)\.(?:0|[1-9][0-9]*)", value):
        raise ValueError("VERSION must contain a stable SemVer")
    if json.loads((ROOT / "bundle.json").read_text(encoding="utf-8"))["version"] != value:
        raise ValueError("VERSION and embedded skill version differ")
    return value


def asset_names():
    binaries = [f"git-workflow-{system}-{arch}" + (".exe" if system == "windows" else "")
                for system, arch in TARGETS]
    return binaries + list(SUPPORT) + [f"git-workflow-source-v{version()}.tar.gz"]


def check_dist(dist):
    names = asset_names()
    entries = {}
    for line in (dist / "SHA256SUMS").read_text(encoding="ascii").splitlines():
        match = re.fullmatch(r"([0-9a-f]{64})  ([A-Za-z0-9_.-]+)", line)
        if not match or match[2] in entries:
            raise ValueError("Malformed or duplicate SHA256SUMS entry")
        entries[match[2]] = match[1]
    if set(entries) != set(names):
        raise ValueError("SHA256SUMS must include all six executables, bootstrap, Python installer and source")
    for name in names:
        path = dist / name
        if path.is_symlink() or not path.is_file() or path.stat().st_size == 0:
            raise ValueError(f"Missing/empty/unsafe release asset: {name}")
        data = path.read_bytes()
        if hashlib.sha256(data).hexdigest() != entries[name]:
            raise ValueError(f"Release asset checksum mismatch: {name}")
        if name.startswith("git-workflow-linux-") and not data.startswith(b"\x7fELF"):
            raise ValueError(f"Invalid Linux executable: {name}")
        if name.startswith("git-workflow-windows-") and not data.startswith(b"MZ"):
            raise ValueError(f"Invalid Windows executable: {name}")
        if name.startswith("git-workflow-darwin-") and data[:4] not in (b"\xcf\xfa\xed\xfe", b"\xfe\xed\xfa\xcf"):
            raise ValueError(f"Invalid macOS executable: {name}")
    return names + ["SHA256SUMS"]


def notes():
    value = version()
    body = render((ROOT / "CHANGELOG.md").read_text(encoding="utf-8"), value)
    rows = []
    for system, arch in TARGETS:
        name = f"git-workflow-{system}-{arch}" + (".exe" if system == "windows" else "")
        label = {"darwin": "macOS", "windows": "Windows", "linux": "Linux"}[system]
        rows.append(f"| {label} {arch} | [{name}]({name}) | CLI portable, tanpa runtime tambahan |")
    rows += ["| Python 3.10+ | [install_git_workflow.py](install_git_workflow.py) | Installer skill; preview default, tambahkan `--apply` |",
             "| Bootstrap | [install.sh](install.sh) / [install.ps1](install.ps1) | Pasang CLI ke PATH pengguna |",
             "| Checksum | [SHA256SUMS](SHA256SUMS) | SHA-256 seluruh aset unduhan |"]
    return (f"## 📋 Apa yang Baru di v{value}?\n\n" + body +
            "\n---\n\n## 📥 Download\n\n| Platform | File | Keterangan |\n| --- | --- | --- |\n" +
            "\n".join(rows) + "\n")


def gh(*args):
    return subprocess.run(["gh", *args], check=True, text=True, capture_output=True).stdout


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("action", choices=("check", "notes", "finalize", "publish"))
    parser.add_argument("--dist", type=Path, default=ROOT / "dist")
    args = parser.parse_args()
    value = version()
    tag = os.environ.get("RELEASE_TAG", "")
    if tag and tag != "v" + value:
        parser.error("Release tag must match VERSION and bundled skill")
    body = notes()
    dist = args.dist
    if args.action == "check":
        print(f"Version, embedded skill and changelog verified: {value}")
        return
    if args.action == "notes":
        print(body, end="")
        return
    if args.action == "finalize":
        dist.mkdir(parents=True, exist_ok=True)
        source = dist / f"git-workflow-source-v{value}.tar.gz"
        subprocess.run(["git", "archive", "--format=tar.gz", "--prefix=git-workflow/",
                        "-o", str(source), "HEAD"], cwd=ROOT, check=True)
        (dist / "SHA256SUMS").write_text("".join(
            hashlib.sha256((dist / name).read_bytes()).hexdigest() + "  " + name + "\n"
            for name in sorted(asset_names())), encoding="ascii")
        check_dist(dist)
        print("All release assets and checksums verified.")
        return
    repo = os.environ.get("RELEASE_REPO", "")
    if not re.fullmatch(r"[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+", repo) or not tag:
        parser.error("publish requires RELEASE_REPO and RELEASE_TAG")
    files = check_dist(dist)
    # Pin links to the verified repository and tag instead of relative release-page URLs.
    for name in asset_names() + ["SHA256SUMS"]:
        body = body.replace(f"]({name})", f"](https://github.com/{repo}/releases/download/{tag}/{name})")
    notes_file = dist / "release-notes.md"
    notes_file.write_text(body, encoding="utf-8")
    pages = json.loads(gh("api", f"repos/{repo}/releases?per_page=100", "--paginate", "--slurp"))
    if any(release["tag_name"] == tag for page in pages for release in page):
        parser.error("A release already exists for this tag; inspect it before replacing or resuming publication")
    # Upload verified assets into a draft. Publish only after GitHub lists every expected asset.
    gh("release", "create", tag, "--repo", repo, "--verify-tag", "--draft",
       "--title", "Git Workflow " + tag, "--notes-file", str(notes_file))
    gh("release", "upload", tag, *[str(dist / name) for name in files], "--repo", repo)
    release = json.loads(gh("api", f"repos/{repo}/releases/tags/{tag}"))
    uploaded = {asset["name"]: asset for asset in release["assets"]}
    if set(uploaded) != set(files):
        raise ValueError("Uploaded asset set differs; release remains a draft")
    for name in files:
        if uploaded[name]["size"] != (dist / name).stat().st_size or uploaded[name]["state"] != "uploaded":
            raise ValueError(f"Incomplete upload {name}; release remains a draft")
        expected_digest = "sha256:" + hashlib.sha256((dist / name).read_bytes()).hexdigest()
        if uploaded[name].get("digest") and uploaded[name]["digest"] != expected_digest:
            raise ValueError(f"Uploaded checksum mismatch {name}; release remains a draft")
    gh("release", "edit", tag, "--repo", repo, "--draft=false")
    published = json.loads(gh("api", f"repos/{repo}/releases/tags/{tag}"))
    if published["draft"] or published["tag_name"] != tag:
        raise ValueError("Release publication could not be verified")
    print(f"Published https://github.com/{repo}/releases/tag/{tag}")


if __name__ == "__main__":
    main()
