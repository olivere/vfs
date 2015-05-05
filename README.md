# VFS - Virtual File System

```go
import "github.com/olivere/vfs"
```

Package vfs implements a virtual file system that gives access to files
below a root directory.

## Why?

When working with uploads and downloads in a web application you have
to ensure that only files inside a certain directory (or its subdirectories)
can be read or written. Package vfs makes this easier by implementing
common file operations available on top of a file system.

The idea is taken from [http.FileSystem](https://golang.org/pkg/net/http/#FileSystem)
and [http.File](https://golang.org/pkg/net/http/#File).

## Example

```go
// Create a file system for a directory
fs := vfs.FileSystem("/u/app1")

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
```
