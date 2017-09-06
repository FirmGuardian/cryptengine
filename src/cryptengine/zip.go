package main

import (
	"fmt"
	"github.com/jhoonb/archivex"
	"os"
)

func archiveFiles(paths []string) string {
	cwd, err := os.Getwd()
	fmt.Println(";;CWD::" + cwd)

	archivePath := ""

	if err != nil {
		archivePath = "./lcsf_secured_files.zip"
	} else {
		archivePath = cwd + "/lcsf_secured_files.zip"
	}

	archive := new(archivex.ZipFile)
	archive.Create(archivePath)
	defer archive.Close()

	// This will hold the paths of any file we skip
	var skippedFiles []string

	for _, path := range paths {
		fileInfo, err := os.Stat(path)
		check(err, errs["fsCantOpenFile"])

		fmode := fileInfo.Mode()
		isDirectory := fmode.IsDir()
		isRegular := fmode.IsRegular()

		if isRegular {
			archive.AddFile(path)
			fmt.Printf("ZIP::Adding file %s\n", path)
		} else if isDirectory {
			archive.AddAll(path, false)
			fmt.Printf("ZIP::Adding directory %s\n", path)
		} else {
			skippedFiles = append(skippedFiles, path)
		}

		for _, spath := range skippedFiles {
			fmt.Printf("SKIPPED::%s\n", spath)
		}
	}

	return archivePath
}
