package main

import (
  "fmt"
  "log"
  dindx "github.com/bmeg/doc-index"
  "flag"
)

func main() {
  log.Printf("Starting Search")
  db_p := flag.String("db", "docs.idx", "Doc Database")
  flag.Parse()

  idx := dindx.Open(*db_p)
  
  if len(flag.Args()) == 0 {
    for a := range idx.Fields() {
      fmt.Printf("%s\n", a)
    }
  }
  if len(flag.Args()) == 1 {
    for a := range idx.Field(flag.Args()[0]).Words() {
      fmt.Printf("%s\n", a)
    }
  }
  if len(flag.Args()) == 2 {
    for a := range idx.Field(flag.Args()[0]).Word(flag.Args()[1]).Keys() {
      fmt.Printf("%s\n", a)
    }
  }
  
}