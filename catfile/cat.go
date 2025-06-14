package catfile

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

func CatFileCommand(hash string) ([]byte, error) {
	objectHash := hash

	if len(objectHash) != 40 {
		fmt.Fprintf(os.Stderr, "usage: mygit cat-file -p <object-hash>\n")
		os.Exit(1)
	}

	dirName := objectHash[0:2]
	fileName := objectHash[2:]
	filePath := fmt.Sprintf(".mygit/objects/%s/%s", dirName, fileName)

	// fmt.Println(filePath)
	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(1)
	}

	b := bytes.NewReader(fileContents)
	r, err := zlib.NewReader(b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decompressing the file: %s\n", err)
		os.Exit(1)
		// panic(err)
	}

	decompressedData, err := io.ReadAll(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading decompressed data: %s\n", err)
		os.Exit(1)
	}
	r.Close()

	nullIndex := bytes.IndexByte(decompressedData, 0)
	if nullIndex == -1 {
		fmt.Fprintf(os.Stderr, "Invalid object format: missing metadata separator\n")
		os.Exit(1)
	}

	content := decompressedData[nullIndex+1:]
	return content, nil
}
