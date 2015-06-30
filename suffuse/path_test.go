package suffuse

import (
  . "gopkg.in/check.v1"
)

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

  // c.Assert(link1.EvalSymlinks(), Equals, link2.EvalSymlinks())


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
