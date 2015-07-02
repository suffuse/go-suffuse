package suffuse

/** Packages confined to this file: syscall, os/signal
 */

import (
  sys "syscall"
  "os"
  "os/signal"
  "time"
  "bazil.org/fuse"
)

const (
  ENOTDIR         = fuse.Errno(sys.ENOTDIR) // not a directory
  EINVAL          = fuse.Errno(sys.EINVAL)  // invalid argument
  MAXNAMELEN      = 255
  MNT_FORCE       = 0x1
  OsModeAnyDevice = os.ModeDevice | os.ModeCharDevice
)

/** Install a signal handler for INT/TERM.
 */
func TrapExit (handler func(os.Signal)) {
  sigs := make(chan os.Signal, 1)
  signal.Notify(sigs, sys.SIGINT, sys.SIGTERM)

  go func() {
    sig := <-sigs
    handler(sig)
  }()
}

func (x Path) SysStatFile() (sys.Stat_t, error) {
  s := sys.Stat_t { }
  err := sys.Stat(x.Path, &s)
  return s, err
}
func (x Path) SysStatLink() (sys.Stat_t, error) {
  s := sys.Stat_t { }
  err := sys.Lstat(x.Path, &s)
  return s, err
}
func (x Path) SysStatfs() (sys.Statfs_t, error) {
  s := sys.Statfs_t { }
  err := sys.Statfs(x.Path, &s)
  return s, err
}
func (x Path) SysUnmount() error {
  return sys.Unmount(x.Path, MNT_FORCE)
}

/** Conversions amongst unix, fuse, and go data structures.
 */

func GoFileInfoToFuseAttr(fi os.FileInfo) fuse.Attr {
  a := fuse.Attr {
     // Valid: cachable,
      Size: uint64(fi.Size()),
      Mode: fi.Mode(),
     Mtime: fi.ModTime(),
  }

  switch sp := fi.Sys().(type) {
    case *sys.Stat_t :
      AssertEq(StatModeToGoMode(uint64(sp.Mode)), a.Mode)
      // System specific attributes
      SetSysAttributes(sp, &a)
      a.Blocks = uint64(sp.Blocks)
      a.Gid    = sp.Gid
      a.Inode  = uint64(sp.Ino)
      a.Nlink  = uint32(sp.Nlink)
      a.Rdev   = uint32(sp.Rdev)
      a.Uid    = sp.Uid
  }

  return a
}

func SysStatfsToFuseStatfs(s sys.Statfs_t) fuse.StatfsResponse {
  return fuse.StatfsResponse {
    Blocks  : s.Blocks,
    Bfree   : s.Bfree,
    Bavail  : s.Bavail,
    Files   : s.Files,
    Ffree   : s.Ffree,
    // Bsize   : s.Bsize, // Ignored by osxfuse, uses mount option iosize
    // Frsize  : s.Bsize, // Ignored, but see http://fuse.996288.n3.nabble.com/statvfs-vs-statfs-td8636.html
    Namelen : MAXNAMELEN,
  }
}

// struct timespec to Go's time.Time.
func TimespecToGoTime(ts sys.Timespec) time.Time {
  return time.Unix(ts.Sec, ts.Nsec)
}

// sys.go:68: cannot use sp.Mode (type uint32) as type uint16 in argument to StatModeToGoMode
// sys.go:70: sp.Atimespec undefined (type *syscall.Stat_t has no field or method Atimespec)
// sys.go:71: sp.Birthtimespec undefined (type *syscall.Stat_t has no field or method Birthtimespec)
// sys.go:72: sp.Ctimespec undefined (type *syscall.Stat_t has no field or method Ctimespec)
// sys.go:74: sp.Flags undefined (type *syscall.Stat_t has no field or method Flags)
// sys.go:92: cannot use s.Bsize (type int64) as type uint32 in field value
// sys.go:93: cannot use s.Bsize (type int64) as type uint32 in field value

// Unix stat/lstat mode bits to Go os.FileMode.
func StatModeToGoMode(bits uint64) os.FileMode {
  // Permission bits passed through.
  mode := os.FileMode(bits & 0777)

  switch bits & sys.S_IFMT {
    case sys.S_IFREG : // nothing
    case sys.S_IFDIR : mode |= os.ModeDir
    case sys.S_IFCHR : mode |= os.ModeCharDevice | os.ModeDevice
    case sys.S_IFBLK : mode |= os.ModeDevice
    case sys.S_IFIFO : mode |= os.ModeNamedPipe
    case sys.S_IFLNK : mode |= os.ModeSymlink
    case sys.S_IFSOCK: mode |= os.ModeSocket
    default          : mode |= os.ModeDevice
  }

  if bits & sys.S_ISUID != 0 { mode |= os.ModeSetuid }
  if bits & sys.S_ISGID != 0 { mode |= os.ModeSetgid }

  return mode
}

// These don't quite match os.FileMode; especially there's an
// explicit unknown, instead of zero value meaning file. They
// are also not quite syscall.DT_*; nothing says the FUSE
// protocol follows those, and even if they were, we don't
// want each fs to fiddle with syscall.
func GoModeToDirentType(mode os.FileMode) fuse.DirentType {
  bits := mode & os.ModeType

  switch {
    case bits                    == 0                : return fuse.DT_File
    case bits & os.ModeDir       == os.ModeDir       : return fuse.DT_Dir
    case bits & os.ModeSymlink   == os.ModeSymlink   : return fuse.DT_Link
    case bits & os.ModeSocket    == os.ModeSocket    : return fuse.DT_Socket
    case bits & os.ModeNamedPipe == os.ModeNamedPipe : return fuse.DT_FIFO
    case bits & OsModeAnyDevice  == OsModeAnyDevice  : return fuse.DT_Char
    case bits & os.ModeDevice    == os.ModeDevice    : return fuse.DT_Block

    default:  return fuse.DT_Unknown
  }
}
