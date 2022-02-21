package utils

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/Vernacular-ai/vcore/errors"

	"github.com/hashicorp/go-getter"
)

var httpClient http.Client

func init() {
	httpClient = http.Client{}
}

// WriteToFile - Create directories/file and write to it
func WriteToFile(stream []byte, toFile string) (file *os.File, err error) {

	// Check if file exists
	_, err = os.Stat(toFile)

	// Create file if not exists
	if os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(toFile), os.ModePerm)
		file, err = os.Create(toFile)
		if err != nil {
			err = errors.NewError("Unable to create file", err, false)
			return
		}
		defer file.Close()
	} else {
		file, err = os.OpenFile(toFile, os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			err = errors.NewError("Unable to open file", err, false)
			return
		}
		defer file.Close()
	}

	file.Write(stream)
	file.Sync()

	return
}

// GetFile - Download file from URL, create directories/file and write to it
// TODO: Allow cancelable requests using contexts https://github.com/hashicorp/go-getter/issues/102
func GetFile(fileURLPath, toFile string) (err error) {
	// Build the client
	client := &getter.Client{
		Src:  fileURLPath,
		Dst:  toFile,
		Mode: getter.ClientModeFile,
	}

	err = client.Get()
	if err != nil {
		return err
	}

	return
}

// ReadCsvFile - Read from csv file and return results as [][]string
func ReadCsvFile(filePath string) ([][]string, error) {
	// Load a csv file.
	f, _ := os.Open(filePath)

	// Create a new reader.
	r := csv.NewReader(bufio.NewReader(f))
	records, err := r.ReadAll()

	return records, err
}

// ReadYamlFile - Read from YAML file and return results
func ReadYamlFile(filePath string, out interface{}) (err error) {
	var yamlFile []byte
	if yamlFile, err = ioutil.ReadFile(filePath); err != nil {
		err = errors.NewError("Unable to read "+filePath, err, true)
	} else {
		if _err := yaml.Unmarshal(yamlFile, out); _err != nil {
			err = errors.NewError("Unable to deserialize "+filePath+" into a struct - ", _err, true)
		}
	}

	return err
}

// DownloadFile - Download from url to a local file
func DownloadFile(ctx context.Context, url, filepath string) (err error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req = req.WithContext(ctx)

	// Download from URL
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// Response status code check
	// 200 <= response status code < 400
	if !(resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusBadRequest) {
		errMessage := fmt.Sprintf("Failed: Error in downloading file.\nRequest Context: filepath => %s and file_download_url => %s\nResponse Context: \nstatus_code => %d\nresponse_text => %s", filepath, url, resp.StatusCode, string(bodyBytes))
		err = errors.NewError(errMessage, nil, false)
		Capture(err, false)
		return
	}
	_, err = WriteToFile(bodyBytes, filepath)
	if err != nil {
		return
	}

	return
}
