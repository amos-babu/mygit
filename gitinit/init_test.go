package gitinit_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/amos-babu/mygit/gitinit"
)

func TestInitCommand(t *testing.T) {
	gitinit.InitCommand()
	//clean-up
	defer os.RemoveAll(".mygit")

	dirs := []string{".mygit", ".mygit/objects", ".mygit/refs"}
	for _, dir := range dirs {
		if _, err := os.Stat(dir); err != nil {
			t.Fatalf("Directory does not exist %v", err)
		}
	}

	if _, err := os.Stat(".mygit/HEAD"); err != nil {
		t.Fatalf("Directory does not exist %v", err)
	}

	expectedMain := fmt.Sprintf("ref: refs/heads/main\n")
	expectedMaster := fmt.Sprintf("ref: refs/heads/master\n")

	actualByte, err := os.ReadFile(".mygit/HEAD")
	if err != nil {
		t.Fatalf("Error reading file %v", err)
	}

	if string(actualByte) != expectedMain && string(actualByte) != expectedMaster {
		t.Errorf("Unexpected HEAD content.\nGot: %q\nExpected: %q or %q", string(actualByte), expectedMain, expectedMaster)
	}
}
