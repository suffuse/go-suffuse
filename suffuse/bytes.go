package suffuse

/** In any other language we could generically leverage the operations
 *  on Range and apply them to Bytes. Ha ha.
 */

type Bytes []byte

var NoBytes = Bytes([]byte {})

func (bs Bytes) ToRange()Range {
  return RangeOffsetLength(0, len(bs))
}
func (bs Bytes) SliceRange(slice Range)Bytes {
  r := bs.ToRange().SliceRange(slice)
  if r.IsEmpty() {
    return NoBytes
  } else {
    return bs[r.StartInt():r.EndInt()]
  }
}
func (bs Bytes) Slice(start, end int)Bytes {
  return bs.SliceRange(RangeStartEnd(start, end))
}
