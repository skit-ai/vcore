package utils

import (
	"os"
	"path/filepath"

	"github.com/Vernacular-ai/vcore/log"

	"github.com/hashicorp/go-getter"
)

// WriteToFile - Create file and write to it
func WriteToFile(stream []uint8, toFile string) *os.File {

	// Check if file exists
	var _, err = os.Stat(toFile)

	// Create file if not exists
	if os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(toFile), os.ModePerm)
		log.Debugf("Writting to file: %s", toFile)
		file, _ := os.Create(toFile)
		defer file.Close()

		file.Write(stream)
		file.Sync()
		log.Debugf("Written to file: %s", toFile)
	}

	return nil
}

// GetFile - Download file from URL
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
