package app

import (
	"net/http"
	"time"

	"git.coding.net/bobxuyang/cy-gateway-BN/utils"
	humanize "github.com/dustin/go-humanize"
	"github.com/tomasen/realip"
)

// wrapper to capture status.
type wrapper struct {
	http.ResponseWriter
	written int
	status  int
}

// capture status.
func (w *wrapper) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// capture written bytes.
func (w *wrapper) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.written += n
	return n, err
}

// NewLoggingMiddle logging middle
func NewLoggingMiddle(logger *utils.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			res := &wrapper{w, 0, 200}
			ip := realip.RealIP(r)

			logger.Infof("[%s] >>> %s %s", ip, r.Method, r.RequestURI)
			h.ServeHTTP(res, r)
			size := humanize.Bytes(uint64(res.written))

			switch {
			case res.status >= 500:
				logger.Errorf("[%s] << %s %s %d (%s) in %s", ip, r.Method, r.RequestURI,
					res.status, size, time.Since(start))
			case res.status >= 400:
				logger.Warningf("[%s] << %s %s %d (%s) in %s", ip, r.Method, r.RequestURI,
					res.status, size, time.Since(start))
			default:
				logger.Infof("[%s] << %s %s %d (%s) in %s", ip, r.Method, r.RequestURI,
					res.status, size, time.Since(start))
			}
		})
	}
}
