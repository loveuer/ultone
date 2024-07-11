package cache

import (
	"fmt"
	"gitea.com/taozitaozi/gredis"
	"github.com/go-redis/redis/v8"
	"net/url"
	"strings"
	"ultone/internal/opt"
	"ultone/internal/tool"
)

func Init() error {

	var (
		err error
	)

	strs := strings.Split(opt.Cfg.Cache.Uri, "::")

	switch strs[0] {
	case "memory":
		gc := gredis.NewGredis(1024 * 1024)
		Client = &_mem{client: gc}
	case "lru":
		if Client, err = newLRUCache(); err != nil {
			return err
		}
	case "redis":
		var (
			ins *url.URL
			err error
		)

		if len(strs) != 2 {
			return fmt.Errorf("cache.Init: invalid cache uri: %s", opt.Cfg.Cache.Uri)
		}

		uri := strs[1]

		if !strings.Contains(uri, "://") {
			uri = fmt.Sprintf("redis://%s", uri)
		}

		if ins, err = url.Parse(uri); err != nil {
			return fmt.Errorf("cache.Init: url parse cache uri: %s, err: %s", opt.Cfg.Cache.Uri, err.Error())
		}

		addr := ins.Host
		username := ins.User.Username()
		password, _ := ins.User.Password()

		var rc *redis.Client
		rc = redis.NewClient(&redis.Options{
			Addr:     addr,
			Username: username,
			Password: password,
		})

		if err = rc.Ping(tool.Timeout(5)).Err(); err != nil {
			return fmt.Errorf("cache.Init: redis ping err: %s", err.Error())
		}

		Client = &_redis{client: rc}
	default:
		return fmt.Errorf("cache type %s not support", strs[0])
	}

	return nil
}
