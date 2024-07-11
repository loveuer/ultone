package es

import (
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/loveuer/esgo2dump/xes/es7"
	"github.com/loveuer/nf/nft/log"
	"net/url"
	"ultone/internal/opt"
	"ultone/internal/tool"
)

var (
	Client *elastic.Client
)

func Init() error {
	ins, err := url.Parse(opt.Cfg.ES.Uri)
	if err != nil {
		return err
	}

	log.Debug("es.InitClient url parse uri: %s, result: %+v", opt.Cfg.ES.Uri, ins)

	if Client, err = es7.NewClient(tool.Timeout(10), ins); err != nil {
		return err
	}

	return nil
}
