package main

import (
  "os"
  "github.com/rveen/ogdl"
  "github.com/paulp/suffuse"
)

func main() {
  for _, file := range os.Args[1:] {
    g := ogdl.ParseFile(file)
    suffuse.Printfln("\n**** %v\n\n%v", file, g.Text())
  }


  // file := os.Args[1]
  // g := ogdl.ParseFile(file)
  // ip,_ := g.GetString("eth0.ip")
  // to,_ := g.GetInt64("eth0.timeout")

  // fmt.Println("ip:",ip,", timeout:",to)
}
