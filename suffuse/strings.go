package suffuse

import (
  "fmt"
  "regexp"
  "strings"
)

type StringMap map[string]string
type Strings struct {
  strings []string
}
type Regex struct {
  *regexp.Regexp
}

var whitespaceRegex = Regexp(`\s+`)

func BytesToLines(bytes []byte) Strings { return SplitLines(TrimSpace(string(bytes)))             }
func SplitLines(s string) Strings       { return Strs(strings.Split(s, "\n")...)                  }
func SplitWords(s string) Strings       { return Strs(whitespaceRegex.Split(TrimSpace(s), -1)...) }
func Strs(xs ...string) Strings         { return Strings{ xs }                                    }
func Regexp(re string) Regex            { return Regex{ regexp.MustCompile(re) }                  }

func Echoerr(format string, a ...interface{})        { fmt.Fprintf(OsStderr(), format + "\n", a...) }
func Printf(format string, a ...interface{})         { fmt.Printf(format, a...)                     }
func Printfln(format string, a ...interface{})       { fmt.Printf(format + "\n", a...)              }
func Println(a ...interface{})                       { fmt.Println(a...)                            }
func Sprintf(format string, a ...interface{}) string { return fmt.Sprintf(format, a...)             }

func Names(names ...string)[]Name {
  xs := make([]Name, len(names))
  for i := range names { xs[i] = Name(names[i]) }
  return xs
}

func (x Strings) Join(sep string) string { return strings.Join(x.strings, sep) }
func (x Strings) JoinWords() string      { return x.TrimAll().Join(" ")        }
func (x Strings) Len() int               { return len(x.strings)               }
func (x Strings) String() string         { return x.Join("\n")                 }

func TrimSpace(s string)string { return strings.TrimSpace(s) }
func (x Strings) TrimAll() Strings {
  return x.Map(func(s string)string { return TrimSpace(s) })
}

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

/** String manipulation, should go on a string type but doesn't yet.
 */
func StripMargin(marginChar rune, s string) string {
  perLine := func(line string)string {
    trimmed := TrimSpace(line)
    index   := strings.IndexRune(trimmed, marginChar)
    if index > -1 {
      return trimmed[index + 1:]
    } else {
      return line
    }
  }
  return SplitLines(TrimSpace(s)).Map(perLine).String()
}
