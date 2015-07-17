package suffuse

import (
  "github.com/docopt/docopt-go"
)

func sfsUsage()string {
  return TrimSpace(`
Usage: suffuse [-v | -vv | -vvv] [-t | -m PATH] [-n VOLNAME] [-c PATH] PATH

Arguments:
  PATH  existing filesystem

Options:
  -c PATH     suffuse config file
  -m PATH     mount point
  -n VOLNAME  volume name (OS X only)
  -t          create scratch directory as mount point
  -v          more verbose logging (also -vv and -vvv)
  `)
}

type sfsOpts struct {
  Config, Mountpoint, VolName string
  Scratch bool
  Verbosity int
  Dir string
}

type SfsConfig struct {
  VolName string
  Config Path
  Mountpoint Path
  LogLevel LogLevel
  Paths []Path
}
type SfsConfigError struct {
  string
}

func CreateSfsConfig(argv []string) *SfsConfig {
  trace("argv: %v", argv)
  opts := optsFromArgs(argv[1:])
  MaybeFatal(validate(opts))
  return configFromOpts(opts)
}

func optsFromArgs(args []string) *sfsOpts {
  argMap, err := docopt.Parse(sfsUsage(), args, /*help*/true, /*version*/"HEAD", /*optionsFirst*/false, /*exit*/false)
  MaybeLog(err)
  trace("argMap: %+v", argMap)

  return &sfsOpts {
    Config: maybeString(argMap["-c"]),
    Mountpoint: maybeString(argMap["-m"]),
    VolName: maybeString(argMap["-n"]),
    Scratch: maybeBool(argMap["-t"]),
    Verbosity: maybeInt(argMap["-v"]),
    Dir: maybeString(argMap["PATH"]),
  }
}

func validate(opts *sfsOpts) *SfsConfigError {
  var err *SfsConfigError

  /** If there's a non-empty path argument to these options,
   *  it should already exist or it's an error.
   */
  ensureExists := func (path string, what string) {
    if path != "" && !Path(path).FileExists() {
      err = opts.newError(Sprintf("%s '%s' does not exist", what, path))
    }
  }
  ensureExists(opts.Mountpoint, "Mountpoint")
  ensureExists(opts.Dir, "Directory")
  ensureExists(opts.Config, "Config file")
  return err
}

func configFromOpts(opts *sfsOpts) *SfsConfig {
  var mountPoint Path
  var level LogLevel

  if opts.Scratch {
    mountPoint = ScratchDir()
  } else {
    mountPoint = Path(opts.Mountpoint)
  }

  switch opts.Verbosity {
    case 0: level = LogWarn
    case 1: level = LogInfo
    case 2: level = LogDebug
    case 3: level = LogTrace
  }

  return &SfsConfig {
    VolName    : opts.VolName,
    Config     : Path(opts.Config),
    Mountpoint : mountPoint,
    Paths      : Paths(opts.Dir),
    LogLevel   : level,
  }
}

func (opts *sfsOpts) newError(msg string) *SfsConfigError { return &SfsConfigError{ msg } }
func (e *SfsConfigError) Error() string { return e.string }
