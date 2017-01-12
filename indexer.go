
package dindx

import (
  //"log"
)

type WordIndex interface {
  Add(key string)
  Set(key string, count int64)
  Remove(key string)
  Keys() chan string
}

type FieldTable interface {
  Words() chan string
  Word(word string) WordIndex
}

type Database interface {
  Fields() chan string
  Field(fieldName string) FieldTable
}


type DocIndex struct {
  db Database
}


func Open(path string) *DocIndex {
  return &DocIndex{db:NewBoltDB(path)}
}

func (docIndex *DocIndex) Fields() chan string {
  return docIndex.db.Fields()
}

func (docIndex *DocIndex) Field(field string) FieldTable {
  return docIndex.db.Field(field)
}

func (docIndex *DocIndex) AddDoc(key string, doc map[string]interface{}) {
  for k, v := range doc {
    switch v.(type) {
    case string:
      docIndex.db.Field(k).Word(v.(string)).Add(key)
    }
  }
}