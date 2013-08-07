package main

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

/*
Markdownd uses pygments, a Python tool, for syntax highlighting. If possible, I want to ship markdownd as a
single binary. To this end, I'm using the following strategy:

1. Tar up the vendor directory (containing pygments) and calculate the md5sum.
2. Use github.com/jteeuwen/go-bindata to write a Go file that embeds the tarball
3. Add the md5 to that file as well
4. When markdownd runs, first check for a vendored pygments colocated with the binary (in case we're running in
   development, for example).
5. If it does not exist, check for ~/.markdownd/.
6. If ~/.markdownd exists, then ~/.markdownd/md5 should contain an md5 checksum. See if it matches.
7. If the checksum doesn't match, or if ~/.markdownd doesn't exist, write out and untar the vendored data into
   ~/.markdownd, along with the current checksum.
8. Use pygments in ~/.markdownd.

1-3 are accomplished by the 'make vendor_data.go' task.
*/

const (
	pygmentPath = "vendor/pygments/pygmentize"
	// Relative to the user's home dir.
	pygmentsCache    = ".markdownd"
	cacheMD5Filename = "checksum"
)

func findPygments() (string, error) {
	// First see if pygments is located alongside the binary.
	exe, err := exec.LookPath(os.Args[0])
	if err == nil {
		pygmentize = filepath.Join(filepath.Dir(exe), pygmentPath)
		if _, err := os.Stat(pygmentize); err == nil {
			dbg.Println("found dev pygments in", filepath.Dir(exe))
			return pygmentize, nil
		}
	}

	// Next see if the cached version exists and is up-to-date.
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	cache := filepath.Join(u.HomeDir, pygmentsCache)
	md5File := filepath.Join(cache, cacheMD5Filename)
	pygmentize = filepath.Join(cache, pygmentPath)
	oldMD5, err := ioutil.ReadFile(md5File)
	if err == nil {
		if string(bytes.TrimSpace(oldMD5)) == VendorMD5 {
			// Up-to-date
			dbg.Println("found up-to-date cached pygments.")
			return pygmentPath, nil
		}
		fmt.Fprintln(os.Stderr, "Updating stale cache in", cache)
	}

	// Need to delete the existing, stale cache version (if it exists) and write out a new one.
	fmt.Fprintln(os.Stderr, "Writing out pygments cache to", cache)
	if err := os.RemoveAll(cache); err != nil {
		return "", err
	}
	if err := os.Mkdir(cache, 0755); err != nil {
		return "", err
	}

	if err := expandTarArchive(cache); err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(md5File, []byte(VendorMD5), 0600); err != nil {
		return "", err
	}

	return pygmentize, nil
}

func expandTarArchive(loc string) error {
	vendorData := bytes.NewBuffer(VendorData())
	reader := tar.NewReader(vendorData)
	for {
		f, err := reader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if f.Typeflag != tar.TypeReg && f.Typeflag != tar.TypeRegA {
			continue
		}
		name := filepath.Join(loc, f.Name)
		if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
			return err
		}
		flags := os.O_WRONLY | os.O_CREATE | os.O_EXCL
		file, err := os.OpenFile(name, flags, os.FileMode(f.Mode))
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := io.Copy(file, reader); err != nil {
			return err
		}
	}
	return nil
}
