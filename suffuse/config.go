package suffuse

import (
  "encoding/json"
)

var RuleFactories = map[string]func() Rule {
  "file_command"   : func() Rule { return &FileCommand{}    },
  "file_conversion": func() Rule { return &FileConversion{} },
}

type ConfigEntry struct {
  Nature string
  Config json.RawMessage
}

func CreateRules(config []byte) ([]Rule, error) {
  var entries []ConfigEntry
  err := json.Unmarshal(config, &entries)
  if err != nil { return nil, err }

  return entriesToRules(entries)
}

func entriesToRules(entries []ConfigEntry) ([]Rule, error) {
  results := make([]Rule, len(entries))

  for i, entry := range entries {
    val := RuleFactories[entry.Nature]()
    err := json.Unmarshal(entry.Config, val)
    if err != nil { return nil, err }
    results[i] = val
  }

  return results, nil
}
