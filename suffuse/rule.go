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
  MetaData (*Inode) *fuse.Attr
  FileData (*Inode) []byte
  DirData  (*Inode) []fuse.Dirent
  LinkData (*Inode) *LinkTarget
}

/** Default implementations. A Rule struct can embed BaseRule
 *  and then not have to implement irrelevant methods. There's
 *  no default for Metadata because a filesystem which never
 *  returns metadata is a useless filesystem.
 */
type BaseRule struct { }
func (*BaseRule) FileData(node *Inode) []byte        { return nil }
func (*BaseRule) DirData (node *Inode) []fuse.Dirent { return nil }
func (*BaseRule) LinkData(node *Inode) *LinkTarget   { return nil }

/** Only marker interfaces at the moment. They will carry data
 *  when filesystems become more complex.
 */
type IdRule struct { BaseRule }

type AttrRule struct { BaseRule }

type FileCommand struct {
  BaseRule
  Separator, Command, Example string
}

type FileConversion struct {
  BaseRule
  From, To, Command, Example string
}

func (*IdRule) MetaData(node *Inode) *fuse.Attr {
  path := node.Path
  fi, err := path.OsLstat()
  if err != nil { return nil }
  attr := GoFileInfoToFuseAttr(fi)
  return &attr
}
func (*IdRule) FileData(node *Inode) []byte {
  path := node.Path
  return path.SlurpBytes()
}
func (*IdRule) DirData(node *Inode) []fuse.Dirent {
  path := node.Path
  names := path.ReadDirnames()
  if names == nil { return nil }
  size := len(names)
  ds := make([]fuse.Dirent, size + 2)

  ds[0] = fuse.Dirent { Name: ".",  Inode: path.Ino(),          Type: fuse.DT_Dir }
  ds[1] = fuse.Dirent { Name: "..", Inode: path.Parent().Ino(), Type: fuse.DT_Dir }

  for i, name := range names { ds[i+2] = childDirent(path, name) }

  return ds
}
func (*IdRule) LinkData(node *Inode) *LinkTarget {
  path := node.Path
  target, err := os.Readlink(string(path))
  if err != nil { return nil }
  res := LinkTarget(target)
  return &res
}

func (*AttrRule) MetaData(x *Inode) *fuse.Attr {
  if x.IsAbsent() { return nil }

  attr := &fuse.Attr {
    Uid  : uint32(x.Uid()),
    Gid  : uint32(x.Gid()),
    Atime: x.Times()[0],
    Mtime: x.Times()[1],
    Mode : x.FuseMode(),
    Nlink: 1,
  }
  switch x.InodeType() {
    case InodeDir  : attr.Nlink = uint32(len(x.DirList()) + 2)
    case InodeLink : attr.Size = uint64(len(x.LinkTarget()))
    case InodeFile : attr.Size = uint64(len(x.Bytes()))
    default        : return nil
  }

  return attr
}
func (*AttrRule) DirData(x *Inode) []fuse.Dirent {
 if x.IsAbsent() {
    return nil
  }
  if !x.IsDir() {
    return nil
  }
  return x.Dirents()
}
func (*AttrRule) FileData(x *Inode) []byte {
  if x.IsAbsent() {
    return nil
  }
  if x.IsDir() {
    return nil
  }
  return x.Bytes()
}
func (*AttrRule) LinkData(x *Inode) *LinkTarget {
  if x.IsAbsent() {
    return nil
  }
  if x.IsLink() {
    link := x.LinkTarget()
    return &link
  }
  return nil
}

func (x *FileCommand) MetaData(node *Inode) *fuse.Attr {
  path := node.Path
  subpath, cmd := x.split(path)
  if cmd == "" { return nil }

  fi, err := subpath.OsStat()
  if err != nil { return nil }

  bytes := x.FileData(node)
  if bytes == nil { return nil }

  attr := GoFileInfoToFuseAttr(fi)
  attr.Size = uint64(len(bytes))
  return &attr
}
func (x *FileCommand) FileData(node *Inode) []byte {
  path := node.Path
  subpath, cmd := x.split(path)
  if cmd == "" { return nil }

  var c = x.Command
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
  }
  return subpath, ""
}

func (x *FileConversion) MetaData(node *Inode) *fuse.Attr {
  path := node.Path
  subpath := x.real(path)

  fi, err := subpath.OsStat()
  if err != nil { return nil }

  bytes := x.FileData(node)
  if bytes == nil { return nil }

  attr := GoFileInfoToFuseAttr(fi)
  attr.Size = uint64(len(bytes))
  return &attr
}
func (x *FileConversion) FileData(node *Inode)[]byte {
  path := node.Path
  subpath := x.real(path)

  c := strings.Replace(x.Command, "$file", string(subpath), -1)

  res := Exec("sh", "-c", c)
  if res.Success() { return res.Stdout }
  return nil
}
func (x *FileConversion) real(p Path) Path {
  return Path(strings.Replace(string(p), "." + x.To, "", 1))
}

func childDirent(x Path, name string) fuse.Dirent {
  child := x.Join(name)

  return fuse.Dirent {
    Inode: child.Ino(),
    Type: direntType(child),
    Name: name,
  }
}

func direntType(x Path) fuse.DirentType {
  fi, err := x.OsLstat()
  if err != nil { return fuse.DT_Unknown }
  return GoModeToDirentType(fi.Mode())
}
