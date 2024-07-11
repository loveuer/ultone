package model

import (
	"fmt"
	"github.com/loveuer/nf/nft/log"
	"gorm.io/gorm"
	"strings"
	"ultone/internal/opt"
	"ultone/internal/sqlType"
)

func Init(db *gorm.DB) error {
	var err error

	if err = initModel(db); err != nil {
		return fmt.Errorf("model.MustInit: init models err=%v", err)
	}

	log.Info("MustInitModels: auto_migrate privilege model success")

	if err = initData(db); err != nil {
		return fmt.Errorf("model.MustInit: init datas err=%v", err)
	}

	return nil
}

func initModel(client *gorm.DB) error {
	if err := client.AutoMigrate(
		&User{},
		&OpLog{},
	); err != nil {
		return err
	}

	log.Info("InitModels: auto_migrate user model success")

	return nil
}

func initData(client *gorm.DB) error {
	var (
		err error
	)

	{
		count := 0

		if err = client.Model(&User{}).Select("count(id)").Take(&count).Error; err != nil {
			return err
		}

		if count < len(initUsers) {
			log.Warn("mustInitDatas: user count = 0, start init...")
			for _, user := range initUsers {
				if err = client.Model(&User{}).Create(user).Error; err != nil {
					if !strings.Contains(err.Error(), "SQLSTATE 23505") {
						return err
					}
				}
			}

			if opt.Cfg.DB.Type == "postgresql" {
				if err = client.Exec(`SELECT setval('users_id_seq', (SELECT MAX(id) FROM users))`).Error; err != nil {
					return err
				}
			}

			log.Info("InitDatas: creat init users success")
		} else {
			ps := make(sqlType.NumSlice[Privilege], 0)
			for _, item := range Privilege(0).All() {
				ps = append(ps, item.(Privilege))
			}
			if err = client.Model(&User{}).Where("id = ?", initUsers[0].Id).
				Updates(map[string]any{
					"privileges": ps,
				}).Error; err != nil {
				return err
			}

			log.Info("initDatas: update init users success")
		}
	}

	return nil
}
