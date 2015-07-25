package suffuse

type InodeGen struct {
  fresh chan InodeNum
}

func NewRoot(startInode uint64) *Inode {
  ichan := make(chan InodeNum)
  // Monotonically increments the counter without any
  // effort to track or reuse.
  go func() {
    count := uint64(startInode)
    for {
      ichan <- InodeNum(count)
      count++
    }
  }()

  gen := InodeGen { ichan }
  return gen.New(InodeDir)
}
func (x *InodeGen) New(tp InodeType) *Inode {
  ino := &Inode { inodeGen: x, AttrMap: AttrMap{} }
  ino.SetAttr(InodeNumKey, <- x.fresh)
  ino.SetAttr(InodeTypeKey, tp)
  ino.SetAttr(PermBitsKey, BasePermBits())
  switch tp {
    case InodeDir : return ino.WithAttr(DirListKey, DirList{})
    case InodeFile: return ino.WithAttr(BytesKey, []byte{})
    case InodeLink: return ino.WithAttr(LinkTargetKey, NoLinkTarget)
    default: return ino
  }
}

func (x *Inode) New(tp InodeType, name Name) (*Inode, error) {
  if x.IsDir() {
    child := x.inodeGen.New(tp)
    x.DirList()[name] = child
    return child, nil
  }
  return nil, NotADir()
}
func (x *Inode) NewDir(name Name) (*Inode, error) {
  return x.New(InodeDir, name)
}
func (x *Inode) NewFile(name Name) (*Inode, error) {
  return x.New(InodeFile, name)
}
func (x *Inode) NewLink(name Name) (*Inode, error) {
  return x.New(InodeLink, name)
}
