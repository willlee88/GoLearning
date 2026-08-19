@echo off
chcp 65001 >nul
cd /d "%~dp0"
echo.
echo  GoLearning static site
echo  Open browser: http://127.0.0.1:4321
echo  Close this window to stop.
echo.
where py >nul 2>&1
if %ERRORLEVEL%==0 (
  py -3 -m http.server 4321
  goto :eof
)
where python >nul 2>&1
if %ERRORLEVEL%==0 (
  python -m http.server 4321
  goto :eof
)
where node >nul 2>&1
if %ERRORLEVEL%==0 (
  npx --yes serve -l 4321 .
  goto :eof
)
echo [ERROR] Need Python 3 or Node.js to serve this folder.
echo Install Python (py launcher) or Node, then run start.bat again.
echo Or upload this folder to any static web host / NAS.
pause
