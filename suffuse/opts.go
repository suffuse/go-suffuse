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
  fs *flag.FlagSet
}

func (opts SfsOpts) UsageAndExit() {
  opts.fs.Usage()
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

  // Exactly one incoming path for the moment.
  // Eventually more than one could mean a union mount.
  if len(opts.Args) != 1 { opts.UsageAndExit() }

  mfs := Sfs {
    Mountpoint : mnt,
    RootNode   : NewIdNode(NewPath(opts.Args[0])),
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
    default                    : opts.UsageAndExit() ; return Path{}
  }
}

func ParseSfsOpts() SfsOpts {
  opts := SfsOpts{}

  fs := flag.NewFlagSet("suffuse", flag.ExitOnError)
  // fs.StringVar(&opts.Config, "c", "", "suffuse config file")
  // fs.StringVar(&opts.ExcludeRegex, "x", "", "name exclusion regex")
  fs.StringVar(&opts.Mountpoint, "m", "", "mount point")
  fs.StringVar(&opts.VolName, "n", "", "volume name (OSX only)")
  fs.BoolVar(&opts.Scratch, "t", false, "create scratch directory as mount point")
  fs.BoolVar(&opts.Verbose, "v", false, "log at INFO level")
  fs.BoolVar(&opts.Debug, "d", false, "log at DEBUG level")
  fs.Usage = func () {
    Echoerr("Usage: %s <options> [path path ...]\n", os.Args[0])
    fs.PrintDefaults()
  }
  fs.Parse(os.Args[1:])
  opts.Args = fs.Args()
  return opts
}
