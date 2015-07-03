package suffuse

import (
  "github.com/suffuse/go-suffuse/suffuse/xattr"
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

func NewIdNode(path Path) *IdNode {
  return &IdNode { Path: path }
}

func (x *IdNode) Access(ctx context.Context, req *f.AccessRequest) error {
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

func (x *IdNode) Attr(ctx context.Context, attr *f.Attr) error {
  logD("Attr", "path", x.Path)

  for _, rule := range Rules {
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
  logD("Lookup", "path", x.Path, "name", name)
  child := x.Path.Join(name)

  var a f.Attr
  err := x.Attr(ctx, &a)
  if err != nil { return nil, err } // f.ENOENT }
  return NewIdNode(child), nil
}

func (x *IdNode) ReadDirAll(ctx context.Context) ([]f.Dirent, error) {
  logD("ReadDirAll", "path", x.Path)

  for _, rule := range Rules {
    children := rule.DirData(x.Path)
    if children != nil { return children, nil }
  }
  return nil, ENOTDIR
}
func (x *IdNode) Readlink(ctx context.Context, req *f.ReadlinkRequest) (string, error) {
  path := x.Path
  logD("Readlink", "path", path)

  for _, rule := range Rules {
    target := rule.LinkData(path)
    if target != nil { return target.Path, nil }
  }
  return "", EINVAL
}
func (x *IdNode) Read(ctx context.Context, req *f.ReadRequest, resp *f.ReadResponse) error {
  path := x.Path
  logD("Read", "path", path, "req", *req)

  for _, rule := range Rules {
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
  logD("ReadAll", "path", path)

  for _, rule := range Rules {
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
  logD("Mkdir", "path", x.Path, "req", *req)

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
  logD("Create", "path", x.Path, "req", *req)
  return nil, nil, f.ENOTSUP
}


/** Xattr ops.
 */
func (x *IdNode) Getxattr(ctx context.Context, req *f.GetxattrRequest, resp *f.GetxattrResponse) error {
  logD("Getxattr", "path", x.Path, "xattr", req.Name)
  bytes := xattr.Get(x.Path.Path, req.Name)
  if bytes != nil { resp.Xattr = bytes }
  return nil
}
func (x *IdNode) Listxattr(ctx context.Context, req *f.ListxattrRequest, resp *f.ListxattrResponse) error {
  logD("Listxattr", "path", x.Path)
  names := xattr.List(x.Path.Path)
  if names != nil { resp.Append(names...) }
  return nil
}
func (x *IdNode) Setxattr(ctx context.Context, req *f.SetxattrRequest) error {
  logD("Setxattr", "path", x.Path, "name", req.Name, "value", string(req.Xattr))
  return xattr.Set(x.Path.Path, req.Name, req.Xattr)
}
func (x *IdNode) Removexattr(ctx context.Context, req *f.RemovexattrRequest) error {
  logD("Removexattr", "path", x.Path, "name", req.Name)
  return xattr.Remove(x.Path.Path, req.Name)
}
