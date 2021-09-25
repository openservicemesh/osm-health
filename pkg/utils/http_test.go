package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	tassert "github.com/stretchr/testify/assert"
)

func TestCheckHTTPResponseCodeEquals(t *testing.T) {
	tests := []struct {
		name               string
		expectedStatusCode int
		argStatusCode      int
		expectedError      bool
	}{
		{
			name:               "response codes are the same",
			expectedStatusCode: http.StatusOK,
			argStatusCode:      http.StatusOK,
			expectedError:      false,
		},
		{
			name:               "response codes are not the same",
			expectedStatusCode: http.StatusOK,
			argStatusCode:      http.StatusOK,
			expectedError:      false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := tassert.New(t)
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.expectedStatusCode)
			}))
			defer ts.Close()

			err := CheckHTTPResponseCodeEquals(ts.URL, test.argStatusCode)
			assert.Equal(test.expectedError, err != nil)
		})
	}
}

func TestGetResponseBody(t *testing.T) {
	expectedResponseBody := "test-response-body"

	tests := []struct {
		name                 string
		statusCode           int
		expectedResponseBody string
		expectedError        error
	}{
		{
			name:                 "StatusOk + Response Body Ok",
			statusCode:           http.StatusOK,
			expectedResponseBody: expectedResponseBody,
			expectedError:        nil,
		},
		{
			name:                 "StatusServiceUnavailable + Empty Response Body",
			statusCode:           http.StatusServiceUnavailable,
			expectedResponseBody: "",
			expectedError:        errors.Errorf("url returned HTTP status code: %d", http.StatusServiceUnavailable),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := tassert.New(t)
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.statusCode)
				_, _ = w.Write([]byte(expectedResponseBody))
			}))
			defer ts.Close()

			respBody, err := GetResponseBody(ts.URL)
			if test.expectedError != nil {
				assert.Equal(test.expectedError.Error(), err.Error())
			} else {
				assert.Equal(test.expectedResponseBody, respBody)
			}
		})
	}
}
