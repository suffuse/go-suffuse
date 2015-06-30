#!/usr/bin/env bash
#

clear && printf '\e[3J' && clear

echo "[Building...]"
go build -v github.com/paulp/suffuse
echo "[Done]"
