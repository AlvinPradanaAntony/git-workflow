package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	if err := loadKit(); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func isolatedHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("LOCALAPPDATA", filepath.Join(home, "Local"))
	t.Setenv("XDG_CACHE_HOME", filepath.Join(home, "cache"))
	return home
}

func git(t *testing.T, root string, args ...string) string {
	t.Helper()
	argv := append([]string{"-C", root}, args...)
	b, err := exec.Command("git", argv...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v: %s", args, err, b)
	}
	return strings.TrimSpace(string(b))
}

func repo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	git(t, root, "init", "-q")
	git(t, root, "config", "user.name", "Test")
	git(t, root, "config", "user.email", "test@example.invalid")
	return root
}

func invoke(t *testing.T, expected int, args ...string) string {
	t.Helper()
	var out, errOut bytes.Buffer
	code := run(args, strings.NewReader(""), &out, &errOut)
	if code != expected {
		t.Fatalf("%v: exit %d wanted %d\n%s\n%s", args, code, expected, out.String(), errOut.String())
	}
	return out.String() + errOut.String()
}

func write(t *testing.T, root, rel, data string) {
	t.Helper()
	file := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}
}

func read(t *testing.T, root, rel string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func missing(t *testing.T, root, rel string) {
	t.Helper()
	if _, err := os.Lstat(filepath.Join(root, filepath.FromSlash(rel))); !os.IsNotExist(err) {
		t.Fatalf("expected absent %s: %v", rel, err)
	}
}

func TestProjectLifecycleAndPreservation(t *testing.T) {
	isolatedHome(t)
	r := repo(t)
	write(t, r, "AGENTS.md", "# Team rules\r\nKeep this content.\r\n")
	write(t, r, "notes.txt", "Unrelated project file\n")
	before := git(t, r, "status", "--porcelain")
	invoke(t, 0, "init", "--project", r, "--dry-run")
	missing(t, r, ".agents/git-workflow-install.json")
	if git(t, r, "status", "--porcelain") != before {
		t.Fatal("dry-run changed repo")
	}
	invoke(t, 0, "init", "--project", r)
	first := read(t, r, ".agents/git-workflow-install.json")
	if !strings.HasPrefix(read(t, r, "AGENTS.md"), "# Team rules\r\nKeep this content.\r\n") {
		t.Fatal("AGENTS content not preserved")
	}
	if !strings.Contains(invoke(t, 0, "init", "--project", r), "Already up to date") {
		t.Fatal("init not idempotent")
	}
	if read(t, r, ".agents/git-workflow-install.json") != first {
		t.Fatal("manifest changed on repeat")
	}
	invoke(t, 0, "doctor", "--project", r)
	invoke(t, 0, "update", "--project", r)
	write(t, r, ".agents/skills/git-workflow/local-notes.txt", "Keep my notes")
	invoke(t, 0, "uninstall", "--project", r, "--dry-run")
	invoke(t, 0, "uninstall", "--project", r)
	missing(t, r, ".agents/git-workflow-install.json")
	missing(t, r, ".agents/skills/git-workflow/SKILL.md")
	if read(t, r, "notes.txt") != "Unrelated project file\n" || read(t, r, ".agents/skills/git-workflow/local-notes.txt") != "Keep my notes" {
		t.Fatal("unrelated files changed")
	}
	if !strings.Contains(read(t, r, "AGENTS.md"), "Keep this content.") {
		t.Fatal("AGENTS content lost on uninstall")
	}
	if read(t, r, "AGENTS.md") != "# Team rules\r\nKeep this content.\r\n" {
		t.Fatal("uninstall did not restore original AGENTS bytes")
	}
	invoke(t, 0, "uninstall", "--project", r)
	invoke(t, 2, "update", "--project", r)
}

func TestCurrentDirectorySubfolderAndLinkedWorktree(t *testing.T) {
	isolatedHome(t)
	r := repo(t)
	write(t, r, "code.txt", "code")
	git(t, r, "add", "code.txt")
	git(t, r, "commit", "-qm", "initial")
	head := git(t, r, "rev-parse", "HEAD")
	worktree := filepath.Join(t.TempDir(), "linked")
	git(t, r, "worktree", "add", "-qb", "test-linked", worktree)
	nested := filepath.Join(worktree, "frontend", "src")
	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(nested)
	invoke(t, 0, "init")
	read(t, worktree, ".agents/git-workflow-install.json")
	missing(t, r, ".agents/git-workflow-install.json")
	missing(t, nested, ".agents")
	if git(t, worktree, "rev-parse", "HEAD") != head {
		t.Fatal("HEAD changed")
	}
	invoke(t, 0, "doctor")
	invoke(t, 0, "uninstall")
}

func TestGlobalSelectionAndPartialUninstall(t *testing.T) {
	home := isolatedHome(t)
	t.Chdir(t.TempDir())
	invoke(t, 2, "init", "-g")
	missing(t, home, ".config/git-workflow/install.json")
	if err := os.Mkdir(filepath.Join(home, ".codex"), 0755); err != nil {
		t.Fatal(err)
	}
	invoke(t, 0, "init", "-g")
	read(t, home, ".agents/skills/git-workflow/SKILL.md")
	missing(t, home, ".gemini/config/skills/git-workflow/SKILL.md")
	invoke(t, 0, "init", "-g", "--agent", "all")
	if strings.Contains(read(t, home, ".gemini/config/skills/git-workflow/SKILL.md"), "Read `.agents/rules/git-workflow.md`") {
		t.Fatal("global skill contains project-relative pointer")
	}
	invoke(t, 0, "doctor", "-g")
	invoke(t, 0, "uninstall", "-g", "--agent", "antigravity")
	missing(t, home, ".gemini/config/skills/git-workflow/SKILL.md")
	read(t, home, ".agents/skills/git-workflow/SKILL.md")
	read(t, home, ".gemini/antigravity-cli/skills/git-workflow/SKILL.md")
	invoke(t, 0, "update", "-g")
	missing(t, home, ".gemini/config/skills/git-workflow/SKILL.md")
	invoke(t, 1, "doctor", "-g", "--agent", "antigravity")
	invoke(t, 0, "uninstall", "-g")
	missing(t, home, ".config/git-workflow/install.json")
}

func TestEditedFileRefusalAndExactBackup(t *testing.T) {
	isolatedHome(t)
	r := repo(t)
	invoke(t, 0, "init", "--project", r)
	rel := ".agents/skills/git-workflow/references/commitmsg.md"
	write(t, r, rel, "My customization\n")
	before := read(t, r, ".agents/git-workflow-install.json")
	invoke(t, 1, "doctor", "--project", r)
	invoke(t, 2, "update", "--project", r)
	invoke(t, 2, "uninstall", "--project", r)
	if read(t, r, rel) != "My customization\n" || read(t, r, ".agents/git-workflow-install.json") != before {
		t.Fatal("refusal changed files")
	}
	output := invoke(t, 0, "update", "--project", r, "--replace")
	var backup string
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, "Exact previous-file backup: ") {
			backup = strings.TrimPrefix(line, "Exact previous-file backup: ")
		}
	}
	if backup == "" || read(t, backup, rel) != "My customization\n" {
		t.Fatal("exact backup missing")
	}
	invoke(t, 0, "doctor", "--project", r)
}

func TestMalformedArgumentsAndUnsafeManifest(t *testing.T) {
	isolatedHome(t)
	r := repo(t)
	invoke(t, 2, "init", "--project", r, "-g")
	invoke(t, 2, "init", "--project", r, "--agent", "unknown")
	invoke(t, 2, "init", "--project", r, "--unknown")
	invoke(t, 2, "commitmsg")
	missing(t, r, ".agents")
	t.Chdir(t.TempDir())
	invoke(t, 2, "init")
	write(t, r, "precious.txt", "preserve")
	m := manifest{Package: "git-workflow", Version: kit.Version, Files: map[string]string{"../precious.txt": digest([]byte("preserve"))}}
	b, _ := json.Marshal(m)
	write(t, r, ".agents/git-workflow-install.json", string(b))
	invoke(t, 2, "uninstall", "--project", r, "--replace")
	if read(t, r, "precious.txt") != "preserve" {
		t.Fatal("unsafe manifest followed")
	}
}

func TestSymlinkDestinationRefused(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink privilege is not available on all Windows hosts")
	}
	isolatedHome(t)
	r := repo(t)
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(r, ".agents")); err != nil {
		t.Fatal(err)
	}
	invoke(t, 2, "init", "--project", r, "--replace")
	missing(t, outside, "skills")
}

func TestLegacyManifestMigrationAndDowngradeRefusal(t *testing.T) {
	isolatedHome(t)
	r := repo(t)
	write(t, r, ".agents/git-workflow.rules.md", "old rules")
	write(t, r, ".agents/workflows/commitmsg.md", "old alias")
	write(t, r, ".agents/skills/commitmsg/SKILL.md", "old skill")
	write(t, r, "AGENTS.md", "# Keep\n"+kit.Blocks[1]+"\n")
	m := manifest{Package: "git-workflow", Version: "2.4.0", Files: map[string]string{
		".agents/git-workflow.rules.md":     digest([]byte("old rules")),
		".agents/workflows/commitmsg.md":    digest([]byte("old alias")),
		".agents/skills/commitmsg/SKILL.md": digest([]byte("old skill")),
	}}
	b, _ := json.Marshal(m)
	write(t, r, ".agents/git-workflow-install.json", string(b))
	invoke(t, 0, "update", "--project", r)
	missing(t, r, ".agents/git-workflow.rules.md")
	missing(t, r, ".agents/workflows/commitmsg.md")
	missing(t, r, ".agents/skills/commitmsg/SKILL.md")
	invoke(t, 0, "doctor", "--project", r)
	current, err := readManifest(r, ".agents/git-workflow-install.json", false)
	if err != nil {
		t.Fatal(err)
	}
	current.Version = "99.0.0"
	b, _ = json.Marshal(current)
	write(t, r, ".agents/git-workflow-install.json", string(b))
	invoke(t, 2, "update", "--project", r, "--replace")
}

func TestChangedPreflightDoesNotWrite(t *testing.T) {
	isolatedHome(t)
	r := repo(t)
	write(t, r, "AGENTS.md", "original")
	if err := loadKit(); err != nil {
		t.Fatal(err)
	}
	plans, err := plan(r, ".agents/git-workflow-install.json", nil, options{Action: "init"}, []string{"codex"})
	if err != nil {
		t.Fatal(err)
	}
	write(t, r, "AGENTS.md", "concurrent edit")
	var out bytes.Buffer
	if err := apply(r, plans, &out); err == nil {
		t.Fatal("changed preflight accepted")
	}
	missing(t, r, ".agents/skills")
	if read(t, r, "AGENTS.md") != "concurrent edit" {
		t.Fatal("concurrent edit lost")
	}
}

func TestUnrelatedMixedNewlinesAndMissingGitDoctor(t *testing.T) {
	home := isolatedHome(t)
	r := repo(t)
	invoke(t, 0, "init", "--project", r)
	write(t, r, "AGENTS.md", read(t, r, "AGENTS.md")+"Unrelated user line.\r\n")
	invoke(t, 0, "doctor", "--project", r)
	invoke(t, 0, "update", "--project", r)
	invoke(t, 0, "uninstall", "--project", r)
	if !strings.Contains(read(t, r, "AGENTS.md"), "Unrelated user line.\r\n") {
		t.Fatal("surrounding mixed newline content was changed")
	}
	invoke(t, 0, "init", "-g", "--agent", "codex")
	t.Setenv("PATH", t.TempDir())
	if !strings.Contains(invoke(t, 1, "doctor", "-g"), "Git: unavailable") {
		t.Fatal("missing Git was not reported")
	}
	read(t, home, ".agents/skills/git-workflow/SKILL.md")
}

func TestEmptyExistingAgentsFileIsPreserved(t *testing.T) {
	isolatedHome(t)
	r := repo(t)
	write(t, r, "AGENTS.md", "")
	invoke(t, 0, "init", "--project", r)
	invoke(t, 0, "update", "--project", r)
	invoke(t, 0, "uninstall", "--project", r)
	if read(t, r, "AGENTS.md") != "" {
		t.Fatal("pre-existing empty AGENTS.md changed")
	}
}
