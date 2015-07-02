package suffuse

import (
  sys "syscall"
  "bazil.org/fuse"
)

func PlatformOptions() []fuse.MountOption {
  return []fuse.MountOption {
  }
}

func SetSysAttributes(sp *sys.Stat_t, a *fuse.Attr) {
}

func (x Path) SysUnmount() error {
  return fuse.Unmount(x.Path)
}
