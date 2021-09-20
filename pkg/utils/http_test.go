package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsHTTPResponseCodeEquals(t *testing.T) {
	tests := []struct {
		name               string
		expectedStatusCode int
	}{
		{
			name:               "StatusOk",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "StatusServiceUnavailable",
			expectedStatusCode: http.StatusServiceUnavailable,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.expectedStatusCode)
			}))
			defer ts.Close()

			err := CheckHTTPResponseCodeEquals(ts.URL, test.expectedStatusCode)
			assert.Nil(t, err)
		})
	}
}
