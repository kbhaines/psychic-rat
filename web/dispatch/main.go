package dispatch

import "net/http"

type (
	MethodSelector map[string]http.HandlerFunc
	URIHandler     struct {
		URI     string
		Handler http.HandlerFunc
	}
)

func ExecHandlerForMethod(selector MethodSelector, writer http.ResponseWriter, request *http.Request) {
	f, exists := selector[request.Method]
	if !exists {
		http.Error(writer, "", http.StatusMethodNotAllowed)
		return
	}
	f(writer, request)
}
