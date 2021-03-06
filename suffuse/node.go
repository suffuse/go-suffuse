package suffuse

import (
  "os"
  f "bazil.org/fuse"
)

// HandleRead handles a read request assuming that data is the entire file content.
// It adjusts the amount returned in resp according to req.Offset and req.Size.
func HandleRead(req *f.ReadRequest, resp *f.ReadResponse, data []byte) {
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

func NewFilePerms(bits uint32) os.FileMode {
  return os.FileMode(bits) & os.ModePerm
}

func direntType(x Path) f.DirentType {
  fi, err := x.OsLstat()
  if err != nil { return f.DT_Unknown }
  return GoModeToDirentType(fi.Mode())
}

func childDirent(x Path, name string) f.Dirent {
  child := x.Join(name)

  return f.Dirent {
    Inode: child.Ino(),
    Type: direntType(child),
    Name: name,
  }
}

func DirChildren(x Path) []f.Dirent {
  names := x.ReadDirnames()
  if names == nil { return nil }
  size := len(names)
  ds := make([]f.Dirent, size + 2)

  ds[0] = f.Dirent { Inode: x.Ino(), Type: f.DT_Dir, Name: "." }
  ds[1] = f.Dirent { Inode: x.Parent().Ino(), Type: f.DT_Dir, Name: ".." }

  for i, name := range names { ds[i+2] = childDirent(x, name) }

  return ds
}
