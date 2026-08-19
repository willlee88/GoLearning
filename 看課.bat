@echo off
chcp 65001 >nul
cd /d "%~dp0site"
if not exist "index.html" (
  echo [ERROR] site\ 裡還沒有 HTML。
  echo 請先在有 Node 的電腦跑: scripts\build-static.ps1
  echo 然後 git commit / push site\ 資料夾。
  pause
  exit /b 1
)
echo.
echo  GoLearning 課程站（git 裡的 site\）
echo  瀏覽器開: http://127.0.0.1:4321
echo  關掉這個視窗就會停止。
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
echo [ERROR] 需要 Python 3 或 Node.js 才能開本機網站。
echo 裝好之後再雙擊 看課.bat。
echo 或者把 site\ 整包丟到任何靜態網站空間／NAS。
pause
