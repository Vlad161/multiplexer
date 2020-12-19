//go:generate mockgen -source=contract.go -package $GOPACKAGE -destination mock_contract_test.go

package multiplexer_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/vlad161/multiplexer/internal/service/multiplexer"
)

const (
	okUrl1  = "ok1"
	okUrl2  = "ok2"
	okUrl3  = "ok3"
	okUrl4  = "ok4"
	okUrl5  = "ok5"
	errUrl1 = "err1"
	errUrl2 = "err2"
	errUrl3 = "err3"
)

var (
	urlData = map[string]map[string]interface{}{
		okUrl1: {"key1": "value1"},
		okUrl2: {"key2": "value2"},
		okUrl3: {"key3": "value3"},
		okUrl4: {"key4": "value4"},
		okUrl5: {"key5": "value5"},
	}

	urlError = map[string]error{
		errUrl1: fmt.Errorf("url1 error"),
		errUrl2: fmt.Errorf("url2 error"),
		errUrl3: fmt.Errorf("url3 error"),
	}
)

func TestService_Urls(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name string
		urls []string

		expected      map[string]map[string]interface{}
		expectedError bool
	}{
		{
			name: "ok",
			urls: []string{okUrl1, okUrl2, okUrl3, okUrl4, okUrl5},

			expected: map[string]map[string]interface{}{
				okUrl1: urlData[okUrl1],
				okUrl2: urlData[okUrl2],
				okUrl3: urlData[okUrl3],
				okUrl4: urlData[okUrl4],
				okUrl5: urlData[okUrl5],
			},
		},
		{
			name: "one url return error",
			urls: []string{okUrl1, okUrl2, errUrl1},

			expectedError: true,
		},
		{
			name: "all urls return errors",
			urls: []string{errUrl1, errUrl2, errUrl3},

			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cl := NewMockClient(ctrl)
			for _, url := range tc.urls {
				var (
					data map[string]interface{}
					err  error
				)

				if _, ok := urlData[url]; ok {
					data = urlData[url]
				} else {
					err = urlError[url]
				}

				cl.EXPECT().Get(ctx, url).Return(data, err).AnyTimes()
			}

			s, errNew := multiplexer.New(multiplexer.WithClient(cl))
			require.NoError(t, errNew)

			result, errUrls := s.Urls(ctx, tc.urls)

			if tc.expectedError {
				require.Error(t, errUrls)
				return
			}

			require.NoError(t, errUrls)
			require.Equal(t, len(tc.expected), len(result))

			for k, v := range tc.expected {
				require.Equal(t, v, result[k])
			}
		})
	}
}
