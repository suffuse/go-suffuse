package suffuse

/** Packages confined to this file: regexp
 */

import (
  "regexp"
  "strings"
)

/** Without generics and with barely usable HOFs we're in hell.
 *  Here we focus on Lines as a collection of Strings and reimplement
 *  swaths of better languages in terms of exactly 1/Nth of the
 *  types of interest.
 */
type Regex struct {
  *regexp.Regexp
}
type Lines struct {
  Strings []string
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

func NewLines(lines ...string) Lines { return Lines { lines                  } }
func NewRegex(re string) Regex       { return Regex { regexp.MustCompile(re) } }

func SplitLines(s string) Lines     { return NewLines(strings.Split(s, "\n")...) }
func Strings(xs ...string) []string { return xs                                  }
func ByteBuf(size int) []byte       { return make([]byte, size)                  }

func (x Lines) Filter(re Regex) Lines    { return filterCommon(x, re, true)    }
func (x Lines) FilterNot(re Regex) Lines { return filterCommon(x, re, false)   }
func (x Lines) IsEmpty() bool            { return x.Len() == 0                 }
func (x Lines) Join(sep string) string   { return strings.Join(x.Strings, sep) }
func (x Lines) JoinWords() string        { return x.TrimAll().Join(" ")        }
func (x Lines) Len() int                 { return len(x.Strings)               }
func (x Lines) String() string           { return x.Join("\n")                 }

func (x Lines) TrimAll() Lines { return x.Map(func(s string)string { return strings.TrimSpace(s) }) }

func filterCommon(x Lines, re Regex, expectTrue bool) Lines {
  xs := x.Strings
  ys := make([]string, 0)

  for _, line := range xs {
    if re.MatchString(line) == expectTrue {
      ys = append(ys, line)
    }
  }
  return NewLines(ys...)
}
func (x Lines) Map(f func(string)string) Lines {
  res := make([]string, x.Len())
  for i, x := range x.Strings { res[i] = f(x) }
  return NewLines(res...)
}
func (x Lines) MapIndex(f func(int, string)string) Lines {
  res := make([]string, x.Len())
  for i, x := range x.Strings { res[i] = f(i, x) }
  return NewLines(res...)
}
func (x Lines) FlatMap(f func(string)[]string) Lines {
  res := make([]string, 0)
  for _, x := range x.Strings { res = append(res, f(x)...) }
  return NewLines(res...)
}
func (x Lines) Fold(f func(string, string)string) string {
  res := ""
  for _, x := range x.Strings { res = f(res, x) }
  return res
}
func BytesToLines(bytes []byte) Lines {
  return SplitLines(strings.TrimSpace(string(bytes)))
}

/** String manipulation, should go on a string type but doesn't yet.
 */
func StripMargin(marginChar rune, s string) string {
  perLine := func(line string)string {
    trimmed := strings.TrimSpace(line)
    if strings.IndexRune(trimmed, marginChar) == 0 {
      return trimmed[1:]
    } else {
      return line
    }
  }
  return SplitLines(strings.TrimSpace(s)).Map(perLine).String()
}
