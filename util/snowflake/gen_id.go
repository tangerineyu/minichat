package snowflake

import (
	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

// Init 初始化雪花算法节点
func Init(machineID int64) error {
	var err error
	node, err = snowflake.NewNode(machineID)
	if err != nil {
		return err
	}
	return nil
}

func GenStringID() string {
	return node.Generate().String()
}
