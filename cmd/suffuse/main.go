package main

import (
  . "github.com/suffuse/go-suffuse/suffuse"
  "os"
)

func main() {
  opts, opts_err := ParseSfsOpts(os.Args)
  if opts_err != nil {
    opts_err.PrintUsage()
    os.Exit(2)
  }

  mfs, err := NewSfs(opts)
  MaybePanic(err)

  err = mfs.Serve()
  MaybeLog(err)

  <- mfs.Connection.Ready
  MaybeLog(mfs.Connection.MountError)
}
