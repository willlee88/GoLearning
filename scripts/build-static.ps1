# Build the teaching site as static HTML (no Docker required).
#
# Outputs:
#   site/                         — tracked in git; other PCs: git pull + 看課.bat
#   release/GoLearning-static/    — local assemble folder (gitignored)
#   release/GoLearning-static.zip — optional copy-anywhere zip (gitignored)

$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
$web = Join-Path $root "web"
$site = Join-Path $root "site"
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

function Publish-StaticFolder([string]$dest) {
  if (Test-Path $dest) {
    Get-ChildItem -Path $dest -Force | Remove-Item -Recurse -Force
  } else {
    New-Item -ItemType Directory -Path $dest -Force | Out-Null
  }
  Copy-Item -Path (Join-Path $web "dist\*") -Destination $dest -Recurse -Force
  Copy-Item (Join-Path $scripts "static-start.bat") (Join-Path $dest "start.bat") -Force
  Copy-Item (Join-Path $scripts "static-start.ps1") (Join-Path $dest "start.ps1") -Force
  Copy-Item (Join-Path $scripts "static-README.txt") (Join-Path $dest "README.txt") -Force
}

Write-Host "==> Publishing site/ (for git pull readers)..."
Publish-StaticFolder $site

Write-Host "==> Assembling release folder + zip..."
New-Item -ItemType Directory -Path (Split-Path $out -Parent) -Force | Out-Null
Publish-StaticFolder $out

$zip = Join-Path $root "release\GoLearning-static.zip"
if (Test-Path $zip) { Remove-Item $zip -Force }
Compress-Archive -Path $out -DestinationPath $zip -Force

Write-Host ""
Write-Host "Done."
Write-Host "  Git-tracked HTML: $site"
Write-Host "  Zip (optional):   $zip"
Write-Host ""
Write-Host "Readers: git pull  then  .\看課.bat"
Write-Host "Authors: commit content/ + site/ together after this script."
