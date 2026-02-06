package snowflake

import (
	"os"
	"strconv"
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	node *snowflake.Node
	once sync.Once
)

// Init 初始化雪花算法节点。
// 一般在 main 启动时调用；如果没调用，Gen* 会使用默认 machineID(1) 做兜底初始化。
func Init(machineID int64) error {
	var err error
	once.Do(func() {
		node, err = snowflake.NewNode(machineID)
	})
	return err
}

func ensureInit() {
	if node != nil {
		return
	}
	// allow env override
	machineID := int64(1)
	if v := os.Getenv("MINICHAT_SNOWFLAKE_MACHINE_ID"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			machineID = n
		}
	}
	_ = Init(machineID)
}

// GenStringID returns a unique snowflake ID in string form.
func GenStringID() string {
	ensureInit()
	return node.Generate().String()
}

// GenInt64ID returns a unique snowflake ID as int64.
func GenInt64ID() int64 {
	ensureInit()
	return node.Generate().Int64()
}
