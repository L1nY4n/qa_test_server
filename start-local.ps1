param(
    [switch]$Lan,
    [string]$LanIP = "",
    [switch]$SkipFrontendBuild
)

$ErrorActionPreference = "Stop"

$projectRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$workspaceRoot = Split-Path -Parent $projectRoot
$tmpDir = Join-Path $projectRoot "tmp"
$dataDir = Join-Path $workspaceRoot "mariadb-data"
$myIni = Join-Path $dataDir "my.ini"
$backendExe = Join-Path $tmpDir "qa_test_server-local.exe"
$frontendDir = Join-Path $projectRoot "qa-test-web"
$templatesDir = Join-Path $projectRoot "templates"
$templatesAssetsDir = Join-Path $templatesDir "assets"
$backendPidFile = Join-Path $tmpDir "backend.pid"
$frontendPidFile = Join-Path $tmpDir "frontend.pid"

$goExe = "C:\Program Files\Go\bin\go.exe"
$mariadbdExe = "C:\Program Files\MariaDB 12.2\bin\mariadbd.exe"
$mysqlExe = "C:\Program Files\MariaDB 12.2\bin\mysql.exe"
$mysqlInstallDbExe = "C:\Program Files\MariaDB 12.2\bin\mysql_install_db.exe"
$pnpmCmd = Join-Path $env:APPDATA "npm\pnpm.cmd"
$npxCmd = "C:\Program Files\nodejs\npx.cmd"
$dbPassword = "L1nFen9.com"

foreach ($path in @($goExe, $mariadbdExe, $mysqlExe, $mysqlInstallDbExe, $pnpmCmd, $npxCmd)) {
    if (-not (Test-Path $path)) {
        throw "Missing required executable: $path"
    }
}

function Test-IsAdmin {
    $identity = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($identity)
    return $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

function Stop-ByPidFile {
    param(
        [string]$PidFile
    )

    if (-not (Test-Path $PidFile)) {
        return
    }

    $procIdRaw = Get-Content $PidFile -ErrorAction SilentlyContinue
    if ($procIdRaw) {
        $procId = [int]$procIdRaw
        if ($procId -ne $PID) {
            Stop-Process -Id $procId -Force -ErrorAction SilentlyContinue
        }
    }
    Remove-Item $PidFile -Force -ErrorAction SilentlyContinue
}

function Stop-LocalServiceProcesses {
    Stop-ByPidFile -PidFile $backendPidFile
    Stop-ByPidFile -PidFile $frontendPidFile

    foreach ($port in @(8080, 5173)) {
        $conns = Get-NetTCPConnection -LocalPort $port -State Listen -ErrorAction SilentlyContinue
        foreach ($conn in $conns) {
            $proc = Get-Process -Id $conn.OwningProcess -ErrorAction SilentlyContinue
            if ($proc -and $proc.ProcessName -in @("node", "qa_test_server", "qa_test_server-local", "go")) {
                Stop-Process -Id $proc.Id -Force -ErrorAction SilentlyContinue
            }
        }
    }
}

function Get-PreferredIPv4 {
    $route = Get-NetRoute -AddressFamily IPv4 -DestinationPrefix "0.0.0.0/0" -ErrorAction SilentlyContinue |
        Sort-Object -Property RouteMetric, InterfaceMetric |
        Select-Object -First 1

    if ($route) {
        $addr = Get-NetIPAddress -AddressFamily IPv4 -InterfaceIndex $route.InterfaceIndex -ErrorAction SilentlyContinue |
            Where-Object { $_.IPAddress -notlike "169.254*" -and $_.IPAddress -ne "127.0.0.1" } |
            Select-Object -First 1
        if ($addr) {
            return $addr.IPAddress
        }
    }

    $fallback = Get-NetIPAddress -AddressFamily IPv4 -ErrorAction SilentlyContinue |
        Where-Object { $_.IPAddress -notlike "169.254*" -and $_.IPAddress -ne "127.0.0.1" } |
        Select-Object -First 1
    if ($fallback) {
        return $fallback.IPAddress
    }

    return ""
}

function Ensure-FirewallProgramRule {
    param(
        [string]$RuleName,
        [string]$ProgramPath,
        [string]$Ports = ""
    )

    if (-not (Test-Path $ProgramPath)) {
        return
    }

    & netsh advfirewall firewall delete rule name="$RuleName" program="$ProgramPath" | Out-Null
    if ([string]::IsNullOrWhiteSpace($Ports)) {
        & netsh advfirewall firewall add rule name="$RuleName" dir=in action=allow program="$ProgramPath" enable=yes profile=any | Out-Null
    } else {
        & netsh advfirewall firewall add rule name="$RuleName" dir=in action=allow program="$ProgramPath" protocol=TCP localport="$Ports" enable=yes profile=any | Out-Null
    }
}

function Ensure-FirewallPortRule {
    param(
        [string]$RuleName,
        [string]$Ports
    )

    if ([string]::IsNullOrWhiteSpace($Ports)) {
        return
    }

    & netsh advfirewall firewall delete rule name="$RuleName" | Out-Null
    & netsh advfirewall firewall add rule name="$RuleName" dir=in action=allow protocol=TCP localport="$Ports" enable=yes profile=any | Out-Null
}

function Ensure-MyIniLocalBind {
    param(
        [string]$IniPath
    )

    if (-not (Test-Path $IniPath)) {
        return
    }

    $content = Get-Content -Path $IniPath -Raw
    $updated = $content -replace '(?im)^\s*bind-address\s*=.*$', 'bind-address=127.0.0.1'
    if ($updated -eq $content) {
        $updated = $content -replace '(?im)^\[mysqld\]\s*$', "[mysqld]`r`nbind-address=127.0.0.1"
    }

    if ($updated -ne $content) {
        Set-Content -Path $IniPath -Value $updated -Encoding ASCII
    }
}

$httpBindAddr = "127.0.0.1:8080"
$tcpBindAddr = "127.0.0.1:4001"
$viteHost = "127.0.0.1"
$displayHost = "localhost"
$frontendEnvFile = Join-Path $projectRoot "qa-test-web\.env.development.local"

if ($Lan) {
    if ([string]::IsNullOrWhiteSpace($LanIP)) {
        $LanIP = Get-PreferredIPv4
    }
    if ([string]::IsNullOrWhiteSpace($LanIP)) {
        throw "Unable to detect LAN IPv4. Please rerun with -LanIP <your-ip>."
    }

    $httpBindAddr = "0.0.0.0:8080"
    $tcpBindAddr = "0.0.0.0:4001"
    $viteHost = "0.0.0.0"
    $displayHost = $LanIP

    Set-Content -Path $frontendEnvFile -Value "VITE_APP_BASE_URL=${LanIP}:8080`r`n" -Encoding ASCII
}

New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null
New-Item -ItemType Directory -Path $dataDir -Force | Out-Null

if (-not (Test-Path $myIni)) {
    $installOut = Join-Path $tmpDir "mariadb-install.out.log"
    $installErr = Join-Path $tmpDir "mariadb-install.err.log"
    & $mysqlInstallDbExe "--datadir=$dataDir" "--port=3306" "--password=$dbPassword" 1>$installOut 2>$installErr
}
Ensure-MyIniLocalBind -IniPath $myIni

$mariaListening = Get-NetTCPConnection -LocalPort 3306 -State Listen -ErrorAction SilentlyContinue
if (-not $mariaListening) {
    $mariaOut = Join-Path $tmpDir "mariadb.out.log"
    $mariaErr = Join-Path $tmpDir "mariadb.err.log"
    $mariaProc = Start-Process -FilePath $mariadbdExe -ArgumentList "--defaults-file=$myIni" -WindowStyle Hidden -RedirectStandardOutput $mariaOut -RedirectStandardError $mariaErr -PassThru
    Set-Content -Path (Join-Path $tmpDir "mariadb.pid") -Value $mariaProc.Id
}

for ($i = 0; $i -lt 15; $i++) {
    if (Get-NetTCPConnection -LocalPort 3306 -State Listen -ErrorAction SilentlyContinue) {
        break
    }
    Start-Sleep -Milliseconds 500
}

& $mysqlExe --no-defaults --protocol=tcp --host=127.0.0.1 --port=3306 --user=root "--password=$dbPassword" -e "CREATE DATABASE IF NOT EXISTS go_test CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;"

# Always restart frontend/backend so 8080 can load the latest compiled templates.
Stop-LocalServiceProcesses

if (-not $SkipFrontendBuild) {
    $frontendBuildOut = Join-Path $tmpDir "frontend-build.out.log"
    $frontendBuildErr = Join-Path $tmpDir "frontend-build.err.log"
    $frontendBuildProc = Start-Process `
        -FilePath $npxCmd `
        -ArgumentList "-y", "node@20.20.1", "node_modules/vite/bin/vite.js", "build" `
        -WorkingDirectory $frontendDir `
        -WindowStyle Hidden `
        -RedirectStandardOutput $frontendBuildOut `
        -RedirectStandardError $frontendBuildErr `
        -PassThru `
        -Wait

    if ($frontendBuildProc.ExitCode -ne 0) {
        throw "Frontend build failed. Check $frontendBuildErr"
    }

    $distDir = Join-Path $frontendDir "dist"
    if (-not (Test-Path (Join-Path $distDir "index.html"))) {
        throw "Frontend build succeeded but dist/index.html is missing."
    }

    if (Test-Path $templatesAssetsDir) {
        Remove-Item $templatesAssetsDir -Recurse -Force
    }
    Copy-Item (Join-Path $distDir "assets") $templatesAssetsDir -Recurse -Force
    Copy-Item (Join-Path $distDir "index.html") (Join-Path $templatesDir "index.html") -Force
}

$env:QA_HTTP_ADDR = $httpBindAddr
$env:QA_TCP_ADDR = $tcpBindAddr
$env:QA_PROXY_FROM_PORT = "0"
$env:QA_PROXY_TO_PORT = "0"
$env:QA_DB_DSN = "root:$dbPassword@tcp(127.0.0.1:3306)/go_test?charset=utf8mb4&parseTime=True&loc=Local"

$buildOut = Join-Path $tmpDir "backend-build.out.log"
$buildErr = Join-Path $tmpDir "backend-build.err.log"
$goBuildProc = Start-Process `
    -FilePath $goExe `
    -ArgumentList "build", "-o", $backendExe, "." `
    -WorkingDirectory $projectRoot `
    -WindowStyle Hidden `
    -RedirectStandardOutput $buildOut `
    -RedirectStandardError $buildErr `
    -PassThru `
    -Wait

if ($goBuildProc.ExitCode -ne 0) {
    throw "Backend build failed. Check $buildErr"
}

if (Test-IsAdmin) {
    Ensure-FirewallProgramRule -RuleName "QA Local Backend 8080" -ProgramPath $backendExe -Ports "8080"
    Ensure-FirewallProgramRule -RuleName "QA Local DB 3306" -ProgramPath $mariadbdExe -Ports "3306"
    if ($Lan) {
        Ensure-FirewallPortRule -RuleName "QA LAN Frontend 5173" -Ports "5173"
        Ensure-FirewallPortRule -RuleName "QA LAN Backend 8080" -Ports "8080"
    }
} else {
    if ($Lan) {
        Write-Host "Tip: run start-lan.ps1 as Administrator once to pre-authorize firewall rules and avoid security prompts."
    } else {
        Write-Host "Tip: run start-local.ps1 as Administrator once to pre-authorize firewall rules and avoid security prompts."
    }
}

$backendOut = Join-Path $tmpDir "backend.out.log"
$backendErr = Join-Path $tmpDir "backend.err.log"
$backendProc = Start-Process -FilePath $backendExe -WorkingDirectory $projectRoot -WindowStyle Hidden -RedirectStandardOutput $backendOut -RedirectStandardError $backendErr -PassThru
Set-Content -Path $backendPidFile -Value $backendProc.Id

$frontendOut = Join-Path $tmpDir "frontend.out.log"
$frontendErr = Join-Path $tmpDir "frontend.err.log"
# Vite 2 works reliably with Node 20 in this project.
$frontendProc = Start-Process `
    -FilePath $npxCmd `
    -ArgumentList "-y", "node@20.20.1", "node_modules/vite/bin/vite.js", "--host", $viteHost, "--port", "5173", "--strictPort" `
    -WorkingDirectory $frontendDir `
    -WindowStyle Hidden `
    -RedirectStandardOutput $frontendOut `
    -RedirectStandardError $frontendErr `
    -PassThru
Set-Content -Path $frontendPidFile -Value $frontendProc.Id

if ($Lan) {
    Write-Host "LAN stack is ready."
    Write-Host "Frontend: http://$displayHost`:5173"
    Write-Host "Backend:  http://$displayHost`:8080/system/info"
    Write-Host "Local:    http://localhost:5173"
} else {
    Write-Host "Local stack is ready."
    Write-Host "Backend:  http://localhost:8080/system/info"
    Write-Host "Frontend: http://localhost:5173"
}
