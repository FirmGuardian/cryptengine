package main

import (
  "crypto/aes"
  "crypto/cipher"
  "io"
  "crypto/rand"
)

const lenAESNonce uint8 = 12

func decryptAES(key []byte, nonce []byte, encryptedData []byte) ([]byte, error) {
  aesCipher, err := aes.NewCipher(key)
  check(err, "Unable to create AES cipher")

  aesgcm, err := cipher.NewGCM(aesCipher)
  check(err, "Unable to create GCM Block")

  decryptedData, err := aesgcm.Open(nil, nonce, encryptedData, nil)
  check(err, "Unable to decrypt data")

  return decryptedData, err
}

func encryptAES(unencryptedData []byte) ([]byte, []byte, []byte, error) {
  rng := rand.Reader

  // Generate AES Session Key; to be RSA encrypted, and used to
  // encrypt input file
  key, err := generateRandomBytes(32)
  check(err, "Unable to generate sessionKey")

  aesCipher, err := aes.NewCipher(key)
  check(err, "Unable to create AES cipher")

  // Generate nonce
  nonce, _ := generateRandomBytes(int(lenAESNonce))

  aesgcm, err := cipher.NewGCM(aesCipher)
  check(err, "Unable to create new AES cipher")

  encryptedBin := aesgcm.Seal(nil, nonce, unencryptedData, nil)

  return encryptedBin, nonce, key, err
}
