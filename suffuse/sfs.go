package suffuse

import (
  "bazil.org/fuse/fs"
  "bazil.org/fuse"
  "os"
)

/** Sfs == Suffuse File System.
 */
type Sfs struct {
  Mountpoint Path
  RootNode fs.Node
  Connection *fuse.Conn
}

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

func NewSfs(conf *SfsConfig) (*Sfs, error) {
  SetLogLevel(conf.LogLevel)
         mnt := conf.Mountpoint
  mount_opts := getFuseMountOptions(conf)

  if !conf.Config.IsEmpty() {
    configFileOpts := readJsonFile(conf.Config)
    Echoerr("%v", configFileOpts)
  }

  c, err := fuse.Mount(string(mnt), mount_opts...)
  if err != nil { return nil, err }

  mfs := &Sfs {
    Mountpoint : mnt,
    RootNode   : NewIdNode(conf.Paths[0]),
    Connection : c,
  }

  /** Start a goroutine which looks for INT/TERM and
   *  calls unmount on the filesystem. Otherwise ctrl-C
   *  leaves us with ghost mounts which outlive the process.
   */
  trap := func (sig os.Signal) {
    Echoerr("Caught %s - forcing unmount(2) of %s\n", sig, mfs.Mountpoint)
    mfs.Unmount()
  }
  TrapExit(trap)
  return mfs, nil
}

func getFuseMountOptions(conf *SfsConfig) []fuse.MountOption {
  mount_opts := []fuse.MountOption { fuse.FSName("suffuse") }
  mount_opts = append(mount_opts, PlatformOptions()...)
  if conf.VolName != "" {
    mount_opts = append(mount_opts,
      fuse.LocalVolume(),
      fuse.VolumeName(conf.VolName),
    )
  }
  return mount_opts
}

func (u *Sfs) Root() (fs.Node, error) { return u.RootNode, nil }

func (u *Sfs) Unmount() error {
  return u.Mountpoint.SysUnmount()
}
func (u *Sfs) Serve() error {
  defer u.Connection.Close()
  return fs.Serve(u.Connection, u)
}
