package suffuse

import (
  "fmt"
  "os"
  "io/ioutil"
  "log"
  "errors"
)

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
