package db

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"ultone/internal/opt"
)

func Init() error {
	strs := strings.Split(opt.Cfg.DB.Uri, "::")

	if len(strs) != 2 {
		return fmt.Errorf("db.Init: opt db uri invalid: %s", opt.Cfg.DB.Uri)
	}

	cli.ttype = strs[0]

	var (
		err error
		dsn = strs[1]
	)

	switch strs[0] {
	case "sqlite":
		opt.Cfg.DB.Type = "sqlite"
		cli.cli, err = gorm.Open(sqlite.Open(dsn))
	case "mysql":
		opt.Cfg.DB.Type = "mysql"
		cli.cli, err = gorm.Open(mysql.Open(dsn))
	case "postgres":
		opt.Cfg.DB.Type = "postgres"
		cli.cli, err = gorm.Open(postgres.Open(dsn))
	default:
		return fmt.Errorf("db type only support: [sqlite, mysql, postgres], unsupported db type: %s", strs[0])
	}

	if err != nil {
		return fmt.Errorf("db.Init: open %s with dsn:%s, err: %w", strs[0], dsn, err)
	}

	return nil
}
