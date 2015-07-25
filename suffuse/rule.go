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
  IsValidFor (Path) bool
  MetaData   (Path) *fuse.Attr
  FileData   (Path) []byte
  DirData    (Path) []fuse.Dirent
  LinkData   (Path) *LinkTarget
}

type RulesNode struct {
  path Path
  rules []Rule
}

func NewRulesNodeRoot(root Path, rules []Rule) *RulesNode {
  return &RulesNode{root, rules}
}

func (x *RulesNode) Lookup(name Name) SuffuseNode {
  path := x.path.Join(string(name))
  for _, rule := range x.rules {
    if rule.IsValidFor(path) {
      return &RulesNode{path, x.rules}
    }
  }
  return nil
}
func (x *RulesNode) MetaData() *fuse.Attr {
  for _, rule := range x.rules {
    a := rule.MetaData(x.path)
    if a != nil {
      return a
    }
  }
  return nil
}
func (x *RulesNode) DirData() []fuse.Dirent {
  for _, rule := range x.rules {
    d := rule.DirData(x.path)
    if d != nil {
      return d
    }
  }
  return nil
}
func (x *RulesNode) FileData() []byte {
  for _, rule := range x.rules {
    f := rule.FileData(x.path)
    if f != nil {
      return f
    }
  }
  return nil
}
func (x *RulesNode) LinkData() *LinkTarget {
  for _, rule := range x.rules {
    l := rule.LinkData(x.path)
    if l != nil {
      return l
    }
  }
  return nil
}

/** Default implementations. A Rule struct can embed BaseRule
 *  and then not have to implement irrelevant methods. There's
 *  no default for Metadata because a filesystem which never
 *  returns metadata is a useless filesystem.
 */
type BaseRule struct { }
func (*BaseRule) FileData(Path) []byte        { return nil }
func (*BaseRule) DirData (Path) []fuse.Dirent { return nil }
func (*BaseRule) LinkData(Path) *LinkTarget   { return nil }

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

func (*IdRule) IsValidFor(path Path) bool {
  _, err := path.OsLstat()
  return err == nil
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
  names := path.ReadDirnames()
  if names == nil { return nil }
  size := len(names)
  ds := make([]fuse.Dirent, size + 2)

  ds[0] = fuse.Dirent { Name: ".",  Inode: path.Ino(),          Type: fuse.DT_Dir }
  ds[1] = fuse.Dirent { Name: "..", Inode: path.Parent().Ino(), Type: fuse.DT_Dir }

  for i, name := range names { ds[i+2] = childDirent(path, name) }

  return ds
}
func (*IdRule) LinkData(path Path) *LinkTarget {
  target, err := os.Readlink(string(path))
  if err != nil { return nil }
  res := LinkTarget(target)
  return &res
}

func (x *FileCommand) IsValidFor(path Path) bool {
  subpath, cmd := x.split(path)
  // fail fast
  if cmd == "" { return false }

   _, err := subpath.OsStat()
  return err == nil

}
func (x *FileCommand) MetaData(path Path) *fuse.Attr {

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
func (x *FileCommand) FileData(path Path) []byte {

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

func (x *FileConversion) IsValidFor(path Path) bool {
  subpath := x.real(path)

  _, err := subpath.OsStat()
  return err == nil
}
func (x *FileConversion) MetaData(path Path) *fuse.Attr {

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
