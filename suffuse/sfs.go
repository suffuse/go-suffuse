package suffuse

import (
  "golang.org/x/net/context"
  "bazil.org/fuse/fs"
  "bazil.org/fuse"
)

const (
  SfsConfigFileName = ".sfs"
)

// There are two approaches to writing a FUSE file system.  The first is to speak
// the low-level message protocol, reading from a Conn using ReadRequest and
// writing using the various Respond methods.  This approach is closest to
// the actual interaction with the kernel and can be the simplest one in contexts
// such as protocol translators.
//
// Servers of synthesized file systems tend to share common
// bookkeeping abstracted away by the second approach, which is to
// call fs.Serve to serve the FUSE protocol using an implementation of
// the service methods in the interfaces FS* (file system), Node* (file
// or directory), and Handle* (opened file or directory).
func (u *Sfs) Root() (fs.Node, error) { return u.RootNode, nil }

func (u *Sfs) Init(ctx context.Context, req *fuse.InitRequest, resp *fuse.InitResponse) error {
  logI("Init", "req", req)
  return nil
}

func (u *Sfs) Destroy(ctx context.Context) {
  logI("Destroy")
}

func (u *Sfs) Statfs(ctx context.Context, req *fuse.StatfsRequest, resp *fuse.StatfsResponse) error {
  logD("Statfs")
  stat, err := u.RootNode.GetPath().SysStatfs()
  if err != nil { return err }
  *resp = SysStatfsToFuseStatfs(stat)
  return nil
}

func (u *Sfs) Unmount() error {
  return u.Mountpoint.SysUnmount()
}
func (u *Sfs) Serve() error {
  defer u.Connection.Close()
  return fs.Serve(u.Connection, u)
}
