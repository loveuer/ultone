package front

import (
	"embed"
	"fmt"
	"github.com/loveuer/nf"
	"net/http"
	"strings"
	"ultone/internal/log"
)

//go:embed dist/front/browser
var DefaultFront embed.FS

func NewFront(ff *embed.FS, basePath string) nf.HandlerFunc {
	var (
		e          error
		indexBytes []byte
		index      string
	)

	index = fmt.Sprintf("%s/index.html", basePath)

	if indexBytes, e = ff.ReadFile(index); e != nil {
		log.Panic(nil, "read index file err: %v", e)
	}

	return func(c *nf.Ctx) error {
		var (
			err  error
			bs   []byte
			path = c.Path()
		)

		if bs, err = ff.ReadFile(basePath + path); err != nil {
			log.Debug(c, "embed read file [%s]%s err: %v", basePath, path, err)
			c.Set("Content-Type", "text/html")
			_, err = c.Write(indexBytes)
			return err
		}

		var dbs []byte
		if len(bs) > 512 {
			dbs = bs[:512]
		} else {
			dbs = bs
		}

		switch {
		case strings.HasSuffix(path, ".js"):
			c.Set("Content-Type", "application/javascript")
		case strings.HasSuffix(path, ".css"):
			c.Set("Content-Type", "text/css")
		default:
			c.Set("Content-Type", http.DetectContentType(dbs))
		}

		_, err = c.Write(bs)
		return err
	}
}
