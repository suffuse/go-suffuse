package suffuse

/** Implementing the fuse read op interface for *Inode.
 *  Trying to keep the logic to a minimum - this should
 *  be limited to issuing errors and extracting/casting
 *  the relevant information from *Inode to fuse.
 */

import (
  "golang.org/x/net/context"
  "bazil.org/fuse"
  "bazil.org/fuse/fs"
)

/** The read-only ops are covered by
 *    Attr, Lookup, ReadAll, ReadDirAll, Readlink
 *  Finer grained ops avoided by the *All variants are
 *    Read, Readdir
 */
func (x *Inode) Attr(ctx context.Context, attrRef *fuse.Attr) error {
  if x.IsAbsent() { return NotExist() }

  attr := fuse.Attr {
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
    default        : return NotExist()
  }

  *attrRef = attr
  return nil
}

func (x *Inode) Lookup(ctx context.Context, name string) (fs.Node, error) {
  if !x.IsDir() {
    return nil, NotADir()
  } else {
    child := x.Child(Name(name)) ; if child != nil {
      return child, nil
    } else {
      return nil, NotExist()
    }
  }
}

func (x *Inode) ReadAll(ctx context.Context) ([]byte, error) {
  if x.IsDir() {
    return nil, IsADir()
  } else {
    return x.Bytes(), nil
  }
}

func (x *Inode) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
  if x.IsAbsent() {
    return nil, NotExist()
  } else if !x.IsDir() {
    return nil, NotADir()
  } else {
    return x.Dirents(), nil
  }
}

func (x *Inode) Readlink(ctx context.Context, req *fuse.ReadlinkRequest) (string, error) {
  if x.IsLink() {
    return string(x.LinkTarget()), nil
  } else {
    return "", NotValidArg()
  }
}
