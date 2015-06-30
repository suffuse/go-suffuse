package suffuse

import (
  "bazil.org/fuse"
  "bazil.org/fuse/fs"
)

var NoPath Path      = NewPath("")
var NoNode *Elem     = Dir(NoPath)
var NoAttr fuse.Attr = fuse.Attr{}

type Vdir interface {
  Children() map[string]Node
}
type Vfile interface {
  Data() []byte
}

type Node interface {
  fs.Node

  GetPath() Path
  GetType() fuse.DirentType
}

type Sfs struct {
  Mountpoint Path
  RootNode Node
  Connection *fuse.Conn
}

type Elem struct {
  fs.NodeRef

  Typ fuse.DirentType
  Path Path
}
func NewElem(tp fuse.DirentType, path Path) *Elem {
  return &Elem {
     Typ: tp,
    Path: path,
  }
}
func Dir(path Path)   *Elem { return NewElem(fuse.DT_Dir, path)      }
func File(path Path)  *Elem { return NewElem(fuse.DT_File, path)     }
func Link(path Path)  *Elem { return NewElem(fuse.DT_Link, path)     }
func Vnode(path Path) *Elem { return NewElem(direntType(path), path) }

func NewUnion(paths ...Path) Node {
  switch len(paths) {
    case 0  : return NoNode
    case 1  : return Dir(paths[0])
    default : panic("No unionfs yet") ; return NoNode // TODO
  }
}
