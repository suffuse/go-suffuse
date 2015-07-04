package suffuse

import (
  "github.com/satori/go.uuid"
  "os"
  "os/exec"
)

func ExampleSed() {
  with(func (in string, out string) {

    execute(in)(`
      seq 1 10 > seq.txt            # Fill a file with digits 1 to 10.
      touch -t 201507100000 seq.txt # Make sure it has the expected date
    `)

    // Mount the directory through suffuse.
    startFuse("suffuse", "-m", out, in)

    execute(out)(`
      ls                   # List the contents through the suffuse mount.
      wc -l seq.txt        # It's a 10 line file, one number to a line.
      cat seq.txt#4,6p     # Via suffuse, a derived file ending with #4,6p is a sed
                           # command executed on the actual file.
      ls -gG seq.txt#5,10p # Arbitrary sed commands, different sized files.
      ls -gG seq.txt#1,3p  # These files are effectively indistinguishable from "real" files.
    `)

    //Set `umask 022` if the output does not match

    //Output:
    // seq.txt
    // 10 seq.txt
    // 4
    // 5
    // 6
    // -rw-r--r-- 1 13 Jul 10  2015 seq.txt#5,10p
    // -rw-r--r-- 1 6 Jul 10  2015 seq.txt#1,3p

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

func execute(dir string) func (string) error {
  return func(script string) error {
    cmd := exec.Command("sh", "-c", script)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Dir = dir
    return cmd.Run()
  }
}
