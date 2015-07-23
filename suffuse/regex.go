package suffuse

import (
  "regexp"
)

type Regex struct {
  *regexp.Regexp
}

var whitespaceRegex = Regexp(`\s+`)

func Regexp(re string) Regex { return Regex{ regexp.MustCompile(re) } }
