package tests

import (
	"context"
	"testing"

	"github.com/skit-ai/vcore/utils"
)

func TestWriteToFile(t *testing.T) {
	text := []byte(`Hello World!!`)

	if _, err := utils.WriteToFile(text, "test.txt"); err != nil {
		t.Error(err)
	}
}

func TestDownloadFile(t *testing.T) {
	if err := utils.DownloadFile(context.TODO(), "https://file-examples.com/wp-content/uploads/2017/11/file_example_MP3_700KB.mp3", "test.mp3"); err != nil {
		t.Error(err)
	}
}

func TestDownloadFileInvalidURL(t *testing.T) {
	if err := utils.DownloadFile(context.TODO(), "", "test"); err == nil {
		t.Error(err)
	}
}

func TestDownloadFileBadStatus(t *testing.T) {
	if err := utils.DownloadFile(context.TODO(), "https://s3.ap-south-1.amazonaws.com/abc/xyz/1645360189.196162.wav", "test"); err == nil {
		t.Error(err)
	}
}

func TestDownloadWithDirPath(t *testing.T) {
	if err := utils.DownloadFile(context.TODO(), "https://file-examples.com/wp-content/uploads/2017/11/file_example_MP3_700KB.mp3", "tests/test.mp3"); err != nil {
		t.Error(err)
	}
}
