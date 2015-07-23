package suffuse

import (
  "io"
  "os"
  "os/exec"
  "github.com/shirou/gopsutil/host"
)

func psutilHostDump () {
  hostInfo, _ := host.HostInfo()
  Println(JsonPretty(*hostInfo))
}

type ExecResult struct {
  Err error
  Stdout []byte
}

func cwd() Path           { return MaybePath(os.Getwd()) }
func OsStderr() io.Writer { return os.Stderr             }

func (x ExecResult) OneLine() string { return x.Lines().JoinWords() }
func (x ExecResult) Lines() Strings  { return BytesToLines(x.Stdout)  }
func (x ExecResult) Slurp() string   { return string(x.Stdout)        }
func (x ExecResult) Success() bool   { return x.Err == nil            }

func Exec(args ...string) ExecResult                { return ExecIn(cwd(), args...)            }
func ExecBash(script string) ExecResult             { return ExecBashIn(cwd(), script)         }
func ExecBashIn(cwd Path, script string) ExecResult { return ExecIn(cwd, "bash", "-c", script) }

func ExecIn(cwd Path, args ...string) ExecResult {
  var res = ExecResult{}

  if len(args) == 0 {
    res.Err = NewErr("no command")
  } else {
    cmd := exec.Command(args[0], args[1:]...)
    cmd.Dir = string(cwd)

    bytes, err := cmd.Output()

    if err != nil {
      res.Err = err
    } else {
      res.Stdout = bytes
    }
  }

  return res
}

func GitWordDiff(p1 Path, p2 Path) ExecResult {
  return Exec("git", "diff", "--no-index", "--word-diff", string(p1), string(p2))
}
