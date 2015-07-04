package main

import (
  "os"
  "fmt"
)

func main() {
  if len(os.Args) != 2 { errorAndExit("This application accepts exactly one file") }

  file, err := os.Open(os.Args[1])
  if err != nil { errorAndExit(err.Error()) }
  defer file.Close()

  fi, err := file.Stat()

  format := "2006-01-02 15:04:05"

  date := fi.ModTime().Format(format)

  fmt.Printf("%s %d %s %s\n",
    fi.Mode().String(),
    fi.Size(),
    date,
    fi.Name(),
  )
}

func errorAndExit(msg string) {
  fmt.Fprintf(os.Stderr, msg + "\n")
  os.Exit(2)
}
