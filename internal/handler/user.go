package handler

import (
	"errors"
	"fmt"
	"github.com/loveuer/nf"
	"github.com/loveuer/nf/nft/resp"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
	"ultone/internal/controller"
	"ultone/internal/database/cache"
	"ultone/internal/database/db"
	"ultone/internal/middleware/oplog"
	"ultone/internal/model"
	"ultone/internal/opt"
	"ultone/internal/sqlType"
	"ultone/internal/tool"
)

func AuthLogin(c *nf.Ctx) error {
	type Req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var (
		err    error
		req    = new(Req)
		target = new(model.User)
		token  string
		now    = time.Now()
	)

	if err = c.BodyParser(req); err != nil {
		return resp.Resp400(c, err.Error())
	}

	if err = db.New(tool.Timeout(3)).
		Model(&model.User{}).
		Where("username = ?", req.Username).
		Where("deleted_at = 0").
		Take(target).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp.Resp400(c, err.Error(), "用户名或密码错误")
		}

		return resp.Resp500(c, err.Error())
	}

	if !tool.ComparePassword(req.Password, target.Password) {
		return resp.Resp400(c, nil, "用户名或密码错误")
	}

	if err = target.IsValid(true); err != nil {
		return resp.Resp401(c, nil, err.Error())
	}

	if err = controller.UserController.CacheUser(c, target); err != nil {
		return resp.RespError(c, err)
	}

	if token, err = target.JwtEncode(); err != nil {
		return resp.Resp500(c, err.Error())
	}

	if err = controller.UserController.CacheToken(c, token, target); err != nil {
		return resp.RespError(c, err)
	}

	if !opt.MultiLogin {
		var (
			last = fmt.Sprintf("%s:user:last_token:%d", opt.CachePrefix, target.Id)
			bs   []byte
		)

		if bs, err = cache.Client.Get(tool.Timeout(3), last); err == nil {
			key := fmt.Sprintf("%s:user:token:%s", opt.CachePrefix, string(bs))
			_ = cache.Client.Del(tool.Timeout(3), key)
		}

		if err = cache.Client.Set(tool.Timeout(3), last, token); err != nil {
			return resp.Resp500(c, err.Error())
		}
	}

	c.Set("Set-Cookie", fmt.Sprintf("%s=%s; Path=/", opt.CookieName, token))
	c.Locals("user", target)
	c.Locals(opt.OpLogLocalKey, &oplog.OpLog{Type: model.OpLogTypeLogin, Content: map[string]any{
		"time": now.UnixMilli(),
		"ip":   c.IP(true),
	}})

	return resp.Resp200(c, nf.Map{"token": token, "user": target})
}

func AuthVerify(c *nf.Ctx) error {
	op, ok := c.Locals("user").(*model.User)
	if !ok {
		return resp.Resp401(c, nil)
	}

	token, ok := c.Locals("token").(string)
	if !ok {
		return resp.Resp401(c, nil)
	}

	return resp.Resp200(c, nf.Map{"token": token, "user": op})
}

func AuthLogout(c *nf.Ctx) error {
	op, ok := c.Locals("user").(*model.User)
	if !ok {
		return resp.Resp401(c, nil)
	}

	_ = controller.UserController.RmUserCache(c, op.Id)

	c.Locals(opt.OpLogLocalKey, &oplog.OpLog{
		Type: model.OpLogTypeLogout,
		Content: map[string]any{
			"time": time.Now().UnixMilli(),
			"ip":   c.IP(),
		},
	})

	return resp.Resp200(c, nil)
}

func ManageUserList(c *nf.Ctx) error {
	type Req struct {
		Page    int    `query:"page"`
		Size    int    `query:"size"`
		Keyword string `query:"keyword"`
	}

	var (
		err   error
		ok    bool
		op    *model.User
		req   = new(Req)
		list  = make([]*model.User, 0)
		total = 0
	)

	if op, ok = c.Locals("user").(*model.User); !ok {
		return resp.Resp401(c, nil)
	}

	if err = c.QueryParser(req); err != nil {
		return resp.Resp400(c, err.Error())
	}

	if req.Size == 0 {
		req.Size = opt.DefaultSize
	}

	if req.Size > opt.MaxSize {
		return resp.Resp400(c, nf.Map{"msg": "size over max", "max": opt.MaxSize})
	}

	if err = op.Role.Where(db.New(tool.Timeout(10)).
		Model(&model.User{}).
		Where("deleted_at = 0")).
		Order("updated_at DESC").
		Offset(req.Page * req.Size).
		Limit(req.Size).
		Find(&list).
		Error; err != nil {
		return resp.Resp500(c, err.Error())
	}

	if err = op.Role.Where(db.New(tool.Timeout(5)).
		Model(&model.User{}).
		Select("COUNT(id)").
		Where("deleted_at = 0")).
		Find(&total).
		Error; err != nil {
		return resp.Resp500(c, err.Error())
	}

	return resp.Resp200(c, nf.Map{"list": list, "total": total})
}

func ManageUserCreate(c *nf.Ctx) error {
	type Req struct {
		Username   string                            `json:"username"`
		Nickname   string                            `json:"nickname"`
		Password   string                            `json:"password"`
		Status     model.Status                      `json:"status"`
		Role       model.Role                        `json:"role"`
		Privileges sqlType.NumSlice[model.Privilege] `json:"privileges"`
		Comment    string                            `json:"comment"`
		ActiveAt   int64                             `json:"active_at"`
		Deadline   int64                             `json:"deadline"`
	}

	var (
		err error
		ok  bool
		op  *model.User
		req = new(Req)
		now = time.Now()
	)

	if op, ok = c.Locals("user").(*model.User); !ok {
		return resp.Resp401(c, nil)
	}

	if err = c.BodyParser(req); err != nil {
		return resp.Resp400(c, err.Error())
	}

	if req.Username == "" || req.Password == "" {
		return resp.Resp400(c, req)
	}

	if err = tool.CheckPassword(req.Password); err != nil {
		return resp.Resp400(c, req, err.Error())
	}

	if req.Nickname == "" {
		req.Nickname = req.Username
	}

	if req.Status.Code() == "unknown" {
		return resp.Resp400(c, req, "用户状态不正常")
	}

	if req.Role == 0 {
		req.Role = model.RoleUser
	}

	if req.ActiveAt == 0 {
		req.ActiveAt = now.UnixMilli()
	}

	if req.Deadline == 0 {
		req.Deadline = now.AddDate(99, 0, 0).UnixMilli()
	}

	newUser := &model.User{
		CreatedAt:     now.UnixMilli(),
		UpdatedAt:     now.UnixMilli(),
		Username:      req.Username,
		Password:      tool.NewPassword(req.Password),
		Status:        req.Status,
		Nickname:      req.Nickname,
		Comment:       req.Comment,
		Role:          req.Role,
		Privileges:    req.Privileges,
		CreatedById:   op.Id,
		CreatedByName: op.CreatedByName,
		ActiveAt:      op.ActiveAt,
		Deadline:      op.Deadline,
	}

	if err = newUser.IsValid(false); err != nil {
		return resp.Resp400(c, newUser, err.Error())
	}

	if !newUser.Role.CanOP(op) {
		return resp.Resp403(c, newUser, "角色不符合")
	}

	if err = db.New(tool.Timeout(5)).
		Create(newUser).
		Error; err != nil {
		return resp.Resp500(c, err.Error())
	}

	c.Locals(opt.OpLogLocalKey, &oplog.OpLog{Type: model.OpLogTypeCreateUser, Content: map[string]any{
		"target_id":       newUser.Id,
		"target_username": newUser.Username,
		"target_nickname": newUser.Nickname,
		"target_status":   newUser.Status.Label(),
		"target_role":     newUser.Role.Label(),
		"target_privileges": lo.Map(newUser.Privileges, func(item model.Privilege, index int) string {
			return item.Label()
		}),
		"target_active_at": op.ActiveAt,
		"target_deadline":  op.Deadline,
	}})

	return resp.Resp200(c, newUser)
}

func ManageUserUpdate(c *nf.Ctx) error {
	type Req struct {
		Id         uint64                            `json:"id"`
		Nickname   string                            `json:"nickname"`
		Password   string                            `json:"password"`
		Status     model.Status                      `json:"status"`
		Comment    string                            `json:"comment"`
		Role       model.Role                        `json:"role"`
		Privileges sqlType.NumSlice[model.Privilege] `json:"privileges"`
		ActiveAt   int64                             `json:"active_at"`
		Deadline   int64                             `json:"deadline"`
	}

	type Change struct {
		Old any `json:"old"`
		New any `json:"new"`
	}

	var (
		ok      bool
		op      *model.User
		target  *model.User
		err     error
		req     = new(Req)
		rm      = make(map[string]any)
		updates = make(map[string]any)
		changes = make(map[string]Change)
	)

	if op, ok = c.Locals("user").(*model.User); !ok {
		return resp.Resp401(c, nil)
	}

	if err = c.BodyParser(req); err != nil {
		return resp.Resp400(c, err.Error())
	}

	if err = c.BodyParser(&rm); err != nil {
		return resp.Resp400(c, err.Error())
	}

	if req.Id == 0 {
		return resp.Resp400(c, "未指定目标用户")
	}

	if target, err = controller.UserController.GetUser(c, req.Id); err != nil {
		return resp.RespError(c, err)
	}

	if op.Role < target.Role || ((op.Role == target.Role) && opt.RoleMustLess) {
		return resp.Resp403(c, req)
	}

	if op.Id == req.Id {
		return resp.Resp403(c, req, "无法更新自己")
	}

	if _, ok = rm["nickname"]; ok {
		if req.Nickname == "" {
			return resp.Resp400(c, req)
		}

		updates["nickname"] = req.Nickname
		changes["昵称"] = Change{Old: target.Nickname, New: req.Nickname}
	}

	if _, ok = rm["password"]; ok {
		if err = tool.CheckPassword(req.Password); err != nil {
			return resp.Resp400(c, err.Error())
		}

		updates["password"] = tool.NewPassword(req.Password)
		changes["密码"] = Change{Old: "******", New: "******"}
	}

	if _, ok = rm["status"]; ok {
		if req.Status.Code() == "unknown" {
			return resp.Resp400(c, req, "用户状态不符合")
		}

		updates["status"] = req.Status
		changes["状态"] = Change{Old: target.Status.Label(), New: req.Status.Label()}
	}

	if _, ok = rm["comment"]; ok {
		updates["comment"] = req.Comment
		changes["备注"] = Change{Old: target.Comment, New: req.Comment}
	}

	if _, ok = rm["role"]; ok {
		if op.Role < req.Role || ((op.Role == req.Role) && opt.RoleMustLess) {
			return resp.Resp400(c, req, "用户角色不符合")
		}

		updates["role"] = req.Role
		changes["角色"] = Change{Old: target.Role.Label(), New: req.Role.Label()}
	}

	if _, ok = rm["privileges"]; ok {
		for _, val := range req.Privileges {
			if lo.IndexOf(op.Privileges, val) < 0 {
				return resp.Resp400(c, req, fmt.Sprintf("权限: %s 不符合", val.Label()))
			}
		}

		changes["权限"] = Change{
			Old: lo.Map(target.Privileges, func(item model.Privilege, index int) string {
				return item.Label()
			}),
			New: lo.Map(req.Privileges, func(item model.Privilege, index int) string {
				return item.Label()
			}),
		}
		updates["privileges"] = req.Privileges
	}

	if _, ok = rm["active_at"]; ok {
		updates["active_at"] = time.UnixMilli(req.ActiveAt).UnixMilli()
		changes["激活时间"] = Change{Old: target.ActiveAt, New: req.ActiveAt}
	}

	if _, ok = rm["deadline"]; ok {
		updates["deadline"] = time.UnixMilli(req.Deadline).UnixMilli()
		changes["到期时间"] = Change{Old: target.Deadline, New: req.Deadline}
	}

	updated := new(model.User)
	if err = db.New(tool.Timeout(5)).
		Model(updated).
		Clauses(clause.Returning{}).
		Where("id = ?", req.Id).
		Updates(updates).
		Error; err != nil {
		return resp.Resp500(c, err.Error())
	}

	if err = controller.UserController.RmUserCache(c, req.Id); err != nil {
		return resp.RespError(c, err)
	}

	c.Locals(opt.OpLogLocalKey, &oplog.OpLog{Type: model.OpLogTypeUpdateUser, Content: map[string]any{
		"target_id":       target.Id,
		"target_username": target.Username,
		"changes":         changes,
	}})

	return resp.Resp200(c, updated)
}

func ManageUserDelete(c *nf.Ctx) error {
	type Req struct {
		Id uint64 `json:"id"`
	}

	var (
		ok     bool
		op     *model.User
		target *model.User
		err    error
		req    = new(Req)
	)

	if err = c.BodyParser(req); err != nil {
		return resp.Resp400(c, err.Error())
	}

	if req.Id == 0 {
		return resp.Resp400(c, req)
	}

	if op, ok = c.Locals("user").(*model.User); !ok {
		return resp.Resp401(c, nil)
	}

	if req.Id == op.Id {
		return resp.Resp400(c, nil, "无法删除自己")
	}

	if target, err = controller.UserController.GetUser(c, req.Id); err != nil {
		return resp.RespError(c, err)
	}

	if op.Role < target.Role || (op.Role == target.Role && opt.RoleMustLess) {
		return resp.Resp403(c, nil)
	}

	if err = controller.UserController.DeleteUser(c, target.Id); err != nil {
		return resp.RespError(c, err)
	}

	c.Locals(opt.OpLogLocalKey, &oplog.OpLog{Type: model.OpLogTypeDeleteUser, Content: map[string]any{
		"target_id":       target.Id,
		"target_username": target.Username,
	}})

	return resp.Resp200(c, nil, "删除成功")
}
