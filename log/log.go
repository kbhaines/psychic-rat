package log

import (
	"fmt"
	"log"
	"net/http"
)

func Logf(r *http.Request, s string, args ...interface{}) {
	rid := r.Context().Value("rid")
	userLog := fmt.Sprintf(s, args...)
	log.Printf("%d %s", rid, userLog)
}
