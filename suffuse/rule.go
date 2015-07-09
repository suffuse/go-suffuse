package suffuse

import (
  "os"
  "strings"
  "bazil.org/fuse"
)

/** An interface encapsulating the essential fuse operations.
 *  Each method returns nil if inapplicable. The fuse operations
 *  use the first non-nil return value found for a given Path.
 */
type Rule interface {
  MetaData (Path) *fuse.Attr
  FileData (Path) []byte
  DirData  (Path) []fuse.Dirent
  LinkData (Path) *Path
}

/** Default implementations. A Rule struct can embed BaseRule
 *  and then not have to implement irrelevant methods. There's
 *  no default for Metadata because a filesystem which never
 *  returns metadata is a useless filesystem.
 */
type BaseRule struct { }
func (*BaseRule) FileData(path Path) []byte        { return nil }
func (*BaseRule) DirData (path Path) []fuse.Dirent { return nil }
func (*BaseRule) LinkData(path Path) *Path         { return nil }

/** Only marker interfaces at the moment. They will carry data
 *  when filesystems become more complex.
 */
type IdRule  struct { BaseRule }

type FileCommand struct {
  BaseRule
  Separator, Command, Example string
}

type FileConversion struct {
  BaseRule
  From, To, Command, Example string
}

func (*IdRule) MetaData(path Path) *fuse.Attr {
  fi, err := path.OsLstat()
  if err != nil { return nil }
  attr := GoFileInfoToFuseAttr(fi)
  return &attr
}
func (*IdRule) FileData(path Path) []byte {
  return path.SlurpBytes()
}
func (*IdRule) DirData(path Path) []fuse.Dirent {
  return DirChildren(path)
}
func (*IdRule) LinkData(path Path) *Path {
  target, err := os.Readlink(string(path))
  if err != nil { return nil }
  res := Path(target)
  return &res
}


func (x *FileCommand) MetaData(path Path)*fuse.Attr {
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
func (x *FileCommand) FileData(path Path)[]byte {
  subpath, cmd := x.split(path)
  if cmd == "" { return nil }

  var c string = x.Command
  c = strings.Replace(c, "$file", string(subpath), -1)
  c = strings.Replace(c, "$args", cmd, -1)

  res := Exec("sh", "-c", c)
  if res.Success() { return res.Stdout }
  return nil
}
func (x *FileCommand) split(p Path) (Path, string) {
  split   := strings.Split(string(p), x.Separator)
  subpath := Path(split[0])
  if len(split) > 1 {
    return subpath, split[1]
  } else {
    return subpath, ""
  }
}

func (x *FileConversion) MetaData(path Path)*fuse.Attr {
  subpath := x.real(path)

  fi, err := subpath.OsStat()
  if err != nil { return nil }

  bytes := x.FileData(path)
  if bytes == nil { return nil }

  attr := GoFileInfoToFuseAttr(fi)
  attr.Size = uint64(len(bytes))
  return &attr
}
func (x *FileConversion) FileData(path Path)[]byte {
  subpath := x.real(path)

  c := strings.Replace(x.Command, "$file", string(subpath), -1)

  res := Exec("sh", "-c", c)
  if res.Success() { return res.Stdout }
  return nil
}
func (x *FileConversion) real(p Path) Path {
  return Path(strings.Replace(string(p), "." + x.To, "", 1))
}
