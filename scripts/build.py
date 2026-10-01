#!/usr/bin/env python3
"""Build native Git Workflow release executables and SHA256SUMS (stdlib only)."""
import argparse
import hashlib
import json
import os
from pathlib import Path
import platform
import shutil
import subprocess

ROOT = Path(__file__).resolve().parents[1]
TARGETS = [('windows', 'amd64'), ('windows', 'arm64'), ('linux', 'amd64'),
           ('linux', 'arm64'), ('darwin', 'amd64'), ('darwin', 'arm64')]


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--go', default='go')
    parser.add_argument('--native-only', action='store_true')
    args = parser.parse_args()
    targets = TARGETS
    if args.native_only:
        goos = {'Windows': 'windows', 'Linux': 'linux', 'Darwin': 'darwin'}[platform.system()]
        goarch = {'x86_64': 'amd64', 'AMD64': 'amd64', 'aarch64': 'arm64', 'arm64': 'arm64', 'ARM64': 'arm64'}[platform.machine()]
        targets = [(goos, goarch)]
    dist = ROOT / 'dist'
    dist.mkdir(exist_ok=True)
    version = json.loads((ROOT / 'bundle.json').read_text(encoding='utf-8'))['version']
    names = []
    for goos, goarch in targets:
        name = f'git-workflow-{goos}-{goarch}' + ('.exe' if goos == 'windows' else '')
        env = dict(os.environ, GOOS=goos, GOARCH=goarch, CGO_ENABLED='0')
        subprocess.run([args.go, 'build', '-buildvcs=false', '-trimpath', '-ldflags=-s -w',
                        '-o', str(dist / name), '.'], cwd=ROOT, env=env, check=True)
        names.append(name)
        print(f'Built {name}', flush=True)
    for name in ('install.sh', 'install.ps1', 'install_git_workflow.py'):
        shutil.copy2(ROOT / name, dist / name)
        names.append(name)
    (dist / 'SHA256SUMS').write_text(''.join(
        hashlib.sha256((dist / name).read_bytes()).hexdigest() + '  ' + name + '\n'
        for name in sorted(names)), encoding='ascii')
    (dist / 'VERSION').write_text(version + '\n', encoding='ascii')


if __name__ == '__main__':
    main()
