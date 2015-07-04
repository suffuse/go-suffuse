Suffuse client
==============

A raw snapshot of generic file virtualization in progress. Written in go.

Prerequisites
=============

[Go](https://golang.org/)

**OSX**
```
brew install go
```

**Linux**

Get and extract the correct version (at least version `1.4.2`) from [The Go Download page](https://golang.org/dl/)

```
# Make sure to export GOROOT to the correct path
export GOROOT="$HOME/opt/go1.4.2"
```

--------------------

[GOPATH environment variable](https://github.com/golang/go/wiki/GOPATH) needs to be correctly set.

An example setup of the paths:

```
export GOPATH="$HOME/workspaces/go"
export PATH="$GOPATH/bin":$PATH
```

Installation
============

**When you see ... it means literally three dots.** It's the go syntax for "all projects under this directory."

```
go get github.com/suffuse/go-suffuse/cmd/suffuse
```

Suffuse has been installed in `$GOPATH/bin`.

A git checkout has been created in `$GOPATH/src/github.com/suffuse/go-suffuse`.


Running
=======

The general steps for suffuse are:

1. Mount using the `suffuse` command
2. Interact with the mounted file system

_Making sure suffuse is available_
```
# Make the suffuse executable available on the path.
% export PATH=$PATH:$GOPATH/bin

# Run suffuse to get an overview of the options.
% suffuse
Usage: suffuse <options> [path path ...]

  -d=false: log at DEBUG level
  -m="": mount point
  -n="": volume name (OSX only)
  -t=false: create scratch directory as mount point
  -v=false: log at INFO level
```

_Preparing a playground_
```
# Create a directory to play in.
% mkdir -p ~/tmp/scratch

# Fill a file with digits 1 to 10.
% seq 1 10 > ~/tmp/scratch/seq.txt
```

_Mounting_
```
# Create a directory that holds the mount.
% mkdir ~/mnt

# Mount a directory through suffuse.
# The `&` runs suffuse in a separate process.
% suffuse -m ~/mnt ~/tmp/scratch &
[1] 9134

# List the contents through the suffuse mount.
% ls ~/mnt
seq.txt
```

_Start playing_
```
# It's a 10 line file, one number to a line.
% wc -l ~/mnt/seq.txt
10 /home/user/mnt/seq.txt

# Via suffuse, a derived file ending with #4,6p is a sed 
# command executed on the actual file.
% cat /mnt/seq.txt#4,6p
4
5
6

# Arbitrary sed commands, different sized files.
% ls -l /mnt/seq.txt#5,10p
-rw-r--r--  1 user  user  13 Jun 30 11:57 /mnt/seq.txt#5,10p

# These files are effectively indistinguishable from "real" files.
% ls -l /mnt/seq.txt#1,3p
-rw-r--r--  1 user  user  6 Jun 30 11:57 /mnt/seq.txt#1,3p
```

_Kill the suffuse instance_
```
% kill %1
```

Development
===========

Make sure that:
- your local checkout is located at: `$GOPATH/src/github.com/suffuse/go-suffuse`
- the `$GOROOT/bin` directory is on your `$PATH`

```
# Install dependencies
go get gopkg.in/check.v1
go get github.com/satori/go.uuid
# Continuous testing
bin/cc.sh
# Docker container works somewhat
bin/docker.sh
```
