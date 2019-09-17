package utils

import (
	"bufio"
	"encoding/csv"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/Vernacular-ai/vcore/errors"
	"github.com/Vernacular-ai/vcore/log"

	"github.com/hashicorp/go-getter"
)

// WriteToFile - Create directories/file and write to it
func WriteToFile(stream []byte, toFile string) (file *os.File, err error) {

	// Check if file exists
	_, err = os.Stat(toFile)

	// Create file if not exists
	if os.IsNotExist(err) {
		log.Debugf("Creating file: %s", toFile)
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

	log.Debugf("Writing to file: %s", toFile)

	file.Write(stream)
	file.Sync()
	log.Debugf("Written to file: %s", toFile)

	return
}

// GetFile - Download file from URL, create directories/file and write to it
func GetFile(fileURLPath, toFile string) (err error) {
	// Build the client
	client := &getter.Client{
		Src:  fileURLPath,
		Dst:  toFile,
		Mode: getter.ClientModeFile,
	}

	log.Debugf("Download file: %s", toFile)
	err = client.Get()
	if err != nil {
		return err
	}
	log.Debugf("Downloaded file: %s", toFile)

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
