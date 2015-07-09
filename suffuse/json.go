package suffuse

import (
  "encoding/json"
  "io/ioutil"
)

func ReadJsonFile(p Path) map[string]interface{} {
  bs  := p.SlurpBytes()
  var f map[string]interface{}
  json.Unmarshal(bs, &f)

  return f
}

func WriteJsonFile(p Path, x interface{}) {
  data, err := json.Marshal(x)
  if err == nil {
    ioutil.WriteFile(string(p), data, 0444)
  }
}

func JsonPretty(x interface{}) string {
  return MaybeByteString(json.MarshalIndent(x, "", "    "))
}
