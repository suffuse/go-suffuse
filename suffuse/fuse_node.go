package suffuse

import (
  "golang.org/x/net/context"
  "bazil.org/fuse"
  "bazil.org/fuse/fs"
)

type FuseNode struct {
  path Path
  node SuffuseNode
}

func NewFuseRoot(root SuffuseNode) *FuseNode {
  return &FuseNode{NoPath, root}
}

func (x *FuseNode) Lookup(ctx context.Context, name string) (fs.Node, error) {
  trace("[%v] Lookup(%v)", x.path, name)

  if node := x.node.Lookup(Name(name)); node != nil {
    return &FuseNode{x.path.Join(name), node}, nil 
  }

  return nil, NotExist()
}

func (x *FuseNode) Attr(ctx context.Context, attrRef *fuse.Attr) error {
  trace("[%v] Attr", x)

  if a := x.node.MetaData(); a != nil { *attrRef = *a; return nil }

  return NotExist()
}

func (x *FuseNode) ReadAll(ctx context.Context) ([]byte, error) {
  trace("[%v] ReadAll", x)

  if bytes := x.node.FileData(); bytes != nil { return bytes, nil }

  return nil, NotExist()
}

func (x *FuseNode) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
  trace("[%v] ReadDirAll", x)

  if children := x.node.DirData(); children != nil { return children, nil }

  return nil, NotADir()
}

func (x *FuseNode) Readlink(ctx context.Context, req *fuse.ReadlinkRequest) (string, error) {
  trace("[%v] Readlink", x)

  if target := x.node.LinkData(); target != nil { return string(*target), nil }

  return "", NotValidArg()
}
