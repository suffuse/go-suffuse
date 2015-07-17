package suffuse

import (
  "fmt"
  "regexp"
  "strings"
  "sort"
)

type StringMap map[string]string
type Lines struct {
  Strings []string
}
type Regex struct {
  *regexp.Regexp
}

var whitespaceRegex = Regexp(`\s+`)

func BytesToLines(bytes []byte) Lines                { return SplitLines(TrimSpace(string(bytes)))             }
func Echoerr(format string, a ...interface{})        { fmt.Fprintf(OsStderr(), format + "\n", a...)            }
func Printf(format string, a ...interface{})         { fmt.Printf(format, a...)                                }
func Printfln(format string, a ...interface{})       { fmt.Printf(format + "\n", a...)                         }
func Println(a ...interface{})                       { fmt.Println(a...)                                       }
func Regexp(re string) Regex                         { return Regex{ regexp.MustCompile(re) }                  }
func SplitLines(s string) Lines                      { return Strs(strings.Split(s, "\n")...)                  }
func SplitWords(s string) Lines                      { return Strs(whitespaceRegex.Split(TrimSpace(s), -1)...) }
func Sprintf(format string, a ...interface{}) string { return fmt.Sprintf(format, a...)                        }
func Strings(xs ...string) []string                  { return xs                                               }
func Strs(xs ...string) Lines                        { return Lines{ xs }                                      }

func Names(names ...string)[]Name {
  xs := make([]Name, len(names))
  for i := range names { xs[i] = Name(names[i]) }
  return xs
}

func TrimSpace(s string)string { return strings.TrimSpace(s) }

func (x Lines) Join(sep string) string { return strings.Join(x.Strings, sep) }
func (x Lines) JoinWords() string      { return x.TrimAll().Join(" ")        }
func (x Lines) JoinLines() string      { return x.Join("\n")                 }
func (x Lines) Len() int               { return len(x.Strings)               }
func (x Lines) String() string         { return x.JoinLines()                }

func (x Lines) SortedStrings() []string {
  dst := make([]string, x.Len())
  copy(dst, x.Strings)
  sort.Strings(dst)
  return dst
}
func (x Lines) TrimAll() Lines {
  return x.Map(func(s string)string { return TrimSpace(s) })
}

func (x Lines) Map(f func(string)string) Lines {
  res := make([]string, x.Len())
  for i, x := range x.Strings { res[i] = f(x) }
  return Strs(res...)
}
func (x Lines) FlatMap(f func(string)[]string) Lines {
  var res []string
  for _, x := range x.Strings { res = append(res, f(x)...) }
  return Strs(res...)
}
func (x Lines) Fold(f func(string, string)string) string {
  res := ""
  for _, x := range x.Strings { res = f(res, x) }
  return res
}

/** String manipulation, should go on a string type but doesn't yet.
 */
func StripMargin(marginChar rune, s string) string {
  perLine := func(line string)string {
    trimmed := TrimSpace(line)
    index   := strings.IndexRune(trimmed, marginChar)
    if index > -1 {
      return trimmed[index + 1:]
    }
    return line
  }
  return SplitLines(TrimSpace(s)).Map(perLine).String()
}
