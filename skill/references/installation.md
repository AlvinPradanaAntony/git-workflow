# Git Workflow installation

## Native CLI (2.13.0)
Use the standalone `git-workflow` executable on PATH for Windows, Linux or macOS, amd64 and arm64. It embeds the skill, needs no Python/Node/Go runtime, and uses an existing Git installation for project discovery. The distribution source is https://github.com/AlvinPradanaAntony/git-workflow. Bootstrap install.ps1/install.sh installs the CLI once in the current user's bin directory; it does not itself install the skill or change a repository. Use adjacent downloaded binaries with SHA256SUMS, --source-dir/-SourceDirectory, or the default distribution repository. --repo/-Repo overrides that repository. Verify release availability before recommending a download URL; a committed workflow does not itself mean a public release exists.

| Terminal command | Behavior |
| --- | --- |
| `git-workflow init` | Install the skill into the Git root resolved from cwd; works in subfolders and linked worktrees. |
| `git-workflow init -g` | Install globally; reuse prior agent targets or detect one clear config. Ask for a choice if unclear; in noninteractive use require --agent. |
| `git-workflow init -g --agent antigravity` | Explicit Antigravity IDE global target. Other choices: codex, antigravity-cli, all. |
| `git-workflow update` | Refresh an existing installation using the skill embedded in the current executable. It does not fetch a newer CLI. |
| `git-workflow update -g` | Refresh the existing global integrations. --agent limits changes and preserves unselected integrations. |
| `git-workflow uninstall` | Back up and remove only manifest-owned project files and the known AGENTS.md pointer. |
| `git-workflow uninstall -g` | Remove managed global files, optionally limited by --agent. |
| `git-workflow doctor` | Check version, target, Git, managed file hashes and project pointer; read-only. Add -g for user scope. |

init/update/uninstall apply immediately. Use --dry-run for a preview and --replace only after reviewing local package edits; it preserves verified exact backups before replacement/removal. --project PATH selects another Git worktree from any folder, including a subfolder. Reject --project with -g. Agent-specific global directories are ~/.agents/skills/git-workflow (Codex), ~/.gemini/config/skills/git-workflow (Antigravity IDE), ~/.gemini/antigravity-cli/skills/git-workflow (Antigravity CLI). Project installs remain one shared .agents/skills/git-workflow skill, .agents/rules/git-workflow.md and a conditional pointer in AGENTS.md. Global installs do not edit repository files or global AGENTS.md.

Download a newer executable or rerun bootstrap to refresh the CLI, then run update at the intended project/user scope. Preserve customized files and stop on invalid manifests, unsafe paths, missing authorization or downgrade attempts. Do not claim an edited file is safe merely because its hash appears in a manifest.

## Migration from Python
The CLI reads existing project/global manifests from Python 2.x and removes obsolete files only when they are recorded as package-owned. Existing user changes stop update/uninstall unless --replace was explicitly requested. Previous backups remain usable. The Python fallback retains preview-first semantics:

```text
python /path/to/install_git_workflow.py --apply
python /path/to/install_git_workflow.py --global --apply
python /path/to/install_git_workflow.py --uninstall --apply
```

Python resolves the Git root from cwd or --project, including subfolders; the installer may be saved anywhere. --global installs all three integrations and retains preview-first semantics; --help and --version are available. The native CLI is the preferred route. CLI installer help is `git-workflow help`; agent command documentation remains `$git-workflow help` / `/git-workflow help`, including --detailed and --cmd. These are different interfaces. Do not run commitmsg/gitrelease as terminal CLI subcommands; pass them to the agent.

The distribution workflow tests native Windows/Linux/macOS targets, builds all six executables, includes install_git_workflow.py and SHA256SUMS, and uses the selected changelog version for release notes. Branch pushes/manual branch runs build artifacts; a matching version tag publishes a Release after complete asset verification. Updating the installer does not authorize publishing a new tag or replacing a Release.
