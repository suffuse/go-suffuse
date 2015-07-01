package main

import . "github.com/suffuse/go-suffuse/suffuse"

var mfs Sfs
var err error

func main() {
  opts := ParseSfsOpts()

  mfs, err = opts.Create()
  MaybePanic(err)

  err = mfs.Serve()
  MaybeLog(err)

  <- mfs.Connection.Ready
  MaybeLog(mfs.Connection.MountError)
}
