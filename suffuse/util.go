package suffuse

import (
  "fmt"
  "os"
  "io"
  "io/ioutil"
  "bufio"
  "log"
  "time"
  "strconv"
  "errors"
  "runtime"
)

// TODO - actually delete on exit.
var deleteOnExit = make([]Path, 0)

func Println(a ...interface{})                       { fmt.Println(a...)                           }
func Printf(format string, a ...interface{})         { fmt.Printf(format, a...)                    }
func Printfln(format string, a ...interface{})       { fmt.Printf(format + "\n", a...)             }
func Sprintf(format string, a ...interface{}) string { return fmt.Sprintf(format, a...)            }
func Echoerr(format string, a ...interface{})        { fmt.Fprintf(os.Stderr, format + "\n", a...) }

func MaybeFatal(e error) { if e != nil { log.Fatal(e) } }
func MaybePanic(e error) { if e != nil { panic(e) } }
func MaybeLog(e error) { if e != nil { Echoerr("ERROR: %v", e) } }
func AssertEq(x interface{}, y interface{}) {
  if (x != y) {
    panic(fmt.Sprintf("%v: %T != %v: %T", x, x, y, y))
  }
}

/** Returns the first non-nil error of all passed,
 *  or nil if they're all nil.
 */
func FindError(errors ...error) error {
  for _, e := range errors {
    if e != nil { return e }
  }
  return nil
}

func NewErr(text string) error { return errors.New(text) }

func GetStack() string {
  buf := make([]byte, 16384)
  bytes := runtime.Stack(buf, true)
  return string(buf[0:bytes])
}
func SplitNanoTime(t time.Time) (sec uint64, nsec uint32) {
  nano := t.UnixNano()
  sec = uint64(nano / 1e9)
  nsec = uint32(nano % 1e9)
  return
}
func ForeachLine(r io.ReadCloser, f func(string)) {
  scanner := bufio.NewScanner(r)
  for scanner.Scan() { f(scanner.Text()) }
}
func StrTake(s string, n int) string {
  if len(s) <= n { return s } else { return s[0:n] }
}
func TruncToAndLeftJust(s string, width int64) string {
  format := "%-" + strconv.FormatInt(width, 10) + "s"
  return fmt.Sprintf(format, TruncTo(s, width))
}
func TruncTo(s string, max int64) string {
  strlen := int64(len(s))
  if (strlen < max) {
    return s;
  } else {
    return s[0:max-3] + "..."
  }
}

/** TODO: better types than "int" for uid/gid etc.
 */
func Getegid() int              { return os.Getegid()     }
func Geteuid() int              { return os.Geteuid()     }
func Getgid() int               { return os.Getgid()      }
func Getgroups() ([]int, error) { return os.Getgroups()   }
func Getpagesize() int          { return os.Getpagesize() }
func Getpid() int               { return os.Getpid()      }
func Getppid() int              { return os.Getppid()     }
func Getuid() int               { return os.Getuid()      }
func Hostname() (string, error) { return os.Hostname()    }

func ScratchDir() Path {
  dir, err := ioutil.TempDir("", "suffuse")
  MaybeFatal(err)
  path := NewPath(dir)
  deleteOnExit = append(deleteOnExit, path)
  return path
}
func ScratchFile() Path {
  fh, err := ioutil.TempFile("", "suffuse")
  MaybeFatal(err)
  path := NewPath(fh.Name())
  deleteOnExit = append(deleteOnExit, path)
  return path
}
