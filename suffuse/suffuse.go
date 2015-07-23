package suffuse

import (
  "os"
)

const (
  defaultFilePerms  = os.FileMode(0644)
)

/** Eventually we'll read this kind of information out of a config
 *  file, but for now it's hardcoded here.
 */
var rules = append(
  []Rule { &IdRule{} },  // The identity fs - passes underlying paths through unchanged
  parseRules()...
)

func parseRules() []Rule {
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
