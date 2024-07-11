package handler

import (
	"github.com/loveuer/nf"
	"github.com/loveuer/nf/nft/resp"
	"ultone/internal/database/db"
	"ultone/internal/log"
	"ultone/internal/model"
	"ultone/internal/opt"
	"ultone/internal/sqlType"
	"ultone/internal/tool"
)

func LogCategories() nf.HandlerFunc {
	return func(c *nf.Ctx) error {
		return resp.Resp200(c, model.OpLogType(0).All())
	}
}

func LogList(c *nf.Ctx) error {
	type Req struct {
		Page    int                               `query:"page"`
		Size    int                               `query:"size"`
		UserIds []uint64                          `query:"user_ids"`
		Types   sqlType.NumSlice[model.OpLogType] `query:"types"`
	}

	var (
		ok    bool
		op    *model.User
		err   error
		req   = new(Req)
		list  = make([]*model.OpLog, 0)
		total int
	)

	if op, ok = c.Locals("user").(*model.User); !ok {
		return resp.Resp401(c, nil)
	}

	if err = c.QueryParser(req); err != nil {
		return resp.Resp400(c, err.Error())
	}

	if req.Size <= 0 {
		req.Size = opt.DefaultSize
	}

	if req.Size > opt.MaxSize {
		return resp.Resp400(c, req, "参数过大")
	}

	txCount := op.Role.Where(db.New(tool.Timeout(3)).
		Model(&model.OpLog{}).
		Select("COUNT(`op_logs`.`id`)").
		Joins("LEFT JOIN users ON `users`.`id` = `op_logs`.`user_id`"))
	txGet := op.Role.Where(db.New(tool.Timeout(10)).
		Model(&model.OpLog{}).
		Joins("LEFT JOIN users ON `users`.`id` = `op_logs`.`user_id`"))

	if len(req.UserIds) != 0 {
		txCount = txCount.Where("op_logs.user_id IN ?", req.UserIds)
		txGet = txGet.Where("op_logs.user_id IN ?", req.UserIds)
	}

	if len(req.Types) != 0 {
		txCount = txCount.Where("op_logs.type IN ?", req.Types)
		txGet = txGet.Where("op_logs.type IN ?", req.Types)
	}

	if err = txCount.
		Find(&total).
		Error; err != nil {
		return resp.Resp500(c, err.Error())
	}

	if err = txGet.
		Offset(req.Page * req.Size).
		Limit(req.Size).
		Order("`op_logs`.`created_at` DESC").
		Find(&list).
		Error; err != nil {
		return resp.Resp500(c, err.Error())
	}

	for _, logItem := range list {
		m := make(map[string]any)
		if err = logItem.Content.Bind(&m); err != nil {
			log.Warn(c, "handler.LogList: log=%d content=%v bind map[string]any err=%v", logItem.Id, logItem.Content, err)
			continue
		}

		if logItem.HTML, err = logItem.Type.Render(m); err != nil {
			log.Warn(c, "handler.LogList: log=%d template=%s render map=%+v err=%v", logItem.Id, logItem.Type.Template(), m, err)
			continue
		}

		log.Debug(c, "handler.LogList: log=%d render map=%+v string=%s", logItem.Id, m, logItem.HTML)
	}

	return resp.Resp200(c, nf.Map{"list": list, "total": total})
}
