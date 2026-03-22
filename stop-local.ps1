$ErrorActionPreference = "SilentlyContinue"

$projectRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$tmpDir = Join-Path $projectRoot "tmp"

function Stop-ByPidFile {
    param(
        [string]$PidFile
    )
    if (Test-Path $PidFile) {
        $procIdRaw = Get-Content $PidFile -ErrorAction SilentlyContinue
        if ($procIdRaw) {
            $procId = [int]$procIdRaw
            if ($procId -ne $PID) {
                Stop-Process -Id $procId -Force -ErrorAction SilentlyContinue
            }
        }
        Remove-Item $PidFile -Force -ErrorAction SilentlyContinue
    }
}

Stop-ByPidFile (Join-Path $tmpDir "frontend.pid")
Stop-ByPidFile (Join-Path $tmpDir "backend.pid")
Stop-ByPidFile (Join-Path $tmpDir "mariadb.pid")

# Clean up stale backend processes started by `go run .` (child qa_test_server.exe),
# which might not be tracked by pid files.
$projectPattern = [Regex]::Escape($projectRoot)
$candidates = Get-CimInstance Win32_Process -ErrorAction SilentlyContinue | Where-Object {
    ($_.Name -in @("qa_test_server.exe", "go.exe")) -and
    $_.CommandLine -and
    $_.CommandLine -match $projectPattern
}
foreach ($proc in $candidates) {
    Stop-Process -Id $proc.ProcessId -Force -ErrorAction SilentlyContinue
}

foreach ($port in @(5173, 8080, 4001, 3306)) {
    $conns = Get-NetTCPConnection -LocalPort $port -State Listen -ErrorAction SilentlyContinue
    foreach ($conn in $conns) {
        $proc = Get-Process -Id $conn.OwningProcess -ErrorAction SilentlyContinue
        if ($proc) {
            if ($proc.ProcessName -in @("node", "qa_test_server", "qa_test_server-local", "mariadbd", "go")) {
                Stop-Process -Id $proc.Id -Force -ErrorAction SilentlyContinue
            }
        }
    }
}

Write-Host "Local stack stopped."
