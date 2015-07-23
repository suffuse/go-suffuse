#!/usr/bin/env bash
#

clear && printf '\e[3J' && clear

function explain {
  echo "                  ^"
  echo "RUN THIS TO FIX: cd \$(go env GOROOT)/src ; GOOS=$1 GOARCH=amd64 ./make.bash --no-clean 2>&1"
  echo ""
}

GOOS=darwin GOARCH=amd64 go build ./...; if [ $? -eq 1 ]; then explain 'darwin'; fi
GOOS=linux  GOARCH=amd64 go build ./...; if [ $? -eq 1 ]; then explain 'linux'; fi

go test -v ./...
