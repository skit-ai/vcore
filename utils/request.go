package utils

import (
	"crypto/tls"
	"crypto/x509"
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

// HTTPGetWithCACrt - Download from url to a bytes array using CA crt
// pathToCaCertFile = Mount path of the CAcrt file
func HTTPGetWithCACrt(pathToCaCertFile string, url string) (payloadBytes []byte, err error) {
	caCert, err := ioutil.ReadFile(pathToCaCertFile)
	if err != nil {
		return
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	// Download from URL
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		payloadBytes = nil
		return
	}
	defer resp.Body.Close()

	payloadBytes, err = ioutil.ReadAll(resp.Body)
	return
}
