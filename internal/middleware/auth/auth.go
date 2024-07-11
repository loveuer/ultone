package auth

import (
	"errors"
	"gitea.com/taozitaozi/gredis"
	"github.com/go-redis/redis/v8"
	"github.com/loveuer/nf"
	"github.com/loveuer/nf/nft/resp"
	"strings"
	"ultone/internal/controller"
	"ultone/internal/log"
	"ultone/internal/opt"
)

var (
	tokenFunc = func(c *nf.Ctx) string {
		token := c.Get("Authorization")
		if token == "" {
			token = c.Cookies(opt.CookieName)
		}

		return token
	}
)

func NewAuth() nf.HandlerFunc {
	return func(c *nf.Ctx) error {
		token := tokenFunc(c)

		if token = strings.TrimPrefix(token, "Bearer "); token == "" {
			return resp.Resp401(c, token)
		}

		log.Debug(c, "middleware.NewAuth: token=%s", token)

		target, err := controller.UserController.GetUserByToken(c, token)
		if err != nil {
			log.Error(c, "middleware.NewAuth: get user by token=%s err=%v", token, err)
			if errors.Is(err, redis.Nil) || errors.Is(err, gredis.ErrKeyNotFound) {
				return resp.Resp401(c, err)
			}

			return resp.RespError(c, err)
		}

		if err = target.IsValid(true); err != nil {
			return resp.Resp401(c, err.Error(), err.Error())
		}

		c.Locals("user", target)
		c.Locals("token", token)

		return c.Next()
	}
}
