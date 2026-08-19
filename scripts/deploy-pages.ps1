# Build with GitHub Pages base (/GoLearning/) and push to gh-pages branch.
# Local site/ (base /) is untouched — use build-static.ps1 for that.

$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
$web = Join-Path $root "web"
$staging = Join-Path $root "release\pages-staging"

Write-Host "==> Building for GitHub Pages (base /GoLearning/)..."
Push-Location $web
try {
  if (-not (Test-Path "node_modules")) { npm install }
  $env:GITHUB_PAGES = "true"
  npm run build
  if ($LASTEXITCODE -ne 0) { throw "astro build failed" }
} finally {
  Remove-Item Env:GITHUB_PAGES -ErrorAction SilentlyContinue
  Pop-Location
}

$dist = Join-Path $web "dist"
if (-not (Test-Path (Join-Path $dist "index.html"))) {
  throw "dist/index.html missing"
}
if (-not (Select-String -Path (Join-Path $dist "index.html") -Pattern "/GoLearning/" -Quiet)) {
  throw "Build missing /GoLearning/ base — aborting deploy"
}

Write-Host "==> Staging..."
if (Test-Path $staging) { Remove-Item $staging -Recurse -Force }
New-Item -ItemType Directory -Path $staging -Force | Out-Null
Copy-Item -Path (Join-Path $dist "*") -Destination $staging -Recurse -Force

# Empty .nojekyll so GitHub Pages won't run Jekyll (underscored folders stay public)
New-Item -ItemType File -Path (Join-Path $staging ".nojekyll") -Force | Out-Null

Write-Host "==> Pushing to gh-pages..."
Push-Location $staging
try {
  git init -b gh-pages | Out-Null
  git checkout -B gh-pages | Out-Null
  git add -A
  $msg = "Deploy site $(Get-Date -Format 'yyyy-MM-dd HH:mm')"
  # Must use a GitHub-linked author; unverified fake emails make Pages builds fail.
  git -c user.name="willlee88" -c user.email="18569620+willlee88@users.noreply.github.com" commit -m $msg | Out-Null
  git remote add origin "https://github.com/willlee88/GoLearning.git"
  git push -f origin gh-pages
  if ($LASTEXITCODE -ne 0) { throw "git push gh-pages failed" }
} finally {
  Pop-Location
}

Write-Host "==> Ensuring Pages serves gh-pages branch..."
# Switch from Actions/workflow mode to branch deploy (no workflow scope needed).
'{"build_type":"legacy","source":{"branch":"gh-pages","path":"/"}}' |
  gh api repos/willlee88/GoLearning/pages -X PUT --input - | Out-Null

# Nudge a build (useful after switching source).
gh api -X POST repos/willlee88/GoLearning/pages/builds | Out-Null

Write-Host ""
Write-Host "Done. Site: https://willlee88.github.io/GoLearning/"
Write-Host "(First enable / cache may take 1–2 minutes.)"
Write-Host ""
Write-Host "Note: web/dist now has Pages base. Re-run build-static.ps1 if you need local site/."
