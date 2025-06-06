package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// Usage: your_program.sh <command> <arg1> <arg2> ...
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Fprintf(os.Stderr, "Logs from your program will appear here!\n")

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	//Initialize git
	case "init":
		if err := initCommand(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("Initialized git directory")

		//Read Blob Objects
	case "cat-file":
		if len(os.Args) != 4 || os.Args[2] != "-p" {
			fmt.Fprintf(os.Stderr, "usage: mygit cat-file -p <object-hash>\n")
			os.Exit(1)
		}

		content, err := catFileCommand(os.Args[3])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Print(string(content))

		//Create Blob Objects
	case "hash-object":
		if len(os.Args) != 4 || os.Args[2] != "-w" {
			fmt.Fprintf(os.Stderr, "usage: mygit hash-object -w <file>")
			os.Exit(1)
		}

		hashedObject, err := hashObjectCommand(os.Args[3])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Print(hashedObject)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}

func initCommand() error {
	for _, dir := range []string{".mygit", ".mygit/objects", ".mygit/refs"} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("Error creating directory: %s: %w\n", dir, err)
		}
	}

	headFileContents := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(".mygit/HEAD", headFileContents, 0644); err != nil {
		return fmt.Errorf("Error writing .mygit/HEAD: %w\n", err)
	}

	return nil
}

func catFileCommand(hash string) ([]byte, error) {
	objectHash := hash

	if len(objectHash) != 40 {
		fmt.Fprintf(os.Stderr, "usage: mygit cat-file -p <object-hash>\n")
		os.Exit(1)
	}

	dirName := objectHash[0:2]
	fileName := objectHash[2:]
	filePath := fmt.Sprintf(".mygit/objects/%s/%s", dirName, fileName)

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

func hashObjectCommand(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("mygit: failed to open file '%s': %w", file, err)
	}
	header := fmt.Sprintf("blob %d \x00 %s\n", len(data), string(data))

	hashedData := sha1.Sum([]byte(header))
	hashedHexString := hex.EncodeToString(hashedData[:])

	dirName := fmt.Sprintf(".mygit/objects/%s", hashedHexString[0:2])
	fileName := fmt.Sprintf(".mygit/objects/%s/%s", hashedHexString[0:2], hashedHexString[2:])

	if err := os.MkdirAll(dirName, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory '%s': %v\n", dirName, err)
	}

	if err := os.WriteFile(fileName, hashedData[:], 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating filename '%s': %v\n", fileName, err)
	}

	return fileName, nil
}
