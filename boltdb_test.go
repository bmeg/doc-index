package dindx

import (
  "testing"
  "log"
  "os"
  "strings"
)

func TestDocIndex(t *testing.T) {
  testFile := "testdb.idx"
  if _, err := os.Stat(testFile); err != nil {
    os.Remove(testFile)
  }
  
  log.Printf("Opening DB")
  idx := Open("testdb.idx")
  
  doc1 := map[string]interface{}{ "field1" : "value1" }
  
  idx.AddDoc("doc1", doc1)
  
  for a := range idx.Field("field1").Word("value1").Keys() {
    if a != "doc1" {
      t.Error("Wrong key")
    }
  }  
  log.Printf("Done")

}


func TestFieldInterface(t *testing.T) {
  testFile := "testdb.idx"
  if _, err := os.Stat(testFile); err != nil {
    os.Remove(testFile)
  }
    
  log.Printf("Opening DB")
  idx := Open("testdb.idx")
  
  log.Printf("Adding Keys")
  idx.Field("field").Word("value").Add("key1")
  idx.Field("field").Word("value").Add("key2")
  idx.Field("field").Word("value").Add("key3")
  idx.Field("field").Word("value").Add("key4")

  log.Printf("Removing Key")  
  idx.Field("field").Word("value").Remove("key1")
  
  count := 0
  for a := range idx.Field("field").Word("value").Keys() {
    if !strings.HasPrefix(a, "key") {
      t.Error("Wrong key")
    }
    count++
  }  
  if count != 3 {
    t.Error("Wrong key count")    
  }
  log.Printf("Done")

}