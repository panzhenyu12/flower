package utils

import (
	"sync"

	"github.com/golang/glog"

	"github.com/panzhenyu12/flower/config"

	"github.com/bwmarrin/snowflake"
)

var idonce sync.Once
var idgenerate *IDGenerate

type IDGenerate struct {
	node *snowflake.Node
}

func GetIDGenerate() *IDGenerate {
	idonce.Do(func() {
		nodeid := config.GetConfig().WorkID
		node, err := snowflake.NewNode(nodeid)
		if err != nil {
			glog.Error(err)
			return
		}
		idgenerate = new(IDGenerate)
		idgenerate.node = node
	})
	return idgenerate
}

//GetID 生成ID
func (generate *IDGenerate) GetID() snowflake.ID {
	return generate.node.Generate()
}
