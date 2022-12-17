$autoFile = "auto.ps1"
& "$PSScriptRoot\$autoFile"
$BUILDname = ".\main.exe"
wails build -o $BUILDname .