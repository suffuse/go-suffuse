package suffuse

import (
  "fmt"
  "runtime"
  "unsafe"
  r "reflect"
)

type CallerInfo struct {
  File Path
  Name string
  Line int
}

// Regex to extract just the function name (and not the module path)
var extractFnName = NewRegex(`^.*\.(.*)$`)


type AnyI interface {         }
type AnyS struct    { AnyI    }
type Type struct    { r.Type  }
type Value struct   { r.Value }

func Any(x AnyI) AnyS    { return AnyS { x } }
func What(x AnyI) string { return fmt.Sprintf("%v: %T", x, x) }

func NewValue(x r.Value) Value { return Value { x } }
func NewType(x r.Type) Type    { return Type  { x } }

func (x AnyS) AnyValue() Value { return NewValue(r.ValueOf(x.AnyI))            }
func (x AnyS) AnyType() Type   { return NewType(r.TypeOf(x.AnyI))              }
func (x AnyS) What() string    { return fmt.Sprintf("%#v: %T", x.AnyI, x.AnyI) }
func (x AnyS) Dump()           { Println(What(x.AnyI))                         }

func (x AnyS) Contains(y interface{})bool { return AnyContains(x, y) }

func (x Type) ChanType() Type      { return NewType(r.ChanOf(r.BothDir, x.Type)) }
func (x Type) SendType() Type      { return NewType(r.ChanOf(r.SendDir, x.Type)) }
func (x Type) RecvType() Type      { return NewType(r.ChanOf(r.RecvDir, x.Type)) }
func (x Type) SliceType() Type     { return NewType(r.SliceOf(x.Type))           }
func (x Type) MapType(v Type) Type { return NewType(r.MapOf(x.Type, v.Type))     }
func (x Type) PtrType() Type       { return NewType(r.PtrTo(x.Type))             }
func (x Value) ValueType() Type    { return NewType(x.Value.Type())              }

func (x Type) MakeChan(buffer int) Value                    { return NewValue(r.MakeChan(x.Type, buffer))    }
func (x Type) MakeFunc(f func([]r.Value) ([]r.Value)) Value { return NewValue(r.MakeFunc(x.Type, f))         }
func (x Type) MakeMap() Value                               { return NewValue(r.MakeMap(x.Type))             }
func (x Type) MakeSlice(len, cap int) Value                 { return NewValue(r.MakeSlice(x.Type, len, cap)) }
func (x Type) New() Value                                   { return NewValue(r.New(x.Type))                 }
func (x Type) NewAt(p unsafe.Pointer) Value                 { return NewValue(r.NewAt(x.Type, p))            }
func (x Value) Get() Value                                  { return NewValue(x.Value.Elem())                }
func (x Value) Addr() Value                                 { return NewValue(x.Value.Addr())                }
func (x Value) Deref() Value                                { return NewValue(r.Indirect(x.Value))           }

func AnyEquals(x, y interface{}) bool   { return equalValues(r.ValueOf(x), r.ValueOf(y))   }
func AnyContains(x, y interface{}) bool { return valueContains(r.ValueOf(x), r.ValueOf(y)) }

// TODO: String as container of runes/bytes.
func valueContains(xs, x r.Value) bool {
  k := xs.Kind()

  if k == r.Array || k == r.Slice {
    for i := 0 ; i < xs.Len() ; i++ {
      if equalValues(xs.Index(i), x) {
        return true
      }
    }
    return false
  } else if k == r.Map {
    keys := xs.MapKeys()
    for i := 0 ; i < len(keys) ; i++ {
      if equalValues(keys[i], x) {
        return true
      }
    }
    return false
  } else {
    return false
  }
}

func equalSlices(x, y r.Value)bool {
  n := x.Len()
  if n != y.Len() { return false }

  for i := 0 ; i < n ; i++ {
    if !equalValues(x.Index(i), y.Index(i)) {
      return false
    }
  }
  return true
}

func equalStructs(x, y r.Value)bool {
  if x.Type() != y.Type() { return false }
  n := x.NumField()
  if n != y.NumField() { return false }

  for i := 0 ; i < n ; i++ {
    if !equalValues(x.Field(i), y.Field(i)) {
      return false
    }
  }
  return true
}

func equalMaps(x, y r.Value)bool {
  keys := x.MapKeys()
  if len(keys) != len(y.MapKeys()) { return false }

  for _, key := range keys {
    if !equalValues(x.MapIndex(key), y.MapIndex(key)) {
      return false
    }
  }
  return true
}

func widenKind(k r.Kind) r.Kind {
  switch k {
    case r.Int, r.Int8, r.Int16, r.Int32, r.Int64      : return r.Int64
    case r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64 : return r.Uint64
    case r.Float32, r.Float64                          : return r.Float64
    case r.Complex64, r.Complex128                     : return r.Complex128
    case r.Array, r.Slice                              : return r.Slice
    default                                            : return k
  }
}

func equalValues(x, y r.Value) bool {
  k1 := widenKind(x.Kind())
  k2 := widenKind(y.Kind())
  if k1 != k2 { return false }

  switch k1 {
    case r.Invalid, r.Chan, r.Func : return false
    case r.Bool                    : return x.Bool() == y.Bool()
    case r.Int64                   : return x.Int() == y.Int()
    case r.Uint64                  : return x.Uint() == y.Uint()
    case r.Float64                 : return x.Float() == y.Float()
    case r.Complex128              : return x.Complex() == y.Complex()
    case r.String                  : return x.String() == y.String()
    case r.Slice                   : return equalSlices(x, y)
    case r.Struct                  : return equalStructs(x, y)
    case r.Map                     : return equalMaps(x, y)
    case r.Interface               : return AnyEquals(x.Interface(), y.Interface())
    case r.Uintptr                 : return equalValues(r.Indirect(x), r.Indirect(y))
    case r.Ptr, r.UnsafePointer    : return equalValues(x.Elem(), y.Elem())
    default                        : return false
  }
}

// Callers fills the slice pc with the return program counters of function invocations on the calling goroutine's stack.
// The argument skip is the number of stack frames to skip before recording in pc, with 0 identifying the frame for
// Callers itself and 1 identifying the caller of Callers. It returns the number of entries written to pc.
//
// Note that since each slice entry pc[i] is a return program counter, looking up the file and line for pc[i] (for
// example, using (*Func).FileLine) will return the file and line number of the instruction immediately following the
// call. To look up the file and line number of the call itself, use pc[i]-1. As an exception to this rule, if pc[i-1]
// corresponds to the function runtime.sigpanic, then pc[i] is the program counter of a faulting instruction and should
// be used without any subtraction.
func FunctionCaller() CallerInfo {
          pc := make([]uintptr, 4)
     /*count := */ runtime.Callers(2, pc) // 0 = Callers, 1 = TestReflect
          fn := runtime.FuncForPC(pc[0])
        name := extractFnName.ReplaceAllString(fn.Name(), "$1")
  path, line := fn.FileLine(pc[0] - 1)

  return CallerInfo {
    File: NewPath(path),
    Name: name,
    Line: line,
  }
}

func (x CallerInfo) String() string {
  return fmt.Sprintf("%v:%v:%v", x.File, x.Name, x.Line)
}
