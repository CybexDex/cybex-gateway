package eventlog

import "cybex-gateway/utils/log"

//Log ... TODO save to db
func Log(name string, event string) {
	log.Infoln(name, event)
}
