package suffuse

import (
  . "gopkg.in/check.v1"
)

type TestEq struct {
  Bar string
  Baz map[string]int
}
type TestEq2 struct {
  Bar string
  Baz map[string]int
}

func (s *Tsfs) TestJsonReadWrite(c *C) {
  f := ScratchFile()
  thing1 := map[string]interface{} { "a" : "5" }
  WriteJsonFile(f, thing1)
  thing2 := ReadJsonFile(f)
  // Echoerr("thing1: %#v: %T", thing1, thing1)
  // Echoerr("thing2: %#v: %T", thing2, thing2)
  c.Assert(thing1, DeepEquals, thing2)
}

func (s *Tsfs) TestAnyContains(c *C) {
  xs := []int { 1, 2, 3 }
  ys := map[int]int { 1: 5, 2: 6, 3: 7 }

  c.Assert(AnyContains(xs, 2), Equals, true)
  c.Assert(AnyContains(xs, 4), Equals, false)
  c.Assert(AnyContains(ys, 2), Equals, true)
  c.Assert(AnyContains(ys, 4), Equals, false)
  c.Assert(AnyContains(xs, 2), Equals, true)
}

func (s *Tsfs) TestAnyEquals(c *C) {
  xs := TestEq { "abc", map[string]int { "a": 6 } }
  ys := TestEq { "abc", map[string]int { "a": 5 } }
  zs := TestEq { "abc", map[string]int { "a": 5 } }
  qs := TestEq2 { "abc", map[string]int { "a": 5 } }

  a1 := 5
  a2 := 5
  a3 := 6

  c.Assert(AnyEquals(xs, zs), Equals, false)
  c.Assert(AnyEquals(ys, zs), Equals, true)
  c.Assert(AnyEquals(&a1, &a2), Equals, true)
  c.Assert(AnyEquals(&a1, &a3), Equals, false)

  a3 = 5
  c.Assert(AnyEquals(&a1, &a3), Equals, true)
  c.Assert(AnyEquals(zs, qs), Equals, false)

  sl1 := []int { 1, 2, 3 }
  sl2 := [3]int{ 1, 2, 3 }
  c.Assert(AnyEquals(sl1, sl2), Equals, true)
}

func (s *Tsfs) TestReflect(c *C) {
  var x0 = Path { "dog" }
  var x1 = Any(x0)
  desc := "suffuse.Path{Path:\"dog\"}: suffuse.Path"

  c.Assert(FunctionCaller().Name, Equals, "TestReflect")
  c.Assert(x1.What(), Equals, desc)
}
