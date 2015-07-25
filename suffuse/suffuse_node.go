package suffuse

import (
  "bazil.org/fuse"
)

type SuffuseNode interface {
  Lookup(Name) SuffuseNode

  MetaData () *fuse.Attr
  FileData () []byte
  DirData  () []fuse.Dirent
  LinkData () *LinkTarget
}

type CompoundNode struct {
  nodes []SuffuseNode
}

func NewCompoundNode(nodes ...SuffuseNode) *CompoundNode {
  return &CompoundNode{nodes}
}

func (x *CompoundNode) Lookup(name Name) SuffuseNode {
  for _, node := range x.nodes {
    result := node.Lookup(name)
    if result != nil { return result }
  }
  return nil
}

func (x *CompoundNode) MetaData() *fuse.Attr {
  for _, node := range x.nodes {
    result := node.MetaData()
    if result != nil { return result }
  }
  return nil
}
func (x *CompoundNode) DirData() []fuse.Dirent {
  for _, node := range x.nodes {
    result := node.DirData()
    if result != nil { return result }
  }
  return nil
}
func (x *CompoundNode) FileData() []byte {
  for _, node := range x.nodes {
    result := node.FileData()
    if result != nil { return result }
  }
  return nil
}
func (x *CompoundNode) LinkData() *LinkTarget {
  for _, node := range x.nodes {
    result := node.LinkData()
    if result != nil { return result }
  }
  return nil
}

