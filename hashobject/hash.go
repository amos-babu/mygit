package hashobject

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
)

func HashObjectCommand(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("mygit: failed to open file '%s': %w", file, err)
	}

	//Format: "blob <size>\x00<content>"
	header := fmt.Sprintf("blob %d\x00", len(data))
	fullData := append([]byte(header), data...)

	//Hashing the Object Data
	hashedData := sha1.Sum([]byte(fullData))
	hashedHexString := hex.EncodeToString(hashedData[:])

	//compress the object using zlib
	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	_, err = w.Write(fullData)
	if err != nil {
		return "", fmt.Errorf("failed to compress: %w", err)
	}
	w.Close()

	// Save to .mygit/objects/xx/yyyyyy...
	dirName := fmt.Sprintf(".mygit/objects/%s", hashedHexString[0:2])
	fileName := fmt.Sprintf(".mygit/objects/%s/%s", hashedHexString[0:2], hashedHexString[2:])

	if err := os.MkdirAll(dirName, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory '%s': %v\n", dirName, err)
	}

	if err := os.WriteFile(fileName, compressed.Bytes(), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating filename '%s': %v\n", fileName, err)
	}

	return hashedHexString, nil
}
