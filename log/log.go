package log

import (
	"fmt"
	"log"
	"net/http"
)

const (
	logLevel   = "  LOG"
	errorLevel = "ERROR"
)

func Logf(r *http.Request, s string, args ...interface{}) {
	logf(r, logLevel, s, args...)
}

func Errorf(r *http.Request, s string, args ...interface{}) {
	logf(r, errorLevel, s, args...)
}

func logf(r *http.Request, level string, s string, args ...interface{}) {
	rid := r.Context().Value("rid")
	uid := r.Context().Value("uid")
	userLog := fmt.Sprintf(s, args...)
	log.Printf("%s: %d %s %s", level, rid, uid, userLog)
}
