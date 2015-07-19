package suffuse

type InodeGen struct {
  fresh chan InodeNum
}

func NewRoot(startInode uint64, path Path) *Inode {
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
  return gen.New(InodeDir, path).WithAttr(DirListKey, DirList{})
}
func (x *InodeGen) New(tp InodeType, path Path) *Inode {
  ino := Inode { inodeGen: x, AttrMap: AttrMap{}, Path: path }
  ino.SetAttr(InodeNumKey, <- x.fresh)
  ino.SetAttr(InodeTypeKey, tp)
  ino.SetAttr(PermBitsKey, BasePermBits())
  return &ino
}

func (x *Inode) New(tp InodeType, name Name) (*Inode, error) {
  if x.IsDir() {
    child := x.inodeGen.New(tp, x.Path.Join(string(name)))
    x.DirList()[name] = child
    return child, nil
  }
  return nil, NotADir()
}
func (x *Inode) NewDir(name Name) (*Inode, error) {
  dir, err := x.New(InodeDir, name)
  if IsError(err) { return nil, err }
  return dir.WithAttr(DirListKey, DirList{}), nil
}
func (x *Inode) NewFile(name Name) (*Inode, error) {
  file, err := x.New(InodeFile, name)
  if IsError(err) { return nil, err }
  return file.WithAttr(BytesKey, []byte{}), nil
}
func (x *Inode) NewLink(name Name) (*Inode, error) {
  link, err := x.New(InodeLink, name)
  if IsError(err) { return nil, err }
  return link.WithAttr(LinkTargetKey, NoLinkTarget), nil
}
