package suffuse

import (
  "fmt"
  "strings"
)

type StringMap map[string]string
type Strings []string

func BytesToStrings(bytes []byte) Strings { return SplitLines(TrimSpace(string(bytes)))             }
func SplitLines(s string) Strings         { return Strings(strings.Split(s, "\n"))                  }
func SplitWords(s string) Strings         { return Strings(whitespaceRegex.Split(TrimSpace(s), -1)) }

func Echoerr(format string, a ...interface{})        { fmt.Fprintf(OsStderr(), format + "\n", a...) }
func Printf(format string, a ...interface{})         { fmt.Printf(format, a...)                     }
func Printfln(format string, a ...interface{})       { fmt.Printf(format + "\n", a...)              }
func Println(a ...interface{})                       { fmt.Println(a...)                            }
func Sprintf(format string, a ...interface{}) string { return fmt.Sprintf(format, a...)             }

func (x Strings) Array() []string        { return []string(x)                  }
func (x Strings) Join(sep string) string { return strings.Join(x.Array(), sep) }
func (x Strings) JoinWords() string      { return x.TrimAll().Join(" ")        }
func (x Strings) Len() int               { return len(x.Array())               }
func (x Strings) String() string         { return x.Join("\n")                 }

func TrimSpace(s string)string { return strings.TrimSpace(s) }
func (x Strings) TrimAll() Strings {
  return x.Map(func(s string)string { return TrimSpace(s) })
}

func (x Strings) Map(f func(string)string) Strings {
  res := make([]string, x.Len())
  for i, x := range x.Array() { res[i] = f(x) }
  return Strings(res)
}
func (x Strings) FlatMap(f func(string)[]string) Strings {
  res := make([]string, 0)
  for _, x := range x.Array() { res = append(res, f(x)...) }
  return Strings(res)
}
func (x Strings) Fold(f func(string, string)string) string {
  res := ""
  for _, x := range x.Array() { res = f(res, x) }
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
