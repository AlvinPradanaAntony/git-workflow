#!/usr/bin/env python3
"""Read-only extraction of one approved changelog version and one predecessor."""

import argparse
from dataclasses import dataclass
from datetime import date
from pathlib import Path
import re
from typing import Optional
from urllib.parse import urlsplit


SEMVER = re.compile(
    r"(?:0|[1-9]\d*)\.(?:0|[1-9]\d*)\.(?:0|[1-9]\d*)"
    r"(?:-(?:0|[1-9]\d*|[0-9A-Za-z-]*[A-Za-z-][0-9A-Za-z-]*)"
    r"(?:\.(?:0|[1-9]\d*|[0-9A-Za-z-]*[A-Za-z-][0-9A-Za-z-]*))*)?"
    r"(?:\+[0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*)?\Z"
)
HEADER = re.compile(
    r"^#{2,3} \[([^\]]+)\](?:\([^\r\n)]*\))? - (\d{4}-\d{2}-\d{2})\s*$"
)


@dataclass(frozen=True)
class Entry:
    version: str
    day: str
    body: str


def parse_entries(content: str) -> list[Entry]:
    # Omit writing-guide comments, and ignore apparent headers in fenced examples.
    content = re.sub(r"<!--[\s\S]*?-->", "", content)
    if "<!--" in content:
        raise ValueError("Unclosed HTML comment in changelog")
    lines = content.splitlines()
    starts = []
    fence: Optional[str] = None
    for number, line in enumerate(lines):
        marker = re.match(r"^\s{0,3}(`{3,}|~{3,})", line)
        if marker:
            token = marker.group(1)
            if fence is None:
                fence = token
            elif token[0] == fence[0] and len(token) >= len(fence):
                fence = None
            continue
        if fence is None:
            match = HEADER.fullmatch(line)
            if match:
                date.fromisoformat(match.group(2))
                starts.append((number, match.group(1), match.group(2)))
    entries = []
    for offset, (number, version, day) in enumerate(starts):
        end = starts[offset + 1][0] if offset + 1 < len(starts) else len(lines)
        body = []
        fence = None
        for line in lines[number + 1:end]:
            marker = re.match(r"^\s{0,3}(`{3,}|~{3,})", line)
            if marker:
                token = marker.group(1)
                if fence is None:
                    fence = token
                elif token[0] == fence[0] and len(token) >= len(fence):
                    fence = None
            elif fence is None and (
                line.strip() == "---" or re.match(r"^\s*</?details\b", line)
            ):
                break
            body.append(line)
        entries.append(Entry(version, day, "\n".join(body).strip()))
    return entries


def select(entries: list[Entry], version: str) -> Entry:
    if not SEMVER.fullmatch(version):
        raise ValueError("Use a SemVer version without the tag prefix")
    matches = [entry for entry in entries if entry.version == version]
    if len(matches) != 1:
        raise ValueError(f"Expected one entry for {version}; found {len(matches)}")
    entry = matches[0]
    if not any(line.strip() and not line.lstrip().startswith("#")
               for line in entry.body.splitlines()):
        raise ValueError(f"No release notes for {version}")
    return entry


def render(content: str, version: str, previous: Optional[str] = None,
           previous_url: Optional[str] = None) -> str:
    entries = parse_entries(content)
    current = select(entries, version)
    if previous_url and not previous:
        raise ValueError("--previous-url requires --previous")
    result = current.body
    if previous:
        if previous == version:
            raise ValueError("Previous version must differ from the release version")
        old = select(entries, previous)
        label = f"[{old.version}]"
        if previous_url:
            url = urlsplit(previous_url)
            if (url.scheme != "https" or not url.netloc or url.username or url.password
                    or url.query or url.fragment or "/releases/tag/" not in url.path
                    or any(char in previous_url for char in "\r\n()<> ")):
                raise ValueError("Use a verified HTTPS GitHub release/tag URL")
            label += f"({previous_url})"
        # Active entries use ### categories; collapsed entries use ####.
        old_body = re.sub(r"(?m)^### (?!\[)", "#### ", old.body)
        result += (
            "\n\n---\n\n<details>\n"
            "<summary><strong>Riwayat versi sebelumnya</strong></summary>\n\n"
            f"### {label} - {old.day}\n\n{old_body}\n\n</details>"
        )
    return result + "\n"


def main(argv=None) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--changelog", required=True, type=Path)
    parser.add_argument("--version", required=True)
    parser.add_argument("--previous")
    parser.add_argument("--previous-url")
    args = parser.parse_args(argv)
    try:
        result = render(args.changelog.read_text(encoding="utf-8-sig"), args.version,
                        args.previous, args.previous_url)
    except (OSError, ValueError) as exc:
        parser.error(str(exc))
    print(result, end="")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
