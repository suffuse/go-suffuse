package suffuse

import (
  "flag"
)

type SfsOpts struct {
  Config, Mountpoint, ExcludeRegex, VolName string
  Scratch, Verbose, Debug bool
  Paths []string

  // stored here because it contains references to arguments used to create SfsOpts
  usage func()

  // stored here so we can generate potential errors during construction
  mountpointPath Path
}

func ParseSfsOpts(args []string) (*SfsOpts, *SfsOptsError) {

  opts, err := optsFromArgs(args)
  if err != nil { return nil, err }

  err = addMountpointPath(opts)
  if err != nil { return nil, err }

  err = validate(opts)

  return opts, err
}

func (opts *SfsOpts) GetMountpoint() Path { return opts.mountpointPath }

func optsFromArgs(args []string) (*SfsOpts, *SfsOptsError) {
  fs := flag.NewFlagSet("suffuse", flag.ContinueOnError)

  fs.Usage = func () {
    Echoerr("Usage: %s <options> [path path ...]\n", args[0])
    fs.PrintDefaults()
  }

  opts := &SfsOpts{}
  opts.usage = fs.Usage

  // fs.StringVar(&opts.Config, "c", "", "suffuse config file")
  // fs.StringVar(&opts.ExcludeRegex, "x", "", "name exclusion regex")
  fs.StringVar(&opts.Mountpoint, "m", "", "mount point")
  fs.StringVar(&opts.VolName, "n", "", "volume name (OSX only)")
  fs.BoolVar(&opts.Scratch, "t", false, "create scratch directory as mount point")
  fs.BoolVar(&opts.Verbose, "v", false, "log at INFO level")
  fs.BoolVar(&opts.Debug, "d", false, "log at DEBUG level")

  err := fs.Parse(args[1:])
  if err != nil {
    // Parse has already printed the problem and usage, there is no way to tell it not to do that
    return nil, &SfsOptsError{err: err.Error(), usage: nil}
  }

  opts.Paths = fs.Args()
  return opts, nil
}

func validate(opts *SfsOpts) *SfsOptsError {
  // Exactly one incoming path for the moment.
  // Eventually more than one could mean a union mount.
  if len(opts.Paths) != 1 { return opts.newError("Suffuse currently accepts a single path") }

  mnt := opts.mountpointPath
  if (!mnt.FileExists()) { return opts.newError(Sprintf("Mountpoint '%s' does not exist", mnt)) }

  for _, path := range opts.Paths {
    if (!NewPath(path).FileExists()) { return opts.newError(Sprintf("Path '%s' does not exist", path)) }
  }

  return nil
}

func (opts *SfsOpts) newError(msg string) *SfsOptsError {
  return &SfsOptsError{msg, opts.usage}
}

func addMountpointPath(opts *SfsOpts) *SfsOptsError {

  determineMountPoint := func() (Path, *SfsOptsError) {
    switch {
      case opts.Scratch          : return ScratchDir(), nil
      case opts.Mountpoint != "" : return NewPath(opts.Mountpoint), nil
      default                    : return Path{}, opts.newError("Mountpoint can not be empty if the scratch flag is false")
    }
  }

  path, err := determineMountPoint()

  opts.mountpointPath = path

  return err
}

type SfsOptsError struct {
  err string
  usage func()
}

func (e *SfsOptsError) Error() string { return e.err }
func (e *SfsOptsError) PrintUsage()   {
  if e.usage != nil {
    Echoerr("%s\n", e.err)
    e.usage()
  }
}
