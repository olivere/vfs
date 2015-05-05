package vfs_test

import (
	"fmt"
	"log"
	"os"

	"github.com/olivere/vfs"
)

func ExampleFileSystem() {
	root := os.TempDir() // assume root = /u/app1 for the rest of the example
	fs := vfs.FileSystem(root)

	fs.Resolve("index.html")     // will return /u/app1/index.html
	fs.Resolve("../secret.html") // will return /u/app1/secret.html
	fs.Resolve("/etc/passwd")    // will return /u/app1/etc/passwd

	// Create a file /u/app1/index.html
	f, err := fs.Create("index.html")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(f, "Hello world")
	f.Close()

	// Read the contents of the new file
	f, err = fs.Open("index.html") // or OpenFile
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	// Stat the file
	fi, err := fs.Stat("index.html")
	if err != nil {
		log.Fatal(err)
	}
	if fi.IsDir() {
		log.Fatal("should not be a directory")
	}

	// Make a subdirectory
	err = fs.Mkdir("subdir", 0755) // or fs.MkdirAll
	if err != nil {
		log.Fatal(err)
	}
	fi, err = fs.Stat("subdir")
	if err != nil {
		log.Fatal(err)
	}
	if !fi.IsDir() {
		log.Fatal("should be a directory")
	}

	// Remove the file
	err = fs.Remove("index.html")
	if err != nil {
		log.Fatal(err)
	}
}
