package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/loveuer/nf"
	"github.com/loveuer/nf/nft/resp"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strings"
	"time"
	"ultone/internal/database/cache"
	"ultone/internal/database/db"
	"ultone/internal/log"
	"ultone/internal/model"
	"ultone/internal/opt"
	"ultone/internal/tool"
)

type userController interface {
	GetUser(c *nf.Ctx, id uint64) (*model.User, error)
	GetUserByToken(c *nf.Ctx, token string) (*model.User, error)
	CacheUser(c *nf.Ctx, user *model.User) error
	CacheToken(c *nf.Ctx, token string, user *model.User) error
	RmUserCache(c *nf.Ctx, id uint64) error
	DeleteUser(c *nf.Ctx, id uint64) error
}

type uc struct{}

var _ userController = (*uc)(nil)

func (u uc) GetUser(c *nf.Ctx, id uint64) (*model.User, error) {
	var (
		err    error
		target = new(model.User)
		key    = fmt.Sprintf("%s:user:id:%d", opt.CachePrefix, id)
		bs     []byte
	)

	if opt.EnableUserCache {
		if bs, err = cache.Client.Get(tool.Timeout(3), key); err != nil {
			log.Warn(c, "controller.GetUser: get user by cache key=%s err=%v", key, err)
			goto ByDB
		}

		if err = json.Unmarshal(bs, target); err != nil {
			log.Warn(c, "controller.GetUser: json unmarshal key=%s by=%s err=%v", key, string(bs), err)
			goto ByDB
		}

		return target, nil
	}

ByDB:
	if err = db.New(tool.Timeout(3)).
		Model(&model.User{}).
		Where("id = ?", id).
		Take(target).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// tips: 公开项目需要考虑击穿处理
			return target, resp.NewError(400, "目标不存在", err, nil)
		}

		return target, resp.NewError(500, "", err, nil)
	}

	if opt.EnableUserCache {
		if err = u.CacheUser(c, target); err != nil {
			log.Warn(c, "controller.GetUser: cache user key=%s err=%v", key, err)
		}
	}

	return target, nil
}

func (u uc) GetUserByToken(c *nf.Ctx, token string) (*model.User, error) {
	strs := strings.Split(token, ".")
	if len(strs) != 3 {
		return nil, fmt.Errorf("controller.GetUserByToken: jwt token invalid, token=%s", token)
	}

	key := fmt.Sprintf("%s:user:token:%s", opt.CachePrefix, strs[2])
	bs, err := cache.Client.Get(tool.Timeout(3), key)
	if err != nil {
		return nil, err
	}

	log.Debug(c, "controller.GetUserByToken: key=%s cache bytes=%s", key, string(bs))

	userId := cast.ToUint64(string(bs))
	if userId == 0 {
		return nil, fmt.Errorf("controller.GetUserByToken: bs=%s cast to uint64 err", string(bs))
	}

	var op *model.User

	if op, err = u.GetUser(c, userId); err != nil {
		return nil, err
	}

	return op, nil
}

func (u uc) CacheUser(c *nf.Ctx, target *model.User) error {
	key := fmt.Sprintf("%s:user:id:%d", opt.CachePrefix, target.Id)
	return cache.Client.Set(tool.Timeout(3), key, target)
}

func (u uc) CacheToken(c *nf.Ctx, token string, user *model.User) error {
	strs := strings.Split(token, ".")
	if len(strs) != 3 {
		return fmt.Errorf("controller.CacheToken: jwt token invalid")
	}

	key := fmt.Sprintf("%s:user:token:%s", opt.CachePrefix, strs[2])
	return cache.Client.SetEx(tool.Timeout(3), key, user.Id, opt.TokenTimeout)
}
func (u uc) RmUserCache(c *nf.Ctx, id uint64) error {
	key := fmt.Sprintf("%s:user:id:%d", opt.CachePrefix, id)
	return cache.Client.Del(tool.Timeout(3), key)
}

func (u uc) DeleteUser(c *nf.Ctx, id uint64) error {
	var (
		err      error
		now      = time.Now()
		username = "CONCAT(username, '@del')"
	)

	if opt.Cfg.DB.Type == "sqlite" {
		username = "username || '@del'"
	}

	if err = db.New(tool.Timeout(5)).
		Model(&model.User{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"deleted_at": now.UnixMilli(),
			"username":   gorm.Expr(username),
		}).Error; err != nil {
		return resp.NewError(500, "", err, nil)
	}

	if opt.EnableUserCache {
		if err = u.RmUserCache(c, id); err != nil {
			log.Warn(c, "controller.DeleteUser: rm user=%d cache err=%v", id, err)
		}
	}

	return nil
}
