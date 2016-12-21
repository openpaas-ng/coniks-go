package kvstorage

import (
	"fmt"

	"github.com/coniks-sys/coniks-go/protocol/extensions/extension"
	"github.com/coniks-sys/coniks-go/storage/kv"
	"github.com/coniks-sys/coniks-go/storage/kv/leveldbkv"
)

func init() {
	extension.RegisterSnapshotStorage(SnapshotKVID, New)
}

const SnapshotKVID = "SnapshotKV"

type SnapshotKV struct {
	db kv.DB
}

var _ extension.SnapshotStorage = (*SnapshotKV)(nil)

func New(path string) (extension.SnapshotStorage, error) {
	db := leveldbkv.OpenDB(path)
	return &SnapshotKV{db: db}, nil
}

func (ss *SnapshotKV) StoreDirectory() {
	fmt.Println("hello")
}
