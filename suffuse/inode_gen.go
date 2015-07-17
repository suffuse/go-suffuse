package suffuse

type InodeGen struct {
  fresh chan InodeNum
}

func NewRoot(startInode uint64)*Inode {
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
  return gen.NewDir()
}

func (x *InodeGen) New(tp InodeType) *Inode {
  ino := Inode { InodeGen: x, AttrMap: AttrMap{} }
  ino.SetAttr(InodeNumKey, <- x.fresh)
  ino.SetAttr(InodeTypeKey, tp)
  ino.SetAttr(PermBitsKey, BasePermBits())
  return &ino
}
func (x *InodeGen) NewDir() *Inode {
  return x.New(InodeDir).WithAttr(DirListKey, DirList{})
}
func (x *InodeGen) NewFile() *Inode {
  return x.New(InodeFile).WithAttr(BytesKey, []byte{})
}
func (x *InodeGen) NewLink() *Inode {
  return x.New(InodeLink).WithAttr(LinkTargetKey, NoLinkTarget)
}
