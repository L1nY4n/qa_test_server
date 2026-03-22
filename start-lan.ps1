param(
    [string]$LanIP = ""
)

$ErrorActionPreference = "Stop"

$projectRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$stopScript = Join-Path $projectRoot "stop-local.ps1"
$startScript = Join-Path $projectRoot "start-local.ps1"

Write-Host "Restarting local stack in LAN mode..."
& $stopScript

$startArgs = @{
    Lan = $true
}
if (-not [string]::IsNullOrWhiteSpace($LanIP)) {
    $startArgs["LanIP"] = $LanIP
}

& $startScript @startArgs
