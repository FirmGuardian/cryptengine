package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"runtime"
	"strings"
)

const (
	dirBin  = "bin"
	dirKeys = "keys"
	dirTmp  = "tmp"
)

// PathInfo for testing and retrieving important info for a given path
type PathInfo struct {
	Base   string
	Clean  string
	Ext    string
	Dir    string
	File   string
	Exists bool
	IsDir  bool
	IsReg  bool
	Size   int64
}

func appRoot() string {
	var p string
	switch runtime.GOOS {
	default:
		panic("Unsupported operating system")
	case "linux":
	case "darwin":
		root := os.Getenv("HOME")
		if root != "" {
			p = path.Join(root, ".LegalCrypt")
		} else {
			// TODO: Do this better
			panic("ERR::Don't know where home is.")
		}
	case "windows":
		root := os.Getenv("AppData")
		if root != "" {
			p = path.Join(root, "LegalCrypt")
		} else {
			// TODO: Do this better
			panic("ERR::Don't know where AppData is")
		}
	}

	return p
}

func binDir() string {
	return path.Join(appRoot(), dirBin)
}

func keyDir() string {
	return path.Join(appRoot(), dirKeys)
}

func outDirDec() string {
	currentUser, err := user.Current()
	if err != nil {
		// TODO: Handle this better
		panic(err)
	}
	return path.Join(currentUser.HomeDir, "Documents", "LegalCrypt", "Received")
}

func outDirEnc() string {
	currentUser, err := user.Current()
	if err != nil {
		// TODO: Handle this better
		panic(err)
	}
	return path.Join(currentUser.HomeDir, "Documents", "LegalCrypt", "Secured")
}

func tmpDir() string {
	return path.Join(appRoot(), dirTmp)
}

// If a path exists, delete it
func nixIfExists(filePath string) {
	if exists, _ := fileExists(filePath); exists {
		check(os.Remove(filePath), errs["fsCantDeleteFile"])
	}
}

func pathInfo(p string) PathInfo {
	d, f := path.Split(p)
	exists, isdir, isreg := false, false, false
	var sz int64
	finfo, err := os.Stat(p)
	if err == nil {
		exists = true
		mode := finfo.Mode()
		isdir = mode.IsDir()
		isreg = mode.IsRegular()
		sz = finfo.Size()
	}

	return PathInfo{
		Base:   path.Base(p),
		Clean:  path.Clean(p),
		Ext:    path.Ext(p),
		Dir:    d,
		File:   f,
		Exists: exists,
		IsDir:  isdir,
		IsReg:  isreg,
		Size:   sz,
	}
}

func scaffoldAppDirs() {
	lcPaths := []string{
		binDir(),
		keyDir(),
		outDirDec(),
		outDirEnc(),
		tmpDir(),
	}

	for _, lcPath := range lcPaths {
		fmt.Printf(";;Scaffolding %v...", lcPath)
		err := os.MkdirAll(lcPath, 0700)
		if err != nil {
			fmt.Println("ERR")
			log.Fatalln(err)
		} else {
			fmt.Println("Done")
		}
	}
}

func pathEndsWith(haystack string, needle string) bool {
	lastIndex := strings.LastIndex(haystack, needle)
	return 0 == len(haystack)-lastIndex-len(needle) // simple, but fragile. consider using the pathInfo.Ext
}

func pathEndsWithLCSF(path string) bool {
	return pathEndsWith(strings.ToLower(path), legalCryptFileExtension)
}

func stripTrailingLCExt(p string) string {
	return strings.Replace(p, legalCryptFileExtension, "", -1)
}

func appendTrailingLCExt(p string) string {
	if strings.HasSuffix(p, ".") {
		p = strings.Replace(p, ".", "", -1)
	}
	return p + legalCryptFileExtension
}
