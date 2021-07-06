package utils

import (
	"io/ioutil"
	"net/http"
)

// HTTPGet - Download from url to a bytes array
func HTTPGet(url string) (payloadBytes []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	// Download from URL
	resp, err := httpClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		payloadBytes = nil
		return
	}
	defer resp.Body.Close()

	payloadBytes, err = ioutil.ReadAll(resp.Body)
	return
}
