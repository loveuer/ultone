package logger

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/loveuer/esgo2dump/log"
	"github.com/loveuer/nf"
	"github.com/loveuer/nf/nft/resp"
	"net/http"
	"strconv"
	"strings"
	"time"
	"ultone/internal/opt"
	"ultone/internal/tool"
)

var (
	Header = http.CanonicalHeaderKey("X-Trace-Id")
)

func New() nf.HandlerFunc {

	return func(c *nf.Ctx) error {
		var (
			now   = time.Now()
			trace = c.Get(Header)
			logFn func(msg string, data ...any)
			ip    = c.IP()
		)

		if trace == "" {
			trace = uuid.Must(uuid.NewV7()).String()
		}

		c.SetHeader(Header, trace)

		traces := strings.Split(trace, "-")
		shortTrace := traces[len(traces)-1]

		c.Locals(opt.LocalTraceKey, trace)
		c.Locals(opt.LocalShortTraceKey, shortTrace)

		err := c.Next()
		status, _ := strconv.Atoi(c.Writer.Header().Get(resp.RealStatusHeader))
		duration := time.Since(now)

		msg := fmt.Sprintf("%s | %15s | %d[%3d] | %s | %6s | %s", shortTrace, ip, c.StatusCode, status, tool.HumanDuration(duration.Nanoseconds()), c.Method(), c.Path())

		switch {
		case status >= 500:
			logFn = log.Error
		case status >= 400:
			logFn = log.Warn
		default:
			logFn = log.Info
		}

		logFn(msg)

		return err
	}
}
