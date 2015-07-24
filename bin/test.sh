#!/usr/bin/env bash
#

clear && printf '\e[3J' && clear

# Eventually we'll want a windows build but that may be a long way off.
# Note that one can "cross-download" windows dependencies with go get
# by setting GOOS, for example:
#
#   GOOS=windows go get github.com/shirou/w32
#
OSARCH="darwin/amd64 linux/amd64"

# I don't see any way in gox to determine if the toolchain has already
# been built. Super annoying. It will rebuild it if one passes -build-toolchain
# when it's already built, so right now this depends on the toolchains being
# built at the same time gox is downloaded.
which gox >/dev/null || {
  echo "go getting gox for cross-build..."
  go get github.com/mitchellh/gox
  echo "Building gox toolchains..."
  gox -osarch="$OSARCH" -build-toolchain
}

gox -osarch="$OSARCH" -output "${TMPDIR:-/tmp}/{{.Dir}}_{{.OS}}_{{.Arch}}" ./...
go test -v ./...
