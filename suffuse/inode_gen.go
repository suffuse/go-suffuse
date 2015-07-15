package suffuse

type InodeGen struct {
  fresh chan InodeNum
}

func NewRootInode(startInode uint64)*Inode {
  ichan := make(chan InodeNum)
  // Monotonically increments the counter without any
  // effort to track or reuse.
  go func() {
    count := uint64(startInode)
    for {
      ichan <- InodeNum(count)
      count += 1
    }
  }()

  gen := InodeGen { ichan }
  return gen.NextDir()
}

func (x *InodeGen) Next(tp InodeType) *Inode {
  ino := Inode { InodeGen: x, AttrMap: AttrMap{} }
  ino.SetAttr(InodeNumKey, <- x.fresh)
  ino.SetAttr(InodeTypeKey, tp)
  ino.SetAttr(PermBitsKey, BasePermBits())
  return &ino
}
func (x *InodeGen) NextDir() *Inode {
  return x.Next(InodeDir).WithAttr(DirListKey, DirList{})
}
func (x *InodeGen) NextFile() *Inode {
  return x.Next(InodeFile).WithAttr(BytesKey, []byte{})
}
func (x *InodeGen) NextLink() *Inode {
  return x.Next(InodeLink).WithAttr(LinkTargetKey, NoLinkTarget)
}
