package jolt

import (
	"github.com/boltdb/bolt"
)

var boltdbrw *bolt.DB
var dbrw = apiConn(boltdbrw)

func Open(dir string) error {
	return dbrw.apiOpen(dir)
}

func Copy(dir string) error {
	return dbrw.apiCopyDB(dir)
}

func Close() {
	dbrw.apiClose()
}

func Save(bucket, key string, data interface{}) error {
	return dbrw.apiSave(bucket, key, data)
}

func List(bucket string) ([]byte, error) {
	return dbrw.apiList(bucket)
}

func ListPrefix(bucket, prefix string) ([]byte, error) {
	return dbrw.apiListPrefix(bucket, prefix)
}

func ListRange(bucket, start, stop string) ([]byte, error) {
	return dbrw.apiListRange(bucket, start, stop)
}

func GetOne(bucket, key string) ([]byte, error) {
	return dbrw.apiGetOne(bucket, key)
}
