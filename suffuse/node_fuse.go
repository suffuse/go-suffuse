package suffuse

import (
  "github.com/suffuse/go-suffuse/suffuse/xattr"
  "os"
  "golang.org/x/net/context"
  "bazil.org/fuse/fs"
  f "bazil.org/fuse"
)

func (x *Elem) Access(ctx context.Context, req *f.AccessRequest) error {
  logD("Access", "path", x.Path)

  // Access checks whether the calling context has permission for
  // the given operations on the receiver. If so, Access should
  // return nil. If not, Access should return EPERM.
  //
  // Note that this call affects the result of the access(2) system
  // call but not the open(2) system call. If Access is not
  // implemented, the Node behaves as if it always returns nil
  // (permission granted), relying on checks in Open instead.
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
  if err != nil { return nil, f.ENOENT }
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
