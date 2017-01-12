
package dindx

import (
  "bytes"
  "fmt"
  //"log"
  "encoding/base64"
  "crypto/sha1"
  "github.com/boltdb/bolt"
	proto "github.com/golang/protobuf/proto"
)


var HASH_LEN = 4
var INDEX_PREFIX = "index:"

type BoltDB struct {
  db *bolt.DB
}

type BoltBucket struct {
  db *bolt.DB
  field string
}

type BoltWordIndex struct {
  db *bolt.DB
  field string
  word string
}

func NewBoltDB(path string) BoltDB {
  db, _ := bolt.Open(path, 0600, nil)
  return BoltDB{db:db}
}

func (self BoltDB) Field(fieldName string) FieldTable {
  err := self.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(fieldName))
    if b == nil {
      return fmt.Errorf("Bucket not found")
    }
    return nil
  })
  if err != nil {
    self.db.Update(func(tx *bolt.Tx) error {
		    _, err = tx.CreateBucket([]byte(fieldName))
        return err
    })
	}
  return BoltBucket{field:fieldName, db:self.db}
}


func (self BoltDB) Fields() chan string {  
  ch := make(chan string, 10)
  go func() {
    self.db.View(func(tx *bolt.Tx) error {
      tx.ForEach( func(name []byte, b *bolt.Bucket) error {
        ch <- string(name)
        return nil
      })
      return nil      
    })
    close(ch)
  } ()
  return ch
}

func ValueKeyHash(value string, key string) string {
  hasher := sha1.New()
  hasher.Write([]byte(key))
  h := hasher.Sum(nil)
  s := base64.URLEncoding.EncodeToString(h)[0:HASH_LEN]
  return fmt.Sprintf("%s%s", value, s)
}


func (self BoltBucket) Words() chan string {
  o := make(chan string, 10)
  go func() {
    self.db.View(func(tx *bolt.Tx) error {
      b := tx.Bucket([]byte(self.field))
      c := b.Cursor()
      last := ""
      for k, _ := c.First(); k != nil; k, _ = c.Next() {
        i := string(k[0:len(k)-HASH_LEN])
        if i != last {
  			  o <- i
          last = i
        }   
  		}
      close(o)
      return nil
    })
  } ()
  return o
}

func (self BoltBucket) Word(word string) WordIndex {
  return BoltWordIndex{ db:self.db, field:self.field, word:word }
}


func (self BoltWordIndex) Keys() chan string {
  o := make(chan string, 10)
  bword := []byte(self.word)
  go func() {
    self.db.View(func(tx *bolt.Tx) error {
      b := tx.Bucket([]byte(self.field))
      c := b.Cursor()
      for k, v := c.Seek(bword); k != nil && bytes.Equal(k[0:len(k)-HASH_LEN], bword); k, v = c.Next() {
          f := FieldIndex{}
          proto.Unmarshal(v, &f)
          for key, _ := range(f.Counts) {
            o <- key
          }
      }
      close(o)
      return nil
    })
  } ()
  return o
}


func (self BoltWordIndex) Add(key string) {
  h := ValueKeyHash(self.word, key)
  self.db.Update(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(self.field))
    v := b.Get([]byte(h))
    var o *FieldIndex
    if (v == nil) {
      o = &FieldIndex{Counts:map[string]int64{}}
    } else {
      o = &FieldIndex{}
      proto.Unmarshal(v, o)
    }
    o.Counts[key] = o.Counts[key] + 1
    d, _ := proto.Marshal(o)
		b.Put([]byte(h), d)
    return nil
  })
}

func (self BoltWordIndex) Set(key string, count int64) {
  h := ValueKeyHash(self.word, key)
  self.db.Update(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(self.field))
    v := b.Get([]byte(h))
    var o *FieldIndex
    if (v == nil) {
      o = &FieldIndex{Counts:map[string]int64{}}
    } else {
      o = &FieldIndex{}
      proto.Unmarshal(v, o)
    }
    o.Counts[key] = count
    d, _ := proto.Marshal(o)
		b.Put([]byte(h), d)
    return nil
  })
}


func (self BoltWordIndex) Remove(key string) {
  h := ValueKeyHash(self.word, key)
  self.db.Update(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(self.field))
    v := b.Get([]byte(h))
    var o *FieldIndex
    if (v != nil) {
      o = &FieldIndex{}
      proto.Unmarshal(v, o)
      delete(o.Counts, key)      
      if len(o.Counts) == 0 {
        b.Delete([]byte(h))
      } else {
        d, _ := proto.Marshal(o)
    		b.Put([]byte(h), d)
      }      
    } 
    return nil
  })
}
