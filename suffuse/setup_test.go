package suffuse

import (
  "fmt"
  "testing"
  "runtime"
  . "gopkg.in/check.v1"
  "math/rand"
  "time"
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
var totalRegex = Regexp(`total \d+`)

func makeTestFs(os string) string {
  xattrPart := ""

  if os == "darwin" {
    xattrPart = "\n" + TrimSpace(`
      xattr    -w suffuse.type link file.txt
      xattr -s -w suffuse.type link flink.txt
    `)
  }
  pre := TrimSpace(`
echo "hello world" > file.txt
ln -s file.txt flink.txt
  `)
  post := TrimSpace(`
mkdir dir
echo "hello dir" > dir/sub.txt
ln -s dir dlink
seq 1 10000 > bigfile.txt
  `)

  return Strings{pre, xattrPart, post}.String()
}

func startFuse(args ...string) {
  conf := CreateSfsConfig(args)
  mfs, err := NewSfs(conf)
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
// "found fmt.Stringer" means you can't pass a string. Oyve.
func AssertString(c *C, found interface{}, expected string) {
  var str string
  switch found := found.(type) {
    case fmt.Stringer : str = found.String()
    case string       : str = found
    default           : panic(Sprintf("Unexpected type %T", found))
  }
  c.Assert(str, Equals, expected)
}

func (s *Tsfs) SetUpSuite(c *C) {
  info("SetUpSuite(%s)\n", *s)
  psutilHostDump()

  // rand.Int() is used by the check library. Without this line it not so random...
  rand.Seed(time.Now().UnixNano())

  s.In  = Path(c.MkDir())
  s.Out = Path(c.MkDir())
  c.Assert(ExecBashIn(s.In, makeTestFs(runtime.GOOS)).Err, IsNil)
  startFuse("suffuse", "-m", string(s.Out), string(s.In))
}

func (s *Tsfs) TearDownSuite(c *C) {
  s.Out.SysUnmount()
  info("TearDownSuite(%s)\n", c.TestName())
}
func (s *Tsfs) SetUpTest(c *C) {
  trace("SetUpTest(%s)\n", c.TestName())
}
func (s *Tsfs) TearDownTest(c *C) {
  trace("TearDownTest(%s)\n", c.GetTestLog())
}

func (x Strings) filter(re Regex) Strings    { return filterCommon(x, re, true)    }
func (x Strings) filterNot(re Regex) Strings { return filterCommon(x, re, false)   }

func filterCommon(x Strings, re Regex, expectTrue bool) Strings {
  var ys []string
  for _, line := range x.Array() {
    if re.MatchString(line) == expectTrue {
      ys = append(ys, line)
    }
  }
  return Strings(ys)
}
