package log

import (
	"context"
	"fmt"
	"log"
)

const (
	logLevel   = "  LOG"
	errorLevel = "ERROR"
)

func Logf(c context.Context, s string, args ...interface{}) {
	logf(c, logLevel, s, args...)
}

func Errorf(c context.Context, s string, args ...interface{}) {
	logf(c, errorLevel, s, args...)
}

func logf(c context.Context, level string, s string, args ...interface{}) {
	rid := c.Value("rid")
	uid := c.Value("uid")
	userLog := fmt.Sprintf(s, args...)
	log.Printf("%s: %d %s %s", level, rid, uid, userLog)
}
