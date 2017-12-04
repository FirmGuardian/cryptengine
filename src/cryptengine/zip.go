package main

import (
	"fmt"
	"log"
	"path"
	"strconv"
	"time"

	"github.com/mholt/archiver"
)

const lcszFileName = "lcsf_secured_files.zip"

func archiveFiles(paths []string) string {
	archivePath := path.Join(tmpDir(), strconv.FormatInt(time.Now().UnixNano(), 36)+"-"+lcszFileName)

	// This will hold the paths of any file that we'll add
	var addedFiles []string
	// This will hold the paths of any file we'll skip
	var skippedFiles []string

	for _, path := range paths {
		fileInfo := pathInfo(path)
		if !fileInfo.Exists {
			//check(errors.New(errs["fsCantOpenFile"].Msg), errs["fsCantOpenFile"])
			skippedFiles = append(skippedFiles, path)
		} else if fileInfo.IsDir || fileInfo.IsReg {
			addedFiles = append(addedFiles, path)
			//archive.AddFile(path)
			//fmt.Printf("ZIP::Adding file %s\n", path)
			//} else if fileInfo.IsDir {
			//
			//	archive.AddAll(path, false)
			//	fmt.Printf("ZIP::Adding directory %s\n", path)
		} else {
			skippedFiles = append(skippedFiles, path)
		}

		for _, spath := range skippedFiles {
			fmt.Printf("SKIPPED::%s\n", spath)
		}
	}

	err := archiver.Zip.Make(archivePath, addedFiles)
	if err != nil {
		// TODO: you know the drill...
		log.Fatalln(err)
	}

	return archivePath
}
