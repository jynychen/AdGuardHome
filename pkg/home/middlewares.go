package home

import (
	"io"
	"net/http"

	"github.com/jynychen/AdGuardHome/pkg/aghio"

	"github.com/AdguardTeam/golibs/log"
)

// middlerware is a wrapper function signature.
type middleware func(http.Handler) http.Handler

// withMiddlewares consequently wraps h with all the middlewares.
func withMiddlewares(h http.Handler, middlewares ...middleware) (wrapped http.Handler) {
	wrapped = h

	for _, mw := range middlewares {
		wrapped = mw(wrapped)
	}

	return wrapped
}

// defaultReqBodySzLim is the default maximum request body size.
const defaultReqBodySzLim = 64 * 1024

// largerReqBodySzLim is the maximum request body size for APIs expecting larger
// requests.
const largerReqBodySzLim = 4 * 1024 * 1024

// expectsLargerRequests shows if this request should use a larger body size
// limit.  These are exceptions for poorly designed current APIs as well as APIs
// that are designed to expect large files and requests.  Remove once the new,
// better APIs are up.
//
// See https://github.com/AdguardTeam/AdGuardHome/issues/2666 and
// https://github.com/AdguardTeam/AdGuardHome/issues/2675.
func expectsLargerRequests(r *http.Request) (ok bool) {
	m := r.Method
	if m != http.MethodPost {
		return false
	}

	p := r.URL.Path
	return p == "/control/access/set" ||
		p == "/control/filtering/set_rules"
}

// limitRequestBody wraps underlying handler h, making it's request's body Read
// method limited.
func limitRequestBody(h http.Handler) (limited http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		var szLim int64 = defaultReqBodySzLim
		if expectsLargerRequests(r) {
			szLim = largerReqBodySzLim
		}

		var reader io.Reader
		reader, err = aghio.LimitReader(r.Body, szLim)
		if err != nil {
			log.Error("limitRequestBody: %s", err)

			return
		}

		// HTTP handlers aren't supposed to call r.Body.Close(), so just
		// replace the body in a clone.
		rr := r.Clone(r.Context())
		rr.Body = io.NopCloser(reader)

		h.ServeHTTP(w, rr)
	})
}
