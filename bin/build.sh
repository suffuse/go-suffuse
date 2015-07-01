#!/usr/bin/env bash
#

clear && printf '\e[3J' && clear

echo "[Building...]"
go build -v github.com/suffuse/go-suffuse/...
echo "[Done]"
