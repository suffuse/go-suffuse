package suffuse

import (
  "time"
  . "gopkg.in/check.v1"
)

func NewDuration(spec string) time.Duration {
  dur, err := time.ParseDuration(spec)
  MaybePanic(err)
  return dur
}
func Hours(n int) time.Duration {
  return NewDuration(Sprintf("%d", n) + "s")
}

func (s *Tsfs) TestPathOsMethods(c *C) {
  root := ScratchDir()
  dir1 := root.Join("a/b/c")

  link1 := root.Join("link1")
  link1.OsSymlink(dir1)

  link2 := root.Join("link2")
  link2.OsSymlink(dir1.Absolute())

  dir1.OsMkdirAll(NewFilePerms(0755))
  dir1 = dir1.EvalSymlinks()
  dir1.OsChdir()

  link3 := NewPath("link3")
  link3.OsSymlink(DotPath)

  dir2 := CwdNoLinks()
  dir3 := NewPath(".").Absolute()
  c.Assert(dir1, Equals, dir2)
  c.Assert(dir1, Equals, dir3)

  AssertSameFile(c, link1, link2)
  AssertSameFile(c, link1, link3)

  file1 := root.Join("file1")
  _, err := file1.OsCreate()
  c.Assert(err, IsNil)

  /** TODO - run these through fuse as well. */
  atime, mtime := file1.OsStatAtimeMtime()
  atime0 := atime.AddDate(-1, -2, -3)
  mtime0 := mtime.AddDate(-4, -5, -6)
  file1.OsChtimes(atime0, mtime0)

  atime1, mtime1 := file1.OsStatAtimeMtime()
  c.Assert(atime0, Equals, atime1)
  c.Assert(mtime0, Equals, mtime1)

  // OsChdir
  // OsChmod
  // OsChown
  // OsChtimes
  // OsMkdir
  // OsMkdirAll
  // OsRemoveAll
  // OsRename
  // OsLstat
  // OsOpen
  // OsOpenFile
  // OsReadLink
  // OsStat
}
