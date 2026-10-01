package main

import (
	"bufio"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

//go:embed bundle.json
var bundleJSON []byte

type bundle struct {
	Version string            `json:"version"`
	Files   map[string]string `json:"files"`
	Rule    string            `json:"rule"`
	Blocks  []string          `json:"agents_blocks"`
}

var kit bundle

var agentBases = map[string]string{
	"codex":           ".agents/skills/git-workflow",
	"antigravity":     ".gemini/config/skills/git-workflow",
	"antigravity-cli": ".gemini/antigravity-cli/skills/git-workflow",
}
var allAgents = []string{"codex", "antigravity", "antigravity-cli"}
var commands = []string{"gitstatus", "commitpln", "commitmsg", "branchname", "prdesc", "gitundo", "gitreset", "gitmergecancel", "gitamend", "gitpushamend", "gitpush", "gitrelease", "commitplan", "git-workflow"}

type options struct {
	Action  string
	Project string
	Global  bool
	Agent   string
	DryRun  bool
	Replace bool
}

func loadKit() error {
	if err := json.Unmarshal(bundleJSON, &kit); err != nil {
		return err
	}
	if kit.Version == "" || len(kit.Blocks) == 0 || len(kit.Files) == 0 {
		return errors.New("invalid embedded package")
	}
	for rel := range kit.Files {
		if !validRel(rel) || !strings.HasPrefix(rel, ".agents/skills/git-workflow/") {
			return fmt.Errorf("invalid embedded package path: %s", rel)
		}
	}
	return nil
}

func usage(w io.Writer) {
	fmt.Fprintf(w, `Git Workflow %s — skill installer for Codex and Antigravity

Usage: git-workflow COMMAND [OPTIONS]

  init       Install into the Git root of the terminal's current directory
  update     Update an existing installation from this CLI's bundled skill
  uninstall  Remove only manifest-owned skill files; keep unrelated files
  doctor     Check executable, Git, installed version and managed-file integrity
  help       Show installer help (agent command help is /git-workflow help)

Options:
  -g, --global        Current user across projects; no repository files changed
  --project PATH      Select another Git repository or one of its subfolders
  --agent NAME        auto (default), codex, antigravity, antigravity-cli, all
  --dry-run           Show proposed changes without writing files
  --replace           Back up and replace/remove edited managed package files
  --version           Print CLI and bundled skill version
  -h, --help          Show this help

Examples:
  git-workflow init
  git-workflow init -g --agent antigravity
  git-workflow update
  git-workflow uninstall -g
  git-workflow doctor

init installs immediately; use --dry-run to preview. update does not download
the latest CLI: rerun the bootstrap or download a new executable first.
The CLI installs skills; it does not run commitmsg or gitrelease.
`, kit.Version)
}

func parse(args []string, w io.Writer) (options, error) {
	var o options
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		if len(args) > 1 {
			return o, errors.New("installer help does not accept agent help options; pass --cmd/--detailed to your agent")
		}
		o.Action = "help"
		return o, nil
	}
	if args[0] == "--version" || args[0] == "version" {
		if len(args) != 1 {
			return o, errors.New("version does not accept arguments")
		}
		o.Action = "version"
		return o, nil
	}
	o.Action = args[0]
	if o.Action != "init" && o.Action != "update" && o.Action != "uninstall" && o.Action != "doctor" {
		return o, fmt.Errorf("unknown installer command %q; use git-workflow help; Git actions run through your agent", o.Action)
	}
	fs := flag.NewFlagSet(o.Action, flag.ContinueOnError)
	fs.SetOutput(w)
	fs.BoolVar(&o.Global, "g", false, "current user")
	fs.BoolVar(&o.Global, "global", false, "current user")
	fs.StringVar(&o.Project, "project", "", "repository path")
	fs.StringVar(&o.Agent, "agent", "auto", "agent selection")
	fs.BoolVar(&o.DryRun, "dry-run", false, "preview")
	fs.BoolVar(&o.Replace, "replace", false, "back up modified files")
	if err := fs.Parse(args[1:]); err != nil {
		return o, err
	}
	if len(fs.Args()) != 0 {
		return o, errors.New("unexpected positional arguments")
	}
	if o.Global && o.Project != "" {
		return o, errors.New("choose either --global or --project")
	}
	if o.Agent != "auto" && o.Agent != "all" {
		if _, ok := agentBases[o.Agent]; !ok {
			return o, fmt.Errorf("unknown agent %q; choose codex, antigravity, antigravity-cli or all", o.Agent)
		}
	}
	if o.Action == "doctor" && (o.DryRun || o.Replace) {
		return o, errors.New("doctor is read-only; --dry-run and --replace do not apply")
	}
	return o, nil
}

func gitRead(args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", args...)
	b, err := cmd.Output()
	if err != nil {
		return "", errors.New("Git could not resolve a working repository; run inside a repo, use --project PATH, or choose -g")
	}
	return strings.TrimRight(string(b), "\r\n"), nil
}

func target(o options) (string, string, error) {
	var root string
	if o.Global {
		p, err := os.UserHomeDir()
		if err != nil {
			return "", "", err
		}
		root = p
	} else {
		p := o.Project
		if p == "" {
			var err error
			p, err = os.Getwd()
			if err != nil {
				return "", "", err
			}
		}
		if p == "~" || strings.HasPrefix(p, "~/") || strings.HasPrefix(p, "~\\") {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", "", err
			}
			p = filepath.Join(home, strings.TrimLeft(p[1:], "/\\"))
		}
		p, err := filepath.Abs(p)
		if err != nil {
			return "", "", err
		}
		root, err = gitRead("-C", p, "rev-parse", "--show-toplevel")
		if err != nil {
			return "", "", err
		}
	}
	root, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", "", err
	}
	root, err = filepath.Abs(root)
	if err != nil {
		return "", "", err
	}
	manifestRel := ".agents/git-workflow-install.json"
	if o.Global {
		manifestRel = ".config/git-workflow/install.json"
	}
	return root, manifestRel, nil
}

func agentForPath(rel string) string {
	for name, base := range agentBases {
		if strings.HasPrefix(rel, base+"/") {
			return name
		}
	}
	return ""
}

func installedAgents(m *manifest, global bool) []string {
	if m == nil {
		return nil
	}
	if !global {
		if len(m.Agents) > 0 {
			return append([]string{}, m.Agents...)
		}
		return []string{"codex", "antigravity"}
	}
	seen := map[string]bool{}
	for rel := range m.Files {
		if a := agentForPath(rel); a != "" {
			seen[a] = true
		}
	}
	var found []string
	for _, name := range allAgents {
		if seen[name] {
			found = append(found, name)
		}
	}
	return found
}

func selectAgents(o options, root string, m *manifest, in io.Reader, out io.Writer) ([]string, error) {
	if o.Agent == "all" {
		return append([]string{}, allAgents...), nil
	}
	if o.Agent != "auto" {
		return []string{o.Agent}, nil
	}
	if existing := installedAgents(m, o.Global); len(existing) > 0 {
		return existing, nil
	}
	if !o.Global {
		return []string{"codex", "antigravity"}, nil // shared portable project location
	}
	var candidates []string
	for _, item := range []struct{ name, config string }{{"codex", ".codex"}, {"antigravity", ".gemini/config"}, {"antigravity-cli", ".gemini/antigravity-cli"}} {
		if info, err := os.Stat(filepath.Join(root, filepath.FromSlash(item.config))); err == nil && info.IsDir() {
			candidates = append(candidates, item.name)
		}
	}
	if len(candidates) == 1 {
		return candidates, nil
	}
	if f, ok := in.(*os.File); ok {
		if info, err := f.Stat(); err != nil || info.Mode()&os.ModeCharDevice == 0 {
			return nil, errors.New("global agent target is unclear; pass --agent codex, antigravity, antigravity-cli or all; no files changed")
		}
	}
	fmt.Fprintf(out, "Global agent target needs a choice. Detected: %s\nChoose codex, antigravity, antigravity-cli or all: ", strings.Join(candidates, ", "))
	line, err := bufio.NewReader(in).ReadString('\n')
	if err != nil && err != io.EOF {
		return nil, err
	}
	line = strings.TrimSpace(line)
	if line == "all" {
		return append([]string{}, allAgents...), nil
	}
	if _, ok := agentBases[line]; !ok {
		return nil, errors.New("no valid agent chosen; no files changed")
	}
	return []string{line}, nil
}

func run(args []string, in io.Reader, out, errOut io.Writer) int {
	if err := loadKit(); err != nil {
		fmt.Fprintln(errOut, "Invalid package:", err)
		return 2
	}
	o, err := parse(args, errOut)
	if errors.Is(err, flag.ErrHelp) {
		usage(out)
		return 0
	}
	if err != nil {
		fmt.Fprintln(errOut, "Stopped:", err)
		return 2
	}
	if o.Action == "help" {
		usage(out)
		return 0
	}
	if o.Action == "version" {
		fmt.Fprintf(out, "git-workflow %s (%s/%s)\n", kit.Version, runtime.GOOS, runtime.GOARCH)
		return 0
	}
	root, manifestRel, err := target(o)
	if err != nil {
		fmt.Fprintln(errOut, "Stopped:", err)
		return 2
	}
	m, err := readManifest(root, manifestRel, o.Global)
	if err != nil {
		fmt.Fprintln(errOut, "Stopped:", err)
		return 2
	}
	if o.Action == "doctor" {
		return doctor(root, manifestRel, m, o, out)
	}
	if m == nil && o.Action != "init" {
		if o.Action == "uninstall" {
			fmt.Fprintln(out, "No managed installation found; no files changed.")
			return 0
		}
		fmt.Fprintln(errOut, "No managed installation found; run git-workflow init first.")
		return 2
	}
	agents, err := selectAgents(o, root, m, in, out)
	if err != nil {
		fmt.Fprintln(errOut, "Stopped:", err)
		return 2
	}
	plans, err := plan(root, manifestRel, m, o, agents)
	if err != nil {
		fmt.Fprintln(errOut, "Stopped:", err, "\nNo files changed.")
		return 2
	}
	fmt.Fprintf(out, "Git Workflow %s: %s; target %s; agents %s\n", kit.Version, o.Action, root, strings.Join(agents, ", "))
	for _, p := range plans {
		verb := "CREATE"
		if p.after == nil {
			verb = "REMOVE"
		} else if p.before.exists {
			verb = "UPDATE"
		}
		fmt.Fprintln(out, verb, p.rel)
	}
	if len(plans) == 0 {
		fmt.Fprintln(out, "Already up to date; no changes needed.")
		return 0
	}
	if o.DryRun {
		fmt.Fprintln(out, "Dry run only; no files changed.")
		return 0
	}
	if err := apply(root, plans, out); err != nil {
		fmt.Fprintln(errOut, "Stopped:", err)
		return 2
	}
	if o.Action == "uninstall" {
		pruneOwnedDirs(root, o.Global, agents)
		fmt.Fprintln(out, "Uninstalled managed files; unrelated files preserved.")
	} else {
		fmt.Fprintln(out, "Installed. Codex: $git-workflow commitmsg; Antigravity: /git-workflow commitmsg.")
	}
	fmt.Fprintln(out, "No Git staging, commits, pushes, hooks or dependency installs were run.")
	return 0
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
