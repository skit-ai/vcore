package tests

import (
	"testing"

	"github.com/Vernacular-ai/vcore/utils"
)

func TestWriteToFile(t *testing.T) {
	text := []byte(`Hello World!!`)

	if _, err := utils.WriteToFile(text, "test.txt"); err != nil {
		t.Error(err)
	}
}
