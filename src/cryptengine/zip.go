package main

import (
	"errors"
	"fmt"
	"github.com/jhoonb/archivex"
	"path"
)

const lcszFileName = "lcsf_secured_files.zip"

func archiveFiles(paths []string) string {
	archivePath := path.Join(tmpDir(), lcszFileName)

	archive := new(archivex.ZipFile)
	archive.Create(archivePath)
	defer archive.Close()

	// This will hold the paths of any file we skip
	var skippedFiles []string

	for _, path := range paths {
		fileInfo := pathInfo(path)
		if !fileInfo.Exists {
			check(errors.New(errs["fsCantOpenFile"].Msg), errs["fsCantOpenFile"])
		}

		if fileInfo.IsReg {
			archive.AddFile(path)
			fmt.Printf("ZIP::Adding file %s\n", path)
		} else if fileInfo.IsDir {
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
