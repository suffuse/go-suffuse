package suffuse

import (
  "github.com/suffuse/go-suffuse/suffuse/xattr"
  "os"
  "time"
  "golang.org/x/net/context"
  "bazil.org/fuse/fs"
  f "bazil.org/fuse"
)

/** Identity node is a wrapper around another Path.
 */
type IdNode struct {
  fs.NodeRef
  Path Path
}

var NoNode = NewIdNode(NoPath)

func NewIdNode(path Path) *IdNode {
  return &IdNode { Path: path }
}

func (x *IdNode) String() string { return string(x.Path) }

// Access checks whether the calling context has permission for
// the given operations on the receiver. If so, Access should
// return nil. If not, Access should return EPERM.
//
// Note that this call affects the result of the access(2) system
// call but not the open(2) system call. If Access is not
// implemented, the Node behaves as if it always returns nil
// (permission granted), relying on checks in Open instead.
// TODO - permissions check.
func (x *IdNode) Access(ctx context.Context, req *f.AccessRequest) error {
  trace("Access", "path", x.Path)
  return nil
}

// There doesn't appear to be any way to set one or the other
// time, you always have to specify two values. So if both are
// valid according to fuse, we pass both, but if only one is,
// we stat the current values and pass what was already there
// for the not-being-set-right-now slot.
func setattrChtimes(path Path, req *f.SetattrRequest) error {
  isAtime := req.Valid.Atime()
  isMtime := req.Valid.Mtime()

  if isAtime || isMtime {
    var atime, mtime time.Time
    if isAtime && isMtime {
      atime = req.Atime
      mtime = req.Mtime
    } else {
      // one of them needs to be preserved
      atime, mtime = path.OsStatAtimeMtime()
      if isAtime { atime = req.Atime }
      if isMtime { mtime = req.Mtime }
    }
    return path.OsChtimes(atime, mtime)
  }
  return nil
}

// man 2 chown:
//   One of the owner or group id's may be left unchanged by specifying it as -1.
func setattrChown(path Path, req *f.SetattrRequest) error {
  isValidUid := req.Valid.Uid()
  isValidGid := req.Valid.Gid()

  if isValidUid || isValidGid {
    uid, gid := -1, -1
    if isValidUid { uid = int(req.Uid) }
    if isValidGid { gid = int(req.Gid) }
    return path.OsChown(uid, gid)
  }
  return nil
}

// TODO - identify which bits besides permission bits can
// theoretically be set, as opposed to being unmodifiable
// consequences of the inode.
func setattrChmod(path Path, req *f.SetattrRequest) error {
  if req.Valid.Mode() {
    return path.OsChmod(req.Mode & os.ModePerm)
  }
  return nil
}

// "Truncate" is kind of a scary name, but it's used to both shrink
// and increase the size of files.
func setattrTruncate(path Path, req *f.SetattrRequest) error {
  if req.Valid.Size() {
    return path.OsTruncate(int64(req.Size))
  }
  return nil
}

// Setattr sets the standard metadata for the receiver.
//
// Note, this is also used to communicate changes in the size of
// the file. Not implementing Setattr causes writes to be unable
// to grow the file (except with OpenDirectIO, which bypasses that
// mechanism).
//
// req.Valid is a bitmask of what fields are actually being set.
// For example, the method should not change the mode of the file
// unless req.Valid.Mode() is true.
func (x *IdNode) Setattr(ctx context.Context, req *f.SetattrRequest, resp *f.SetattrResponse) error {
  trace("Setattr", "path", x.Path, "req", *req)
  // Not Yet Implemented:
  // req.Valid.{ Handle, LockOwner, Crtime, Chgtime, Bkuptime, Flags }

  return FindError(
    setattrChtimes(x.Path, req),
    setattrChown(x.Path, req),
    setattrChmod(x.Path, req),
    setattrTruncate(x.Path, req),
  )
}

func (x *IdNode) Attr(ctx context.Context, attr *f.Attr) error {
  trace("Attr", "path", x.Path)

  for _, rule := range rules {
    a := rule.MetaData(x.Path)
    if a != nil {
      *attr = *a
      return nil
    }
  }

  _, err := x.Path.OsStat()
  return err
}

func (x *IdNode) Lookup(ctx context.Context, name string) (fs.Node, error) {
  trace("Lookup", "path", x.Path, "name", name)
  child := x.Path.Join(name)

  var a f.Attr
  err := x.Attr(ctx, &a)

  MaybeLog(err)

  if err != nil { return nil, err } // f.ENOENT }
  return NewIdNode(child), nil
}

func (x *IdNode) ReadDirAll(ctx context.Context) ([]f.Dirent, error) {
  trace("ReadDirAll", "path", x.Path)

  for _, rule := range rules {
    children := rule.DirData(x.Path)
    if children != nil { return children, nil }
  }
  return nil, ENOTDIR
}
func (x *IdNode) Readlink(ctx context.Context, req *f.ReadlinkRequest) (string, error) {
  path := x.Path
  trace("Readlink", "path", path)

  for _, rule := range rules {
    target := rule.LinkData(path)
    if target != nil { return string(*target), nil }
  }
  return "", EINVAL
}
func (x *IdNode) Read(ctx context.Context, req *f.ReadRequest, resp *f.ReadResponse) error {
  path := x.Path
  trace("Read", "path", path, "req", *req)

  for _, rule := range rules {
    bytes := rule.FileData(path)
    if bytes != nil {
      HandleRead(req, resp, bytes)
      return nil
    }
  }
  return nil
}
func (x *IdNode) ReadAll(ctx context.Context) ([]byte, error) {
  path := x.Path
  trace("ReadAll", "path", path)

  for _, rule := range rules {
    bytes := rule.FileData(path)
    if bytes != nil {
      return bytes, nil
    }
  }
  return nil, f.ENOENT
}

/** Write ops.
 */

func (x *IdNode) Mkdir(ctx context.Context, req *f.MkdirRequest) (fs.Node, error) {
  trace("Mkdir", "path", x.Path, "req", *req)

  name := req.Name
  mode := req.Mode
  umask := req.Umask
  path := x.Path.Join(name)
  err := path.OsMkdir(mode & umask)

  MaybeLog(err)

  if err != nil { return nil, err }
  return NewIdNode(path), nil
}
func (x *IdNode) Create(ctx context.Context, req *f.CreateRequest, resp *f.CreateResponse) (fs.Node, fs.Handle, error) {
  trace("Create", "path", x.Path, "req", *req)
  return nil, nil, f.ENOTSUP
}

/** Xattr ops.
 */
func (x *IdNode) Getxattr(ctx context.Context, req *f.GetxattrRequest, resp *f.GetxattrResponse) error {
  trace("Getxattr", "path", x.Path, "xattr", req.Name)
  bytes := xattr.Get(string(x.Path), req.Name)
  if bytes != nil { resp.Xattr = bytes }
  return nil
}
func (x *IdNode) Listxattr(ctx context.Context, req *f.ListxattrRequest, resp *f.ListxattrResponse) error {
  trace("Listxattr", "path", x.Path)
  names := xattr.List(string(x.Path))
  if names != nil { resp.Append(names...) }
  return nil
}
func (x *IdNode) Setxattr(ctx context.Context, req *f.SetxattrRequest) error {
  trace("Setxattr", "path", x.Path, "name", req.Name, "value", string(req.Xattr))
  return xattr.Set(string(x.Path), req.Name, req.Xattr)
}
func (x *IdNode) Removexattr(ctx context.Context, req *f.RemovexattrRequest) error {
  trace("Removexattr", "path", x.Path, "name", req.Name)
  return xattr.Remove(string(x.Path), req.Name)
}
