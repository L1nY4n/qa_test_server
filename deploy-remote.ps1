[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$ServerIp,

    [string]$User = "root",

    [Parameter(Mandatory = $true)]
    [string]$Password,

    [string]$HostKey = "",

    [string]$RemoteAppDir = "/root/qa_test_server",
    [string]$RemoteArchive = "/root/qa_test_server_deploy.tar.gz",
    [string]$RemoteCoreScript = "/root/remote_deploy_and_heal.sh",
    [string]$ServiceName = "qa_test_server",

    [int]$HttpPort = 8080,
    [int]$TcpPort = 4001,
    [int]$ProxyFromPort = 4003,
    [int]$ProxyToPort = 7777,

    [string]$DbName = "go_test",
    [string]$DbUser = "qa_user",
    [string]$DbPassword = "L1nFen9.com",

    [string]$AuthSecret = "qa_test_server_prod_secret_change_me",
    [string]$AdminUser = "admin",
    [string]$AdminPassword = "Admin@123456",

    [int]$MaxHealRounds = 3,
    [switch]$EnableNginxProxy,
    [switch]$SkipExternalCheck
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
if (Get-Variable PSNativeCommandUseErrorActionPreference -ErrorAction SilentlyContinue) {
    $PSNativeCommandUseErrorActionPreference = $false
}

function Write-Stage {
    param([string]$Message)
    $now = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    Write-Host "[$now] $Message"
}

function Escape-BashSingleQuoted {
    param([string]$Value)
    $replacement = "'" + '"' + "'" + '"' + "'"
    return $Value.Replace("'", $replacement)
}

function Ensure-FileDownloaded {
    param(
        [string]$Path,
        [string]$Url
    )
    if (Test-Path $Path) {
        return
    }
    Write-Stage "downloading $(Split-Path $Path -Leaf)"
    Invoke-WebRequest -Uri $Url -OutFile $Path
}

function Invoke-Tar {
    param(
        [string]$ArchivePath,
        [string]$WorkspaceRoot,
        [string]$ProjectName
    )

    $tarCmd = Get-Command tar -ErrorAction Stop
    $args = @(
        "-czf", $ArchivePath,
        "--exclude=$ProjectName/.git",
        "--exclude=$ProjectName/qa-test-web/node_modules",
        "--exclude=$ProjectName/tmp",
        "--exclude=$ProjectName/device/__debug_bin",
        "-C", $WorkspaceRoot,
        $ProjectName
    )
    & $tarCmd.Source @args
    if ($LASTEXITCODE -ne 0) {
        throw "tar create archive failed"
    }
}

function Resolve-HostKeyFromOutput {
    param([string]$Text)
    $m = [regex]::Match($Text, "SHA256:[A-Za-z0-9+/=]+")
    if ($m.Success) {
        return $m.Value
    }
    return ""
}

$script:ResolvedHostKey = $HostKey

function Invoke-NativeCapture {
    param(
        [string]$FilePath,
        [string[]]$Arguments
    )

    $stdoutPath = [System.IO.Path]::GetTempFileName()
    $stderrPath = [System.IO.Path]::GetTempFileName()
    try {
        $proc = Start-Process -FilePath $FilePath -ArgumentList $Arguments -NoNewWindow -Wait -PassThru -RedirectStandardOutput $stdoutPath -RedirectStandardError $stderrPath
        $stdoutText = ""
        $stderrText = ""

        if (Test-Path $stdoutPath) {
            $stdoutText = Get-Content -Path $stdoutPath -Raw -ErrorAction SilentlyContinue
        }
        if (Test-Path $stderrPath) {
            $stderrText = Get-Content -Path $stderrPath -Raw -ErrorAction SilentlyContinue
        }

        $parts = @()
        if ($stdoutText) {
            $parts += $stdoutText.TrimEnd("`r", "`n")
        }
        if ($stderrText) {
            $parts += $stderrText.TrimEnd("`r", "`n")
        }

        return [pscustomobject]@{
            ExitCode = $proc.ExitCode
            Text     = ($parts -join "`n")
        }
    }
    finally {
        Remove-Item -Force -ErrorAction SilentlyContinue $stdoutPath, $stderrPath
    }
}

function Invoke-Plink {
    param(
        [string]$PlinkPath,
        [string]$Target,
        [string]$PasswordText,
        [string]$InlineCommand = "",
        [string]$CommandFile = "",
        [switch]$NoHostKeyRetry
    )

    $args = @("-batch", "-ssh", $Target, "-pw", $PasswordText)
    if ($script:ResolvedHostKey) {
        $args += @("-hostkey", $script:ResolvedHostKey)
    }

    if ($CommandFile) {
        $args += @("-m", $CommandFile)
    } elseif ($InlineCommand) {
        $args += @($InlineCommand)
    } else {
        throw "Invoke-Plink requires InlineCommand or CommandFile"
    }

    $result = Invoke-NativeCapture -FilePath $PlinkPath -Arguments $args
    $exitCode = $result.ExitCode
    $text = $result.Text
    $output = @()
    if ($text) {
        $output = $text -split "\r?\n"
    }

    if ($exitCode -eq 0) {
        return ,$output
    }

    if (-not $script:ResolvedHostKey -and -not $NoHostKeyRetry) {
        $discovered = Resolve-HostKeyFromOutput -Text $text
        if ($discovered) {
            $script:ResolvedHostKey = $discovered
            Write-Stage "discovered host key: $discovered"
            return Invoke-Plink -PlinkPath $PlinkPath -Target $Target -PasswordText $PasswordText -InlineCommand $InlineCommand -CommandFile $CommandFile -NoHostKeyRetry
        }
    }

    throw "plink failed (exit $exitCode):`n$text"
}

function Invoke-PscpUpload {
    param(
        [string]$PscpPath,
        [string]$Target,
        [string]$PasswordText,
        [string]$LocalPath,
        [string]$RemotePath
    )

    $args = @("-batch", "-pw", $PasswordText)
    if ($script:ResolvedHostKey) {
        $args += @("-hostkey", $script:ResolvedHostKey)
    }
    $args += @($LocalPath, "$Target`:$RemotePath")

    $result = Invoke-NativeCapture -FilePath $PscpPath -Arguments $args
    if ($result.ExitCode -ne 0) {
        $text = $result.Text
        throw "pscp upload failed:`n$text"
    }
}

function New-RemoteRunnerFile {
    param(
        [string]$Path,
        [string]$Mode,
        [int]$EnableProxyValue
    )

    $kv = @{}
    $kv["APP_DIR"] = $RemoteAppDir
    $kv["ARCHIVE_PATH"] = $RemoteArchive
    $kv["SERVICE_NAME"] = $ServiceName
    $kv["HTTP_PORT"] = [string]$HttpPort
    $kv["TCP_PORT"] = [string]$TcpPort
    $kv["PROXY_FROM_PORT"] = [string]$ProxyFromPort
    $kv["PROXY_TO_PORT"] = [string]$ProxyToPort
    $kv["DB_NAME"] = $DbName
    $kv["DB_USER"] = $DbUser
    $kv["DB_PASS"] = $DbPassword
    $kv["AUTH_SECRET"] = $AuthSecret
    $kv["ADMIN_USER"] = $AdminUser
    $kv["ADMIN_PASS"] = $AdminPassword
    $kv["ENABLE_NGINX_PROXY"] = [string]$EnableProxyValue
    $kv["HEAL_MODE"] = $Mode

    $lines = @("#!/usr/bin/env bash", "set -euo pipefail")
    foreach ($key in $kv.Keys) {
        $escaped = Escape-BashSingleQuoted -Value $kv[$key]
        $lines += "export $key='$escaped'"
    }
    $scriptEscaped = Escape-BashSingleQuoted -Value $RemoteCoreScript
    $lines += "bash '$scriptEscaped'"
    Set-Content -Path $Path -Value $lines -Encoding ASCII
}

function Get-HttpCode {
    param([string]$Url)
    $out = & curl.exe -s -o NUL -w "%{http_code}" --max-time 8 $Url 2>$null
    if ($LASTEXITCODE -ne 0) {
        return 0
    }
    $code = 0
    if ([int]::TryParse($out, [ref]$code)) {
        return $code
    }
    return 0
}

if ($MaxHealRounds -lt 1) {
    throw "MaxHealRounds must be >= 1"
}

$projectRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$workspaceRoot = Split-Path -Parent $projectRoot
$projectName = Split-Path $projectRoot -Leaf
$remoteScriptLocal = Join-Path $projectRoot "scripts\remote_deploy_and_heal.sh"

if (-not (Test-Path $remoteScriptLocal)) {
    throw "missing script: $remoteScriptLocal"
}

$plinkPath = Join-Path $workspaceRoot "plink.exe"
$pscpPath = Join-Path $workspaceRoot "pscp.exe"
Ensure-FileDownloaded -Path $plinkPath -Url "https://the.earth.li/~sgtatham/putty/latest/w64/plink.exe"
Ensure-FileDownloaded -Path $pscpPath -Url "https://the.earth.li/~sgtatham/putty/latest/w64/pscp.exe"

$target = "$User@$ServerIp"
$tempDir = Join-Path $env:TEMP ("qa_deploy_" + [guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $tempDir -Force | Out-Null

try {
    Write-Stage "checking SSH connectivity"
    Invoke-Plink -PlinkPath $plinkPath -Target $target -PasswordText $Password -InlineCommand "echo connected" | Out-Null

    $archivePath = Join-Path $tempDir ("${projectName}_deploy.tar.gz")
    Write-Stage "packing project archive"
    Invoke-Tar -ArchivePath $archivePath -WorkspaceRoot $workspaceRoot -ProjectName $projectName

    Write-Stage "uploading archive"
    Invoke-PscpUpload -PscpPath $pscpPath -Target $target -PasswordText $Password -LocalPath $archivePath -RemotePath $RemoteArchive

    Write-Stage "uploading remote deploy core script"
    Invoke-PscpUpload -PscpPath $pscpPath -Target $target -PasswordText $Password -LocalPath $remoteScriptLocal -RemotePath $RemoteCoreScript

    $deploySucceeded = $false
    $useProxy = if ($EnableNginxProxy.IsPresent) { 1 } else { 0 }

    for ($round = 1; $round -le $MaxHealRounds; $round++) {
        $mode = if ($round -eq 1) { "deploy" } else { "heal" }
        Write-Stage "round $round/${MaxHealRounds}: running remote mode=$mode (proxy=$useProxy)"

        $runnerLocal = Join-Path $tempDir ("runner_${mode}_${round}.sh")
        New-RemoteRunnerFile -Path $runnerLocal -Mode $mode -EnableProxyValue $useProxy
        $runnerRemote = "/root/runner_${mode}_${round}.sh"
        Invoke-PscpUpload -PscpPath $pscpPath -Target $target -PasswordText $Password -LocalPath $runnerLocal -RemotePath $runnerRemote

        $prepCmd = "bash -lc 'chmod +x $RemoteCoreScript $runnerRemote'"
        Invoke-Plink -PlinkPath $plinkPath -Target $target -PasswordText $Password -InlineCommand $prepCmd | Out-Null

        try {
            $out = Invoke-Plink -PlinkPath $plinkPath -Target $target -PasswordText $Password -CommandFile $runnerLocal
            if ($out) {
                Write-Host ($out -join "`n")
            }
        } catch {
            Write-Stage "round $round failed: $($_.Exception.Message)"
            if ($round -ge $MaxHealRounds) {
                throw
            }
            continue
        }

        Write-Stage "verifying remote health"
        $verifyLocal = Join-Path $tempDir ("runner_verify_${round}.sh")
        New-RemoteRunnerFile -Path $verifyLocal -Mode "verify" -EnableProxyValue $useProxy
        $verifyOut = Invoke-Plink -PlinkPath $plinkPath -Target $target -PasswordText $Password -CommandFile $verifyLocal
        if ($verifyOut) {
            Write-Host ($verifyOut -join "`n")
        }

        $deploySucceeded = $true
        break
    }

    if (-not $deploySucceeded) {
        throw "deployment did not pass remote verification"
    }

    if (-not $SkipExternalCheck.IsPresent) {
        $directUrl = "http://${ServerIp}:${HttpPort}/"
        Write-Stage "checking external endpoint: $directUrl"
        $directCode = Get-HttpCode -Url $directUrl

        if ($directCode -eq 200) {
            Write-Stage "external direct check passed (HTTP 200)"
        } else {
            Write-Stage "external direct check failed (HTTP $directCode)"
            if (-not $EnableNginxProxy.IsPresent) {
                Write-Stage "trying one more recursive heal with nginx proxy on port 80"
                $runnerProxy = Join-Path $tempDir "runner_heal_proxy.sh"
                New-RemoteRunnerFile -Path $runnerProxy -Mode "heal" -EnableProxyValue 1
                Invoke-Plink -PlinkPath $plinkPath -Target $target -PasswordText $Password -CommandFile $runnerProxy | Out-Null

                $proxyUrl = "http://${ServerIp}/"
                $proxyCode = Get-HttpCode -Url $proxyUrl
                if ($proxyCode -eq 200) {
                    Write-Stage "external proxy check passed via $proxyUrl"
                } else {
                    Write-Warning "service is healthy on server, but external access still failed (direct:$directCode, proxy:$proxyCode). Check cloud security-group inbound rules for TCP 80/8080."
                }
            } else {
                Write-Warning "service is healthy on server, but external direct access failed (HTTP $directCode). Check cloud security-group inbound rules."
            }
        }
    }

    Write-Stage "deployment completed"
}
finally {
    if (Test-Path $tempDir) {
        Remove-Item -Recurse -Force $tempDir
    }
}
