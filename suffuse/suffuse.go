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
var Rules = append(
  []Rule { &IdRule{} },  // The identity fs - passes underlying paths through unchanged
  rules()...
)

func rules() []Rule {
  conf := []byte(`[
    {
      "nature" : "file_command",
      "config" : {
        "separator" : "#",
        "command"   : "sed -ne $args $file",
        "example"   : "seq.txt#2,3p"
      }
    },
    {
      "nature" : "file_conversion",
      "config" : {
        "from"      : "json",
        "to"        : "yaml",
        "command"   : "json2yaml $file",
        "example"   : "jsonFile.yaml"
      }
    }
  ]`)

  rules, err := CreateRules(conf)
  if err != nil { Echoerr(err.Error()) }

  return rules
}

var NoNode = NewIdNode(NoPath)

var NoPath     = Path("")
var RootPath   = Path("/")
var DotPath    = Path(".")
var DotDotPath = Path("..")

// TODO - actually delete on exit.
var deleteOnExit = make([]Path, 0)
