package suffuse

/** Adding a small amount of type safety to indices and ranges.
 *  Index is in a signed int because otherwise there's no way to
 *  encode the None value. If we could control the type we'd
 *  enforce -1 to be the only invalid index, but you can't really
 *  do much in go so don't count on that being true.
 */
type Offset uint
type Length uint

type Range struct {
  Offset
  Length
}

var EmptyRange Range = Range { Offset(0), Length(0) }

func (x Range) IsEmpty()bool     { return x.LengthInt() == 0           }
func (x Range) StartInt()int     { return int(x.Offset)                }
func (x Range) EndInt()int       { return x.StartInt() + x.LengthInt() }
func (x Range) LengthInt()int    { return int(x.Length)                }
func (x Range) EndOffset()Offset { return Offset(x.EndInt())           }

func (x Range) String()string {
  if x.IsEmpty() {
    return "0:0"
  } else {
    return Sprintf("%v:%v", x.StartInt(), x.LengthInt())
  }
}
func min(x int, y int)int {
  if x <= y { return x }
  return y
}

func RangeOffsetLength(offset int, length int)Range {
  if offset < 0 {
    return RangeOffsetLength(0, length)
  } else if length <= 0 {
    return EmptyRange
  } else {
    return Range { Offset(offset), Length(length) }
  }
}
func RangeStartEnd(start int, end int)Range {
  if start < 0 {
    return RangeStartEnd(0, end)
  } else if end <= start || end <= 0 {
    return EmptyRange
  } else {
    return RangeOffsetLength(start, end - start)
  }
}

func (x Range) Tail()Range { return x.Drop(1) }
func (x Range) Init()Range { return x.DropRight(1) }

func (x Range) Drop(n int)Range {
  if n <= 0 { return x }
  return RangeStartEnd(x.StartInt() + n, x.EndInt())
}
func (x Range) DropRight(n int)Range {
  if n <= 0 { return x }
  return RangeStartEnd(x.StartInt(), x.EndInt() - n)
}
func (x Range) Take(n int)Range {
  length := min(n, x.LengthInt())
  return RangeOffsetLength(x.StartInt(), length)
}
func (x Range) TakeRight(n int)Range {
  length := min(n, x.LengthInt())
  return RangeOffsetLength(x.EndInt() - length, length)
}
func (x Range) SliceRange(r Range)Range {
  return x.Drop(r.StartInt()).Take(r.LengthInt())
}
func (x Range) Slice(start, end int)Range {
  return x.SliceRange(RangeStartEnd(start, end))
}
