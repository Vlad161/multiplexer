//go:generate mockgen -source=contract.go -package $GOPACKAGE -destination mock_contract_test.go

package handler_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/vlad161/multiplexer/internal/http/handler"
	"github.com/vlad161/multiplexer/internal/logger"
)

const (
	url1 = "http://example1.com"
	url2 = "http://example2.com"
	url3 = "http://example3.com"
)

var (
	tooMuchUrls = []string{"http://example1.com", "http://example2.com", "http://example3.com", "http://example4.com",
		"http://example5.com", "http://example6.com", "http://example7.com", "http://example8.com", "http://example9.com",
		"http://example10.com", "http://example11.com", "http://example12.com", "http://example13.com", "http://example14.com",
		"http://example15.com", "http://example16.com", "http://example17.com", "http://example18.com", "http://example19.com",
		"http://example20.com", "http://example21.com"}
)

func TestHandler_ServeHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	log := logger.New()

	tests := []struct {
		name string

		method  string
		urls    []string
		reqBody string

		mpData  map[string]map[string]interface{}
		mpError error

		expectedCode int
	}{
		{
			name: "ok",

			method: http.MethodPost,
			urls:   []string{url1, url2, url3},
			mpData: map[string]map[string]interface{}{
				url1: {"key1": "value1"},
				url2: {"key2": "value2"},
				url3: {"key3": "value3"},
			},

			expectedCode: http.StatusOK,
		},
		{
			name: "empty url",

			method: http.MethodPost,
			urls:   []string{},

			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid request body",
			method:       http.MethodPost,
			urls:         []string{url1},
			reqBody:      "[",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "method get",

			method: http.MethodGet,

			expectedCode: http.StatusNotFound,
		},
		{
			name: "too much urls",

			method: http.MethodPost,
			urls:   tooMuchUrls,

			expectedCode: http.StatusRequestEntityTooLarge,
		},
		{
			name:         "multiplexer service error",
			method:       http.MethodPost,
			urls:         []string{url1},
			mpError:      fmt.Errorf("multiplexer service error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if len(tc.reqBody) == 0 {
				tc.reqBody = reqBodyByUrls(tc.urls)
			}

			mp := NewMockMultiplexer(ctrl)
			mp.EXPECT().
				Urls(gomock.Any(), tc.urls).
				Return(tc.mpData, tc.mpError).
				AnyTimes()

			reqBody := ioutil.NopCloser(strings.NewReader(tc.reqBody))
			defer reqBody.Close()
			rr := httptest.NewRecorder()

			req, errNewReq := http.NewRequest(tc.method, "/multiplexer", reqBody)
			require.NoError(t, errNewReq)
			req = req.WithContext(logger.WithContext(req.Context(), log))

			handler.NewMultiplexer(mp).ServeHTTP(rr, req)

			require.Equal(t, tc.expectedCode, rr.Code)
			if rr.Code != http.StatusOK {
				return
			}

			var respBody map[string]map[string]interface{}
			errDecode := json.NewDecoder(rr.Body).Decode(&respBody)
			require.NoError(t, errDecode)
			require.Equal(t, len(tc.mpData), len(respBody))

			for k, v := range tc.mpData {
				require.Equal(t, v, respBody[k])
			}
		})
	}
}

func reqBodyByUrls(urls []string) string {
	b := strings.Builder{}
	b.WriteRune('[')
	for i, url := range urls {
		if i > 0 {
			b.WriteString(`, `)
		}

		b.WriteString(fmt.Sprintf(`"%s"`, url))
	}
	b.WriteRune(']')
	return b.String()
}
