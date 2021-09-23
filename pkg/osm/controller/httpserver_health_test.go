package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	tassert "github.com/stretchr/testify/assert"
)

func TestCheckControllerHealth(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		shouldError bool
	}{
		{
			name:        "StatusOk",
			statusCode:  http.StatusOK,
			shouldError: false,
		},
		{
			name:        "StatusServiceUnavailable",
			statusCode:  http.StatusServiceUnavailable,
			shouldError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := tassert.New(t)
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.statusCode)
			}))
			defer ts.Close()

			err := checkControllerHealthReadiness(ts.URL)
			assert.Equal(test.shouldError, err != nil)
			err = checkControllerHealthLiveness(ts.URL)
			assert.Equal(test.shouldError, err != nil)
		})
	}
}
