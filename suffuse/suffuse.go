package suffuse

import (
  "os"
)

const (
  SfsConfigFileName = ".sfs"
  DefaultFilePerms  = os.FileMode(0644)
  DefaultDirPerms   = os.FileMode(0755)
)

/** Eventually we'll read this kind of information out of a config
 *  file, but for now it's hardcoded here.
 */
var Rules = []Rule {
  IdRule{},  // The identity fs - passes underlying paths through unchanged
  SedRule{}, // The sed fs - parses a '#...' suffix and applies it to the prefix path
}

var NoNode = NewIdNode(NoPath)

var NoPath     = Path("")
var RootPath   = Path("/")
var DotPath    = Path(".")
var DotDotPath = Path("..")

// TODO - actually delete on exit.
var deleteOnExit = make([]Path, 0)
