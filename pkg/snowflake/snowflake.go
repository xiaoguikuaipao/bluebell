package snowflake

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func Init(startTime string, machinID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2023-05-03", startTime)
	if err != nil {
		return err
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(machinID)
	return err
}

func GenID() int64 {
	return node.Generate().Int64()
}
