package cmd

import (
	"flag"
	"time"
	"ultone/internal/opt"
)

func init() {
	time.Local = time.FixedZone("CST", 8*3600)

	flag.StringVar(&filename, "c", "etc/config.json", "config json file path")
	flag.BoolVar(&opt.Debug, "debug", false, "")

	flag.Parse()
}
