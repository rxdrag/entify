package utils

import (
	"github.com/bwmarrin/snowflake"
	"rxdrag.com/entity-engine/config"
)

func CreateId() snowflake.ID {
	node, err := snowflake.NewNode(config.NODE_ID)
	if err != nil {
		panic("Create Id error:" + err.Error())
	}

	// Generate a snowflake ID.
	return node.Generate()
}