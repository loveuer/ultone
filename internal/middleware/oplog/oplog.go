package oplog

import "ultone/internal/model"

type OpLog struct {
	Type    model.OpLogType
	Content map[string]any
}
