package suffuse

import (
  "syscall"
  "time"
  "bazil.org/fuse"
)

func PlatformOptions() []fuse.MountOption {
  return []fuse.MountOption {
    // fuse.LocalVolume(),
    // fuse.VolumeName("Suffuse"),
    // fuse.ArbitraryMountOption("noappledouble", ""),
    // fuse.FuseMountFlag("noapplexattr"),
    // fuse.FuseMountFlag("negative_vncache"),
    // fuse.FuseMountFlag("kernel_cache"),
    // fuse.FuseMountFlag("debug"),
  }
}

func SetSysAttributes(sp *syscall.Stat_t, a *fuse.Attr) {
  a.Atime  = TimespecToGoTime(sp.Atimespec)
  a.Crtime = TimespecToGoTime(sp.Birthtimespec) // time of creation (OS X only)
  a.Ctime  = TimespecToGoTime(sp.Ctimespec)     // time of last inode change
  a.Flags  = sp.Flags
}

func SysAtimeMtime(sp *syscall.Stat_t) (atime time.Time, mtime time.Time) {
  return TimespecToGoTime(sp.Atimespec), TimespecToGoTime(sp.Mtimespec)
}
