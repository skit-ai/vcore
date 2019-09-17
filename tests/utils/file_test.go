package tests

import (
	"testing"

	"github.com/Vernacular-ai/vcore/log"
	"github.com/Vernacular-ai/vcore/utils"
)

func TestWriteToFile(t *testing.T) {
	log.SetLevel(log.DEBUG)
	text := []byte(`Hello World!!`)

	if _, err := utils.WriteToFile(text, "test.txt"); err != nil {
		t.Error(err)
	}
}
