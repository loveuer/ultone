package model

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/spf13/cast"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	"html/template"
	"time"
	"ultone/internal/sqlType"
)

var (
	FuncMap = template.FuncMap{
		"time_format": func(mil any, format string) string {
			return time.UnixMilli(cast.ToInt64(mil)).Format(format)
		},
	}
)

var (
	_ OpLogger = (*OpLogType)(nil)
)

type OpLogType uint64

const (
	OpLogTypeLogin OpLogType = iota + 1
	OpLogTypeLogout
	OpLogTypeCreateUser
	OpLogTypeUpdateUser
	OpLogTypeDeleteUser

	// todo: 添加自己的操作日志 分类
)

func (o OpLogType) Value() int64 {
	return int64(o)
}

func (o OpLogType) Code() string {
	switch o {
	case OpLogTypeLogin:
		return "login"
	case OpLogTypeLogout:
		return "logout"
	case OpLogTypeCreateUser:
		return "create_user"
	case OpLogTypeUpdateUser:
		return "update_user"
	case OpLogTypeDeleteUser:
		return "delete_user"
	default:
		return "unknown"
	}
}

func (o OpLogType) Label() string {
	switch o {
	case OpLogTypeLogin:
		return "登入"
	case OpLogTypeLogout:
		return "登出"
	case OpLogTypeCreateUser:
		return "创建用户"
	case OpLogTypeUpdateUser:
		return "修改用户"
	case OpLogTypeDeleteUser:
		return "删除用户"
	default:
		return "未知"
	}
}

func (o OpLogType) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"value": o.Value(),
		"code":  o.Code(),
		"label": o.Label(),
	})
}

func (o OpLogType) All() []Enum {
	return []Enum{
		OpLogTypeLogin,
		OpLogTypeLogout,
		OpLogTypeCreateUser,
		OpLogTypeUpdateUser,
		OpLogTypeDeleteUser,
	}
}

func _trimHTML(v []byte) string {
	return base64.StdEncoding.EncodeToString(v)
}

var (
	_mini = minify.New()
)

func init() {
	_mini.AddFunc("text/html", html.Minify)
}

func (o OpLogType) Render(content map[string]any) (string, error) {
	var (
		err    error
		render *template.Template
		buf    bytes.Buffer
		bs     []byte
	)

	if render, err = template.New(o.Code()).
		Funcs(FuncMap).
		Parse(o.Template()); err != nil {
		return "", err
	}

	if err = render.Execute(&buf, content); err != nil {
		return "", err
	}

	if bs, err = _mini.Bytes("text/html", buf.Bytes()); err != nil {
		return "", err
	}

	return _trimHTML(bs), nil
}

const (
	oplogTemplateLogin = `
		<div class="nf-op-log">
			用户
			<span 
				class="nf-op-log-user nf-op-log-keyword"
				nf-op-log-user="{{ .user_id }}"
			>{{ .username }}
			</span>
			于
			<span 
				class="nf-op-log-time nf-op-log-keyword"
				nf-op-log-time="{{ .time }}"
			>{{ time_format .time "2006-01-02 15:04:05" }}
			</span>	
			在
			<span
				class="nf-op-log-ip nf-op-log-keyword"
			>{{ .ip }}
			</span>
			上
			<span
				class="nf-op-log-op nf-op-log-keyword"
			>
			登入
			</span>
			了系统
		</div>
	`
	oplogTemplateLogout = `
		<div class="nf-op-log">
			用户
			<span 
				class="nf-op-log-user nf-op-log-keyword" 
				nf-op-log-user="{{ .user_id }}"
			>{{ .username }}
			</span>
			于
			<span 
				class="nf-op-log-time nf-op-log-keyword"
				nf-op-log-time="{{ .time }}"
			>{{ time_format .time "2006-01-02 15:04:05" }}
			</span>	
			在
			<span
				class="nf-op-log-ip nf-op-log-keyword"
			>{{ .ip }}
			</span>
			上
			<span
				class="nf-op-log-op nf-op-log-keyword"
			>
			登出
			</span>
			了系统
		</div>
`
	oplogTemplateCreateUser = `
		<div class="nf-op-log">
			用户
			<span 
				class="nf-op-log-user nf-op-log-keyword" 
				nf-op-log-user="{{ .user_id }}"
			>{{ .username }}
			</span>
			于
			<span 
				class="nf-op-log-time nf-op-log-keyword"
				nf-op-log-time="{{ .time }}"
			>{{ time_format .time "2006-01-02 15:04:05" }}
			</span>	
			<span class="nf-op-log-keyword">
			创建
			</span>
			了用户
			<span
				class="nf-op-log-target nf-op-log-keyword"
				nf-op-log-target="{{ .target_id }}"
			>{{ .target_username }}
			</span>
		</div>
`
	oplogTemplateUpdateUser = `
		<div class="nf-op-log">
			用户
			<span 
				class="nf-op-log-user nf-op-log-keyword" 
				nf-op-log-user='{{ .user_id }}'
			>{{ .username }}
			</span>
			于
			<span 
				class="nf-op-log-time nf-op-log-keyword"
				nf-op-log-time='{{ .time }}'
			>{{ time_format .time "2006-01-02 15:04:05" }}
			</span>	
			<span class="nf-op-log-keyword">
			编辑	
			</span>
			了用户
			<span
				class="nf-op-log-target nf-op-log-keyword"
				nf-op-log-target="{{ .target_id }}"
			>{{ .target_username }}
			</span>
		</div>
`
	oplogTemplateDeleteUser = `
		<div class="nf-op-log">
			用户
			<span 
				class="nf-op-log-user nf-op-log-keyword" 
				nf-op-log-user="{{ .user_id }}"
			>{{ .username }}
			</span>
			于
			<span 
				class="nf-op-log-time nf-op-log-keyword"
				nf-op-log-time="{{ .time }}"
			>{{ time_format .time "2006-01-02 15:04:05" }}
			</span>	
			<span class="nf-op-log-keyword">
			删除	
			</span>
			了用户
			<span
				class="nf-op-log-target nf-op-log-keyword"
				nf-op-log-target="{{ .target_id }}"
			>{{ .target_username }}
			</span>
		</div>
`
)

func (o OpLogType) Template() string {
	switch o {
	case OpLogTypeLogin:
		return oplogTemplateLogin
	case OpLogTypeLogout:
		return oplogTemplateLogout
	case OpLogTypeCreateUser:
		return oplogTemplateCreateUser
	case OpLogTypeUpdateUser:
		return oplogTemplateUpdateUser
	case OpLogTypeDeleteUser:
		return oplogTemplateDeleteUser
	default:
		return `<div>错误的日志类型</div>`
	}
}

type OpLog struct {
	Id        uint64 `json:"id" gorm:"primaryKey;column:id"`
	CreatedAt int64  `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime:milli"`
	DeletedAt int64  `json:"deleted_at" gorm:"index;column:deleted_at;default:0"`

	UserId   uint64        `json:"user_id" gorm:"column:user_id"`
	Username string        `json:"username" gorm:"column:username;varchar(128)"`
	Type     OpLogType     `json:"type" gorm:"column:type;type:varchar(128)"`
	Content  sqlType.JSONB `json:"content" gorm:"column:content;type:jsonb"`
	HTML     string        `json:"html" gorm:"-"`
}
