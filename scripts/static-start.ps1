Set-Location $PSScriptRoot
Write-Host "GoLearning -> http://127.0.0.1:4321  (Ctrl+C to stop)"
if (Get-Command py -ErrorAction SilentlyContinue) {
  py -3 -m http.server 4321
} elseif (Get-Command python -ErrorAction SilentlyContinue) {
  python -m http.server 4321
} elseif (Get-Command node -ErrorAction SilentlyContinue) {
  npx --yes serve -l 4321 .
} else {
  Write-Error "Need Python 3 or Node.js to serve this folder."
}
