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

  trace("[%v] Attr", x)

  for _, rule := range rules {
    a := rule.MetaData(x)
    if a != nil {
      *attrRef = *a
      return nil
    }
  }

  return NotExist()
}

func (x *Inode) Lookup(ctx context.Context, name string) (fs.Node, error) {
  trace("[%v] Lookup(%v)", x.Path, name)

  // TODO: move the next four lines to a separate function, maybe something like: GetRealNode
  fi, err := x.Path.Join(name).OsLstat()
  if IsNilError(err) {
    return x.New(FileModeToInodeType(fi.Mode()), Name(name))
  }

  if !x.IsDir() {
    return nil, NotADir()
  }
  child := x.Child(Name(name))
  if child != nil {
    return child, nil
  }

  return x.New(InodeNone, Name(name))
}

func (x *Inode) ReadAll(ctx context.Context) ([]byte, error) {
  trace("[%v] ReadAll", x)

  for _, rule := range rules {
    bytes := rule.FileData(x)
    if bytes != nil {
      return bytes, nil
    }
  }
  return nil, NotExist()
}

func (x *Inode) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
  trace("[%v] ReadDirAll", x)

  for _, rule := range rules {
    children := rule.DirData(x)
    if children != nil { return children, nil }
  }
  return nil, NotADir()
}

func (x *Inode) Readlink(ctx context.Context, req *fuse.ReadlinkRequest) (string, error) {
  trace("[%v] Readlink", x)

  for _, rule := range rules {
    target := rule.LinkData(x)
    if target != nil { return string(*target), nil }
  }
  return "", NotValidArg()
}

func (x *Inode) Forget() {
  // Echoerr("Forget %v", *x)
}
