package suffuse

import (
  . "gopkg.in/check.v1"
)

// Brute force testing all the obvious boundary conditions.
// Ranges are printed as <start>:<length>
// All empty ranges are 0:0 so <start> cannot accidentally
// be treated as meaningful on a zero length range.
func (s *Tsfs) TestRange(c *C) {
  r := RangeStartEnd(1, 5) // [ 1, 2, 3, 4 ]
  AssertString(c, r, "1:4")

  // negative arg
  AssertString(c, r.Drop(-1), "1:4")
  AssertString(c, r.DropRight(-1), "1:4")
  AssertString(c, r.Take(-1), "0:0")
  AssertString(c, r.TakeRight(-1), "0:0")

  // zero arg
  AssertString(c, r.Drop(0), "1:4")
  AssertString(c, r.DropRight(0), "1:4")
  AssertString(c, r.Take(0), "0:0")
  AssertString(c, r.TakeRight(0), "0:0")

  // positive arg
  AssertString(c, r.Drop(1), "2:3")
  AssertString(c, r.DropRight(1), "1:3")
  AssertString(c, r.Take(1), "1:1")
  AssertString(c, r.TakeRight(1), "4:1")

  // length - 1
  AssertString(c, r.Drop(3), "4:1")
  AssertString(c, r.DropRight(3), "1:1")
  AssertString(c, r.Take(3), "1:3")
  AssertString(c, r.TakeRight(3), "2:3")

  // length
  AssertString(c, r.Drop(4), "0:0")
  AssertString(c, r.DropRight(4), "0:0")
  AssertString(c, r.Take(4), "1:4")
  AssertString(c, r.TakeRight(4), "1:4")

  // length + 1
  AssertString(c, r.Drop(5), "0:0")
  AssertString(c, r.DropRight(5), "0:0")
  AssertString(c, r.Take(5), "1:4")
  AssertString(c, r.TakeRight(5), "1:4")

  // slice is implemented in terms of take and drop so is
  // primarily assured by the above passing.
  AssertString(c, r.Slice(0, 0), "0:0")
  AssertString(c, r.Slice(0, 1), "1:1")
  AssertString(c, r.Slice(0, 2), "1:2")
  AssertString(c, r.Slice(1, 2), "2:1")
  AssertString(c, r.Slice(2, 1), "0:0")
  AssertString(c, r.Slice(-2, -1), "0:0")
  AssertString(c, r.Slice(-2, 2), "1:2")

}
