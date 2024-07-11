package oplog

import (
	"context"
	"github.com/loveuer/nf"
	"github.com/loveuer/nf/nft/log"
	"sync"
	"time"
	"ultone/internal/database/db"
	"ultone/internal/model"
	"ultone/internal/opt"
	"ultone/internal/sqlType"
	"ultone/internal/tool"
)

var (
	_once = &sync.Once{}
	lc    = make(chan *model.OpLog, 1024)
)

// NewOpLog
//
// * 记录操作日志的 中间件使用方法如下:
//
//	 app := nf.New()
//	 app.Post("/login", oplog.NewOpLog(ctx), HandleLog)
//
//	 func HandleLog(c *nf.Ctx) error {
//		// 你的操作逻辑
//		c.Local(opt.OpLogLocalKey, &oplog.OpLog{})
//		// 剩下某些逻辑
//		// return xxx
//	 }
func NewOpLog(ctx context.Context) nf.HandlerFunc {

	_once.Do(func() {
		go func() {
			var (
				err    error
				ticker = time.NewTicker(time.Duration(opt.OpLogWriteDurationSecond) * time.Second)
				list   = make([]*model.OpLog, 0, 1024)

				write = func() {
					if len(list) == 0 {
						return
					}

					if err = db.New(tool.Timeout(10)).
						Model(&model.OpLog{}).
						Create(&list).
						Error; err != nil {
						log.Error("middleware.NewOpLog: write logs err=%v", err)
					}

					list = list[:0]
				}
			)

		Loop:
			for {
				select {
				case <-ctx.Done():
					break Loop
				case <-ticker.C:
					write()
				case item, ok := <-lc:
					if !ok {
						return
					}

					list = append(list, item)

					if len(list) >= 100 {
						write()
					}
				}
			}

			write()
		}()
	})

	return func(c *nf.Ctx) error {
		now := time.Now()

		err := c.Next()

		op, ok := c.Locals("user").(*model.User)

		opv := c.Locals(opt.OpLogLocalKey)
		logItem, ok := opv.(*OpLog)
		if !ok {
			log.Warn("middleware.NewOpLog: %s - %s local '%s' to [*OpLog] invalid", c.Method(), c.Path(), opt.OpLogLocalKey)
			return err
		}

		logItem.Content["time"] = now.UnixMilli()
		logItem.Content["user_id"] = op.Id
		logItem.Content["username"] = op.Username
		logItem.Content["created_at"] = now.UnixMilli()

		select {
		case lc <- &model.OpLog{
			CreatedAt: now.UnixMilli(),
			UpdatedAt: now.UnixMilli(),
			UserId:    op.Id,
			Username:  op.Username,
			Type:      logItem.Type,
			Content:   sqlType.NewJSONB(logItem.Content),
		}:
		case <-tool.Timeout(3).Done():
			log.Warn("middleware.NewOpLog: %s - %s log -> chan timeout[3s]", c.Method, c.Path())
		}

		return err
	}
}
