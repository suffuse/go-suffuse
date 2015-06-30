package suffuse

import (
  "os"
  "fmt"
  "runtime"
  . "gopkg.in/check.v1"
)

/** Mount a fuse filesystem and check that functions
 *  produce the same result against the underlying filesystem
 *  and the fuse filesystem.
 */
func (s *Tsfs) TestSameWithLs(c *C) {
  // no -@ option to ls except on osx.
  if runtime.GOOS == "darwin" {
    AssertSame(s, c,
      func(p Path) string {
        return ExecIn(p, "ls", "-laR@").Lines().FilterNot(totalRegex).String()
      },
    )
  }
}
func (s *Tsfs) TestSameWithFind(c *C) {
  AssertSame(s, c,
    func(p Path) string {
      return ExecIn(p, "find", ".").Lines().FilterNot(totalRegex).String()
    },
  )
}
func (s *Tsfs) TestSameWithWalk(c *C) {
  AssertSame(s, c,
    func(root Path) string {
      return root.WalkCollect(
        func(path string, info os.FileInfo) string {
          return fmt.Sprintf("[%-6d] %s\n", info.Size(), path)
        },
      ).String()
    },
  )
}
func (s *Tsfs) TestSameWithDiff(c *C) {
  c.Assert(GitWordDiff(s.In, s.Out).Slurp(), Equals, "")
  c.Assert(Exec("diff", "-qr", s.In.Path, s.Out.Path).Err, IsNil)
}

func (s *Tsfs) TestXattr(c *C) {
  // xattr doesn't work inside a container.
  switch runtime.GOOS {
    case "darwin" :
      Printfln("Running xattr test on darwin")
      AssertExecSuccess(c, "xattr", "-w", "dog", "pony", s.Out.Path)
      // c.Assert(Exec("xattr", "-w", "dog", "pony", s.Out.Path).Success(), Equals, true)
      c.Assert(Exec("xattr", s.Out.Path).OneLine(), Equals, "dog")
      c.Assert(Exec("xattr", "-l", s.Out.Path).OneLine(), Equals, "dog: pony")
      c.Assert(Exec("xattr", "-p", "dog", s.Out.Path).OneLine(), Equals, "pony")
      c.Assert(Exec("xattr", "-d", "dog", s.Out.Path).Success(), Equals, true)
      c.Assert(Exec("xattr", "-l", s.Out.Path).OneLine(), Equals, "")
    default:
      Printfln("Skipping xattr test on %s", runtime.GOOS)
  }
}

func (s *Tsfs) TestSedSuffix(c *C) {
  s1 := ExecIn(s.In, "sed", "-ne", "11,13p", "bigfile.txt").OneLine()
  s2 := ExecIn(s.Out, "cat", "bigfile.txt#11,13p").OneLine()
  expected := "11 12 13"

  c.Assert(s1, Equals, expected)
  c.Assert(s2, Equals, expected)
}
