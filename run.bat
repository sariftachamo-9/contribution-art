@echo off
set "GO_BIN=C:\Users\97798\go-sdk\go\bin\go.exe"
if exist "%GO_BIN%" (
    "%GO_BIN%" run main.go %*
) else (
    go run main.go %*
)
