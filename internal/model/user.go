package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/loveuer/nf/nft/log"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"strings"
	"time"
	"ultone/internal/opt"
	"ultone/internal/sqlType"
	"ultone/internal/tool"
)

var (
	initUsers = []*User{
		{
			Id:       1,
			Username: "admin",
			Password: tool.NewPassword("123456"),
			Nickname: "admin",
			Role:     RoleAdmin,
			Privileges: lo.Map(Privilege(0).All(), func(item Enum, index int) Privilege {
				return item.(Privilege)
			}),
			CreatedById:   1,
			CreatedByName: "admin",
			ActiveAt:      time.Now().UnixMilli(),
			Deadline:      time.Now().AddDate(100, 0, 0).UnixMilli(),
		},
	}

	_ Enum = Status(0)
)

type Status uint64

const (
	StatusNormal Status = iota
	StatusFrozen
)

func (s Status) Value() int64 {
	return int64(s)
}

func (s Status) Code() string {
	switch s {
	case StatusNormal:
		return "normal"
	case StatusFrozen:
		return "frozen"
	default:
		return "unknown"
	}
}

func (s Status) Label() string {
	switch s {
	case StatusNormal:
		return "正常"
	case StatusFrozen:
		return "冻结"
	default:
		return "异常"
	}
}

func (s Status) All() []Enum {
	return []Enum{
		StatusNormal,
		StatusFrozen,
	}
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"value": s.Value(),
		"code":  s.Code(),
		"label": s.Label(),
	})
}

type User struct {
	Id        uint64 `json:"id" gorm:"primaryKey;column:id"`
	CreatedAt int64  `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
	DeletedAt int64  `json:"deleted_at" gorm:"index;column:deleted_at;default:0"`

	Username string `json:"username" gorm:"column:username;type:varchar(64);unique"`
	Password string `json:"-" gorm:"column:password;type:varchar(256)"`

	Status Status `json:"status" gorm:"column:status;default:0"`

	Nickname string `json:"nickname" gorm:"column:nickname;type:varchar(64)"`
	Comment  string `json:"comment" gorm:"column:comment"`

	Role       Role                        `json:"role" gorm:"column:role"`
	Privileges sqlType.NumSlice[Privilege] `json:"privileges" gorm:"column:privileges;type:bigint[]"`

	CreatedById   uint64 `json:"created_by_id" gorm:"column:created_by_id"`
	CreatedByName string `json:"created_by_name" gorm:"column:created_by_name;type:varchar(64)"`

	ActiveAt int64 `json:"active_at" gorm:"column:active_at"`
	Deadline int64 `json:"deadline" gorm:"column:deadline"`

	LoginAt int64 `json:"login_at" gorm:"-"`
}

func (u *User) CheckStatus(mustOk bool) error {
	switch u.Status {
	case StatusNormal:
	case StatusFrozen:
		if mustOk {
			return errors.New("用户被冻结")
		}
	default:
		return errors.New("用户状态未知")
	}

	return nil
}

func (u *User) IsValid(mustOk bool) error {
	now := time.Now()

	if now.UnixMilli() >= u.Deadline {
		return errors.New("用户已过期")
	}

	if now.UnixMilli() < u.ActiveAt {
		return errors.New("用户未启用")
	}

	if u.DeletedAt > 0 {
		return errors.New("用户不存在")
	}

	return u.CheckStatus(mustOk)
}

func (u *User) JwtEncode() (token string, err error) {

	now := time.Now()

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id":       u.Id,
		"username": u.Username,
		"status":   u.Status,
		"deadline": u.Deadline,
		"login_at": now.UnixMilli(),
	})

	if token, err = jwtToken.SignedString([]byte(opt.JwtTokenSecret)); err != nil {
		err = fmt.Errorf("JwtEncode: jwt token signed secret err: %v", err)
		log.Error(err.Error())
		return "", nil
	}

	return
}

func (u *User) FromJwt(token string) *User {
	var (
		ok     bool
		err    error
		pt     *jwt.Token
		claims jwt.MapClaims
	)

	token = strings.TrimPrefix(token, "Bearer ")

	if pt, err = jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok = t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(opt.JwtTokenSecret), nil
	}); err != nil {
		log.Error("jwt parse err: %v", err)
		return nil
	}

	if !pt.Valid {
		log.Warn("parsed jwt invalid")
		return nil
	}

	if claims, ok = pt.Claims.(jwt.MapClaims); !ok {
		log.Error("convert jwt claims err")
		return nil
	}

	u.Id = cast.ToUint64(claims["user_id"])
	u.Username = cast.ToString(claims["username"])
	u.Status = Status(cast.ToInt64(claims["status"]))
	u.Deadline = cast.ToInt64(claims["deadline"])
	u.LoginAt = cast.ToInt64(claims["login_at"])

	return u
}

func (u User) MarshalBinary() ([]byte, error) {
	return json.Marshal(map[string]any{
		"id":         u.Id,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
		"deleted_at": u.DeletedAt,
		"username":   u.Username,
		"status":     u.Status.Value(),
		"nickname":   u.Nickname,
		"comment":    u.Comment,
		"role":       uint8(u.Role),
		"privileges": lo.Map(u.Privileges, func(item Privilege, index int) int64 {
			return item.Value()
		}),
		"created_by_id":   u.CreatedById,
		"created_by_name": u.CreatedByName,
		"active_at":       u.ActiveAt,
		"deadline":        u.Deadline,
		"login_at":        u.LoginAt,
	})
}
