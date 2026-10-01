# Install the native CLI once; then run git-workflow init in any repository.
[CmdletBinding()]
param(
    [string]$Repo = $env:GIT_WORKFLOW_REPO,
    [string]$Version = 'latest',
    [string]$SourceDirectory,
    [string]$BinDirectory = (Join-Path $env:LOCALAPPDATA 'Programs\GitWorkflow\bin'),
    [switch]$NoPath,
    [switch]$Replace
)
$ErrorActionPreference = 'Stop'
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

function Get-GitWorkflowHash([string]$Path) {
    $stream = [IO.File]::OpenRead($Path)
    $sha = [Security.Cryptography.SHA256]::Create()
    try {
        return [BitConverter]::ToString($sha.ComputeHash($stream)).Replace('-', '').ToLowerInvariant()
    } finally {
        $sha.Dispose()
        $stream.Dispose()
    }
}

function Save-GitWorkflowDownload([string]$Uri, [string]$Path) {
    $client = [Net.WebClient]::new()
    try { $client.DownloadFile($Uri, $Path) } finally { $client.Dispose() }
}

if ($env:OS -ne 'Windows_NT') { throw 'Use install.sh on Linux/macOS.' }
$cpu = $env:PROCESSOR_ARCHITEW6432
if (-not $cpu) { $cpu = $env:PROCESSOR_ARCHITECTURE }
switch ($cpu.ToUpperInvariant()) {
    'AMD64' { $arch = 'amd64' }
    'ARM64' { $arch = 'arm64' }
    default { throw "Unsupported architecture: $cpu" }
}
$asset = "git-workflow-windows-$arch.exe"
if (-not $SourceDirectory -and -not $Repo) {
    if ((Test-Path -LiteralPath (Join-Path $PSScriptRoot $asset)) -and
        (Test-Path -LiteralPath (Join-Path $PSScriptRoot 'SHA256SUMS'))) {
        $SourceDirectory = $PSScriptRoot
    } else {
        $Repo = 'AlvinPradanaAntony/git-workflow'
    }
}
if ($SourceDirectory -and $Repo) { throw 'Choose -SourceDirectory or -Repo.' }
if (-not [IO.Path]::IsPathRooted($BinDirectory)) { throw 'BinDirectory must be an absolute path.' }
$destination = Join-Path $BinDirectory 'git-workflow.exe'
$receipt = Join-Path $BinDirectory '.git-workflow-cli.sha256'
foreach ($item in @($BinDirectory, $destination, $receipt)) {
    if (Test-Path -LiteralPath $item) {
        if ((Get-Item -LiteralPath $item -Force).Attributes -band [IO.FileAttributes]::ReparsePoint) {
            throw "Refusing a reparse-point destination: $item"
        }
    }
}
if ($Repo) {
    if ($Repo -notmatch '^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$') { throw 'Expected GitHub OWNER/REPOSITORY.' }
    if ($Version -eq 'latest') { $base = "https://github.com/$Repo/releases/latest/download" }
    elseif ($Version -match '^v[0-9][A-Za-z0-9_.-]*$') { $base = "https://github.com/$Repo/releases/download/$Version" }
    else { throw 'Use latest or a version tag such as v2.13.0.' }
}

$tempDir = Join-Path ([IO.Path]::GetTempPath()) ('git-workflow-' + [Guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $tempDir | Out-Null
try {
    $binary = Join-Path $tempDir $asset
    $sums = Join-Path $tempDir 'SHA256SUMS'
    if ($SourceDirectory) {
        Copy-Item -LiteralPath (Join-Path $SourceDirectory $asset) -Destination $binary
        Copy-Item -LiteralPath (Join-Path $SourceDirectory 'SHA256SUMS') -Destination $sums
    } else {
        Save-GitWorkflowDownload "$base/$asset" $binary
        Save-GitWorkflowDownload "$base/SHA256SUMS" $sums
    }
    $pattern = '^([0-9a-f]{64})\s+' + [Regex]::Escape($asset) + '$'
    $found = @(Get-Content -LiteralPath $sums | Where-Object { $_ -cmatch $pattern })
    if ($found.Count -ne 1) { throw "Missing/duplicate SHA256SUMS entry for $asset." }
    $expected = [Regex]::Match($found[0], $pattern).Groups[1].Value
    $actual = Get-GitWorkflowHash $binary
    if ($expected -ne $actual) { throw 'Executable checksum mismatch; nothing installed.' }

    New-Item -ItemType Directory -Force -Path $BinDirectory | Out-Null
    $installNeeded = $true
    if (Test-Path -LiteralPath $destination) {
        if (-not (Test-Path -LiteralPath $destination -PathType Leaf)) { throw 'Destination is not a regular file.' }
        $oldHash = Get-GitWorkflowHash $destination
        if ($oldHash -eq $expected) { $installNeeded = $false }
        if ($oldHash -ne $expected) {
            $recorded = ''
            if (Test-Path -LiteralPath $receipt) { $recorded = (Get-Content -Raw -LiteralPath $receipt).Trim() }
            if ($recorded -ne $oldHash -and -not $Replace) {
                throw 'Existing executable is unmanaged/edited; review it or use -Replace.'
            }
        }
    }
    if ($installNeeded) {
        $tempBinary = Join-Path $BinDirectory ('.git-workflow-install-' + [Guid]::NewGuid().ToString('N') + '.exe')
        Copy-Item -LiteralPath $binary -Destination $tempBinary
        try {
            if (Test-Path -LiteralPath $destination) {
                $backup = Join-Path $BinDirectory ('git-workflow.backup.' + [Guid]::NewGuid().ToString('N') + '.exe')
                [IO.File]::Replace($tempBinary, $destination, $backup)
                Write-Output "Previous executable backup: $backup"
            } else {
                [IO.File]::Move($tempBinary, $destination)
            }
        } finally {
            if (Test-Path -LiteralPath $tempBinary) { Remove-Item -LiteralPath $tempBinary -Force }
        }
    } else {
        Write-Output 'Executable is already up to date.'
    }
    [IO.File]::WriteAllText($receipt, $expected + [Environment]::NewLine, [Text.Encoding]::ASCII)
    if (-not $NoPath) {
        $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
        $entries = @($userPath -split ';' | Where-Object { $_ })
        $normal = $BinDirectory.TrimEnd('\')
        $present = @($entries | Where-Object { $_.TrimEnd('\') -ieq $normal }).Count -gt 0
        if (-not $present) {
            $newPath = (($entries + $BinDirectory) -join ';')
            [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
        }
        if (@($env:Path -split ';' | Where-Object { $_.TrimEnd('\') -ieq $normal }).Count -eq 0) {
            $env:Path = $BinDirectory + ';' + $env:Path
        }
        Write-Output 'PATH configured for this user and the current PowerShell session.'
    }
    & $destination --version
    if ($LASTEXITCODE -ne 0) { throw 'Installed executable verification failed.' }
    Write-Output "CLI installed at $destination. In a repository run: git-workflow init"
} finally {
    Remove-Item -LiteralPath $tempDir -Recurse -Force
}
