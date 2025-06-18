package hashobject_test

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/amos-babu/mygit/hashobject"
)

func TestHashObjectCommand(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := []byte("hello world\n")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Expected SHA-1 of blob header + content
	header := fmt.Sprintf("blob %d\x00", len(content))
	fullData := append([]byte(header), content...)
	expectedSha := sha1.Sum(fullData)
	expectedShaHex := hex.EncodeToString(expectedSha[:])

	// Run your command
	actualSha, err := hashobject.HashObjectCommand(tmpFile.Name())
	if err != nil {
		t.Fatalf("HashObjectCommand returned error: %v", err)
	}

	if actualSha != expectedShaHex {
		t.Errorf("SHA mismatch: got %s, expected %s", actualSha, expectedShaHex)
	}

	// Check if the compressed object file exists
	objectPath := filepath.Join(".mygit", "objects", actualSha[:2], actualSha[2:])
	data, err := os.ReadFile(objectPath)
	if err != nil {
		t.Fatalf("failed to read object file: %v", err)
	}

	defer os.RemoveAll(".mygit")

	// Decompress the file
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decompress object file: %v", err)
	}
	defer r.Close()

	decompressed, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read decompressed data: %v", err)
	}

	if !bytes.Equal(decompressed, fullData) {
		t.Errorf("decompressed content mismatch.\nExpected: %q\nGot: %q", fullData, decompressed)
	}

}
