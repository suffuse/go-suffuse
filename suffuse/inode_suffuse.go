package suffuse

/** Implementing the Suffuse read op interface for *Inode.
 *  Trying to keep the logic to a minimum - this should
 *  be limited to issuing errors and extracting/casting
 *  the relevant information from *Inode to fuse.
 */

import (
  "bazil.org/fuse"
)

func (x *Inode) MetaData() *fuse.Attr {
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

func (x *Inode) Lookup(name Name) SuffuseNode {
  if !x.IsDir() {
    return nil
  }
  child := x.Child(Name(name))
  // returning child here fails future nil tests
  if child == nil { return nil }
  return child
}

func (x *Inode) FileData() []byte {
  if x.IsAbsent() {
    return nil
  }
  if x.IsDir() {
    return nil
  }
  return x.Bytes()
}

func (x *Inode) DirData() []fuse.Dirent {
 if x.IsAbsent() {
    return nil
  }
  if !x.IsDir() {
    return nil
  }
  return x.Dirents()
}

func (x *Inode) LinkData() *LinkTarget {
  if x.IsAbsent() {
    return nil
  }
  if x.IsLink() {
    link := x.LinkTarget()
    return &link
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
