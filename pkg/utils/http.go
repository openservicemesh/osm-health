package utils

import (
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// GetResponseBody returns the respo
func GetResponseBody(url string) (string, error) {
	// #nosec G107: Potential HTTP request made with variable url
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.Errorf("error fetching (GET) url %s: %s", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("url returned http status code: %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Errorf("error rendering HTTP response: %s", err)
	}

	return string(respBody), nil
}
