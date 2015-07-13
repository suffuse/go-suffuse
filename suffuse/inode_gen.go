package suffuse

type InodeGen struct {
  count uint64
}

func NewRootInode()*Inode {
  gen := InodeGen { 2 }
  return gen.NextDir()
}

func (x *InodeGen) Next() *Inode {
  ino := Inode { InodeGen: x, AttrMap: AttrMap{} }
  ino.SetAttr(InodeNumKey, InodeNum(x.count))
  ino.SetAttr(PermBitsKey, BasePermBits())
  x.count += 1
  return &ino
}
func (x *InodeGen) NextDir() *Inode {
  ino := x.Next()
  ino.SetAttr(InodeTypeKey, InodeDir)
  ino.SetAttr(DirListKey, DirList{})
  return ino
}
func (x *InodeGen) NextFile() *Inode {
  ino := x.Next()
  ino.SetAttr(InodeTypeKey, InodeFile)
  ino.SetAttr(BytesKey, NoBytes)
  return ino
}
func (x *InodeGen) NextLink() *Inode {
  ino := x.Next()
  ino.SetAttr(InodeTypeKey, InodeLink)
  return ino
}
