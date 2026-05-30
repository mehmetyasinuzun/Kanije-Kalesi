@echo off
REM ============================================================
REM  Kanije Kalesi - Cift tikla kurulum
REM  Bu dosyaya cift tiklayin. install.ps1 kendini yonetici
REM  olarak yukseltir, gerekli her seyi yapar ve gizli baslatir.
REM ============================================================
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0install.ps1" %*
