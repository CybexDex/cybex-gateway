package app

import (
	"net/http"

	"github.com/jinzhu/gorm"

	u "git.coding.net/bobxuyang/cy-gateway-BN/utils"
)

//NewLogEventMiddle ...
func NewLogEventMiddle(db *gorm.DB) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//start := time.Now()
			//ip := realip.RealIP(r)
			//
			//logger.Infof("[%s] >>> %s %s", ip, r.Method, r.RequestURI)
			//h.ServeHTTP(w, r)
			//size := humanize.Bytes(uint64(res.written))
			//
			//switch {
			//case res.status >= 500:
			//	logger.Errorf("[%s] << %s %s %d (%s) in %s", ip, r.Method, r.RequestURI,
			//		res.status, size, time.Since(start))
			//case res.status >= 400:
			//	//logger.Warningf("[%s] << %s %s %d (%s) in %s", ip, r.Method, r.RequestURI,
			//		res.status, size, time.Since(start))
			//default:
			//	logger.Infof("[%s] << %s %s %d (%s) in %s", ip, r.Method, r.RequestURI,
			//		res.status, size, time.Since(start))
			//}
			u.Debugln(111)
			h.ServeHTTP(w, r)
			u.Debugln(222)
			id := r.Context().Value("UserID")
			u.Infof("Login UserID: %d", id)
		})
	}
}
