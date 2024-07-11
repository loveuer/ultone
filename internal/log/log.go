package log

import (
	"fmt"
	"github.com/loveuer/nf"
	ulog "github.com/loveuer/nf/nft/log"
	"ultone/internal/opt"
)

func _mix(c *nf.Ctx, msg string) string {
	if c == nil {
		return msg
	}

	return fmt.Sprintf("%v | %s", c.Locals(opt.LocalShortTraceKey), msg)
}

func Debug(c *nf.Ctx, msg string, data ...any) {
	ulog.Debug(_mix(c, msg), data...)
}

func Info(c *nf.Ctx, msg string, data ...any) {
	ulog.Info(_mix(c, msg), data...)
}

func Warn(c *nf.Ctx, msg string, data ...any) {
	ulog.Warn(_mix(c, msg), data...)
}

func Error(c *nf.Ctx, msg string, data ...any) {
	ulog.Error(_mix(c, msg), data...)
}

func Panic(c *nf.Ctx, msg string, data ...any) {
	ulog.Panic(_mix(c, msg), data...)
}

func Fatal(c *nf.Ctx, msg string, data ...any) {
	ulog.Fatal(_mix(c, msg), data...)
}
