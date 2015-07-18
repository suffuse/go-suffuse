package suffuse

import (
  "errors"
  "reflect"
)

/** If you have a *Foo type which implements error and you pass it
 *  to a method accepting error, when you pass a nil pointer it
 *  does NOT test as nil. It's apparently a non-nil interface with
 *  a nil value. That's what the second test is for.
 */
func IsNilError(e error) bool {
  return e == nil || reflect.ValueOf(e).IsNil()
}
func IsError(e error) bool {
  return !IsNilError(e)
}
func NoError(_ interface{}, err error) bool {
  return IsNilError(err)
}

func TypedString(x interface{}) string {
  return Sprintf("%v %T", x, x)
}

func MaybeFatal(e error) { if IsError(e) { sfsLogger.Fatal(e) } }
func MaybePanic(e error) { if IsError(e) { sfsLogger.Panic(e) } }
func MaybeLog(e error)   { if IsError(e) { sfsLogger.Println(TypedString(e)) } }

func AssertEq(x interface{}, y interface{}) {
  if (x != y) {
    panic(Sprintf("%v != %v", TypedString(x), TypedString(y)))
  }
}

/** Returns the first non-nil error of all passed,
 *  or nil if they're all nil.
 */
func FindError(errors ...error) error {
  for _, e := range errors {
    if IsError(e) { return e }
  }
  return nil
}

func NewErr(text string) error { return errors.New(text) }

func maybeByteString(result []byte, err error) string {
  if err != nil { return "" }
  return string(result)
}
func maybeBytes(result []byte, err error) []byte {
  if err != nil { return nil }
  return result
}
func maybeString(x interface{}) string {
  switch x := x.(type) {
    case string: return x
    default: return ""
  }
}
func maybeBool(x interface{}) bool {
  switch x := x.(type) {
    case bool: return x
    default: return false
  }
}
func maybeInt(x interface{}) int {
  switch x := x.(type) {
    case int: return x
    default: return 0
  }
}
