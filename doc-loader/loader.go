package main

import (
  "log"
  dindx "github.com/bmeg/doc-index"
  "flag"
  "os"
  "bufio"
  "encoding/json"
)

func ReadLines(path string) (chan []byte, error) {
  out := make(chan []byte, 100)
  if file, err := os.Open(path); err == nil {
    go func() {
      reader := bufio.NewReaderSize(file, 102400)
      var isPrefix bool = true
      var err error = nil
      var line, ln []byte
      for err == nil {
        line, isPrefix, err = reader.ReadLine()
        ln = append(ln, line...)
        if !isPrefix {
          out <- ln
          ln = []byte{}
        }
      }
      close(out)
   } ()
   return out, nil
  } else {
    return out, err
  }
}

func main() {
  log.Printf("Starting Load")
  db_p := flag.String("db", "docs.idx", "Doc Database")
  key_p := flag.String("key", "_id", "Key Field")
  flag.Parse()

  idx := dindx.Open(*db_p)
  ch := make(chan *map[string]interface{}, 10)
  go func() {
    for _, f := range flag.Args() {
      log.Printf("Loading %s", f)      
      if reader, err := ReadLines(f); err == nil {
        for line := range(reader) {
          o := map[string]interface{}{}
          if err := json.Unmarshal(line, &o); err == nil {
            ch <- &o
          } else {
            log.Printf("Error parsing: %s", err)
          }
        }
        //file.Close()
      } else {
        log.Printf("Error parsing: %s", err)
      }
    }
    close(ch)
  } ()
  
  count := 0
  skipped := 0
  for d := range ch {
    if k, ok := (*d)[*key_p]; ok {
      if ks, ok := k.(string); ok {
        idx.AddDoc(ks, *d)
        count += 1
      } else {
        skipped += 1
      }
    }
  }
  log.Printf("%d docs indexed (%d skipped)", count, skipped)
  
}