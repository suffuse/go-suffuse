package suffuse

/** Go slices fail hard if the inputs don't make sense, but
 *  slices are much more convenient to work with if they do the
 *  sensible thing when you overshoot, e.g. slice[0:10] if there
 *  are only 5 elements would give you the 5. Range superimposes
 *  the scala-like forgiving slice logic on top of native slices.
 */
type Offset uint
type Length uint

type Range struct {
  Offset
  Length
}

/** All empty ranges are defined as 0:0. We don't want the
 *  start value of an empty range to be considered meaningful
 *  or stable under any operations.
 */
var EmptyRange Range = Range { Offset(0), Length(0) }

func (x Range) IsEmpty() bool  { return x.LengthInt() == 0           }
func (x Range) StartInt() int  { return int(x.Offset)                }
func (x Range) EndInt() int    { return x.StartInt() + x.LengthInt() }
func (x Range) LengthInt() int { return int(x.Length)                }

func (x Range) String() string {
  if x.IsEmpty() {
    return "0:0"
  } else {
    return Sprintf("%v:%v", x.StartInt(), x.LengthInt())
  }
}

func RangeOffsetLength(offset int, length int) Range {
  if offset < 0 {
    return RangeOffsetLength(0, length)
  } else if length <= 0 {
    return EmptyRange
  } else {
    return Range { Offset(offset), Length(length) }
  }
}
func RangeStartEnd(start int, end int) Range {
  if start < 0 {
    return RangeStartEnd(0, end)
  } else if end <= start || end <= 0 {
    return EmptyRange
  } else {
    return RangeOffsetLength(start, end - start)
  }
}

func (x Range) Drop(n int) Range {
  if n <= 0 { return x }
  return RangeStartEnd(x.StartInt() + n, x.EndInt())
}
func (x Range) DropRight(n int) Range {
  if n <= 0 { return x }
  return RangeStartEnd(x.StartInt(), x.EndInt() - n)
}
func (x Range) Take(n int) Range {
  length := min(n, x.LengthInt())
  return RangeOffsetLength(x.StartInt(), length)
}
func (x Range) TakeRight(n int) Range {
  length := min(n, x.LengthInt())
  return RangeOffsetLength(x.EndInt() - length, length)
}

func (x Range) SliceRange(r Range) Range {
  return x.Drop(r.StartInt()).Take(r.LengthInt())
}
func (x Range) Slice(start, end int) Range {
  return x.SliceRange(RangeStartEnd(start, end))
}

func min(x int, y int) int {
  if x <= y { return x }
  return y
}
