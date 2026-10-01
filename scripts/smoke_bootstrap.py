#!/usr/bin/env python3
"""Exercise bootstrap installation with disposable user profiles and release files."""
import hashlib
import os
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[1]


def main():
    if os.name == 'nt':
        print('PowerShell bootstrap native test')
        with tempfile.TemporaryDirectory(prefix='git-workflow-bootstrap-') as temp:
            dest = Path(temp) / 'bin with spaces'
            cmd = ['powershell', '-NoProfile', '-File', str(ROOT / 'install.ps1'),
                   '-SourceDirectory', str(ROOT / 'dist'), '-BinDirectory', str(dest), '-NoPath']
            for _ in range(2):
                subprocess.run(cmd, check=True)
                subprocess.run([str(dest / 'git-workflow.exe'), '--version'], check=True)
            bad = Path(temp) / 'bad-source'
            shutil.copytree(ROOT / 'dist', bad)
            for item in bad.glob('git-workflow-windows-*.exe'):
                item.write_bytes(b'altered executable')
            before = hashlib.sha256((dest / 'git-workflow.exe').read_bytes()).hexdigest()
            bad_cmd = cmd[:]
            bad_cmd[bad_cmd.index('-SourceDirectory') + 1] = str(bad)
            result = subprocess.run(bad_cmd, capture_output=True, text=True)
            assert result.returncode != 0 and 'checksum mismatch' in result.stderr
            assert hashlib.sha256((dest / 'git-workflow.exe').read_bytes()).hexdigest() == before
        return
    print('Shell bootstrap native test')
    with tempfile.TemporaryDirectory(prefix='git-workflow-bootstrap-') as temp:
        home = Path(temp)
        dest = home / "bin with spaces and ' quote"
        env = dict(os.environ, HOME=str(home), SHELL='/bin/bash')
        cmd = ['sh', str(ROOT / 'install.sh'), '--source-dir', str(ROOT / 'dist'), '--bin-dir', str(dest)]
        for _ in range(2):
            subprocess.run(cmd, env=env, check=True)
            subprocess.run([str(dest / 'git-workflow'), '--version'], env=env, check=True)
        profile = (home / '.bashrc').read_text()
        assert profile.count('# git-workflow:path:begin') == 1
        result = subprocess.run(['sh', '-c', '. "$1"; command -v git-workflow', 'sh', str(home / '.bashrc')],
                                env=env, capture_output=True, text=True, check=True)
        assert result.stdout.strip() == str(dest / 'git-workflow'), result.stdout
        # An altered binary must fail checksum validation before destination/profile writes.
        bad = home / 'bad-source'
        bad.mkdir()
        for name in os.listdir(ROOT / 'dist'):
            item = ROOT / 'dist' / name
            if item.is_file():
                shutil.copy2(item, bad / name)
        for item in bad.glob('git-workflow-*'):
            item.write_bytes(b'altered executable')
        before = hashlib.sha256((dest / 'git-workflow').read_bytes()).hexdigest()
        result = subprocess.run(['sh', str(ROOT / 'install.sh'), '--source-dir', str(bad), '--bin-dir', str(dest)],
                                env=env, capture_output=True, text=True)
        assert result.returncode != 0 and 'checksum mismatch' in result.stderr
        assert hashlib.sha256((dest / 'git-workflow').read_bytes()).hexdigest() == before
    print('Bootstrap checks passed')


if __name__ == '__main__':
    main()
