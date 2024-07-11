package opt

import "time"

const (
	// todo: 可以替换自己生生成的 secret
	JwtTokenSecret = "7^D+UW3BPB2Mnz)bY3uVrAUyv&dj8Kdz"

	// todo: 是否打开 gorm 的 debug 打印 (开发和 dev 环境时可以打开)
	DBDebug = true

	// todo: 是否加载默认的前端用户管理界面
	EnableFront = false

	// todo: 同一个账号是否可以多 client 登录
	MultiLogin = false

	// todo: 用户量不大的情况, 并没有缓存用户具体信息, 如果需要可以打开
	EnableUserCache = true

	// todo: 缓存时, key 的前缀
	CachePrefix = "ultone"

	// todo: 登录颁发的 cookie 的 name
	CookieName = "utlone-token"

	// todo: 用户列表,日志列表 size 参数
	DefaultSize, MaxSize = 20, 200

	// todo: 操作用户时, role 相等时能否操作: 包括 列表, 能否新建,修改,删除同样 role 的用户
	RoleMustLess = false

	// todo: 通过 c.Local() 存入 oplog 时的 key 值
	OpLogLocalKey = "oplog"

	// todo: 操作日志 最多延迟多少秒写入（最多缓存多少秒的日志，然后 bulk 写入)
	OpLogWriteDurationSecond = 5

	LocalTraceKey      = "X-Trace"
	LocalShortTraceKey = "X-Short-Trace"
)

var (
	// todo: 颁发的 token, (cookie) 在缓存中存在的时间 (每次请求该时间也会被刷新)
	TokenTimeout = time.Duration(3600*12) * time.Second
)
