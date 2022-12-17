#################################
# Set up build
#################################
$WASMpath = ".\wasm\main.wasm"
$WASMname = "main.wasm"
$PREWASMname = "."
$DESTpath = ".\frontend\src\wasm\main.wasm"

#################################
# Build the program
#################################
cd .\wasm\
# Set GOOS and GOARCH to build the WASM file.
$Env:GOOS = "js"; $Env:GOARCH = "wasm"
go build -tags=wails -o $WASMname $PREWASMname
cd ..
# Copy file to destination path, then remove the WASM file.
Copy-Item $WASMpath $DESTpath -Force
Remove-Item $WASMpath
# Set GOOS and GOARCH back to normal, and build the program
$env:GOOS = "windows"; $env:GOARCH = "amd64";
