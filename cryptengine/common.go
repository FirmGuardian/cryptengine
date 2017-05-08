package main

import (
  "crypto/rand"
  "fmt"
  "os"
)

const maxInputFileSize uint64 = 1024 * 1024 * 1024 * 2 // 2GB; uint64 to support >= 4GB

func check(e error, msg string) {
  if e != nil {
    fmt.Fprintln(os.Stderr, "ERR::" + msg)
    panic(e)
  }
}

func fileExists(filePath string) (bool) {
  if _, err := os.Stat(filePath); err == nil {
    return true
  }

  return false
}

func generateRandomBytes(s int) ([]byte, error) {
  b := make([]byte, s)
  n, err := rand.Read(b)
  if n != len(b) || err != nil {
    return nil, fmt.Errorf("Unable to successfully read from the system CSPRNG (%v)", err)
  }

  return b, nil
}

func nixIfExists(filePath string) {
  if _, err := os.Stat(filePath); err == nil {
    check(os.Remove(filePath), "Unable to remove existing file")
  }
}
