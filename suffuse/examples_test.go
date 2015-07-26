package suffuse

import (
  "github.com/satori/go.uuid"
  "os"
  "os/exec"
  "bytes"
  . "testing"
  "strings"
)

func TestSed(t *T) {
  with(func (in string, out string) {

    execute(in)(`
      seq 1 10 > seq.txt            # Fill a file with digits 1 to 10.
      touch -t 201507100000 seq.txt # Make sure it has the expected date
    `)

    // Mount the directory through suffuse.
    startFuse("suffuse", "-m", out, "-c", "suffuse_tests.conf", in)

    result, _ := execute(out)(`
      ls                   # List the contents through the suffuse mount.
      wc -l seq.txt        # It's a 10 line file, one number to a line.
      cat seq.txt#4,6p     # Via suffuse, a derived file ending with #4,6p is a sed
                           # command executed on the actual file.

      # Get a file utility to display coherent cross platform information
      go get github.com/suffuse/go-suffuse/cmd/fileinfo

      fileinfo seq.txt#5,10p # Arbitrary sed commands, different sized files.
      fileinfo seq.txt#1,3p  # These files are effectively indistinguishable from "real" files.
    `)

    //Set `umask 022` if the output does not match

    compareIgnoringWhitespace(t, result,
      `|seq.txt
       |10 seq.txt
       |4
       |5
       |6
       |-rw-r--r-- 13 2015-07-10 00:00:00 seq.txt#5,10p
       |-rw-r--r-- 6 2015-07-10 00:00:00 seq.txt#1,3p`)
  })
}

func TestJsonToYaml(t *T) {
  with(func (in string, out string) {

    installConverterFunctions()

    execute(in)(`
      echo '{"test": "file", "list" : [1, 2, 3]}' > jsonFile
    `)

    // Mount the directory through suffuse.
    startFuse("suffuse", "-m", out, "-c", "suffuse_tests.conf", in)

    result, _ := execute(out)(`
      cat jsonFile.yaml
    `)

    compare(t, result,
      `|list:
       |- 1
       |- 2
       |- 3
       |test: file`)
  })
}

func with(f func(string, string)) {
  in := newTmpDir()
  out := newTmpDir()

  f(in, out)

  Path(out).SysUnmount()

  os.Remove(in)
  os.Remove(out)
}

func newTmpDir() string {
  path := Sprintf("%s%csuffuse_example-%s", os.TempDir(), os.PathSeparator, uuid.NewV4())
  os.Mkdir(path, 0700)
  return path
}

func execute(dir string) func (string) (string, error) {
  return func(script string) (string, error) {
    stdout := new(bytes.Buffer)
    errout := new(bytes.Buffer)

    cmd := exec.Command("bash", "-c", script)
    cmd.Stdout = stdout
    cmd.Stderr = errout
    cmd.Dir = dir
    err := cmd.Run()

    if err := errout.String(); len(err) > 0 { Echoerr(err) }
    return stdout.String(), err
  }
}

func compareIgnoringWhitespace(t *T, result string, expected string) {
  compareWithOptions(t, "-w", result, expected)
}

func compare(t *T, result string, expected string) {
  compareWithOptions(t, "", result, expected)
}

func compareWithOptions(t *T, opts string, result string, expected string) {
  script := "diff " + opts + " -y " + echo(result) + " " + echo(StripMargin('|', expected))

  result, err := execute("")(script)

  if err != nil { t.Error(result) }
}

func echo(s string) string {
  trimmed := strings.TrimSpace(s)
  escaped := strings.Replace(trimmed, "\n", "\\n", -1)
  return `<(echo -e "\n` + escaped + `")`
}

func installConverterFunctions() {
  execute(string(cwd()))(`
    go get github.com/dbohdan/remarshal

    for if in toml yaml json; do
      for of in toml yaml json; do
        FILE="$GOPATH/bin/${if}2${of}"
        if [ ! -f $FILE ]; then ln -s "$GOPATH/bin/remarshal" $FILE; fi
      done
    done
  `)
}
