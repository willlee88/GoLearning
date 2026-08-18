# Build a copy-anywhere static teaching site (no Docker required).
# Output: release/GoLearning-static/  and  release/GoLearning-static.zip

$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
$web = Join-Path $root "web"
$out = Join-Path $root "release\GoLearning-static"
$scripts = $PSScriptRoot

Write-Host "==> Building Astro site..."
Push-Location $web
try {
  if (-not (Test-Path "node_modules")) {
    npm install
  }
  Remove-Item Env:GITHUB_PAGES -ErrorAction SilentlyContinue
  npm run build
  if ($LASTEXITCODE -ne 0) { throw "astro build failed" }
} finally {
  Pop-Location
}

Write-Host "==> Assembling release folder..."
if (Test-Path $out) { Remove-Item $out -Recurse -Force }
New-Item -ItemType Directory -Path $out -Force | Out-Null
Copy-Item -Path (Join-Path $web "dist\*") -Destination $out -Recurse -Force
Copy-Item (Join-Path $scripts "static-start.bat") (Join-Path $out "start.bat") -Force
Copy-Item (Join-Path $scripts "static-start.ps1") (Join-Path $out "start.ps1") -Force
Copy-Item (Join-Path $scripts "static-README.txt") (Join-Path $out "README.txt") -Force

$zip = Join-Path $root "release\GoLearning-static.zip"
if (Test-Path $zip) { Remove-Item $zip -Force }
Write-Host "==> Zipping..."
Compress-Archive -Path $out -DestinationPath $zip -Force

Write-Host ""
Write-Host "Done."
Write-Host "  Folder: $out"
Write-Host "  Zip:    $zip"
Write-Host "Copy the zip to the other PC, extract, run start.bat"
