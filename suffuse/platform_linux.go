package suffuse

import (
  "time"
  "syscall"
  "bazil.org/fuse"
)

func PlatformOptions() []fuse.MountOption {
  return []fuse.MountOption {
  }
}

func SetSysAttributes(sp *syscall.Stat_t, a *fuse.Attr) {
}

func SysAtimeMtime(sp *syscall.Stat_t) (atime time.Time, mtime time.Time) {
  return TimespecToGoTime(sp.Atim), TimespecToGoTime(sp.Mtim)
}
