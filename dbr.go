package joltDB

import (
	"github.com/boltdb/bolt"
)

var boltdbr *bolt.DB
var dbr = apiConn(boltdbr)

func OpenReadOnly(dir string) error {
	return dbr.apiOpenReadOnly(dir)
}

func CloseReadOnly() {
	dbr.apiClose()
}

func ListReadOnly(bucket string) ([]byte, error) {
	return dbr.apiList(bucket)
}

func ListPrefixReadOnly(bucket, prefix string) ([]byte, error) {
	return dbrw.apiListPrefix(bucket, prefix)
}

func ListRangeReadOnly(bucket, start, stop string) ([]byte, error) {
	return dbrw.apiListRange(bucket, start, stop)
}

func GetOneReadOnly(bucket, key string) ([]byte, error) {
	return dbr.apiGetOne(bucket, key)
}
