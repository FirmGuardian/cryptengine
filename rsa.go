package main

import (
  "errors"
  "crypto/rand"
  "crypto/x509"
  "io/ioutil"
  "crypto/rsa"
  "encoding/pem"
  "encoding/base64"
  "golang.org/x/crypto/sha3"
  "bufio"
  "os"
)

func decryptRSA(filePath string) {
  fileInfo, _ := os.Stat(filePath)
  fileSize := fileInfo.Size()

  inFile, err := os.Open(filePath)
  check(err, "Unable to open encrypted file for reading")
  r := bufio.NewReaderSize(inFile, int(fileSize))

  publicKeyHash64 := make([]byte, 88)
  _, err = r.Read(publicKeyHash64)
  check(err, "Error reading publicKeyHash")

  encryptedKey64 := make([]byte, 684)
  _, err = r.Read(encryptedKey64)
  check(err, "Error reading encrypted key")
  encryptedKey, err := base64.StdEncoding.DecodeString(string(encryptedKey64))

  nonce64 := make([]byte, 16)
  _, err = r.Read(nonce64)
  nonce, _ := base64.StdEncoding.DecodeString(string(nonce64))

  szEncryptedData := uint64(r.Buffered())

  if szEncryptedData > maxInputFileSize + 4096 {
    check(errors.New("Encrypted file is larger than maximum!"), "Encrypted file is larger than maximum!")
  }

  encryptedData64 := make([]byte, szEncryptedData) // max encryptable filesize + pad
  _, err = r.Read(encryptedData64)
  inFile.Close()
  check(err, "Error reading encryptedData")

  encryptedData, err := base64.StdEncoding.DecodeString(string(encryptedData64))
  check(err, "Unable to deserialize encrypted data")

  keySlurp, err :=ioutil.ReadFile("./id_rsa")
  check(err, "Unable to read private key")
  privateBlock, _ := pem.Decode(keySlurp)
  if privateBlock == nil || privateBlock.Type != "RSA PRIVATE KEY" {
    check(errors.New("Failed to decode PEM block containing private key"), "Failed to decode PEM block containing private key")
  }

  der, err := x509.DecryptPEMBlock(privateBlock, []byte(constPassphrase))
  check(err, "Unable to decrypt private block")

  privatePKCS, err := x509.ParsePKCS1PrivateKey(der)
  check(err, "Unable to parse decrypted private key")

  hash := sha3.New512()
  rng := rand.Reader
  sessionKey, err := rsa.DecryptOAEP(hash, rng, privatePKCS, encryptedKey, []byte(""))
  check(err, "Unable to decrypt cipherkey")

  // BEGIN AES DECRYPT (sessionKey, nonce, encryptedData)
  decryptedData, err := decryptAES(sessionKey, nonce, encryptedData)
  check(err, "Unable to decrypt file")
  // END AES DECRYPT

  outFilePath, err := getDecryptedFilename(filePath)
  nixIfExists(outFilePath)
  outFile, err := os.Create(outFilePath)
  check(err, "Unable to create output file")
  defer outFile.Close()
  w := bufio.NewWriter(outFile)

  w.Write(decryptedData)
  w.Flush()
}

func encryptRSA(filePath string) (error) {
  outFilePath := getEncryptedFilename(filePath)
  nixIfExists(outFilePath)

  // Create output file, and Writer
  outFile, err := os.Create(outFilePath)
  defer outFile.Close()
  w := bufio.NewWriter(outFile)

  // Slurp and parse public key to encrypt AES Session Key
  keySlurp, err := ioutil.ReadFile("./id_rsa.pub")
  check(err, "Unable to read public key")
  publicBlock, _ := pem.Decode(keySlurp)
  if publicBlock == nil || publicBlock.Type != "PUBLIC KEY" {
    check(errors.New("Failed to decode PEM block containing public key"), "Failed to decode PEM block containing public key")
  }

  publicKey, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
  check(err, "Unable to parse public key")

  publicKeyHash := sha3.New512()
  publicKeyHash.Write(publicBlock.Bytes)

  w.WriteString(base64.StdEncoding.EncodeToString(publicKeyHash.Sum(nil)))

  hash := sha3.New512()
  rng := rand.Reader

  // Slurp file to be encrypted
  fSlurp, err := ioutil.ReadFile(filePath)
  check(err, "Unable to read file to be secured.")

  // BEGIN AES ENCRYPTION
  encryptedBin, nonce, sessionKey, err := encryptAES(fSlurp)

  encryptedSessionKey, err := rsa.EncryptOAEP(hash, rng, publicKey.(*rsa.PublicKey), sessionKey, []byte(""))
  encryptedSessionKey64 := base64.StdEncoding.EncodeToString(encryptedSessionKey)

  w.Write([]byte(encryptedSessionKey64))

  nonce64 := base64.StdEncoding.EncodeToString(nonce)
  w.Write([]byte(nonce64))

  encrypted64 := base64.StdEncoding.EncodeToString(encryptedBin)
  w.Write([]byte(encrypted64))

  w.Flush()

  return nil
}

func generateRSA4096(secret string) {
  privateFilename := "./id_rsa"
  publicFilename  := privateFilename + ".pub"
  if fileExists(privateFilename) && fileExists(publicFilename) {
    return
  }

  rng := rand.Reader

  // private1 *rsa.PrivateKey;
  // err error;
  private1, err := rsa.GenerateKey(rng, 4096)
  check(err, "Failed to generate private key.")

  err = private1.Validate()
  check(err, "Validation failed.")

  privateDer := x509.MarshalPKCS1PrivateKey(private1)

  // pem.Block
  // blk pem.Block
  private2 := pem.Block{
    Type:    "RSA PRIVATE KEY",
    Headers: nil,
    Bytes:   privateDer,
  }

  // Encrypt the pem
  private3, err := x509.EncryptPEMBlock(rng, private2.Type, private2.Bytes, []byte(secret), x509.PEMCipherAES256)
  check(err, "Failed to encrypt PEM private block")

  // Resultant private key in PEM format.
  // priv_pem string
  privatePem := pem.EncodeToMemory(private3)

  // Public Key generation
  publicDer, err := x509.MarshalPKIXPublicKey(&private1.PublicKey)
  check(err, "Failed to get der format for PublicKey.")

  public2 := pem.Block{
    Type:    "PUBLIC KEY",
    Headers: nil,
    Bytes:   publicDer,
  }
  publicPem := pem.EncodeToMemory(&public2)

  _ = ioutil.WriteFile(privateFilename, privatePem, 0400)
  _ = ioutil.WriteFile(publicFilename,  publicPem, 0644)
}

