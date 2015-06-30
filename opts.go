package suffuse

import (
  "flag"
  "os"
  "bazil.org/fuse"
)

type SfsOpts struct {
  Config, Mountpoint, ExcludeRegex, VolName string
  Scratch, Verbose, Debug bool
  Args []string
}

func UsageAndExit() {
  flag.Usage()
  os.Exit(2)
}

func (opts SfsOpts) Create() (Sfs, error) {
  SetLogLevel(opts.GetLogLevel())
         mnt := opts.GetMountpoint()
  mount_opts := []fuse.MountOption { fuse.FSName("suffuse") }

  mount_opts = append(mount_opts, PlatformOptions()...)
  if opts.VolName != "" {
    mount_opts = append(mount_opts,
      fuse.LocalVolume(),
      fuse.VolumeName(opts.VolName),
    )
  }
  if opts.Config != "" {
    configFileOpts := ReadJsonFile(NewPath(opts.Config))
    Echoerr("%v", configFileOpts)
  }

  c, err := fuse.Mount(mnt.Path, mount_opts...)
  if err != nil { return Sfs{}, err }

  mfs := Sfs {
    Mountpoint : mnt,
    RootNode   : NewUnion(opts.GetPaths()...),
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

func (opts SfsOpts) GetMountpoint() Path {
  switch {
    case opts.Scratch          : return ScratchDir()
    case opts.Mountpoint != "" : return NewPath(opts.Mountpoint)
    default                    : UsageAndExit() ; return Path{}
  }
}

func (opts SfsOpts) GetPaths() []Path {
  if len(opts.Args) == 0 { UsageAndExit() }
  xs := make([]Path, len(opts.Args))
  for i, p := range opts.Args { xs[i] = NewPath(p) }
  return xs
}

func ParseSfsOpts() SfsOpts {
  opts := SfsOpts{}

  flag.StringVar(&opts.Config, "c", "", "suffuse config file")
  flag.StringVar(&opts.Mountpoint, "m", "", "mount point")
  flag.StringVar(&opts.ExcludeRegex, "x", "", "name exclusion regex")
  flag.StringVar(&opts.VolName, "n", "", "volume name (OSX only)")
  flag.BoolVar(&opts.Scratch, "t", false, "create scratch directory as mount point")
  flag.BoolVar(&opts.Verbose, "v", false, "log at INFO level")
  flag.BoolVar(&opts.Debug, "d", false, "log at DEBUG level")
  flag.Usage = func () {
    Echoerr("Usage: %s <options> [path path ...]\n", os.Args[0])
    flag.PrintDefaults()
  }
  flag.Parse()
  opts.Args = flag.Args()
  return opts
}
