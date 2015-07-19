package suffuse

import (
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

func (x *IdNode) Attr(ctx context.Context, attr *f.Attr) error {
  trace("[%v] Attr", x.Path)

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
  trace("[%v] Lookup(%v)", x.Path, name)
  child := x.Path.Join(name)

  var a f.Attr
  err := x.Attr(ctx, &a)

  MaybeLog(err)

  if err != nil { return nil, err } // f.ENOENT }
  return NewIdNode(child), nil
}

func (x *IdNode) ReadDirAll(ctx context.Context) ([]f.Dirent, error) {
  trace("[%v] ReadDirAll", x.Path)

  for _, rule := range rules {
    children := rule.DirData(x.Path)
    if children != nil { return children, nil }
  }
  return nil, NotADir()
}
func (x *IdNode) Readlink(ctx context.Context, req *f.ReadlinkRequest) (string, error) {
  path := x.Path
  trace("[%v] Readlink", path)

  for _, rule := range rules {
    target := rule.LinkData(path)
    if target != nil { return string(*target), nil }
  }
  return "", NotValidArg()
}
func (x *IdNode) Read(ctx context.Context, req *f.ReadRequest, resp *f.ReadResponse) error {
  path := x.Path
  trace("[%v] Read(%+v)", path, *req)

  for _, rule := range rules {
    bytes := rule.FileData(path)
    if bytes != nil {
      handleRead(req, resp, bytes)
      return nil
    }
  }
  return nil
}
func (x *IdNode) ReadAll(ctx context.Context) ([]byte, error) {
  path := x.Path
  trace("[%v] ReadAll", path)

  for _, rule := range rules {
    bytes := rule.FileData(path)
    if bytes != nil {
      return bytes, nil
    }
  }
  return nil, NotExist()
}

// HandleRead handles a read request assuming that data is the entire file content.
// It adjusts the amount returned in resp according to req.Offset and req.Size.
func handleRead(req *f.ReadRequest, resp *f.ReadResponse, data []byte) {
  if req.Offset >= int64(len(data)) {
    data = nil
  } else {
    data = data[req.Offset:]
  }
  if len(data) > req.Size {
    data = data[:req.Size]
  }
  n := copy(resp.Data[:req.Size], data)
  resp.Data = resp.Data[:n]
}
