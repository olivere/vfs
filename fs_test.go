package vfs

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func cp(dst, src string) error {
	srcf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcf.Close()
	dstf, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(dstf, srcf); err != nil {
		dstf.Close()
		return err
	}
	return dstf.Close()
}

func TestFileSystemResolve(t *testing.T) {
	root, err := ioutil.TempDir("", "vfs")
	if err != nil {
		t.Fatal(err)
	}

	fs := FileSystem(root)
	if fs == nil {
		t.Fatalf("expected file system; got: %v", fs)
	}
	if fs.root != root {
		t.Fatalf("expected root %q; got: %q", root, fs.root)
	}

	tests := []struct {
		Name     string
		Expected string
	}{
		{"", root},
		{".", root},
		{"..", root},
		{"/", root},
		{"file.gif", root + "/file.gif"},
		{"files/file.gif", root + "/files/file.gif"},
		{"files\\file.gif", root + "/files\\file.gif"},
		{"/files/file.gif", root + "/files/file.gif"},
		{"/etc/password", root + "/etc/password"},
		{"../../../../../../etc/password", root + "/etc/password"},
		{"//file.gif", root + "/file.gif"},
		{"///file.gif", root + "/file.gif"},
		{"file.gif//", root + "/file.gif"},
		{"file.gif//.", root + "/file.gif"},
		{`\\host\share\test`, root + `/\\host\share\test`},
		{`\\host\share\..\test`, root + `/\\host\share\..\test`},
	}
	for _, test := range tests {
		got := fs.Resolve(test.Name)
		if got != test.Expected {
			t.Errorf("expected %q; got: %q", test.Expected, got)
		}
	}
}

func TestFileSystemJoin(t *testing.T) {
	root, err := ioutil.TempDir("", "storev2")
	if err != nil {
		t.Fatal(err)
	}

	fs := FileSystem(root)
	if fs == nil {
		t.Fatalf("expected file system; got: %v", fs)
	}
	if fs.root != root {
		t.Fatalf("expected root %q; got: %q", root, fs.root)
	}

	tests := []struct {
		Parts    []string
		Expected string
	}{
		{[]string{""}, root},
		{[]string{"."}, root},
		{[]string{".."}, root},
		{[]string{"/"}, root},
		{[]string{"file.gif"}, root + "/file.gif"},
		{[]string{"files/file.gif"}, root + "/files/file.gif"},
		{[]string{"files", "file.gif"}, root + "/files/file.gif"},
		{[]string{"a", "b", "c", "file.gif"}, root + "/a/b/c/file.gif"},
		{[]string{"..", "..", "..", "etc", "password"}, root + "/etc/password"},
		{[]string{"..", "..", "..", "/", "etc", "password"}, root + "/etc/password"},
		{[]string{"/", "etc", "password"}, root + "/etc/password"},
	}
	for _, test := range tests {
		got := fs.Join(test.Parts...)
		if got != test.Expected {
			t.Errorf("expected %q; got: %q", test.Expected, got)
		}
	}
}

func TestFileSystemOpen(t *testing.T) {
	root, err := ioutil.TempDir("", "storev2")
	if err != nil {
		t.Fatal(err)
	}

	fs := FileSystem(root)
	if fs == nil {
		t.Fatalf("expected file system; got: %v", fs)
	}
	if fs.root != root {
		t.Fatalf("expected root %q; got: %q", root, fs.root)
	}

	// Setup a directory structure for the project
	os.MkdirAll(path.Join(root, "a", "1"), 0775)
	os.MkdirAll(path.Join(root, "b", "1"), 0775)
	os.MkdirAll(path.Join(root, "c", "1"), 0775)
	if err := cp(path.Join(root, "image1.jpg"), path.Join("testdata", "routercat.jpg")); err != nil {
		t.Fatal(err)
	}
	if err := cp(path.Join(root, "b", "1", "image-b-1.jpg"), path.Join("testdata", "routercat.jpg")); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		Name    string
		Success bool
	}{
		{"", false},
		{".", false},
		{"..", false},
		{"/", false},
		{"a", false},
		{"a/", false},
		{"a/.", false},
		{"a/..", false},
		{"d/", false},
		{"image1.jpg", true},
		{"/image1.jpg", true},
		{"//image1.jpg", true},
		{"image1.jpg//", true},
		{"//image1.jpg//", true},
		{"../image1.jpg", false},
		{"/b/1/image-b-1.jpg", true},
		{"/b//1/image-b-1.jpg", true},
		{"/b/image-b-1.jpg", false},
		{"/image-b-1.jpg", false},
		{"image-b-1.jpg", false},
	}
	for _, test := range tests {
		f, err := fs.Open(test.Name)
		if err != nil {
			// Expected error
			if test.Success {
				t.Fatalf("expected success; got: %v", err)
			}
		} else {
			// Expected success
			if f == nil {
				t.Fatalf("expected file; got: %v", f)
			}
			f.Close()
		}
	}
}

func TestFileSystemCreate(t *testing.T) {
	root, err := ioutil.TempDir("", "storev2")
	if err != nil {
		t.Fatal(err)
	}

	fs := FileSystem(root)
	if fs == nil {
		t.Fatalf("expected file system; got: %v", fs)
	}
	if fs.root != root {
		t.Fatalf("expected root %q; got: %q", root, fs.root)
	}

	// Setup a directory structure for the project
	os.MkdirAll(path.Join(root, "a"), 0775)
	os.MkdirAll(path.Join(root, "b"), 0775)
	if err := cp(path.Join(root, "image1.jpg"), path.Join("testdata", "routercat.jpg")); err != nil {
		t.Fatal(err)
	}
	if err := cp(path.Join(root, "a", "image1.jpg"), path.Join("testdata", "routercat.jpg")); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		Name    string
		Success bool
	}{
		{"", false},
		{".", false},
		{"..", false},
		{"/", false},
		{"image1.jpg", true},
		{"image2.jpg", false},
		{"a/image1.jpg", true},
	}
	for _, test := range tests {
		f, err := fs.Create(test.Name)
		if err != nil {
			// Expected error
			if test.Success {
				t.Fatalf("expected success; got: %v", err)
			}
		} else {
			// Expected success
			if f == nil {
				t.Fatalf("expected file; got: %v", f)
			}
			f.Close()
		}
	}
}

func TestFileSystemStat(t *testing.T) {
	root, err := ioutil.TempDir("", "storev2")
	if err != nil {
		t.Fatal(err)
	}

	fs := FileSystem(root)
	if fs == nil {
		t.Fatalf("expected file system; got: %v", fs)
	}
	if fs.root != root {
		t.Fatalf("expected root %q; got: %q", root, fs.root)
	}

	// Setup a directory structure for the project
	os.MkdirAll(path.Join(root, "a"), 0775)
	os.MkdirAll(path.Join(root, "b"), 0775)
	if err := cp(path.Join(root, "image1.jpg"), path.Join("testdata", "routercat.jpg")); err != nil {
		t.Fatal(err)
	}
	if err := cp(path.Join(root, "image2.jpg"), path.Join("testdata", "routercat.jpg")); err != nil {
		t.Fatal(err)
	}
	if err := cp(path.Join(root, "a", "image1.jpg"), path.Join("testdata", "routercat.jpg")); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		Name    string
		Success bool
		IsDir   bool
	}{
		{"", true, true},
		{".", true, true},
		{"..", true, true},
		{"/", true, true},
		{"image1.jpg", true, false},
		{"image2.jpg", true, false},
		{"a", true, true},
		{"a/", true, true},
		{"a/.", true, true},
		{"a/..", true, true},
		{"a/image1.jpg", true, false},
		{"b", false, true},
	}
	for _, test := range tests {
		fi, err := fs.Stat(test.Name)
		if err != nil {
			// Expected error
			if test.Success {
				t.Fatalf("expected success; got: %v", err)
			}
		} else {
			// Expected success
			if fi == nil {
				t.Fatalf("expected file info; got: %v", fi)
			}
			if test.IsDir && !fi.IsDir() {
				t.Errorf("expected %q to be a dir", test.Name)
			}
			if !test.IsDir && fi.IsDir() {
				t.Errorf("expected %q to not be a dir", test.Name)
			}
		}
	}
}

func TestFileSystemRemove(t *testing.T) {
	root, err := ioutil.TempDir("", "storev2")
	if err != nil {
		t.Fatal(err)
	}

	fs := FileSystem(root)
	if fs == nil {
		t.Fatalf("expected file system; got: %v", fs)
	}
	if fs.root != root {
		t.Fatalf("expected root %q; got: %q", root, fs.root)
	}

	// Setup a directory structure for the project
	os.MkdirAll(path.Join(root, "a"), 0775)
	if err := cp(path.Join(root, "image1.jpg"), path.Join("testdata", "routercat.jpg")); err != nil {
		t.Fatal(err)
	}
	if err := cp(path.Join(root, "image2.jpg"), path.Join("testdata", "routercat.jpg")); err != nil {
		t.Fatal(err)
	}
	if err := cp(path.Join(root, "a", "image1.jpg"), path.Join("testdata", "routercat.jpg")); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		Name    string
		Success bool
	}{
		{"", false},
		{".", false},
		{"..", false},
		{"/", false},
		{"image1.jpg", true},
		{"image1.jpg", false},
		{"image2.jpg", true},
		{"image2.jpg", false},
		{"a", false},
		{"a/", false},
		{"a/.", false},
		{"a/..", false},
		{"a/image1.jpg", true},
		{"a/image1.jpg", false},
		{"b", false},
	}
	for _, test := range tests {
		err = fs.Remove(test.Name)
		if err == nil && !test.Success {
			t.Fatalf("expected success for %q; got: %v", test.Name, err)
		}
		if err != nil && test.Success {
			t.Fatalf("expected failure for %q; got: %v", test.Name, err)
		}
	}
}
