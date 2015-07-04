package suffuse

import (
  "github.com/satori/go.uuid"
  "os"
  "os/exec"
  "time"
)

func ExampleSed() {
  with(func (in string, out string) {
    execOut := executeIn(out)

    // Fill a file with digits 1 to 10.
    pipe(command("seq", "1", "10"), in + "/seq.txt")

    // Mount a directory through suffuse.
    startFuse("suffuse", "-m", out, in)

    //List the contents through the suffuse mount.
    execOut("ls")

    // It's a 10 line file, one number to a line.
    execOut("wc", "-l", "seq.txt")

    // Via suffuse, a derived file ending with #4,6p is a sed
    // command executed on the actual file.
    execOut("cat", "seq.txt#4,6p")

    // Arbitrary sed commands, different sized files.
    execOut("ls", "-g", "-G", "seq.txt#5,10p")

    // These files are effectively indistinguishable from "real" files.
    execOut("ls", "-g", "-G", "seq.txt#1,3p")

    //Output:
    // seq.txt
    // 10 seq.txt
    // 4
    // 5
    // 6
    // -rw-rw-r-- 1 13 Jul 10  2015 seq.txt#5,10p
    // -rw-rw-r-- 1 6 Jul 10  2015 seq.txt#1,3p

  })
}

func with(f func(string, string)) {
  in := newTmpDir()
  out := newTmpDir()

  f(in, out)

  NewPath(out).SysUnmount()

  os.Remove(in)
  os.Remove(out)
}

func newTmpDir() string {
  path := Sprintf("%s%csuffuse_example-%s", os.TempDir(), os.PathSeparator, uuid.NewV4())
  os.Mkdir(path, 0700)
  return path
}

func executeIn(dir string) func (string, ...string) error {
  return func(name string, args ...string) error {
    cmd := command(name, args...)
    cmd.Stdout = os.Stdout
    cmd.Dir = dir
    return cmd.Run()
  }
}

func command(name string, args ...string) *exec.Cmd {
  cmd := exec.Command(name, args...)
  cmd.Stderr = os.Stderr
  return cmd
}

func pipe(cmd *exec.Cmd, file string) error {
  outfile, _ := os.Create(file)
  cmd.Stdout = outfile
  cmd.Run()
  time := time.Date(2015, 7, 10, 0, 0, 0, 0, time.UTC)
  os.Chtimes(file, time, time)
  return outfile.Close()
}
