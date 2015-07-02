package suffuse

import (
  "github.com/suffuse/go-suffuse/suffuse/xattr"
  "os"
  "time"
  "golang.org/x/net/context"
  "bazil.org/fuse/fs"
  f "bazil.org/fuse"
)

// Access checks whether the calling context has permission for
// the given operations on the receiver. If so, Access should
// return nil. If not, Access should return EPERM.
//
// Note that this call affects the result of the access(2) system
// call but not the open(2) system call. If Access is not
// implemented, the Node behaves as if it always returns nil
// (permission granted), relying on checks in Open instead.
func (x *Elem) Access(ctx context.Context, req *f.AccessRequest) error {
  logD("Access", "path", x.Path)
  // TODO - permissions check.
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
func (x *Elem) Setattr(ctx context.Context, req *f.SetattrRequest, resp *f.SetattrResponse) error {
  logD("Setattr", "path", x.Path, "req", *req)
  // TODO - permissions check on ops which require it (see also Access.)
  // TODO - accumulate and report errors somehow.

  // There doesn't appear to be any way to set one or the other
  // time, you always have to specify two values. So if both are
  // valid according to fuse, we pass both, but if only one is,
  // we stat the current values and pass what was already there
  // for the not-being-set-right-now slot.
  if req.Valid.Atime() || req.Valid.Mtime() {
    var atime, mtime time.Time
    if req.Valid.Atime() && req.Valid.Mtime() {
      atime = req.Atime
      mtime = req.Mtime
    } else {
      // one of them needs to be preserved
      atime, mtime = x.Path.OsStatAtimeMtime()
      if req.Valid.Atime() { atime = req.Atime }
      if req.Valid.Mtime() { mtime = req.Mtime }
    }
    x.Path.OsChtimes(atime, mtime)
  }
  // TODO - identify which bits besides permission bits can
  // theoretically be set, as opposed to being unmodifiable
  // consequences of the inode.
  if req.Valid.Mode() {
    perms := req.Mode & os.ModePerm
    x.Path.OsChmod(perms)
  }
  if req.Valid.Uid() || req.Valid.Gid() {
    uid, gid := -1, -1
    if req.Valid.Uid() { uid = int(req.Uid) }
    if req.Valid.Gid() { gid = int(req.Gid) }
    x.Path.OsChown(uid, gid)
  }
  // "Truncate" is kind of a scary name, but it's used to both shrink
  // and increase the size of files.
  if req.Valid.Size() {
    x.Path.OsTruncate(int64(req.Size))
  }

  /** TODO - or at least to investigate. **/
  // if req.Valid.Handle() {
  // }
  // if req.Valid.LockOwner() {
  // }
  // if req.Valid.Crtime() {
  // }
  // if r.Valid.Chgtime() {
  // }
  // if r.Valid.Bkuptime() {
  // }
  // if r.Valid.Flags() {
  // }

  return nil
}

func (x *Elem) Attr(ctx context.Context, attr *f.Attr) error {
  logD("Attr", "path", x.Path)

     path := x.Path
  fi, err := path.OsLstat()
  var a f.Attr

  if err != nil {
    subpath, cmd := splitSedSuffix(path)
    if cmd == "" {
      return err
    }
    fi, err = subpath.OsStat()
    a = GoFileInfoToFuseAttr(fi)
    a.Size = uint64(len(slurpSedSuffix(path)))
  } else {
    a = GoFileInfoToFuseAttr(fi)
  }

  *attr = a
  return nil
}

func (x *Elem) Lookup(ctx context.Context, name string) (fs.Node, error) {
  logD("Lookup", "path", x.Path, "name", name)
  child := x.Path.Join(name)

  var a f.Attr
  err := x.Attr(ctx, &a)
  MaybeLog(err)

  if err != nil { return nil, err }
  return Vnode(child), nil
}

func (x *Elem) ReadDirAll(ctx context.Context) (dirents []f.Dirent, err error) {
  logD("ReadDirAll", "path", x.Path)

  dirents = DirChildren(x.Path)
  if dirents == nil { err = ENOTDIR }
  return
}
func (x *Elem) Readlink(ctx context.Context, req *f.ReadlinkRequest) (s string, err error) {
  logD("Readlink", "path", x.Path)

  if x.IsLink() {
    return os.Readlink(x.Path.Path)
  } else {
    s = ""
    err = EINVAL
  }
  return //"", EINVAL
}
func (x *Elem) Read(ctx context.Context, req *f.ReadRequest, resp *f.ReadResponse) error {
  logD("Read", "path", x.Path, "req", *req)
  bytes := FindBytes(x.Path)
  HandleRead(req, resp, bytes)
  return nil
}
func (x *Elem) ReadAll(ctx context.Context) (bytes []byte, err error) {
  logD("ReadAll", "path", x.Path)
  bytes = FindBytes(x.Path)
  if bytes == nil { err = f.ENOENT }
  return
}

/** Write ops.
 */

func (x *Elem) Mkdir(ctx context.Context, req *f.MkdirRequest) (fs.Node, error) {
  logD("Mkdir", "path", x.Path, "req", *req)

  name := req.Name
  mode := req.Mode
  umask := req.Umask
  path := x.Path.Join(name)
  err := path.OsMkdir(mode & umask)

  MaybeLog(err)

  if err != nil { return nil, err }
  return Dir(path), nil
}
func (x *Elem) Create(ctx context.Context, req *f.CreateRequest, resp *f.CreateResponse) (fs.Node, fs.Handle, error) {
  logD("Create", "path", x.Path, "req", *req)
  return nil, nil, f.ENOTSUP
}

/** Xattr ops.
 */
func (x *Elem) Getxattr(ctx context.Context, req *f.GetxattrRequest, resp *f.GetxattrResponse) error {
  logD("Getxattr", "path", x.Path, "xattr", req.Name)
  bytes := xattr.Get(x.Path.Path, req.Name)
  if bytes != nil { resp.Xattr = bytes }
  return nil
}
func (x *Elem) Listxattr(ctx context.Context, req *f.ListxattrRequest, resp *f.ListxattrResponse) error {
  logD("Listxattr", "path", x.Path)
  names := xattr.List(x.Path.Path)
  if names != nil { resp.Append(names...) }
  return nil
}
func (x *Elem) Setxattr(ctx context.Context, req *f.SetxattrRequest) error {
  logD("Setxattr", "path", x.Path, "name", req.Name, "value", string(req.Xattr))
  return xattr.Set(x.Path.Path, req.Name, req.Xattr)
}
func (x *Elem) Removexattr(ctx context.Context, req *f.RemovexattrRequest) error {
  logD("Removexattr", "path", x.Path, "name", req.Name)
  return xattr.Remove(x.Path.Path, req.Name)
}
