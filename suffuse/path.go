package suffuse

import (
  "os"
  "time"
  "strings"
  "path/filepath"
  "io/ioutil"
)

type Path struct {
  Path string
}
type SfsWalkFunc func(string, os.FileInfo) string

func NewPath(path string) Path {
  return Path { path }
}
func NewPaths(paths ...string) []Path {
  xs := make([]Path, len(paths))
  for i, p := range paths { xs[i] = NewPath(p) }
  return xs
}
func Cwd() Path        { return MaybePath(os.Getwd()) }
func CwdNoLinks() Path { return Cwd().EvalSymlinks()  }

var RootPath         = NewPath("/")
var DotPath          = NewPath(".")
var DotDotPath       = NewPath("..")
var DefaultFilePerms = os.FileMode(0644)
var DefaultDirPerms  = os.FileMode(0755)

func (x Path) MaybeNewPath(result string, err error) Path {
  if err != nil { return x }
  return NewPath(result)
}
func MaybePath(path string, err error) Path {
  if err != nil { return NoPath }
  return NewPath(path)
}
func MaybeString(result string, err error) string {
  if err != nil { return "" }
  return result
}
func MaybeByteString(result []byte, err error) string {
  if err != nil { return "" }
  return string(result)
}

func (x Path) GoAbs() (string, error)          { return filepath.Abs(x.Path)           }
func (x Path) GoBase() string                  { return filepath.Base(x.Path)          }
func (x Path) GoClean() string                 { return filepath.Clean(x.Path)         }
func (x Path) GoDir() string                   { return filepath.Dir(x.Path)           }
func (x Path) GoEvalSymlinks() (string, error) { return filepath.EvalSymlinks(x.Path)  }
func (x Path) GoExt() string                   { return filepath.Ext(x.Path)           }
func (x Path) GoRel(abs Path) (string, error)  { return filepath.Rel(x.Path, abs.Path) }
func (x Path) GoSplitList() []string           { return filepath.SplitList(x.Path)     }

func (x Path) Name() string           { return x.GoBase()                    }
func (x Path) Extension() string      { return x.GoExt()                     }
func (x Path) Absolute() Path         { return MaybePath(x.GoAbs())          }
func (x Path) Relative(abs Path) Path { return MaybePath(x.GoRel(abs))       }
func (x Path) Clean() Path            { return NewPath(x.GoClean())          }
func (x Path) Parent() Path           { return NewPath(x.GoDir())            }
func (x Path) Segments() []string     { return x.GoSplitList()               }
func (x Path) EvalSymlinks() Path     { return MaybePath(x.GoEvalSymlinks()) }

func (x Path) IsAbs() bool        { return filepath.IsAbs(x.Path)     }
func (x Path) IsEmpty() bool      { return x.Path == ""               }

func (x Path) OsChdir() error                         { return os.Chdir(x.Path)                 }
func (x Path) OsChmod(mode os.FileMode) error         { return os.Chmod(x.Path, mode)           }
func (x Path) OsChown(uid, gid int) error             { return os.Chown(x.Path, uid, gid)       }
func (x Path) OsChtimes(atime, mtime time.Time) error { return os.Chtimes(x.Path, atime, mtime) }
func (x Path) OsLchown(uid, gid int) error            { return os.Lchown(x.Path, uid, gid)      }
func (x Path) OsLink(to Path) error                   { return os.Link(x.Path, to.Path)         }
func (x Path) OsMkdir(perm os.FileMode) error         { return os.Mkdir(x.Path, perm)           }
func (x Path) OsMkdirAll(perm os.FileMode) error      { return os.MkdirAll(x.Path, perm)        }
func (x Path) OsRemove() error                        { return os.Remove(x.Path)                }
func (x Path) OsRemoveAll() error                     { return os.RemoveAll(x.Path)             }
func (x Path) OsRename(to Path) error                 { return os.Rename(x.Path, to.Path)       }
func (x Path) OsSymlink(target Path) error            { return os.Symlink(x.Path, target.Path)  }
func (x Path) OsTruncate(size int64) error            { return os.Truncate(x.Path, size)        }

func (x Path) OsCreate() (*os.File, error)                             { return os.Create(x.Path)               }
func (x Path) OsLstat() (os.FileInfo, error)                           { return os.Lstat(x.Path)                }
func (x Path) OsOpen() (*os.File, error)                               { return os.Open(x.Path)                 }
func (x Path) OsOpenFile(flag int, perm os.FileMode) (*os.File, error) { return os.OpenFile(x.Path, flag, perm) }
func (x Path) OsReadLink() (string, error)                             { return os.Readlink(x.Path)             }
func (x Path) OsStat() (os.FileInfo, error)                            { return os.Stat(x.Path)                 }

func (x Path) OsStatAtimeMtime() (atime time.Time, mtime time.Time) {
  // The go os.FileMode distillation of the contents of Stat_t doesn't
  // expose Atimespec, so we have to plunge into the platform-dependent part.
  stat, err := x.SysStatFile()
  if err != nil { return }
  return SysAtimeMtime(&stat)
}

func (x Path) IndexOfByte(b byte) int { return strings.IndexByte(x.Path, '#') }

func (x Path) SplitAround(idx int) (Path, string) {
  max := len(x.Path) - 1
  if (idx > 0 && idx < max) {
    return NewPath(x.Path[:idx]), x.Path[idx+1:]
  }
  return x, ""
}
func (x Path) Split() (Path, string) {
  dir, file := filepath.Split(x.Path)
  return Path { dir }, file
}
func (x Path) Glob(glob string) []Path {
  ms, err := filepath.Glob(glob)
  if err != nil { return NewPaths() }
  return NewPaths(ms...)
}

func (x Path) WriteBytes(bs []byte) error        { return ioutil.WriteFile(x.Path, bs, DefaultFilePerms) }
func (x Path) WriteString(text string) error     { return x.WriteBytes([]byte(text))                     }
func (x Path) IoReadDir() ([]os.FileInfo, error) { return ioutil.ReadDir(x.Path)                         }
func (x Path) IoReadFile() ([]byte, error)       { return ioutil.ReadFile(x.Path)                        }

func (x Path) FollowOnce() Path { return x.MaybeNewPath(x.OsReadLink()) }

func (x Path) Join(elem ...string) Path {
  buf := make([]string, len(elem) + 1)
  buf = append(buf, x.Path)
  buf = append(buf, elem...)
  return NewPath(filepath.Join(buf...))
}

func (x Path) Slurp() string {
  bytes, err := ioutil.ReadFile(x.Path)
  if err != nil { return "" }
  return string(bytes)
}
func (x Path) SlurpBytes() []byte {
  bytes, err := ioutil.ReadFile(x.Path)
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

func (p Path) String() string { return p.Path }

func (x Path) FollowAll() Path {
  next, err := x.OsReadLink()
  if err != nil { return x }
  return NewPath(next).FollowAll()
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
  return filepath.Walk(x.Path, walkFn)
}
func (x Path) WalkCollect(f SfsWalkFunc) Lines {
  res := make([]string, 0)
  x.Walk(
    func(path string, info os.FileInfo, err error) error {
      rel, err := filepath.Rel(x.Path, path)
      if err != nil {
        x := f(rel, info)
        if x != "" { res = append(res, x) }
      }
      return nil
    },
  )
  return NewLines(res...)
}
