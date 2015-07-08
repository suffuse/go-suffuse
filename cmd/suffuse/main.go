package main

import (
  . "github.com/suffuse/go-suffuse/suffuse"
  "os"
)

func main() {
  conf, conf_err := CreateSfsConfig(os.Args)
  if conf_err != nil {
    conf_err.PrintUsage()
    os.Exit(2)
  }

  mfs, err := NewSfs(conf)
  MaybePanic(err)

  err = mfs.Serve()
  MaybeLog(err)

  <- mfs.Connection.Ready
  MaybeLog(mfs.Connection.MountError)
}
