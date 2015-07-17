package suffuse

/** Packages confined to this file: regexp
 */

import (
  "strings"
)

type Lines struct {
  Strings []string
}

func NewLines(lines ...string) Lines { return Lines { lines } }

func BytesToLines(bytes []byte) Lines {
  return SplitLines(strings.TrimSpace(string(bytes)))
}

func (x Lines) Join(sep string) string   { return strings.Join(x.Strings, sep) }
func (x Lines) JoinWords() string        { return x.TrimAll().Join(" ")        }
func (x Lines) Len() int                 { return len(x.Strings)               }
func (x Lines) String() string           { return x.Join("\n")                 }

func (x Lines) TrimAll() Lines { return x.Map(func(s string)string { return strings.TrimSpace(s) }) }

func (x Lines) Map(f func(string)string) Lines {
  res := make([]string, x.Len())
  for i, x := range x.Strings { res[i] = f(x) }
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

func SplitLines(s string) Lines { return NewLines(strings.Split(s, "\n")...) }

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
