package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"ultone/internal/opt"
)

type _role struct {
	Value uint8  `json:"value"`
	Code  string `json:"code"`
	Label string `json:"label"`
}

type Role uint8

var _ Enum = Role(0)

func (u Role) MarshalJSON() ([]byte, error) {
	m := _role{
		Value: uint8(u),
		Code:  u.Code(),
		Label: u.Label(),
	}
	return json.Marshal(m)
}

const (
	RoleRoot  Role = 255
	RoleAdmin Role = 254
	RoleUser  Role = 100
)

func (u Role) Code() string {
	switch u {
	case RoleRoot:
		return "root"
	case RoleAdmin:
		return "admin"
	case RoleUser:
		return "user"
	default:
		return "unknown"
	}
}

func (u Role) Label() string {
	switch u {
	case RoleRoot:
		return "根用户"
	case RoleAdmin:
		return "管理员"
	case RoleUser:
		return "用户"
	default:
		return "未知"
	}
}

func (u Role) Value() int64 {
	return int64(u)
}

func (u Role) All() []Enum {
	return []Enum{
		RoleAdmin,
		RoleUser,
	}
}

func (u Role) Where(db *gorm.DB) *gorm.DB {
	if opt.RoleMustLess {
		return db.Where("users.role < ?", u.Value())
	} else {
		return db.Where("users.role <= ?", u.Value())
	}
}

func (u Role) CanOP(op *User) bool {
	if opt.RoleMustLess {
		return op.Role > u
	}

	return op.Role >= u
}
