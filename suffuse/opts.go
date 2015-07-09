package suffuse

import (
  "flag"
  lg "gopkg.in/inconshreveable/log15.v2"
)

type sfsOpts struct {
  printUsage func()
  Config, Mountpoint, VolName string
  Scratch, Verbose, Debug bool
  Args []string
}

type SfsConfig struct {
  VolName string
  Config Path
  Mountpoint Path
  LogLevel lg.Lvl
  Paths []Path
}

func CreateSfsConfig(args []string) (*SfsConfig, *SfsConfigError) {

  opts, err := optsFromArgs(args)
  if err != nil { return nil, err }

  err = validate(opts)
  if err != nil { return nil, err }

  return configFromOpts(opts)
}

func optsFromArgs(args []string) (*sfsOpts, *SfsConfigError) {
  fs := flag.NewFlagSet("suffuse", flag.ContinueOnError)

  fs.Usage = func () {
    Echoerr("Usage: %s <options> [path path ...]\n", args[0])
    fs.PrintDefaults()
  }

  opts := &sfsOpts{}
  opts.printUsage = fs.Usage

  fs.StringVar(&opts.Config, "c", "", "suffuse config file")
  fs.StringVar(&opts.Mountpoint, "m", "", "mount point")
  fs.StringVar(&opts.VolName, "n", "", "volume name (OSX only)")
  fs.BoolVar(&opts.Scratch, "t", false, "create scratch directory as mount point")
  fs.BoolVar(&opts.Verbose, "v", false, "log at INFO level")
  fs.BoolVar(&opts.Debug, "d", false, "log at DEBUG level")

  err := fs.Parse(args[1:])
  if err != nil {
    // Parse has already printed the problem and usage, there is no way to tell it not to do that
    return nil, &SfsConfigError{err: err.Error(), usage: nil}
  }

  opts.Args = fs.Args()
  return opts, nil
}

func validate(opts *sfsOpts) *SfsConfigError {
  exists := func (path string, msg string) *SfsConfigError {
    if (!Path(path).FileExists()) { return opts.newError(Sprintf(msg, path)) }
    return nil
  }

  // Exactly one incoming path for the moment.
  // Eventually more than one could mean a union mount.
  if len(opts.Args) != 1 { return opts.newError("Suffuse currently accepts a single path") }

  if err := exists(opts.Mountpoint, "Mountpoint '%s' does not exist" ); err != nil { return err }
  if err := exists(opts.Config    , "Config file '%s' does not exist"); err != nil && opts.Config != "" { return err }

  for _, path := range opts.Args {
    if err := exists(path, "Path '%s' does not exist"); err != nil { return err }
  }

  return nil
}

func configFromOpts(opts *sfsOpts) (*SfsConfig, *SfsConfigError) {

  determineMountPoint := func() (Path, *SfsConfigError) {
    switch {
      case opts.Scratch          : return ScratchDir(), nil
      case opts.Mountpoint != "" : return Path(opts.Mountpoint), nil
      default                    : return NoPath, opts.newError("Mountpoint can not be empty if the scratch flag is false")
    }
  }

  determineLogLevel := func() lg.Lvl {
    switch {
      case opts.Debug   : return lg.LvlDebug
      case opts.Verbose : return lg.LvlInfo
      default           : return lg.LvlWarn
    }
  }

  mnt, err := determineMountPoint()
  if err != nil { return nil, err }

  return &SfsConfig {
    VolName    : opts.VolName,
    Config     : Path(opts.Config),
    Mountpoint : mnt,
    Paths      : Paths(opts.Args...),
    LogLevel   : determineLogLevel(),
  }, nil
}

func (opts *sfsOpts) newError(msg string) *SfsConfigError {
  return &SfsConfigError{msg, opts.printUsage}
}

type SfsConfigError struct {
  err string
  usage func()
}

func (e *SfsConfigError) Error() string { return e.err }
func (e *SfsConfigError) PrintUsage()   {
  if e.usage != nil {
    Echoerr("%s\n", e.err)
    e.usage()
  }
}
