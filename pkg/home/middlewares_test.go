package home

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jynychen/AdGuardHome/pkg/aghio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLimitRequestBody(t *testing.T) {
	errReqLimitReached := &aghio.LimitReachedError{
		Limit: defaultReqBodySzLim,
	}

	testCases := []struct {
		name    string
		body    string
		want    []byte
		wantErr error
	}{{
		name:    "not_so_big",
		body:    "somestr",
		want:    []byte("somestr"),
		wantErr: nil,
	}, {
		name:    "so_big",
		body:    string(make([]byte, defaultReqBodySzLim+1)),
		want:    make([]byte, defaultReqBodySzLim),
		wantErr: errReqLimitReached,
	}, {
		name:    "empty",
		body:    "",
		want:    []byte(nil),
		wantErr: nil,
	}}

	makeHandler := func(t *testing.T, err *error) http.HandlerFunc {
		t.Helper()

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var b []byte
			b, *err = io.ReadAll(r.Body)
			_, werr := w.Write(b)
			require.NoError(t, werr)
		})
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			handler := makeHandler(t, &err)
			lim := limitRequestBody(handler)

			req := httptest.NewRequest(http.MethodPost, "https://www.example.com", strings.NewReader(tc.body))
			res := httptest.NewRecorder()

			lim.ServeHTTP(res, req)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, res.Body.Bytes())
		})
	}
}
