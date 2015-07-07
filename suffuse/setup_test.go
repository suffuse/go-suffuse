package suffuse

import (
  "testing"
  "strings"
  "runtime"
  . "gopkg.in/check.v1"
)

var _ = Suite(&Tsfs{})

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type Tsfs struct { In Path ; Out Path }

/** A bug in either osxfuse, fuse, or the tools which read statfs
 *  structures means that we cannot avoid different outcome for ls.
 *  In particular the "total NN" line differs. The details are
 *  not especially interesting. We just filter out the tottal line.
 */
var totalRegex = NewRegex(`total \d+`)

func make_test_fs(os string) string {
  xattrPart := ""

  if os == "darwin" {
    xattrPart = "\n" + strings.TrimSpace(`
      xattr    -w suffuse.type link file.txt
      xattr -s -w suffuse.type link flink.txt
    `)
  }
  pre := strings.TrimSpace(`
echo "hello world" > file.txt
ln -s file.txt flink.txt
  `)
  post := strings.TrimSpace(`
mkdir dir
echo "hello dir" > dir/sub.txt
ln -s dir dlink
seq 1 10000 > bigfile.txt
  `)

  return NewLines(pre, xattrPart, post).String()
}

func startFuse(args ...string) {
  opts, opts_err := ParseSfsOpts(args)
  if opts_err != nil { return }

  mfs, err := NewSfs(opts)
  MaybePanic(err)

  go func() {
    defer mfs.Unmount()
    MaybeLog(mfs.Serve())
    <- mfs.Connection.Ready
    MaybeLog(mfs.Connection.MountError)
  }()
}

/** Same result string for both given Paths.
 */
func AssertSame(s *Tsfs, c *C, f func(Path)string) {
  c.Assert(f(s.In), Equals, f(s.Out))
}
func AssertExecLine(c *C, expect string, args ...string) {
  c.Assert(Exec(args...).OneLine(), Equals, expect)
}
func AssertExecSuccess(c *C, args ...string) {
  c.Assert(Exec(args...).Success(), Equals, true)
}
func AssertSameFile(c *C, p1, p2 Path) {
  c.Assert(p1.EvalSymlinks(), Equals, p2.EvalSymlinks())
}

func (s *Tsfs) SetUpSuite(c *C) {
  logI("SetUpSuite(%s)\n", *s)

  PsutilHostDump()
  // fmt.Println(GetStack())

  s.In  = NewPath(c.MkDir())
  s.Out = NewPath(c.MkDir())
  c.Assert(ExecBashIn(s.In, make_test_fs(runtime.GOOS)).Err, IsNil)
  startFuse("suffuse", "-m", s.Out.Path, s.In.Path)
}

func (s *Tsfs) TearDownSuite(c *C) {
  Exec("umount", "-f", s.Out.Path)
  logI("TearDownSuite(%s)\n", c.TestName())
}
func (s *Tsfs) SetUpTest(c *C) {
  logD("SetUpTest(%s)\n", c.TestName())
}
func (s *Tsfs) TearDownTest(c *C) {
  logD("TearDownTest(%s)\n", c.GetTestLog())
}
