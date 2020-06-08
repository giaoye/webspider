package snowflake

import (
  "github.com/bwmarrin/snowflake"
  "sync"
)

const NODEN int64 = 111

var (
  once sync.Once
  node *snowflake.Node
)
func NewNode() (*snowflake.Node, error) {
  return snowflake.NewNode(NODEN)
}

func GetNode() (*snowflake.Node, error) {
  var err error
  once.Do(func() {
    node, err = NewNode()
  })
  return node, err
}