# [bazil.org/fuse](https://github.com/bazil/fuse) data types

```go
// A Conn represents a connection to a mounted FUSE file system.
type Conn struct {
  // Ready is closed when the mount is complete or has failed.
  Ready <-chan struct{}
  // MountError stores any error from the mount process. Only valid
  // after Ready is closed.
  MountError error
  // File handle for kernel communication. Only safe to access if
  // rio or wio is held.
  dev *os.File
  buf []byte
  wio sync.Mutex
  rio sync.RWMutex
}

// An Attr is the metadata for a single file or directory.
type Attr struct {
  Valid time.Duration // how long Attr can be cached
  Inode  uint64      // inode number
  Size   uint64      // size in bytes
  Blocks uint64      // size in blocks
  Atime  time.Time   // time of last access
  Mtime  time.Time   // time of last modification
  Ctime  time.Time   // time of last inode change
  Crtime time.Time   // time of creation (OS X only)
  Mode   os.FileMode // file mode
  Nlink  uint32      // number of links
  Uid    uint32      // owner uid
  Gid    uint32      // group gid
  Rdev   uint32      // device numbers
  Flags  uint32      // chflags(2) flags (OS X only)
}

// A Header describes the basic information sent in every request.
type Header struct {
  Conn *Conn     `json:"-"` // connection this request was received on
  ID   RequestID // unique ID for request
  Node NodeID    // file or directory the request is about
  Uid  uint32    // user ID of process making request
  Gid  uint32    // group ID of process making request
  Pid  uint32    // process ID of process making request
  // for returning to reqPool
  msg *message
}

type Dirent struct {
  Inode uint64
  Type DirentType (alias for uint32)
  Name string
}
```

# [golang](http://golang.org/) data types

```go
// package os
type FileInfo interface {
  Name() string       // base name of the file
  Size() int64        // length in bytes for regular files; system-dependent for others
  Mode() FileMode     // file mode bits
  ModTime() time.Time // modification time
  IsDir() bool        // abbreviation for Mode().IsDir()
  Sys() interface{}   // underlying data source (can return nil)
}

// package os FileMode constants
const (
  // The single letters are the abbreviations
  // used by the String methods formatting.
  ModeDir        FileMode = 1 << (32 - 1 - iota) // d: is a directory
  ModeAppend                                     // a: append-only
  ModeExclusive                                  // l: exclusive use
  ModeTemporary                                  // T: temporary file (not backed up)
  ModeSymlink                                    // L: symbolic link
  ModeDevice                                     // D: device file
  ModeNamedPipe                                  // p: named pipe (FIFO)
  ModeSocket                                     // S: Unix domain socket
  ModeSetuid                                     // u: setuid
  ModeSetgid                                     // g: setgid
  ModeCharDevice                                 // c: Unix character device, when ModeDevice is set
  ModeSticky                                     // t: sticky
  // Mask for the type bits. For regular files, none will be set.
  ModeType = ModeDir | ModeSymlink | ModeNamedPipe | ModeSocket | ModeDevice
  ModePerm FileMode = 0777 // permission bits
)

// package syscall
// This is actually the linux variation, because the go documentation
// always gives you the linux structures without ever telling you that.
type Stat_t struct {
  Dev       uint64
  Ino       uint64
  Nlink     uint64
  Mode      uint32
  Uid       uint32
  Gid       uint32
  X__pad0   int32
  Rdev      uint64
  Size      int64
  Blksize   int64
  Blocks    int64
  Atim      Timespec
  Mtim      Timespec
  Ctim      Timespec
  X__unused [3]int64
}
```
