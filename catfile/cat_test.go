package catfile_test

import (
	"os"
	"testing"

	"github.com/amos-babu/mygit/catfile"
	"github.com/amos-babu/mygit/hashobject"
)

func TestCatFileCommand(t *testing.T) {
	//create a temp file
	tempFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatalf("Failed to create a temp file: %v", err)
	}

	defer os.Remove(tempFile.Name())

	content := []byte("Hello Amos")

	if _, err := tempFile.Write(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	tempFile.Close()

	hash, err := hashobject.HashObjectCommand(tempFile.Name())
	if err != nil {
		t.Errorf("Failed to hash the file %v", err)
	}

	defer os.RemoveAll(".mygit")

	contentByte, err := catfile.CatFileCommand(hash)
	if err != nil {
		t.Errorf("Failed to read the hash %v", err)
	}

	if string(contentByte) != string(content) {
		t.Fatalf("Actual file is not same as the generated file")
	}

}
