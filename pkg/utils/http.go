package utils

import (
<<<<<<< HEAD
=======
	"io/ioutil"
>>>>>>> aece626 (check(controller): Check controller http server health probes.)
	"net/http"

	"github.com/pkg/errors"
)

// IsHTTPResponseCodeEquals checks whether the returned response from the url matches the status code.
func IsHTTPResponseCodeEquals(url string, statusCode int) error {
	// #nosec G107: Potential HTTP request made with variable url
	resp, err := http.Get(url)
	if err != nil {
		return errors.Errorf("error fetching (GET) url %s: %s", url, err)
	}

	if resp.StatusCode != statusCode {
		return errors.Errorf("checking for HTTP status code: %d, but url returned HTTP status code: %d", statusCode, resp.StatusCode)
	}

	return nil
}

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
