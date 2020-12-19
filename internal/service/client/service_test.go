//go:generate mockgen -source=contract.go -package $GOPACKAGE -destination mock_contract_test.go

package client_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/vlad161/multiplexer/internal/service/client"
)

const (
	respBodyOk        = `{"key": "value"}`
	respBodyIncorrect = `{"key":`
)

func TestService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name string
		url  string

		respCode int
		respBody string
		respErr  error

		expected      map[string]interface{}
		expectedError bool
	}{
		{
			name: "ok",
			url:  "http://example.com/",

			respCode: 200,
			respBody: respBodyOk,

			expected: map[string]interface{}{"key": "value"},
		},
		{
			name: "incorrect body",
			url:  "http://example.com/",

			respCode: 200,
			respBody: respBodyIncorrect,

			expectedError: true,
		},
		{
			name: "response code 404",
			url:  "http://example.com/",

			respCode: 404,

			expectedError: true,
		},
		{
			name: "transport error",
			url:  "http://example.com/",

			respErr: fmt.Errorf("transport error"),

			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, errNewRequest := http.NewRequestWithContext(ctx, http.MethodGet, tc.url, nil)
			require.NoError(t, errNewRequest)

			respBody := ioutil.NopCloser(strings.NewReader(tc.respBody))
			defer respBody.Close()

			transport := NewMockTransport(ctrl)
			transport.EXPECT().
				Do(req).
				Return(&http.Response{StatusCode: tc.respCode, Body: respBody}, tc.respErr).
				Times(1)

			cl, errCreateCl := client.New(client.WithTransport(transport))
			require.NoError(t, errCreateCl)

			result, errGet := cl.Get(ctx, tc.url)

			if tc.expectedError {
				require.Error(t, errGet)
				return
			}

			require.NoError(t, errGet)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestService_Create_Request(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		url string

		expectedError bool
	}{
		{
			url: "http://example.com",

			expectedError: false,
		},
		{
			url: ":",

			expectedError: true,
		},
		{
			url: "://",

			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.url, func(t *testing.T) {
			req, errNewRequest := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", nil)
			require.NoError(t, errNewRequest)

			respBody := ioutil.NopCloser(strings.NewReader(respBodyOk))
			defer respBody.Close()

			transport := NewMockTransport(ctrl)
			transport.EXPECT().
				Do(req).
				Return(&http.Response{StatusCode: 200, Body: respBody}, nil).
				AnyTimes()

			cl, errCreateCl := client.New(client.WithTransport(transport))
			require.NoError(t, errCreateCl)

			_, errGet := cl.Get(ctx, tc.url)
			require.Equal(t, errGet != nil, tc.expectedError)
		})
	}
}
