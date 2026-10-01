import hashlib
import importlib.util
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import unittest
from unittest.mock import patch

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT / "scripts"))
import package
import release


def load_installer(path):
    spec = importlib.util.spec_from_file_location("test_installer", path)
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


class DistributionTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory(prefix="git-workflow-tests-")
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.home = self.root / "user"
        self.home.mkdir()
        self.env = dict(os.environ, HOME=str(self.home), USERPROFILE=str(self.home),
                        LOCALAPPDATA=str(self.home / "Local"), XDG_CACHE_HOME=str(self.home / "cache"))

    def repo(self):
        target = self.root / "repo with spaces"
        target.mkdir()
        subprocess.run(["git", "init", "-q", str(target)], check=True, env=self.env)
        return target

    def invoke(self, cwd, *args, script=None, code=0):
        result = subprocess.run([sys.executable, str(script or ROOT / "install_git_workflow.py"), *args],
                                cwd=cwd, env=self.env, text=True, capture_output=True)
        self.assertEqual(result.returncode, code, result.stdout + result.stderr)
        return result

    def test_generated_files_and_icon_are_consistent(self):
        for path, text in package.outputs(ROOT / "skill").items():
            self.assertEqual(path.read_text(encoding="utf-8"), text)
        installer = load_installer(ROOT / "install_git_workflow.py")
        payload = installer.load_payload()
        bundle = json.loads((ROOT / "bundle.json").read_text(encoding="utf-8"))
        self.assertEqual(payload, bundle["files"])
        self.assertIn(".agents/skills/git-workflow/assets/icon.svg", payload)
        self.assertEqual(installer.VERSION, release.version())

    def test_python_cwd_subfolder_idempotence_and_uninstall(self):
        target = self.repo()
        nested = target / "frontend" / "src"
        nested.mkdir(parents=True)
        agents = target / "AGENTS.md"
        agents.write_bytes(b"# Team\r\nKeep this.\r\n")
        self.invoke(nested)
        self.assertFalse((target / ".agents").exists())
        self.invoke(nested, "--apply")
        self.assertFalse((nested / ".agents").exists())
        self.assertTrue(agents.read_bytes().startswith(b"# Team\r\nKeep this.\r\n"))
        manifest = target / ".agents/git-workflow-install.json"
        before = manifest.read_bytes()
        self.invoke(nested, "--apply")
        self.assertEqual(manifest.read_bytes(), before)
        extra = target / ".agents/skills/git-workflow/my-notes.txt"
        extra.write_text("preserve", encoding="utf-8")
        self.invoke(nested, "--uninstall", "--apply")
        self.assertFalse(manifest.exists())
        self.assertEqual(extra.read_text(), "preserve")
        self.assertEqual(agents.read_bytes(), b"# Team\r\nKeep this.\r\n")

    def test_existing_empty_agents_file_and_added_user_rule_survive_uninstall(self):
        target = self.repo()
        agents = target / "AGENTS.md"
        agents.write_bytes(b"")
        self.invoke(target, "--apply")
        self.invoke(target, "--apply")
        self.invoke(target, "--uninstall", "--apply")
        self.assertTrue(agents.is_file())
        self.assertEqual(agents.read_bytes(), b"")
        agents.write_bytes(b"Team rule.\n")
        self.invoke(target, "--apply")
        with agents.open("ab") as stream:
            stream.write(b"New user rule.\r\n")
        self.invoke(target, "--uninstall", "--apply")
        self.assertEqual(agents.read_bytes(), b"Team rule.\nNew user rule.\r\n")

    def test_explicit_project_resolves_git_root_and_requires_git(self):
        target = self.repo()
        nested = target / "nested"
        nested.mkdir()
        self.invoke(self.home, "--project", str(nested), "--apply")
        self.assertTrue((target / ".agents/git-workflow-install.json").is_file())
        self.assertFalse((nested / ".agents").exists())
        self.invoke(self.home, "--apply", code=2)
        self.assertFalse((self.home / ".agents").exists())

    def test_upgrade_from_python_2123_preserves_user_files(self):
        target = self.repo()
        self.invoke(target, "--project", str(target), "--apply", script=ROOT / "scripts/legacy-2.12.3.py")
        agents = target / "AGENTS.md"
        with agents.open("ab") as stream:
            stream.write(b"My rule.\n")
        self.invoke(target, "--apply")
        self.assertEqual(json.loads((target / ".agents/git-workflow-install.json").read_text())["version"], release.version())
        self.assertIn(b"My rule.", agents.read_bytes())
        self.assertTrue((target / ".agents/skills/git-workflow/references/installation.md").exists())

    def test_edited_file_blocks_update_and_replace_backs_up(self):
        target = self.repo()
        self.invoke(target, "--apply")
        rel = ".agents/skills/git-workflow/references/commitmsg.md"
        file = target / rel
        file.write_text("My edited file", encoding="utf-8")
        result = self.invoke(target, "--apply", code=2)
        self.assertIn("No files changed", result.stderr)
        self.assertEqual(file.read_text(), "My edited file")
        result = self.invoke(target, "--apply", "--replace")
        path = next(line.split(": ", 1)[1] for line in result.stdout.splitlines()
                    if line.startswith("Exact previous-file backup:"))
        backup = Path(path)
        self.addCleanup(shutil.rmtree, backup, True)
        self.assertEqual((backup / rel).read_text(), "My edited file")

    def test_unsafe_manifest_and_downgrade_are_rejected(self):
        target = self.repo()
        self.invoke(target, "--apply")
        path = target / ".agents/git-workflow-install.json"
        manifest = json.loads(path.read_text())
        original = dict(manifest)
        manifest["version"] = "99.0.0"
        path.write_text(json.dumps(manifest))
        self.invoke(target, "--apply", code=2)
        manifest = original
        manifest["files"]["../precious.txt"] = "0" * 64
        path.write_text(json.dumps(manifest))
        self.invoke(target, "--apply", "--replace", code=2)

    def test_global_installs_and_uninstalls_all_integrations(self):
        self.invoke(self.home, "--global", "--apply")
        paths = [self.home / base / "git-workflow/SKILL.md" for base in
                 (".agents/skills", ".gemini/config/skills", ".gemini/antigravity-cli/skills")]
        for path in paths:
            self.assertTrue(path.is_file())
            self.assertNotIn("Read `.agents/rules/git-workflow.md`", path.read_text())
        self.assertFalse((self.home / "AGENTS.md").exists())
        self.invoke(self.home, "--global", "--uninstall", "--apply")
        for path in paths:
            self.assertFalse(path.exists())

    def test_release_notes_are_scoped_and_tag_must_match(self):
        notes = release.notes()
        self.assertIn("2.13.0", notes)
        self.assertNotIn("## [2.12.3]", notes)
        self.assertIn("install_git_workflow.py", notes)
        env = dict(self.env, RELEASE_TAG="v99.0.0")
        result = subprocess.run([sys.executable, str(ROOT / "scripts/release.py"), "check"],
                                cwd=ROOT, env=env, text=True, capture_output=True)
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("must match", result.stderr)

    def test_missing_tampered_or_duplicate_release_assets_stop_publication(self):
        dist = self.root / "dist"
        dist.mkdir()
        for name in release.asset_names():
            magic = b"data"
            if name.startswith("git-workflow-linux-"):
                magic = b"\x7fELF"
            elif name.startswith("git-workflow-windows-"):
                magic = b"MZ"
            elif name.startswith("git-workflow-darwin-"):
                magic = b"\xcf\xfa\xed\xfe"
            (dist / name).write_bytes(magic + b" fixture")
        sums = "".join(hashlib.sha256((dist / name).read_bytes()).hexdigest() + "  " + name + "\n"
                       for name in release.asset_names())
        (dist / "SHA256SUMS").write_text(sums)
        self.assertEqual(len(release.check_dist(dist)), 11)
        file = dist / "install_git_workflow.py"
        original = file.read_bytes()
        file.write_bytes(b"tampered")
        with self.assertRaisesRegex(ValueError, "checksum mismatch"):
            release.check_dist(dist)
        file.write_bytes(original)
        (dist / "SHA256SUMS").write_text(sums + sums.splitlines()[0] + "\n")
        with self.assertRaisesRegex(ValueError, "duplicate"):
            release.check_dist(dist)
        (dist / "SHA256SUMS").write_text(sums)
        file.unlink()
        with self.assertRaisesRegex(ValueError, "Missing/empty"):
            release.check_dist(dist)

    def test_publish_stays_draft_until_all_expected_assets_are_uploaded(self):
        dist = self.root / "publication"
        dist.mkdir()
        names = release.asset_names() + ["SHA256SUMS"]
        for name in names:
            (dist / name).write_bytes(b"fixture asset")
        uploaded = [{"name": name, "size": (dist / name).stat().st_size, "state": "uploaded"}
                    for name in names]
        for incomplete in (False, True):
            calls = []
            def fake_gh(*args):
                calls.append(args)
                if args[0] == "api" and "?per_page=" in args[1]:
                    return "[[]]"
                if args[0] == "api":
                    is_draft = not any(call[:2] == ("release", "edit") for call in calls)
                    return json.dumps({"tag_name": "v2.13.0", "draft": is_draft,
                                       "assets": uploaded[:-1] if incomplete else uploaded})
                return ""
            with patch.dict(os.environ, RELEASE_TAG="v2.13.0", RELEASE_REPO="example/fixture"), \
                    patch.object(sys, "argv", ["release.py", "publish", "--dist", str(dist)]), \
                    patch.object(release, "check_dist", return_value=names), \
                    patch.object(release, "gh", side_effect=fake_gh):
                if incomplete:
                    with self.assertRaisesRegex(ValueError, "remains a draft"):
                        release.main()
                else:
                    release.main()
            edits = [call for call in calls if call[:2] == ("release", "edit")]
            self.assertEqual(len(edits), 0 if incomplete else 1)
            create = next(call for call in calls if call[:2] == ("release", "create"))
            self.assertIn("--draft", create)
            if not incomplete:
                self.assertIn("https://github.com/example/fixture/releases/download/v2.13.0/", (dist / "release-notes.md").read_text())


if __name__ == "__main__":
    unittest.main()
