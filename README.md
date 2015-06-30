Suffuse client
==============

A raw snapshot of generic file virtualization in progress. Written in go.

Installation
============

`go install github.com/paulp/suffuse/...`

Development
===========

Continuous testing.
```
go get github.com/paulp/suffuse
cd "$(go list -f '{{.Dir}}' github.com/paulp/suffuse)"
bin/cc.sh
```

Docker container works somewhat.
```
bin/docker.sh
```
