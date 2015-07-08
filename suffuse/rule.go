package suffuse

import (
  "os"
  "bazil.org/fuse"
)

/** An interface encapsulating the essential fuse operations.
 *  Each method returns nil if inapplicable. The fuse operations
 *  use the first non-nil return value found for a given Path.
 */
type Rule interface {
  MetaData(Path)*fuse.Attr
  FileData(Path)[]byte
  DirData(Path)[]fuse.Dirent
  LinkData(Path)*Path
}

/** Default implementations. A Rule struct can embed BaseRule
 *  and then not have to implement irrelevant methods. There's
 *  no default for Metadata because a filesystem which never
 *  returns metadata is a useless filesystem.
 */
type BaseRule struct { }
func (BaseRule) FileData(path Path)[]byte       { return nil }
func (BaseRule) DirData(path Path)[]fuse.Dirent { return nil }
func (BaseRule) LinkData(path Path)*Path        { return nil }

/** Only marker interfaces at the moment. They will carry data
 *  when filesystems become more complex.
 */
type IdRule  struct { BaseRule }
type SedRule struct { BaseRule }

func (IdRule) MetaData(path Path)*fuse.Attr {
  fi, err := path.OsLstat()
  if err != nil { return nil }
  attr := GoFileInfoToFuseAttr(fi)
  return &attr
}
func (IdRule) FileData(path Path)[]byte {
  return path.SlurpBytes()
}
func (IdRule) DirData(path Path)[]fuse.Dirent {
  return DirChildren(path)
}
func (IdRule) LinkData(path Path)*Path {
  target, err := os.Readlink(string(path))
  if err != nil { return nil }
  res := Path(target)
  return &res
}

func (x SedRule) MetaData(path Path)*fuse.Attr {
  subpath, cmd := x.split(path)
  if cmd == "" { return nil }

  fi, err := subpath.OsStat()
  if err != nil { return nil }

  bytes := x.FileData(path)
  if bytes == nil { return nil }

  attr := GoFileInfoToFuseAttr(fi)
  attr.Size = uint64(len(bytes))
  return &attr
}
func (x SedRule) FileData(path Path)[]byte {
  subpath, cmd := x.split(path)
  if cmd == "" { return nil }

  res := Exec("sed", "-ne", cmd, string(subpath))
  if res.Success() { return res.Stdout }
  return nil
}
func (SedRule) split(p Path) (Path, string) {
  return p.SplitAround(p.IndexOfByte('#'))
}
