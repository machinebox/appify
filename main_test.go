package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/matryer/is"
)

// NOTE: Run these tests with `make test`

func Test(t *testing.T) {
	is := is.New(t)
	out, err := exec.Command("./appify",
		"-name", "Test",
		"-icon", "testdata/machina-square.png",
		"testdata/app").CombinedOutput()
	t.Logf("%q", string(out))
	is.NoErr(err)
	defer os.RemoveAll("Test.app")
	actualAppHash := filehash(t, "testdata/app")
	type file struct {
		path string
		perm string
		hash string
	}
	for _, f := range []file{
		{path: "Test.app", perm: "drwxr-xr-x"},
		{path: "Test.app/Contents", perm: "drwxr-xr-x"},
		{path: "Test.app/Contents/MacOS", perm: "drwxr-xr-x"},
		{path: "Test.app/Contents/MacOS/Test.app", perm: "-rwxr-xr-x", hash: actualAppHash},
		{path: "Test.app/Contents/Info.plist", perm: "-rw-r--r--", hash: "0cd092b7b884e87617648dbdadb6a804"},
		{path: "Test.app/Contents/README", perm: "-rw-r--r--", hash: "afeb10df47c7f189b848ae44a54e7e06"},
		{path: "Test.app/Contents/Resources", perm: "drwxr-xr-x"},
		{path: "Test.app/Contents/Resources/icon.icns", perm: "-rw-r--r--", hash: "23bdc36475094ed8886f319811d3a182"},
	} {
		t.Run(f.path, func(t *testing.T) {
			is := is.New(t)
			info, err := os.Stat(f.path)
			is.NoErr(err)
			is.Equal(info.Mode().String(), f.perm) // perm
			if f.hash != "" {
				actual := filehash(t, f.path)
				is.Equal(actual, f.hash) // hash
			}
		})
	}
}

// filehash gets an md5 hash of the file at path.
func filehash(t *testing.T, path string) string {
	is := is.New(t)
	f, err := os.Open(path)
	is.NoErr(err)
	defer f.Close()
	h := md5.New()
	_, err = io.Copy(h, f)
	is.NoErr(err)
	return fmt.Sprintf("%x", h.Sum(nil))
}
