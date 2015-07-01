Suffuse client
==============

A raw snapshot of generic file virtualization in progress. Written in go.

Prerequisites
=============

- [Go](https://golang.org/).
- The [GOPATH environment variable](https://github.com/golang/go/wiki/GOPATH) correctly set.

Installation
============

**When you see ... it means literally three dots.** It's the go syntax for "all projects under this directory."

The easy 'go' way
```
go get github.com/suffuse/go-suffuse/cmd/suffuse
```

The harder 'go' way
```
go get github.com/suffuse/go-suffuse/...
cd "$(go list -f '{{.Dir}}' github.com/suffuse/go-suffuse)"
go install ./...
```

The not-so-'go' way
```
git clone https://github.com/suffuse/go-suffuse
cd go-suffuse
go install ./...
```

Running
=======

After install, suffuse has been installed in $GOROOT/bin.
```
% suffuse
Usage: suffuse <options> [path path ...]

  -d=false: log at DEBUG level
  -m="": mount point
  -n="": volume name (OSX only)
  -t=false: create scratch directory as mount point
  -v=false: log at INFO level

% suffuse -m /mnt $HOME &
% ls -l /mnt
%
```

Development
===========

```
# Continuous testing
bin/cc.sh
# Docker container works somewhat
bin/docker.sh
```

Taster
======

```
% seq 1 10 > /scratch/seq.txt
% suffuse -m /mnt /scratch &
# It's a 10 line file, one number to a line.
% wc -l /mnt/seq.txt
      10 /mnt/seq.txt
# Via suffuse, a derived file: 4,6p is a sed command.
% cat /mnt/seq.txt#4,6p
4
5
6
# Arbitrary sed commands, different sized files.
% ls -l /mnt/seq.txt#5,10p
-rw-r--r--  1 paulp  wheel  13 Jun 30 11:57 /mnt/seq.txt#5,10p
# These files are effectively indistinguishable from "real" files.
% ls -l /mnt/seq.txt#1,3p
-rw-r--r--  1 paulp  wheel  6 Jun 30 11:57 /mnt/seq.txt#1,3p
```
