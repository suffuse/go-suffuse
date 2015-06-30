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
    ioutil.WriteFile(p.Path, data, 0444)
  }
}
