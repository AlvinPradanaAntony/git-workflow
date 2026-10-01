# Git Workflow development

The terminal CLI installs the bundled skill. Git actions such as commitmsg and gitrelease run through an agent and are not CLI subcommands.

`VERSION` is the distribution version. `skill/` contains the exported skill source. After editing it or the Python installer template, run `python scripts/package.py` to refresh `bundle.json` and `install_git_workflow.py`. Keep the existing 12 Git commands and help behavior intact.

Validate with `python scripts/package.py --check`, `python -m unittest discover -s tests -v`, `go test ./...`, and `go vet ./...`. Build with `python scripts/build.py`; native bootstrap smoke tests use `python scripts/build.py --native-only` and `python scripts/smoke_bootstrap.py`. GitHub Actions runs native tests on Windows/Linux/macOS, amd64 and arm64.

Preserve unrelated user files during install/update/uninstall. Do not disable manifest checks, backups, symlink checks, or rollback. Release publication is handled only by the tag workflow after all assets pass validation. Do not publish or replace a release while implementing a code change unless the user requests it.
