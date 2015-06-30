#!/usr/bin/env bash
#

set -x
set -e

# Xs are necessary for linux's mktemp, ignored by OSX's. Madness.
workdir="$(mktemp -d -t suffuseXXXXXX)"
package="github.com/paulp/suffuse"
target="$workdir/go/src/$package"
docker_image="suffuse"
docker_run_opts="--rm --cap-add SYS_ADMIN --device /dev/fuse -ti $docker_image"

export GOPATH="$workdir"
mkdir -p "$target"
rsync -av --relative cmd suffuse "$target"
cp Dockerfile "$workdir"

cd "$workdir"
# We can go-get some dependencies from outside the container, but not
# all of them, because it won't get the linux-specific parts on OSX.
set +e && ( cd "$target" && go get -t -d -v ) && set -e
docker build -t $docker_image .
docker run $docker_run_opts
