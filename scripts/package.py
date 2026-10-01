#!/usr/bin/env python3
"""Generate the embedded CLI payload and standalone Python installer from one skill."""
import argparse
import base64
import importlib.util
import json
from pathlib import Path
import re
import zlib

ROOT = Path(__file__).resolve().parents[1]
PORTABLE = ("The portable project copy is managed by the standalone installer. "
            "Read `.agents/rules/git-workflow.md` and the matching command reference before acting.")


def outputs(skill):
    version = (ROOT / "VERSION").read_text(encoding="ascii").strip()
    if not re.fullmatch(r"(?:0|[1-9][0-9]*)\.(?:0|[1-9][0-9]*)\.(?:0|[1-9][0-9]*)", version):
        raise ValueError("VERSION must contain a stable SemVer")
    template = ROOT / "scripts/installer_template.py"
    spec = importlib.util.spec_from_file_location("installer_template", template)
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    paths = [skill / "SKILL.md", skill / "agents/openai.yaml"]
    paths += sorted((skill / "references").glob("*.md"))
    paths += sorted((skill / "assets").glob("*.svg"))
    paths += [skill / "scripts/release_notes.py"]
    files = {}
    for item in paths:
        rel = item.relative_to(skill).as_posix()
        content = item.read_text(encoding="utf-8")
        if rel == "SKILL.md":
            content += "\n" + PORTABLE + "\n"
        files[".agents/skills/git-workflow/" + rel] = content
    bundle = {"version": version, "files": files, "rule": module.RULE,
              "agents_blocks": [module.AGENTS_BLOCK, module.OLD_AGENTS_BLOCK, module.PREVIOUS_AGENTS_BLOCK]}
    bundle_text = json.dumps(bundle, ensure_ascii=False, indent=2) + "\n"
    payload = base64.b64encode(zlib.compress(json.dumps(files, ensure_ascii=False).encode("utf-8"))).decode("ascii")
    installer = template.read_text(encoding="utf-8")
    installer = installer.replace('VERSION = "__VERSION__"', f'VERSION = "{version}"')
    installer = installer.replace('PAYLOAD_B64 = "__PAYLOAD__"', "PAYLOAD_B64 = " + repr(payload))
    if "__VERSION__" in installer or "__PAYLOAD__" in installer:
        raise ValueError("Installer generation left unresolved placeholders")
    return {ROOT / "bundle.json": bundle_text, ROOT / "install_git_workflow.py": installer}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--skill-dir", type=Path, default=ROOT / "skill")
    parser.add_argument("--check", action="store_true", help="Fail if generated files are out of date")
    args = parser.parse_args()
    generated = outputs(args.skill_dir.resolve())
    for path, text in generated.items():
        if args.check:
            if not path.exists() or path.read_text(encoding="utf-8") != text:
                parser.error(f"{path.name} is stale; run python scripts/package.py")
        else:
            path.write_text(text, encoding="utf-8", newline="\n")
    print("Embedded skill and Python installer are consistent." if args.check else "Generated bundle.json and install_git_workflow.py.")


if __name__ == "__main__":
    main()
