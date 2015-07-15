package suffuse

import (
  "encoding/json"
)

func readJsonFile(p Path) map[string]interface{} {
  bs  := p.SlurpBytes()
  var f map[string]interface{}
  json.Unmarshal(bs, &f)

  return f
}

func JsonPretty(x interface{}) string {
  return maybeByteString(json.MarshalIndent(x, "", "    "))
}
