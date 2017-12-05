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
	antiCollision := strconv.FormatInt(time.Now().UnixNano(), 36) + "-"
	archivePath := path.Join(tmpDir(), antiCollision+lcszFileName)

	var addedFiles []string
	var skippedFiles []string

	for _, path := range paths {
		fileInfo := pathInfo(path)
		if !fileInfo.Exists {
			skippedFiles = append(skippedFiles, path)
		} else if fileInfo.IsDir || fileInfo.IsReg {
			addedFiles = append(addedFiles, path)
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
