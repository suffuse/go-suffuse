package suffuse

import (
  "time"
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

func SysAtimeMtime(sp *sys.Stat_t) (atime time.Time, mtime time.Time) {
  return TimespecToGoTime(sp.Atim), TimespecToGoTime(sp.Mtim)
}
