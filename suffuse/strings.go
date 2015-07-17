package suffuse

import (
  "strings"
)

type Strings struct {
  strings []string
}

func Strs(lines ...string) Strings { return Strings { lines } }

func BytesToLines(bytes []byte) Strings {
  return SplitLines(strings.TrimSpace(string(bytes)))
}

func (x Strings) Join(sep string) string { return strings.Join(x.strings, sep) }
func (x Strings) JoinWords() string      { return x.TrimAll().Join(" ")        }
func (x Strings) Len() int               { return len(x.strings)               }
func (x Strings) String() string         { return x.Join("\n")                 }

func (x Strings) TrimAll() Strings { return x.Map(func(s string)string { return strings.TrimSpace(s) }) }

func (x Strings) Map(f func(string)string) Strings {
  res := make([]string, x.Len())
  for i, x := range x.strings { res[i] = f(x) }
  return Strs(res...)
}
func (x Strings) FlatMap(f func(string)[]string) Strings {
  res := make([]string, 0)
  for _, x := range x.strings { res = append(res, f(x)...) }
  return Strs(res...)
}
func (x Strings) Fold(f func(string, string)string) string {
  res := ""
  for _, x := range x.strings { res = f(res, x) }
  return res
}

// TODO - strings methods to expose on a more convenient String type
//
// Contains(s, substr string) bool
// Index(s, sep string) int
// LastIndex(s, sep string) int
// Repeat(s string, count int) string
// Replace(s, old, new string, n int) string
// Split(s, sep string) []string
// ToLower(s string) string
// ToTitle(s string) string
// ToUpper(s string) string
// Trim(s string, cutset string) string
// TrimPrefix(s, prefix string) string
// TrimSuffix(s, suffix string) string

func SplitLines(s string) Strings { return Strs(strings.Split(s, "\n")...) }

/** String manipulation, should go on a string type but doesn't yet.
 */
func StripMargin(marginChar rune, s string) string {
  perLine := func(line string)string {
    trimmed := strings.TrimSpace(line)
    index   := strings.IndexRune(trimmed, marginChar)
    if index > -1 {
      return trimmed[index + 1:]
    } else {
      return line
    }
  }
  return SplitLines(strings.TrimSpace(s)).Map(perLine).String()
}
