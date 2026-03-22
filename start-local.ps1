param(
    [switch]$Lan,
    [string]$LanIP = ""
)

$ErrorActionPreference = "Stop"

$projectRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$workspaceRoot = Split-Path -Parent $projectRoot
$tmpDir = Join-Path $projectRoot "tmp"
$dataDir = Join-Path $workspaceRoot "mariadb-data"
$myIni = Join-Path $dataDir "my.ini"
$backendExe = Join-Path $tmpDir "qa_test_server-local.exe"

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

$env:QA_HTTP_ADDR = $httpBindAddr
$env:QA_TCP_ADDR = $tcpBindAddr
$env:QA_PROXY_FROM_PORT = "0"
$env:QA_PROXY_TO_PORT = "0"
$env:QA_DB_DSN = "root:$dbPassword@tcp(127.0.0.1:3306)/go_test?charset=utf8mb4&parseTime=True&loc=Local"

$buildOut = Join-Path $tmpDir "backend-build.out.log"
$buildErr = Join-Path $tmpDir "backend-build.err.log"
& $goExe build -o $backendExe . 1>$buildOut 2>$buildErr

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

$backendListening = Get-NetTCPConnection -LocalPort 8080 -State Listen -ErrorAction SilentlyContinue
if (-not $backendListening) {
    $backendOut = Join-Path $tmpDir "backend.out.log"
    $backendErr = Join-Path $tmpDir "backend.err.log"
    $backendProc = Start-Process -FilePath $backendExe -WorkingDirectory $projectRoot -WindowStyle Hidden -RedirectStandardOutput $backendOut -RedirectStandardError $backendErr -PassThru
    Set-Content -Path (Join-Path $tmpDir "backend.pid") -Value $backendProc.Id
}

$frontendListening = Get-NetTCPConnection -LocalPort 5173 -State Listen -ErrorAction SilentlyContinue
if (-not $frontendListening) {
    $frontendOut = Join-Path $tmpDir "frontend.out.log"
    $frontendErr = Join-Path $tmpDir "frontend.err.log"
    # Vite 2 works reliably with Node 20 in this project.
    $frontendProc = Start-Process `
        -FilePath $npxCmd `
        -ArgumentList "-y", "node@20.20.1", "node_modules/vite/bin/vite.js", "--host", $viteHost, "--port", "5173", "--strictPort" `
        -WorkingDirectory (Join-Path $projectRoot "qa-test-web") `
        -WindowStyle Hidden `
        -RedirectStandardOutput $frontendOut `
        -RedirectStandardError $frontendErr `
        -PassThru
    Set-Content -Path (Join-Path $tmpDir "frontend.pid") -Value $frontendProc.Id
}

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
