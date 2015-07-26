package suffuse

import (
  "os"
)

const (
  defaultFilePerms  = os.FileMode(0644)
)

var defaultRules = []Rule {
  &IdRule{}, // The identity fs - passes underlying paths through unchanged
}
