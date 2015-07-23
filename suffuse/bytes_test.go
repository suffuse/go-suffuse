package suffuse

import (
  . "gopkg.in/check.v1"
)

func (s *Tsfs) TestBytesSliceRange(c *C) {
  // byte slices.
  bs := Bytes([]byte("0123456789"))
  s1 := RangeOffsetLength(4, 6)
  s2 := RangeOffsetLength(0, 4)
  s3 := s1.SliceRange(s2)
  s4 := s2.SliceRange(s1)

  AssertString(c, string(bs.SliceRange(s1)), "456789")
  AssertString(c, string(bs.SliceRange(s2)), "0123")
  AssertString(c, string(bs.SliceRange(s3)), "4567")
  AssertString(c, string(bs.SliceRange(s4)), "")

  AssertString(c, string(bs.SliceRange(s1).SliceRange(s2)), "4567")
  AssertString(c, string(bs.SliceRange(s2).SliceRange(s1)), "")
  AssertString(c, string(bs.SliceRange(s1).Slice(0, 4)), "4567")
  AssertString(c, string(bs.SliceRange(s2).Slice(4, 6)), "")
}
