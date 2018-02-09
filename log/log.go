package log

import (
	"fmt"
	"log"
	"net/http"
)

func Logf(r *http.Request, s string, args ...interface{}) {
	rid := r.Context().Value("rid")
	uid := r.Context().Value("uid")
	userLog := fmt.Sprintf(s, args...)
	log.Printf("%d %s %s", rid, uid, userLog)
}
