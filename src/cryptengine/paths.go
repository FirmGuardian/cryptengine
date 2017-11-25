package main

import (
	"fmt"
	"os"
	"path"
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
			segments:    []string{root, "private", "bin"},
			permissions: 0700,
		},
		{
			segments:    []string{root, "private", "keys"},
			permissions: 0600,
		},
		{
			segments:    []string{root, "private", "tmp"},
			permissions: 0600,
		},
	}

	for _, lcPath := range lcPaths {
		testPath := path.Join(lcPath.segments...)
		fmt.Printf("%v...", testPath)
		err := os.MkdirAll(testPath, lcPath.permissions)
		if err != nil {
			fmt.Println("ERR")
		} else {
			fmt.Println("Done")
		}
	}
}
