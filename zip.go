package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/mholt/archiver"
)

const lcszFileName = "lcsf_secured_files"

func archiveFiles(paths []string) string {
	antiCollision := "-" + strconv.FormatInt(time.Now().UnixNano(), 36)
	archivePath := path.Join(tmpDir(), lcszFileName+antiCollision+".zip")

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

func unarchiveFiles(zipFile string) error {
	zInfo := pathInfo(zipFile)
	if !zInfo.Exists || !zInfo.IsReg {
		// TODO: You know what to do, here
		log.Fatalf("Cannot extract! %v is not regular, or file does not exist!\n", zInfo.Clean)
	}
	outPath := path.Join(outDirDec(), strings.Replace(zInfo.File, zInfo.Ext, "", -1))
	os.MkdirAll(outPath, 0700)
	err := archiver.Zip.Open(zInfo.Clean, outPath)
	if err != nil {
		return err
	}

	return os.Remove(zInfo.Clean)
}
