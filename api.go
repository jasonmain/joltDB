package jolt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

type BoltDB struct {
	db   *bolt.DB
	open bool
}

func apiConn(db *bolt.DB) *BoltDB {
	return &BoltDB{db, false}
}

func (boltdb *BoltDB) apiOpenReadOnly(dir string) error {
	var err error
	config := &bolt.Options{ReadOnly: true, Timeout: 1 * time.Second}
	boltdb.db, err = bolt.Open(dir, 0666, config)
	if err != nil {
		return fmt.Errorf("db open error: %s.", err)
	}
	boltdb.open = true
	return nil
}

func (boltdb *BoltDB) apiOpen(dir string) error {
	var err error
	config := &bolt.Options{Timeout: 1 * time.Second}
	boltdb.db, err = bolt.Open(dir, 0600, config)
	if err != nil {
		return fmt.Errorf("db open error: %s.", err)
	}
	boltdb.open = true
	return nil
}

func (boltdb *BoltDB) apiClose() {
	boltdb.open = false
	boltdb.db.Close()
}

func (boltdb *BoltDB) apiCopyDB(dir string) error {
	err := boltdb.db.View(func(tx *bolt.Tx) error {
		err := tx.CopyFile(dir, 0666)
		if err != nil {
			return fmt.Errorf("Copy db error: %s.", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Could not get copy db.")
	}
	return nil
}

func (boltdb *BoltDB) apiSave(bucket, key string, data interface{}) error {
	if !boltdb.open {
		return fmt.Errorf("db must be opened before saving!")
	}
	err := boltdb.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("Create bucket error: %s.", err)
		}
		enc, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("Could not encode to json. Key %s: %s.", key, err)
		}
		err = b.Put([]byte(key), enc)
		if err != nil {
			return fmt.Errorf("Could not save: %s to %s. Error: %s", key, bucket, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Could not save: %s to %s. Error: %s", key, bucket, err)
	}
	return nil
}

func fmtToJsonArr(s []byte) []byte {
	s = bytes.Replace(s, []byte("{"), []byte("[{"), 1)
	s = bytes.Replace(s, []byte("}"), []byte("},"), -1)
	s = bytes.TrimSuffix(s, []byte(","))
	s = append(s, []byte("]")...)
	return s
}

func (boltdb *BoltDB) apiList(bucket string) ([]byte, error) {
	if !boltdb.open {
		return nil, fmt.Errorf("db must be opened before saving!")
	}
	var jsonSlice []byte
	err := boltdb.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(bucket)).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			jsonSlice = append(jsonSlice, v...)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Could not get Bucket: %s", bucket)
	}
	jsonSlice = fmtToJsonArr(jsonSlice)
	return jsonSlice, nil
}

func (boltdb *BoltDB) apiListPrefix(bucket, prefix string) ([]byte, error) {
	if !boltdb.open {
		return nil, fmt.Errorf("db must be opened before saving!")
	}
	var jsonSlice []byte
	err := boltdb.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(bucket)).Cursor()
		p := []byte(prefix)
		for k, v := c.Seek(p); bytes.HasPrefix(k, p); k, v = c.Next() {
			jsonSlice = append(jsonSlice, v...)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Could not get Bucket: %s", bucket)
	}
	jsonSlice = fmtToJsonArr(jsonSlice)
	return jsonSlice, nil
}

func (boltdb *BoltDB) apiListRange(bucket, start, stop string) ([]byte, error) {
	if !boltdb.open {
		return nil, fmt.Errorf("db must be opened before saving!")
	}
	var jsonSlice []byte
	err := boltdb.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(bucket)).Cursor()
		min := []byte(start)
		max := []byte(stop)
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			jsonSlice = append(jsonSlice, v...)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Could not get Bucket: %s. Error: %s", bucket, err)
	}
	jsonSlice = fmtToJsonArr(jsonSlice)
	return jsonSlice, nil
}

func (boltdb *BoltDB) apiGetOne(bucket, key string) ([]byte, error) {
	if !boltdb.open {
		return nil, fmt.Errorf("db must be opened before saving!")
	}
	var data []byte
	err := boltdb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		k := []byte(key)
		data = b.Get(k)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Could not get Bucket: %s, Key: %s. Error: %s", bucket, key, err)
	}
	return data, nil
}
