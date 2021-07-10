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
  "errors"
)

const lcszFileName = "lcsf_secured_files"

func archiveFiles(paths []string) string {
  if len(paths) == 0 {
    return ""
  }

	antiCollision := "-" + strconv.FormatInt(time.Now().UnixNano(), 36)
	archivePath := path.Join(tmpDir(), lcszFileName + antiCollision + ".zip")

	var addedFiles []string
	var action string

	for _, path := range paths {
		fileInfo := pathInfo(path)
		if fileInfo.Exists && (fileInfo.IsDir || fileInfo.IsReg) {
			addedFiles = append(addedFiles, path)
			action = "ADDED"
		} else {
      action = "SKIPPED"
    }
    fmt.Printf("%s::%s\n", action, path)
	}

	if len(addedFiles) == 0 {
	  return ""
  }

  // We are required to make sure output path exists otherwise archive fails to zip
  // This could change? https://github.com/mholt/archiver/issues/61
  os.MkdirAll(tmpDir(), 700)
	err := archiver.Zip.Make(archivePath, addedFiles)
	if err != nil {
		// TODO: you know the drill...
		log.Fatalln(err)
	}

	return archivePath
}
// unarchiveFiles will unzip a file and place it in the cryptengine output directory, and then delete the zip file.
// The path of the unzipped file is returned as well as any errors that may have occurred.
func unarchiveFiles(zipFile string) (string, error) {
	zInfo := pathInfo(zipFile)
	if !zInfo.Exists || !zInfo.IsReg {
		// TODO: You know what to do, here
		return "", errors.New(fmt.Sprintf("Cannot extract! %v is not regular, or file does not exist!\n", zInfo.Clean))
	}
	outPath := path.Join(outDirDec(), strings.Replace(zInfo.File, zInfo.Ext, "", -1))
	os.MkdirAll(outPath, 0700)
	err := archiver.Zip.Open(zInfo.Clean, outPath)
	if err != nil {
		return "", err
	}

	return outPath, os.Remove(zInfo.Clean)
}
