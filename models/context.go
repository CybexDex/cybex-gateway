package models

import "bitbucket.org/woyoutlz/bbb-gateway/utils/log"

// Context ...
type Context struct {
	logger  *log.Logger
	traceID string
}
