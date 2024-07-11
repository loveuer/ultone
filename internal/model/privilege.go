package model

import "encoding/json"

type Privilege uint64

type _privilege struct {
	Value int64  `json:"value"`
	Code  string `json:"code"`
	Label string `json:"label"`
}

const (
	PrivilegeUserManage Privilege = iota + 1
	PrivilegeOpLog
)

func (p Privilege) Value() int64 {
	return int64(p)
}

func (p Privilege) Code() string {
	switch p {
	case PrivilegeUserManage:
		return "user_manage"
	case PrivilegeOpLog:
		return "oplog"
	default:
		return "unknown"
	}
}

func (p Privilege) Label() string {
	switch p {
	case PrivilegeUserManage:
		return "用户管理"
	case PrivilegeOpLog:
		return "操作日志"
	default:
		return "未知"
	}
}

func (p Privilege) MarshalJSON() ([]byte, error) {
	_p := &_privilege{
		Value: int64(p),
		Code:  p.Code(),
		Label: p.Label(),
	}

	return json.Marshal(_p)
}

func (p Privilege) All() []Enum {
	return []Enum{
		PrivilegeUserManage,
		PrivilegeOpLog,
	}
}

var _ Enum = (*Privilege)(nil)
