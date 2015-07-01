#!/usr/bin/env bash
#

bindir="$(cd "$(dirname "$0")" && pwd -P)"
goroot="$(dirname "$bindir")"

what=${1:-test}
shift

which spy || go get github.com/jpillora/spy

spy --dir "$goroot" "$bindir/$what.sh" "$@"
