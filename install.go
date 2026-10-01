package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type manifest struct {
	Package       string            `json:"package"`
	Version       string            `json:"version"`
	Files         map[string]string `json:"files"`
	Agents        []string          `json:"agents,omitempty"`
	AgentVersions map[string]string `json:"agent_versions,omitempty"`
}

type state struct {
	exists bool
	data   []byte
	mode   os.FileMode
}

type change struct {
	rel    string
	before state
	after  []byte // nil means removal
}

func digest(data []byte) string {
	x := sha256.Sum256(data)
	return hex.EncodeToString(x[:])
}

func validRel(rel string) bool {
	return rel != "" && rel != "." && !path.IsAbs(rel) && path.Clean(rel) == rel &&
		!strings.ContainsAny(rel, "\\:\x00\r\n") && !strings.HasPrefix(rel, "../") && rel != ".."
}

func owned(rel string, global bool) bool {
	if !validRel(rel) {
		return false
	}
	if global {
		return agentForPath(rel) != ""
	}
	if rel == ".agents/rules/git-workflow.md" || rel == ".agents/git-workflow.rules.md" {
		return true
	}
	for _, name := range commands {
		if strings.HasPrefix(rel, ".agents/skills/"+name+"/") ||
			rel == ".agents/workflows/"+name+".md" || rel == ".agent/workflows/"+name+".md" {
			return true
		}
	}
	return false
}

func safePath(root, rel string) (string, error) {
	rootInfo, rootErr := os.Lstat(root)
	if rootErr != nil {
		return "", rootErr
	}
	if !rootInfo.IsDir() || rootInfo.Mode()&os.ModeSymlink != 0 {
		return "", errors.New("installation root is no longer a regular directory")
	}
	if !validRel(rel) {
		return "", fmt.Errorf("unsafe package path: %q", rel)
	}
	parts := strings.Split(rel, "/")
	p := root
	for i, part := range parts {
		p = filepath.Join(p, part)
		info, err := os.Lstat(p)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return "", err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("refusing symlink destination: %s", p)
		}
		if i < len(parts)-1 && !info.IsDir() {
			return "", fmt.Errorf("parent is not a directory: %s", p)
		}
		if i == len(parts)-1 && !info.Mode().IsRegular() {
			return "", fmt.Errorf("destination is not a regular file: %s", p)
		}
	}
	return p, nil
}

func readState(root, rel string) (state, error) {
	p, err := safePath(root, rel)
	if err != nil {
		return state{}, err
	}
	info, err := os.Lstat(p)
	if errors.Is(err, os.ErrNotExist) {
		return state{}, nil
	}
	if err != nil {
		return state{}, err
	}
	if info.Size() > 8<<20 {
		return state{}, fmt.Errorf("managed file too large to safely process: %s", rel)
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return state{}, err
	}
	return state{true, b, info.Mode().Perm()}, nil
}

func readManifest(root, rel string, global bool) (*manifest, error) {
	s, err := readState(root, rel)
	if err != nil || !s.exists {
		return nil, err
	}
	var m manifest
	if err := json.Unmarshal(s.data, &m); err != nil {
		return nil, fmt.Errorf("invalid install manifest: %w", err)
	}
	if m.Package != "git-workflow" || m.Version == "" || m.Files == nil {
		return nil, errors.New("unrecognized install manifest; review it before continuing")
	}
	for file, hash := range m.Files {
		decoded, err := hex.DecodeString(hash)
		if !owned(file, global) || err != nil || len(decoded) != 32 || strings.ToLower(hash) != hash {
			return nil, fmt.Errorf("unsafe install manifest entry: %s", file)
		}
	}
	for _, a := range m.Agents {
		if _, ok := agentBases[a]; !ok {
			return nil, fmt.Errorf("invalid manifest agent: %s", a)
		}
	}
	return &m, nil
}

const begin = "<!-- git-workflow:begin -->"
const end = "<!-- git-workflow:end -->"

func agentsContent(old []byte, remove, replace bool) ([]byte, error) {
	text := string(old)
	newline := "\n"
	if strings.Contains(text, "\r\n") {
		newline = "\r\n"
	}
	block := strings.ReplaceAll(kit.Blocks[0], "\n", newline)
	if !strings.Contains(text, begin) && !strings.Contains(text, end) {
		if remove {
			return old, nil
		}
		if text != "" && !strings.HasSuffix(text, newline+newline) {
			text += newline + newline
		}
		return []byte(text + block + newline), nil
	}
	if strings.Count(text, begin) != 1 || strings.Count(text, end) != 1 || strings.Index(text, end) < strings.Index(text, begin) {
		return nil, errors.New("malformed/duplicate Git Workflow block in AGENTS.md; review manually")
	}
	start, finish := strings.Index(text, begin), strings.Index(text, end)+len(end)
	// Existing pointer style belongs to that block, independent of surrounding user text.
	existingBlock := text[start:finish]
	newline = "\n"
	if strings.Contains(existingBlock, "\r\n") {
		newline = "\r\n"
	}
	block = strings.ReplaceAll(kit.Blocks[0], "\n", newline)
	known := false
	for _, candidate := range kit.Blocks {
		if text[start:finish] == strings.ReplaceAll(candidate, "\n", newline) {
			known = true
			break
		}
	}
	if !known && !replace {
		return nil, errors.New("Git Workflow block in AGENTS.md was edited; review it or use --replace with backup")
	}
	if remove {
		block = ""
	}
	result := []byte(text[:start] + block + text[finish:])
	if remove && len(bytes.TrimSpace(result)) == 0 {
		return nil, nil
	}
	return result, nil
}

func hasAgent(agents []string, name string) bool {
	for _, a := range agents {
		if a == name {
			return true
		}
	}
	return false
}

func semverCompare(a, b string) (int, error) {
	parse := func(s string) ([3]int, error) {
		var result [3]int
		parts := strings.Split(s, ".")
		if len(parts) != 3 {
			return result, fmt.Errorf("unsupported installer version: %s", s)
		}
		for i := range parts {
			n, err := strconv.Atoi(parts[i])
			if err != nil || n < 0 {
				return result, fmt.Errorf("invalid installer version: %s", s)
			}
			result[i] = n
		}
		return result, nil
	}
	x, err := parse(a)
	if err != nil {
		return 0, err
	}
	y, err := parse(b)
	if err != nil {
		return 0, err
	}
	for i := range x {
		if x[i] < y[i] {
			return -1, nil
		}
		if x[i] > y[i] {
			return 1, nil
		}
	}
	return 0, nil
}

func desiredFiles(global bool, agents []string) map[string][]byte {
	want := map[string][]byte{}
	for rel, data := range kit.Files {
		if !global {
			want[rel] = []byte(data)
			continue
		}
		suffix := strings.TrimPrefix(rel, ".agents/skills/git-workflow/")
		if suffix == "SKILL.md" {
			data = strings.Replace(data,
				"The portable project copy is managed by the standalone installer. Read `.agents/rules/git-workflow.md` and the matching command reference before acting.",
				"This copy applies to the current user's projects. Read `references/common.md` and the matching command reference before acting.", 1)
		}
		for _, a := range agents {
			want[agentBases[a]+"/"+suffix] = []byte(data)
		}
	}
	if !global {
		want[".agents/rules/git-workflow.md"] = []byte(kit.Rule)
	}
	return want
}

func plan(root, manifestRel string, old *manifest, o options, agents []string) ([]change, error) {
	if o.Action != "uninstall" && old != nil {
		compare, err := semverCompare(old.Version, kit.Version)
		if err != nil {
			return nil, err
		}
		if compare > 0 {
			return nil, errors.New("installed skill is newer than this CLI; update the executable first")
		}
	}
	want := desiredFiles(o.Global, agents)
	if o.Action == "uninstall" {
		want = map[string][]byte{}
	}
	newManifest := manifest{Package: "git-workflow", Version: kit.Version, Files: map[string]string{}, Agents: agents, AgentVersions: map[string]string{}}
	if old != nil {
		for rel, hash := range old.Files {
			newManifest.Files[rel] = hash
		}
		for a, ver := range old.AgentVersions {
			newManifest.AgentVersions[a] = ver
		}
		for _, a := range installedAgents(old, o.Global) {
			if newManifest.AgentVersions[a] == "" {
				newManifest.AgentVersions[a] = old.Version
			}
		}
	}
	var changes []change
	for _, rel := range sortedKeys(want) {
		before, err := readState(root, rel)
		if err != nil {
			return nil, err
		}
		after := want[rel]
		newManifest.Files[rel] = digest(after)
		if before.exists && bytes.Equal(before.data, after) {
			continue
		}
		if before.exists && !o.Replace && (old == nil || digest(before.data) != old.Files[rel]) {
			return nil, fmt.Errorf("edited/unmanaged package file %s; review it or use --replace", rel)
		}
		changes = append(changes, change{rel, before, after})
	}
	if old != nil {
		for _, rel := range sortedKeys(old.Files) {
			if _, keep := want[rel]; keep {
				continue
			}
			if o.Global && !hasAgent(agents, agentForPath(rel)) {
				continue // preserve unselected global integrations
			}
			before, err := readState(root, rel)
			if err != nil {
				return nil, err
			}
			if before.exists {
				if digest(before.data) != old.Files[rel] && !o.Replace {
					return nil, fmt.Errorf("edited package file %s; review it or use --replace", rel)
				}
				changes = append(changes, change{rel, before, nil})
			}
			delete(newManifest.Files, rel)
		}
	}
	if !o.Global {
		before, err := readState(root, "AGENTS.md")
		if err != nil {
			return nil, err
		}
		after, err := agentsContent(before.data, o.Action == "uninstall", o.Replace)
		if err != nil {
			return nil, err
		}
		if !(before.exists && bytes.Equal(before.data, after)) && !(after == nil && !before.exists) {
			changes = append(changes, change{"AGENTS.md", before, after})
		}
	}
	for _, a := range agents {
		if o.Action == "uninstall" {
			delete(newManifest.AgentVersions, a)
		} else {
			newManifest.AgentVersions[a] = kit.Version
		}
	}
	if o.Global {
		newManifest.Agents = installedAgents(&newManifest, true)
	}
	before, err := readState(root, manifestRel)
	if err != nil {
		return nil, err
	}
	var after []byte
	if len(newManifest.Files) > 0 {
		after, err = json.MarshalIndent(newManifest, "", "  ")
		if err != nil {
			return nil, err
		}
		after = append(after, '\n')
	}
	if !(before.exists && bytes.Equal(before.data, after)) && !(after == nil && !before.exists) {
		changes = append(changes, change{manifestRel, before, after}) // commit manifest last
	}
	return changes, nil
}

func sameState(a, b state) bool {
	return a.exists == b.exists && (!a.exists || bytes.Equal(a.data, b.data) && a.mode == b.mode)
}

func atomicWrite(root, rel string, data []byte, mode os.FileMode) error {
	p, err := safePath(root, rel)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return err
	}
	if _, err := safePath(root, rel); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(p), ".git-workflow-*")
	if err != nil {
		return err
	}
	temp := f.Name()
	defer os.Remove(temp)
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err := f.Chmod(mode); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if _, err := safePath(root, rel); err != nil {
		return err
	}
	return os.Rename(temp, p)
}

func backupChanges(plans []change, out io.Writer) error {
	any := false
	for _, p := range plans {
		any = any || p.before.exists
	}
	if !any {
		return nil
	}
	cache, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	base := filepath.Join(cache, "git-workflow", "backups")
	if err := os.MkdirAll(base, 0700); err != nil {
		return err
	}
	backup, err := os.MkdirTemp(base, "installation-")
	if err != nil {
		return err
	}
	type saved struct {
		SHA256 string `json:"sha256"`
		Mode   uint32 `json:"mode"`
	}
	meta := map[string]saved{}
	for _, p := range plans {
		if !p.before.exists {
			continue
		}
		if err := atomicWrite(backup, p.rel, p.before.data, 0600); err != nil {
			return err
		}
		check, err := readState(backup, p.rel)
		if err != nil || !bytes.Equal(check.data, p.before.data) {
			return errors.New("backup verification failed; no installation changes applied")
		}
		meta[p.rel] = saved{digest(p.before.data), uint32(p.before.mode)}
	}
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	if err := atomicWrite(backup, "backup-index.json", append(data, '\n'), 0600); err != nil {
		return err
	}
	fmt.Fprintln(out, "Exact previous-file backup:", backup)
	return nil
}

func apply(root string, plans []change, out io.Writer) error {
	for _, p := range plans {
		now, err := readState(root, p.rel)
		if err != nil {
			return err
		}
		if !sameState(now, p.before) {
			return fmt.Errorf("file changed during preflight: %s", p.rel)
		}
	}
	if err := backupChanges(plans, out); err != nil {
		return err
	}
	completed := []change{}
	rollback := func(cause error) error {
		var failures []string
		for i := len(completed) - 1; i >= 0; i-- {
			p := completed[i]
			now, err := readState(root, p.rel)
			if err != nil || (p.after == nil && now.exists) || (p.after != nil && (!now.exists || !bytes.Equal(now.data, p.after))) {
				failures = append(failures, p.rel)
				continue
			}
			if p.before.exists {
				err = atomicWrite(root, p.rel, p.before.data, p.before.mode)
			} else {
				var path string
				path, err = safePath(root, p.rel)
				if err == nil {
					err = os.Remove(path)
				}
			}
			if err != nil {
				failures = append(failures, p.rel)
			}
		}
		if len(failures) > 0 {
			return fmt.Errorf("%w; rollback needs review for: %s", cause, strings.Join(failures, ", "))
		}
		return fmt.Errorf("%w; completed file changes rolled back", cause)
	}
	for _, p := range plans {
		now, err := readState(root, p.rel)
		if err != nil {
			return rollback(err)
		}
		if !sameState(now, p.before) {
			return rollback(fmt.Errorf("file changed during installation: %s", p.rel))
		}
		if p.after == nil {
			var name string
			name, err = safePath(root, p.rel)
			if err == nil {
				err = os.Remove(name)
			}
		} else {
			mode := os.FileMode(0644)
			if p.before.exists {
				mode = p.before.mode
			}
			err = atomicWrite(root, p.rel, p.after, mode)
		}
		if err != nil {
			return rollback(err)
		}
		completed = append(completed, p)
	}
	return nil
}

func pruneOwnedDirs(root string, global bool, agents []string) {
	bases := []string{}
	if global {
		for _, a := range agents {
			bases = append(bases, agentBases[a])
		}
	} else {
		for _, name := range commands {
			bases = append(bases, ".agents/skills/"+name)
		}
	}
	var empty func(string)
	empty = func(p string) {
		info, err := os.Lstat(p)
		if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return
		}
		entries, err := os.ReadDir(p)
		if err != nil {
			return
		}
		for _, entry := range entries {
			if entry.IsDir() {
				empty(filepath.Join(p, entry.Name()))
			}
		}
		_ = os.Remove(p) // succeeds only for an empty directory
	}
	for _, base := range bases {
		empty(filepath.Join(root, filepath.FromSlash(base)))
	}
}

func doctor(root, manifestRel string, m *manifest, o options, out io.Writer) int {
	failures := 0
	fmt.Fprintf(out, "CLI/bundled skill: %s (%s/%s)\nTarget: %s\nManifest: %s\n", kit.Version, runtime.GOOS, runtime.GOARCH, root, manifestRel)
	if exe, err := os.Executable(); err == nil {
		fmt.Fprintln(out, "Executable:", exe)
	}
	if command, err := exec.LookPath("git-workflow"); err == nil {
		fmt.Fprintln(out, "PATH command:", command)
	} else {
		fmt.Fprintln(out, "PATH command: not found; run the bootstrap or use the executable's full path.")
	}
	if git, err := gitRead("--version"); err == nil {
		fmt.Fprintln(out, "Git:", git)
	} else {
		fmt.Fprintln(out, "Git: unavailable")
		failures++
	}
	if m == nil {
		fmt.Fprintln(out, "NOT INSTALLED: run git-workflow init with the same scope.")
		return 1
	}
	fmt.Fprintln(out, "Installed package:", m.Version)
	agents := installedAgents(m, o.Global)
	fmt.Fprintln(out, "Agents:", strings.Join(agents, ", "))
	for _, rel := range sortedKeys(m.Files) {
		if o.Global && o.Agent != "auto" && o.Agent != "all" && agentForPath(rel) != o.Agent {
			continue
		}
		s, err := readState(root, rel)
		if err != nil {
			fmt.Fprintln(out, "UNSAFE:", rel, err)
			failures++
		} else if !s.exists {
			fmt.Fprintln(out, "MISSING:", rel)
			failures++
		} else if digest(s.data) != m.Files[rel] {
			fmt.Fprintln(out, "MODIFIED:", rel)
			failures++
		}
	}
	if o.Global && o.Agent != "auto" && o.Agent != "all" && !hasAgent(agents, o.Agent) {
		fmt.Fprintln(out, "NOT INSTALLED agent:", o.Agent)
		failures++
	}
	if !o.Global {
		s, err := readState(root, "AGENTS.md")
		if err != nil || !s.exists {
			fmt.Fprintln(out, "MISSING/UNSAFE: AGENTS.md")
			failures++
		} else if changed, err := agentsContent(s.data, false, false); err != nil || !bytes.Equal(changed, s.data) {
			fmt.Fprintln(out, "INVALID: Git Workflow pointer in AGENTS.md")
			failures++
		}
	}
	if compare, err := semverCompare(m.Version, kit.Version); err == nil && compare < 0 {
		fmt.Fprintln(out, "UPDATE AVAILABLE FROM THIS CLI: run git-workflow update with the same scope.")
	}
	for _, a := range agents {
		if ver := m.AgentVersions[a]; ver != "" && ver != kit.Version {
			fmt.Fprintf(out, "Agent %s uses %s; CLI bundles %s.\n", a, ver, kit.Version)
		}
	}
	if failures > 0 {
		fmt.Fprintf(out, "Check failed: %d issue(s); no files changed.\n", failures)
		return 1
	}
	fmt.Fprintf(out, "OK: managed-file checks passed; no files changed.\n")
	return 0
}
