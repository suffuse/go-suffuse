package suffuse

import (
  "os"
  "strings"
  "bazil.org/fuse"
  . "gopkg.in/check.v1"
)

/** TODO de-dup the logic from this and other fuse startups.
 */
func SfsFromMap(mountPoint Path, nodes map[Path]interface{})*Sfs {
  conn, err := fuse.Mount(string(mountPoint))
  if err != nil { return nil }

  root := NewRootInode()
  root.AddNodeMap(nodes)
  vfs := &Sfs {
    Mountpoint: mountPoint,
    RootNode: root,
    Connection: conn,
  }
  trap := func (sig os.Signal) {
    Echoerr("Caught %v - forcing unmount(2) of %v\n", sig, mountPoint)
    vfs.Unmount()
  }
  TrapExit(trap)
  return vfs
}

/** TODO put these big expected strings in text files and slurp the files.
 */
var expectedTree string = strings.TrimSpace(`
.
|-- 02/
|   |-- 02/
|   |   \-- 4*
|   |-- 03/
|   |   \-- 6*
|   |-- 05/
|   |   \-- 10*
|   |-- 07/
|   |   \-- 14*
|   \-- 11/
|       \-- 22*
|-- 03/
|   |-- 02/
|   |   \-- 6 -> ../../02/03/6*
|   |-- 03/
|   |   \-- 9*
|   |-- 05/
|   |   \-- 15*
|   |-- 07/
|   |   \-- 21*
|   \-- 11/
|       \-- 33*
|-- 05/
|   |-- 02/
|   |   \-- 10 -> ../../02/05/10*
|   |-- 03/
|   |   \-- 15 -> ../../03/05/15*
|   |-- 05/
|   |   \-- 25*
|   |-- 07/
|   |   \-- 35*
|   \-- 11/
|       \-- 55*
|-- 07/
|   |-- 02/
|   |   \-- 14 -> ../../02/07/14*
|   |-- 03/
|   |   \-- 21 -> ../../03/07/21*
|   |-- 05/
|   |   \-- 35 -> ../../05/07/35*
|   |-- 07/
|   |   \-- 49*
|   \-- 11/
|       \-- 77*
\-- 11/
    |-- 02/
    |   \-- 22 -> ../../02/11/22*
    |-- 03/
    |   \-- 33 -> ../../03/11/33*
    |-- 05/
    |   \-- 55 -> ../../05/11/55*
    |-- 07/
    |   \-- 77 -> ../../07/11/77*
    \-- 11/
        \-- 121*

30 directories, 25 files
`)

var expectedFiles = strings.TrimSpace(`
 2 *  2 = 4
 2 *  3 = 6
 2 *  5 = 10
 2 *  7 = 14
 2 * 11 = 22
 3 *  3 = 9
 3 *  5 = 15
 3 *  7 = 21
 3 * 11 = 33
 5 *  5 = 25
 5 *  7 = 35
 5 * 11 = 55
 7 *  7 = 49
 7 * 11 = 77
11 * 11 = 121
`)

var expectedLinks = strings.TrimSpace(`
 2 *  3 = 6
 2 *  5 = 10
 3 *  5 = 15
 2 *  7 = 14
 3 *  7 = 21
 5 *  7 = 35
 2 * 11 = 22
 3 * 11 = 33
 5 * 11 = 55
 7 * 11 = 77
`)

// Translate backticks to backslashes in the output of tree.
func removeBackticks(s string)string {
  return NewRegex("[`]").ReplaceAllLiteralString(s, "\\")
}

// Exercising basic filesystem functions.
func primeFilesystemMap()map[Path]interface{} {
  numbers := []int { 2, 3, 5, 7, 11 }
  mapping := make(map[Path]interface{}, 0)

  for _, n1 := range numbers {
    for _, n2 := range numbers {
      path := Path(Sprintf("%02v/%02v/%v", n1, n2, n1 * n2))
      if n1 <= n2 {
        mapping[path] = Bytes(Sprintf("%2v * %2v = %v\n", n1, n2, n1 * n2))
      } else {
        mapping[path] = LinkTarget(Sprintf("../../%02v/%02v/%v", n2, n1, n1 * n2))
      }
    }
  }

  return mapping
}

func (s *Tsfs) TestNewVfs(c *C) {
  mnt := ScratchDir()
  mapping := primeFilesystemMap()
  mfs := SfsFromMap(mnt, mapping)

  defer mfs.Unmount()

  go func() {
    MaybeLog(mfs.Serve())
    <- mfs.Connection.Ready
    MaybeLog(mfs.Connection.MountError)
  }()

  // Without a sleep here the tests non-deterministically fail with no output.
  // Not sure what we're supposed to wait on.
  SleepMillis(50)

  // Use "find" to select a subset of virtual nodes and "cat"
  // to obtain their contents.
  getpaths := func(arg string)[]string {
    paths := ExecIn(mnt, "find", ".", "-type", arg).Lines().SortedStrings()
    return append(Strings("cat"), paths...)
  }

  treeOut := removeBackticks(strings.TrimSpace(ExecIn(mnt, "tree", "-F", "-n", "--charset", "US-ASCII").Slurp()))
  AssertString(c, treeOut, expectedTree)

  // find . -type f
  findFiles := strings.TrimSpace(ExecIn(mnt, getpaths("f")...).Slurp())
  AssertString(c, findFiles, expectedFiles)

  // find . -type l
  findLinks := strings.TrimSpace(ExecIn(mnt, getpaths("l")...).Slurp())
  AssertString(c, findLinks, expectedLinks)
}
