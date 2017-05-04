package store

import (
	"github.com/google/btree"
	"github.com/tonjun/pubsub"
)

type treeItem struct {
	Key   uint64 // Connection ID
	Value pubsub.Conn
}

func (a treeItem) Less(b btree.Item) bool {
	return a.Key < b.(treeItem).Key
}
