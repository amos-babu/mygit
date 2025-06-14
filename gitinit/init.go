package gitinit

import (
	"fmt"
	"os"
)

func InitCommand() error {
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
