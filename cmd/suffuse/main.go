package main

import (
  . "github.com/suffuse/go-suffuse/suffuse"
  "os"
)

func main() {
  conf := CreateSfsConfig(os.Args)
  mfs, err := NewSfs(conf)
  MaybePanic(err)

  err = mfs.Serve()
  MaybeLog(err)

  <- mfs.Connection.Ready
  MaybeLog(mfs.Connection.MountError)
}
