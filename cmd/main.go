package main

import (
	"fmt"
	"os"

	"github.com/amos-babu/mygit/catfile"
	"github.com/amos-babu/mygit/gitinit"
	"github.com/amos-babu/mygit/hashobject"
)

// Usage: your_program.sh <command> <arg1> <arg2> ...
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	//Initialize git
	case "init":
		if err := gitinit.InitCommand(); err != nil {
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

		content, err := catfile.CatFileCommand(os.Args[3])
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
		hashedObject, err := hashobject.HashObjectCommand(os.Args[3])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Printf("%s\n", hashedObject)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
