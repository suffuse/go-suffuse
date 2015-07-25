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
    if result := node.Lookup(name); result != nil { return result }
  }
  return nil
}

func (x *CompoundNode) MetaData() *fuse.Attr {
  for _, node := range x.nodes {
    if result := node.MetaData(); result != nil { return result }
  }
  return nil
}
func (x *CompoundNode) DirData() []fuse.Dirent {
  for _, node := range x.nodes {
    if result := node.DirData(); result != nil { return result }
  }
  return nil
}
func (x *CompoundNode) FileData() []byte {
  for _, node := range x.nodes {
    if result := node.FileData(); result != nil { return result }
  }
  return nil
}
func (x *CompoundNode) LinkData() *LinkTarget {
  for _, node := range x.nodes {
    if result := node.LinkData(); result != nil { return result }
  }
  return nil
}

