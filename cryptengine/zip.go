package main

import (
  "github.com/jhoonb/archivex"
  "os"
  "fmt"
)

func archiveFiles(paths []string) string {
  zipPath := "./lcsf_secured_files.zip"

  zip := new(archivex.ZipFile)
  zip.Create(zipPath)

  var skipped []string

  for _, path := range paths {
    fileInfo, err := os.Stat(path)
    check(err, "Something's fucky with " + path)

    fmode := fileInfo.Mode()
    isdir := fmode.IsDir()
    isreg := fmode.IsRegular()

    if isreg {
      zip.AddFile(path)
      fmt.Printf("ZIP::Adding file %s\n", path)
    } else if isdir {
      zip.AddAll(path, false)
      fmt.Printf("ZIP::Adding directory %s\n", path)
    } else {
      skipped = append(skipped, path)
    }

    zip.Close()

    for _, spath := range skipped {
      fmt.Printf("SKIPPED::%s\n", spath)
    }
  }

  return zipPath
}
