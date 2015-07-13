package suffuse

import (
  "strings"
  . "gopkg.in/check.v1"
)

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
var treeOptions = Strings("tree", "-F", "-n", "--charset", "US-ASCII")

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

func (s *Tsfs) TestPrimesFilesystem(c *C) {
  WithSfs(func(mnt Path, root *Inode) {
    runSlurp := func(args ...string)string {
      return strings.TrimSpace(ExecIn(mnt, args...).Slurp())
    }
    // Use "find" to select a subset of virtual nodes and "cat"
    // to obtain their contents.
    getpaths := func(arg string)[]string {
      return append(Strings("cat"), ExecIn(mnt, "find", ".", "-type", arg).Lines().SortedStrings()...)
    }

    // Populate the filesystem.
    root.AddNodeMap(primeFilesystemMap())

    // `tree .`
    AssertString(c, removeBackticks(runSlurp(treeOptions...)), expectedTree)

    // `find . -type f | xargs cat`
    AssertString(c, runSlurp(getpaths("f")...), expectedFiles)

    // `find . -type l | xargs cat`
    AssertString(c, runSlurp(getpaths("l")...), expectedLinks)
  })
}
