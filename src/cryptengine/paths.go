package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
)

// LCPath makes for os-independent path declarations for scaffolding
type LCPath struct {
	segments    []string
	permissions os.FileMode
}

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

func keyDir() string {
	return appRoot() + pathSeparator() + "keys"
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
	finfo, err := os.Stat(p)
	if err == nil {
		exists = true
		mode := finfo.Mode()
		isdir = mode.IsDir()
		isreg = mode.IsRegular()
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
	}
}

func scaffoldAppDirs() {
	root := appRoot()

	lcPaths := []LCPath{
		{
			segments:    []string{root, "bin"},
			permissions: 0700,
		},
		{
			segments:    []string{root, "keys"},
			permissions: 0600,
		},
		{
			segments:    []string{root, "tmp"},
			permissions: 0600,
		},
	}

	for _, lcPath := range lcPaths {
		p := path.Join(lcPath.segments...)
		fmt.Printf("%v...", p)
		err := os.MkdirAll(p, lcPath.permissions)
		if err != nil {
			fmt.Println("ERR")
		} else {
			fmt.Println("Done")
		}
	}
}

func pathEndsWith(haystack string, needle string) bool {
	lastIndex := strings.LastIndex(haystack, needle)
	return 0 == len(haystack)-lastIndex-len(needle)
}

func pathEndsWithLCSF(path string) bool {
	return pathEndsWith(path, legalCryptFileExtension)
}

//func pathEndsWithSeparator(path string) bool {
//	return pathEndsWith(path, pathSeparator())
//}

// TODO: uncomment to implement zipping in tmp after outPath implemented
//func tmpDir() string {
//	return appRoot() + pathSeparator() + "tmp"
//}

func pathSeparator() string {
	return string(os.PathSeparator)
}
