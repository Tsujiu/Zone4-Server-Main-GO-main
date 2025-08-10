@echo off
chcp 65001 >nul
setlocal
cd /d "%~dp0"

title 🌐 Zone4 - Start Manager + All Channels
color 0A

echo [INFO] Iniciando Channel Manager...
start "Zone4-Manager" cmd /c "go run ./cmd/manager"

REM ==== espera o manager responder ====
set /a retries=40
:WAIT_LOOP
  >nul 2>&1 curl -fsS http://127.0.0.1:8080/status
  if %errorlevel%==0 goto :UP
  set /a retries-=1
  if %retries% LEQ 0 (
    echo [ERRO] Manager nao respondeu em tempo habil. Abortei.
    goto :END
  )
  timeout /t 1 >nul
goto :WAIT_LOOP

:UP
echo [INFO] Manager online. Disparando START ALL...
>nul 2>&1 curl -fsS -X POST http://127.0.0.1:8080/start-all
if not %errorlevel%==0 (
  echo [AVISO] Falha ao iniciar via POST. Tentando GET...
  >nul 2>&1 curl -fsS http://127.0.0.1:8080/start-all
)

echo [INFO] Abrindo painel...
start "" http://127.0.0.1:8080/

:END
echo.
pause
endlocal
