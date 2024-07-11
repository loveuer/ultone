package opt

import (
	"encoding/json"
	"fmt"
	"github.com/loveuer/nf/nft/log"
	"os"
	"ultone/internal/tool"
)

type listen struct {
	Http string `json:"http"`
	Grpc string `json:"grpc"`
	Unix string `json:"unix"`
}

type db struct {
	Type string `json:"-"` // postgres, mysql, sqlite
	Uri  string `json:"uri"`
}

type cache struct {
	Uri string `json:"uri"`
}

type es struct {
	Uri   string `json:"uri"`
	Index struct {
		Staff string `json:"staff"`
	} `json:"index"`
}

type Nebula struct {
	Uri      string `json:"uri"`
	Username string `json:"username"`
	Password string `json:"password"`
	Space    string `json:"space"`
}

type config struct {
	Name   string `json:"name"`
	Listen listen `json:"listen"`
	DB     db     `json:"db"`
	Cache  cache  `json:"cache"`
	ES     es     `json:"es"`
	Nebula Nebula `json:"nebula"`
}

var (
	Debug bool
	Cfg   = &config{}
)

func Init(filename string) error {

	var (
		err error
		bs  []byte
	)

	log.Info("opt.Init: start reading config file: %s", filename)

	if bs, err = os.ReadFile(filename); err != nil {
		return fmt.Errorf("opt.Init: read config file=%s err=%v", filename, err)
	}

	if err = json.Unmarshal(bs, Cfg); err != nil {
		return fmt.Errorf("opt.Init: json marshal config=%s err=%v", string(bs), err)
	}

	if Debug {
		log.SetLogLevel(log.LogLevelDebug)
	}

	tool.TablePrinter(Cfg)

	return nil
}
