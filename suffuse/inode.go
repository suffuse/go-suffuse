package suffuse

import (
  "strings"
  "sort"
  "bazil.org/fuse"
)

/** Every directory inode has to have a way to obtain fresh
 *  inodes to respond to e.g. Mkdir and Create calls. We could
 *  reduce the footprint by splitting out Dir inodes from the
 *  others, but for now every Inode carries the pointer.
 */
type Inode struct {
  *InodeGen
  AttrMap
}

func (x *Inode) AttrKeys()[]AttrKey {
  keys := make([]AttrKey, 0, len(x.AttrMap))
  for k := range x.AttrMap { keys = append(keys, k) }
  return keys
}
func (x *Inode) ChildNames()[]Name {
  ds := x.DirList()
  var buf []string
  for k := range ds { buf = append(buf, string(k)) }
  sort.Strings(buf)
  return Names(buf...)
}

func (x *Inode) String()string {
  return Sprintf("%v[%v]", x.InodeType(), x.InodeNum())
}
func (x *Inode) TreeString()string {
  return Strs(expand(0, NoName, x)...).JoinLines()
}

func (x *Inode) SetAttr(k AttrKey, v AttrValue) {
  x.AttrMap[k] = v
}
func (x *Inode) WithAttr(k AttrKey, v AttrValue)*Inode {
  x.SetAttr(k, v)
  return x
}

func (x *Inode) AttrOr(key AttrKey, alt AttrValue)(r AttrValue) {
  r = x.AttrMap[key]
  if r == nil { r = alt }
  return
}

func (x *Inode) Text()string   { return string(x.Bytes())          }

func (x *Inode) AddChildDir(name Name)(*Inode, error) {
  if x.IsDir() {
    child := x.NewDir()
    x.DirList()[name] = child
    return child, nil
  }
  return nil, NotADir()
}
func (x *Inode) AddChild(name Name, ino *Inode)error {
  if x.IsDir() {
    x.DirList()[name] = ino
    return nil
  }
  return NotADir()
}
func (x *Inode) Child(name Name)*Inode {
  if x.IsDir() {
    return x.DirList()[name]
  }
  return nil
}
func (x *Inode) Dirents()[]fuse.Dirent {
  var res []fuse.Dirent
  for _, name := range x.ChildNames() {
    child := x.Child(name)
    res = append(res, child.FuseDirent(name))
  }
  return res
}

func (x *Inode) FuseDirent(name Name)fuse.Dirent {
  return fuse.Dirent {
    Inode: uint64(x.InodeNum()),
    Type: x.InodeType().ToFuseType(),
    Name: string(name),
  }
}

func expand(level int, name Name, node *Inode)[]string {
  indent := strings.Repeat(" ", level * 4)
  line := Sprintf("%v[%v]%v", indent, node.InodeNum(), name)

  switch node.InodeType() {
    case InodeDir:
      res := []string { line + "/" }
      for _, name := range node.ChildNames() {
        res = append(res, expand(level + 1, name, node.Child(name))...)
      }
      return res
    case InodeLink:
      return []string { Sprintf("%v -> %q", line, node.LinkTarget()) }
    case InodeFile:
      return []string { Sprintf("%v: %q", line, node.Text()) }
    default:
      return []string { line }
  }
}

func (x *Inode) AddNodeMap(nodes map[Path]interface{})error {
  for path, node := range nodes {
    err := x.AddPath(path, node)
    if err != nil { return err }
  }
  return nil
}
func (x *Inode) AddNode(segs []Name, node *Inode)error {
  switch len(segs) {
    case 0: return NotExist()
    case 1: return x.AddChild(segs[0], node)
    default:
      next := x.Child(segs[0])
      if next == nil {
        next = x.NewDir()
        x.AddChild(segs[0], next)
      }
      return next.AddNode(segs[1:], node)
  }
}
func (x *Inode) AddPath(path Path, data interface{})error {
  switch data := data.(type) {
    case Bytes:
      ino := x.NewFile()
      ino.SetAttr(BytesKey, data)
      return x.AddNode(path.Segments(), ino)
    case LinkTarget:
      ino := x.NewLink()
      ino.SetAttr(LinkTargetKey, data)
      return x.AddNode(path.Segments(), ino)
    default:
      Echoerr("Failed: %v %T", data, data)
      return NotImplemented()
  }
}
