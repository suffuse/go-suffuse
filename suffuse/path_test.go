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
  return NewDuration(Sprintf("%dh", n))
}

/** TODO - run all path methods through fuse as well. */
func (s *Tsfs) TestPathTimeManipulation(c *C) {
  f := ScratchFile()
  // Read the times on a scratch file.
  atime0, mtime0 := f.OsStatAtimeMtime()
  // Calculate recognizably different times from those.
  atime1 := atime0.AddDate(-1, -2, -3)
  mtime1 := mtime0.AddDate(-4, -5, -6)
  // Set the scratch file times to the different ones.
  f.OsChtimes(atime1, mtime1)
  // Read the file times again.
  atime2, mtime2 := f.OsStatAtimeMtime()
  // And hope they look like the ones we set.
  c.Assert(atime1, Equals, atime2)
  c.Assert(mtime1, Equals, mtime2)
}

func (s *Tsfs) TestPathSymlinkLogic(c *C) {
  root := ScratchDir()
  d := root.Join("a/b/c")

  // Create symlink "root/link1" -> "root/a/b/c"
  link1 := root.Join("link1")
  link1.OsSymlink(d)

  // Create symlink "root/link2" -> "a/b/c"
  link2 := root.Join("link2")
  link2.OsSymlink(Path("a/b/c"))

  // Create a/b/c, resolve all links, and chdir into it.
  d.OsMkdirAll(NewFilePerms(0755))
  d = d.EvalSymlinks()
  d.OsChdir()

  // Create symlink "link3" -> ".", where "." is now $d.
  link3 := Path("link3")
  link3.OsSymlink(DotPath)

  // Ensure that $d, $(pwd), and "." all resolve to the same path.
  c.Assert(d, Equals, CwdNoLinks())
  c.Assert(d, Equals, DotPath.Absolute())

  // Ensure that link1, link2, and link3 all resolve to the same path.
  AssertSameFile(c, link1, link2)
  AssertSameFile(c, link1, link3)
}

func (s *Tsfs) TestPathFileCreation(c *C) {
  root := ScratchDir()
  f := root.Join("file")
  _, err := f.OsCreate()
  c.Assert(err, IsNil)
}

// Still to Test.
//
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
