Suffuse client
==============

A raw snapshot of generic file virtualization in progress. Written in go.

Prerequisites
=============

[Go](https://golang.org/).

Installation
============

The easy 'go' way
```
go install github.com/paulp/suffuse/...
```

The harder 'go' way
```
go get github.com/paulp/suffuse
cd "$(go list -f '{{.Dir}}' github.com/paulp/suffuse)"
go install ./...
```

The not-so-'go' way
```
git clone https://github.com/paulp/suffuse
cd suffuse
go install ./...
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
% wc -l /mnt/seq.txt
      10 /mnt/seq.txt
% cat /mnt/seq.txt#4,6p
4
5
6
% ls -l /mnt/seq.txt#5,10p
-rw-r--r--  1 paulp  wheel  13 Jun 30 11:57 /mnt/seq.txt#5,10p
% ls -l /mnt/seq.txt#1,3p
-rw-r--r--  1 paulp  wheel  6 Jun 30 11:57 /mnt/seq.txt#1,3p
```
