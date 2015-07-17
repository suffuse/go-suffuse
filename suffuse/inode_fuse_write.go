package suffuse

/** The fuse write operations.
 *  See also inode_fuse_read.go.
 */

import (
  "golang.org/x/net/context"
  "bazil.org/fuse"
  "bazil.org/fuse/fs"
)

/** The basic writable filesystem ops are
 *    Setattr, Create, Mkdir, Link
 */
func (x *Inode) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
  // TODO: req.Mode, req.Umask
  return x.AddChildDir(Name(req.Name))
}
func (x *Inode) Mknod(ctx context.Context, req *fuse.MknodRequest) (fs.Node, error) {
  Echoerr("%v.Mknod(*v)", *x, *req)
  return nil, NotImplemented()
}

func (x *Inode) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
  Echoerr("[Create]\n  %+v\n  %+v", *x, *req)

  switch FuseTypeToInodeType(GoModeToDirentType(req.Mode)) {
    case InodeFile:
      ino := x.NewFile()
      x.AddChild(Name(req.Name), ino)
      return ino, ino, nil
    default:
      return nil, nil, NotSupported()
  }
  return nil, nil, NotSupported()
}
func (x *Inode) Symlink(ctx context.Context, req *fuse.SymlinkRequest) (fs.Node, error) {
  return nil, NotSupported()
}

func (x *Inode) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
  start := Offset(req.Offset)
  length := Length(len(req.Data))
  rng := Range { start, length }
  // start := Index(req.Offset)
  // end := start + Index(len(req.Data))
  // rng := Range { start, end }
  Echoerr("[Write] %v[%v] = %q", *x, rng, string(req.Data))
  return NotSupported()
}
func (x *Inode) Setattr(ctx context.Context, req *fuse.SetattrRequest, resp *fuse.SetattrResponse) error {
  Echoerr("[Setattr]\n  %+v\n  %+v", *x, *req)
  return NotSupported()
}

/** Some relevant request/response structures.
 */

// type MkdirRequest struct {
//   Name   string
//   Mode   os.FileMode
//   Umask  os.FileMode
// }
// type MknodRequest struct {
//   Name   string
//   Mode   os.FileMode
//   Rdev   uint32
//   Umask  os.FileMode
// }
// type CreateRequest struct {
//   Name   string
//   Flags  OpenFlags
//   Mode   os.FileMode
//   Umask  os.FileMode
// }
// type WriteRequest struct {
//   Handle    HandleID
//   Offset    int64
//   Data      []byte
//   Flags     WriteFlags
//   LockOwner uint64
//   FileFlags OpenFlags
// }

// type CreateResponse struct {
//   LookupResponse
//   OpenResponse
// }
// type LookupResponse struct {
//   Node       NodeID
//   Generation uint64
//   EntryValid time.Duration
//   Attr       Attr
// }
// type OpenResponse struct {
//   Handle HandleID
//   Flags  OpenResponseFlags
// }
