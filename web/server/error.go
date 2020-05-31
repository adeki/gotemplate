package server

import (
	"net/http"

	"github.com/adeki/go-utils/web/errcode"
)

// errorHandler(badRequest, errcode.InvalidArguments)
func errorHandler(errfn func(http.ResponseWriter, errcode.Code), code errcode.Code) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		errfn(w, code)
	}
	return http.HandlerFunc(fn)
}

func renderError(w http.ResponseWriter, status int, c errcode.Code) {
	renderJSON(w, status, c.Struct())
}

func notFound(w http.ResponseWriter, c errcode.Code) {
	renderError(w, http.StatusNotFound, c)
}

func badRequest(w http.ResponseWriter, c errcode.Code) {
	renderError(w, http.StatusBadRequest, c)
}

func notAllowed(w http.ResponseWriter, c errcode.Code) {
	renderError(w, http.StatusMethodNotAllowed, c)
}

func internalError(w http.ResponseWriter, c errcode.Code) {
	renderError(w, http.StatusInternalServerError, c)
}
