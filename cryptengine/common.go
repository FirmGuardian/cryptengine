package main

import (
  "crypto/rand"
  "fmt"
  "strings"
  "io/ioutil"
  "os"
)

func check(e error, msg string) {
  if e != nil {
    fmt.Println("ERR::" + msg)
    panic(e)
  }
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

func writeKeyPair(privatePem []byte, publicPem []byte, encType string) {
  privateFilename := "./id_" + strings.ToLower(encType)
  publicFilename := "./id_" + strings.ToLower(encType) + ".pub"

  _ = ioutil.WriteFile(privateFilename, privatePem, 0400)
  _ = ioutil.WriteFile(publicFilename,  publicPem, 0644)
}

