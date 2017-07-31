package main

import (
  "github.com/jhoonb/archivex"
  "os"
  "fmt"
)

func archiveFiles(paths []string) string {
  cwd, err := os.Getwd()
  fmt.Println(";;CWD::" + cwd)

  archivePath := ""

  if (err != nil) {
    archivePath = "./lcsf_secured_files.zip"
  } else {
    archivePath = cwd + "/lcsf_secured_files.zip"
  }

  archive := new(archivex.ZipFile)
  archive.Create(archivePath)
  defer archive.Close()

  var skipped []string

  for _, path := range paths {
    fileInfo, err := os.Stat(path)
    check(err, "Something's fucky with " + path)

    fmode := fileInfo.Mode()
    isdir := fmode.IsDir()
    isreg := fmode.IsRegular()

    if isreg {
      archive.AddFile(path)
      fmt.Printf("ZIP::Adding file %s\n", path)
    } else if isdir {
      archive.AddAll(path, false)
      fmt.Printf("ZIP::Adding directory %s\n", path)
    } else {
      skipped = append(skipped, path)
    }

    for _, spath := range skipped {
      fmt.Printf("SKIPPED::%s\n", spath)
    }
  }

  return archivePath
}
