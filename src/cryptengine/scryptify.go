package main

import (
  "encoding/base64"
  "golang.org/x/crypto/scrypt"
)

func scryptify(str string, keyLen int) []byte {
  passwd := []byte(str)
  salt := []byte("liam@storskegg.org")
  N := 512 * 1024
  r := 32
  p := 2

  key, _ := scrypt.Key(passwd, salt, N, r, p, keyLen)

  return key
}

func scryptify64(str string, keyLen int) string {
  return base64.StdEncoding.EncodeToString(scryptify(str, keyLen))
}
