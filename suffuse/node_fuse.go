package suffuse

import (
  "github.com/paulp/suffuse/suffuse/xattr"
  "os"
  "golang.org/x/net/context"
  "bazil.org/fuse/fs"
  f "bazil.org/fuse"
)

func (x Elem) Attr(ctx context.Context, attr *f.Attr) error {
  logD("Attr", "path", x.Path)

  a := x.GetVattr()
  switch a.(type) {
    case VattrOk: *attr = a.Attr()
  }
  return a.Err()
}
func (x Elem) Lookup(ctx context.Context, name string) (fs.Node, error) {
  logD("Lookup", "path", x.Path, "name", name)

  if x.IsDir() {
    return Vnode(x.Path.Join(name)), nil
  }
  return nil, ENOTDIR
}
func (x Elem) ReadDirAll(ctx context.Context) (dirents []f.Dirent, err error) {
  logD("ReadDirAll", "path", x.Path)

  dirents = DirChildren(x.Path)
  if dirents == nil { err = ENOTDIR }
  return
}
func (x Elem) Readlink(ctx context.Context, req *f.ReadlinkRequest) (s string, err error) {
  logD("Readlink", "path", x.Path)

  if x.IsLink() {
    return os.Readlink(x.Path.Path)
  } else {
    s = ""
    err = EINVAL
  }
  return //"", EINVAL
}
func (x Elem) Read(ctx context.Context, req *f.ReadRequest, resp *f.ReadResponse) error {
  logD("Read", "path", x.Path, "req", *req)
  bytes := FindBytes(x.Path)
  HandleRead(req, resp, bytes)
  return nil
}
func (x Elem) ReadAll(ctx context.Context) (bytes []byte, err error) {
  logD("ReadAll", "path", x.Path)
  bytes = FindBytes(x.Path)
  if bytes == nil { err = f.ENOENT }
  return
}
func (x Elem) Getxattr(ctx context.Context, req *f.GetxattrRequest, resp *f.GetxattrResponse) error {
  logD("Getxattr", "path", x.Path, "xattr", req.Name)
  bytes := xattr.Get(x.Path.Path, req.Name)
  if bytes != nil { resp.Xattr = bytes }
  return nil
}
func (x Elem) Listxattr(ctx context.Context, req *f.ListxattrRequest, resp *f.ListxattrResponse) error {
  logD("Listxattr", "path", x.Path)
  names := xattr.List(x.Path.Path)
  if names != nil { resp.Append(names...) }
  return nil
}
func (x Elem) Setxattr(ctx context.Context, req *f.SetxattrRequest) error {
  logD("Setxattr", "path", x.Path, "name", req.Name, "value", string(req.Xattr))
  return xattr.Set(x.Path.Path, req.Name, req.Xattr)
}
func (x Elem) Removexattr(ctx context.Context, req *f.RemovexattrRequest) error {
  logD("Removexattr", "path", x.Path, "name", req.Name)
  return xattr.Remove(x.Path.Path, req.Name)
}
