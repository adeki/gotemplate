package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gorilla/csrf"
)

/*
	Wrapper of gorilla/csrf.
	This middleware is inject a token into a form tag automatically.
*/

type csrfResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (crw *csrfResponseWriter) Write(p []byte) (int, error) {
	return crw.buf.Write(p)
}

func CSRF(authKey []byte, args ...CSRFOption) func(http.Handler) http.Handler {
	opts := csrfOptions{
		cookieName: "csrf",
		fieldName:  "auth_token",
		secure:     true,
		errorHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "CSRF ERROR.", http.StatusForbidden)
		}),
	}
	for _, o := range args {
		o(&opts)
	}
	return func(next http.Handler) http.Handler {
		csrfProtect := csrf.Protect(
			authKey,
			csrf.CookieName(opts.cookieName),
			csrf.FieldName(opts.fieldName),
			csrf.Secure(opts.secure),
			csrf.ErrorHandler(opts.errorHandler),
		)
		fn := func(w http.ResponseWriter, r *http.Request) {
			crw := &csrfResponseWriter{
				ResponseWriter: w,
				buf:            &bytes.Buffer{},
			}
			next.ServeHTTP(crw, r)

			ptn := `(<form\s*[^>]*method=["'](?i)(?:post)["'][^>]*>)`
			field := string(csrf.TemplateField(r))
			reg := regexp.MustCompile(ptn)
			b := reg.ReplaceAll(crw.buf.Bytes(), []byte("$1\n"+field))
			fmt.Fprint(w, string(b))
		}

		return csrfProtect(http.HandlerFunc(fn))
	}
}

//
// options
//

type csrfOptions struct {
	cookieName   string
	fieldName    string
	secure       bool
	errorHandler http.Handler
}

type CSRFOption func(*csrfOptions)

func CSRFCookieName(s string) CSRFOption {
	return func(opts *csrfOptions) {
		if s != "" {
			opts.cookieName = s
		}
	}
}

func CSRFFieldName(s string) CSRFOption {
	return func(opts *csrfOptions) {
		if s != "" {
			opts.fieldName = s
		}
	}
}

func CSRFSecure(b bool) CSRFOption {
	return func(opts *csrfOptions) {
		opts.secure = b
	}
}

func CSRFErrorHandler(h http.Handler) CSRFOption {
	return func(opts *csrfOptions) {
		opts.errorHandler = h
	}
}
