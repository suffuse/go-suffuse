package suffuse

import (
  "encoding/json"
)

var ruleFactories = map[string]func() Rule {
  "file_command"   : func() Rule { return &FileCommand{}    },
  "file_conversion": func() Rule { return &FileConversion{} },
}

type configEntry struct {
  Nature string
  Config json.RawMessage
}

func CreateRules(config []byte) ([]Rule, error) {
  var entries []configEntry
  err := json.Unmarshal(config, &entries)
  if err != nil { return nil, err }

  return entriesToRules(entries)
}

func entriesToRules(entries []configEntry) ([]Rule, error) {
  results := make([]Rule, len(entries))

  for i, entry := range entries {
    val := ruleFactories[entry.Nature]()
    err := json.Unmarshal(entry.Config, val)
    if err != nil { return nil, err }
    results[i] = val
  }

  return results, nil
}
