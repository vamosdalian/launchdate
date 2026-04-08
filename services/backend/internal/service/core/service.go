package core

import (
	"github.com/bwmarrin/snowflake"
	"github.com/vamosdalian/launchdate-backend/internal/db"
)

type MainService struct {
	mc *db.MongoDB
	sn *snowflake.Node
}

func NewMainService(mc *db.MongoDB) *MainService {
	node, _ := snowflake.NewNode(0)
	return &MainService{
		mc: mc,
		sn: node,
	}
}
