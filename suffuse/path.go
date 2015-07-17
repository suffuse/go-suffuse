package suffuse

import (
  "os"
  "time"
  "strings"
  "path/filepath"
  "io/ioutil"
)

type Path string

var NoPath = Path("")

func Paths(paths ...string) []Path {
  xs := make([]Path, len(paths))
  for i, p := range paths { xs[i] = Path(p) }
  return xs
}

func MaybePath(path string, err error) Path {
  if err != nil { return NoPath }
  return Path(path)
}

func (x Path) MaybePath(result string, err error) Path {
  if err != nil { return x }
  return Path(result)
}

func (x Path) Name() string           { return filepath.Base(string(x))                        }
func (x Path) Extension() string      { return filepath.Ext(string(x))                         }
func (x Path) Absolute() Path         { return MaybePath(filepath.Abs(string(x)))              }
func (x Path) Relative(abs Path) Path { return MaybePath(filepath.Rel(string(x), string(abs))) }
func (x Path) Clean() Path            { return Path(filepath.Clean(string(x)))                 }
func (x Path) Parent() Path           { return Path(filepath.Dir(string(x)))                   }
func (x Path) Segments() []string     { return filepath.SplitList(string(x))                   }
func (x Path) EvalSymlinks() Path     { return MaybePath(filepath.EvalSymlinks(string(x)))     }

func (x Path) IsAbs() bool   { return filepath.IsAbs(string(x)) }
func (x Path) IsEmpty() bool { return string(x) == ""            }

func (x Path) OsChdir() error                         { return os.Chdir(string(x))                 }
func (x Path) OsChmod(mode os.FileMode) error         { return os.Chmod(string(x), mode)           }
func (x Path) OsChown(uid, gid int) error             { return os.Chown(string(x), uid, gid)       }
func (x Path) OsChtimes(atime, mtime time.Time) error { return os.Chtimes(string(x), atime, mtime) }
func (x Path) OsLchown(uid, gid int) error            { return os.Lchown(string(x), uid, gid)      }
func (x Path) OsLink(to Path) error                   { return os.Link(string(x), string(to))        }
func (x Path) OsMkdir(perm os.FileMode) error         { return os.Mkdir(string(x), perm)           }
func (x Path) OsMkdirAll(perm os.FileMode) error      { return os.MkdirAll(string(x), perm)        }
func (x Path) OsRemove() error                        { return os.Remove(string(x))                }
func (x Path) OsRemoveAll() error                     { return os.RemoveAll(string(x))             }
func (x Path) OsRename(to Path) error                 { return os.Rename(string(x), string(to))      }
func (x Path) OsSymlink(target Path) error            { return os.Symlink(string(x), string(target)) }
func (x Path) OsTruncate(size int64) error            { return os.Truncate(string(x), size)        }

func (x Path) OsCreate() (*os.File, error)                             { return os.Create(string(x))               }
func (x Path) OsLstat() (os.FileInfo, error)                           { return os.Lstat(string(x))                }
func (x Path) OsOpen() (*os.File, error)                               { return os.Open(string(x))                 }
func (x Path) OsOpenFile(flag int, perm os.FileMode) (*os.File, error) { return os.OpenFile(string(x), flag, perm) }
func (x Path) OsReadLink() (string, error)                             { return os.Readlink(string(x))             }
func (x Path) OsStat() (os.FileInfo, error)                            { return os.Stat(string(x))                 }

func (x Path) OsStatAtimeMtime() (atime time.Time, mtime time.Time) {
  // The go os.FileMode distillation of the contents of Stat_t doesn't
  // expose Atimespec, so we have to plunge into the platform-dependent part.
  stat, err := x.SysStatFile()
  if err != nil { return }
  return SysAtimeMtime(&stat)
}

func (x Path) IndexOfByte(b byte) int    { return strings.IndexByte(string(x), b) }
func (x Path) Slice(start, end int) Path { return Path(string(x)[start:end])      }

func (x Path) SplitAround(idx int) (Path, string) {
  max := len(string(x)) - 1
  if (idx > 0 && idx < max) {
    return x.Slice(0, idx), string(x)[idx+1:]
  }
  return x, ""
}
func (x Path) Split() (Path, string) {
  dir, file := filepath.Split(string(x))
  return Path(dir), file
}
func (x Path) Glob(glob string) []Path {
  ms, err := filepath.Glob(glob)
  if err != nil { return Paths() }
  return Paths(ms...)
}

func (x Path) WriteBytes(bs []byte) error        { return ioutil.WriteFile(string(x), bs, defaultFilePerms) }
func (x Path) WriteString(text string) error     { return x.WriteBytes([]byte(text))                      }
func (x Path) IoReadDir() ([]os.FileInfo, error) { return ioutil.ReadDir(string(x))                         }
func (x Path) IoReadFile() ([]byte, error)       { return ioutil.ReadFile(string(x))                        }

func (x Path) FollowOnce() Path { return x.MaybePath(x.OsReadLink()) }

func (x Path) Join(elem ...string) Path {
  buf := make([]string, len(elem) + 1)
  buf = append(buf, string(x))
  buf = append(buf, elem...)
  return Path(filepath.Join(buf...))
}

func (x Path) Slurp() string {
  bytes, err := ioutil.ReadFile(string(x))
  if err != nil { return "" }
  return string(bytes)
}
func (x Path) SlurpBytes() []byte {
  bytes, err := ioutil.ReadFile(string(x))
  if err != nil { bytes = nil }
  return bytes
}
func (x Path) Ino() uint64 {
  stat, err := x.SysStatLink()
  if err != nil { return 0 }
  return stat.Ino
}
func (x Path) Size() uint64 {
  fi, err := x.OsStat()
  if err != nil { return 0 }
  return uint64(fi.Size())
}

func (x Path) FollowAll() Path {
  next, err := x.OsReadLink()
  if err != nil { return x }
  return Path(next).FollowAll()
}

func (x Path) ReadDirnames() []string {
  fd, err := x.OsOpen()
  if err != nil { return nil }
  names, err := fd.Readdirnames(0)
  if err != nil { return nil }
  return names
}
func (x Path) ReadDirnodes() []os.FileInfo {
  fd, err := x.OsOpen()
  if err != nil { return nil }
  names, err := fd.Readdir(0)
  if err != nil { return nil }
  return names
}

func (x Path) Walk(walkFn filepath.WalkFunc) error {
  return filepath.Walk(string(x), walkFn)
}

// https://github.com/golang/go/issues/1312
func (x Path) FileExists() bool {
  _, err := x.OsStat()
  return err == nil
}
