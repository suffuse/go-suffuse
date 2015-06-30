package suffuse

import (
  "syscall"
  "bazil.org/fuse"
)

func PlatformOptions() []fuse.MountOption {
  return []fuse.MountOption {
  }
}

func SetSysAttributes(sp *syscall.Stat_t, a *fuse.Attr) {
}
