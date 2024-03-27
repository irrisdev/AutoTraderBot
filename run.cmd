@echo off

REM Set CGO_ENABLED to 1
set CGO_ENABLED=1

REM Run your Go program with CGO enabled
go run cmd/main.go
