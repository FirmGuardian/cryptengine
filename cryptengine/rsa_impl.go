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
  "crypto/aes"
  "io"
  "crypto/cipher"
)

func decryptRSA(filePath string) {
  inFile, err := os.Open(filePath)
  check(err, "Unable to open encrypted file for reading")
  r := bufio.NewReader(inFile)

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

  encryptedData64 := make([]byte, r.Buffered()) // max encryptable filesize + pad
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

  aesCipher, err := aes.NewCipher(sessionKey)
  check(err, "Unable to create AES cipher")

  aesgcm, err := cipher.NewGCM(aesCipher)
  check(err, "Unable to create GCM Block")

  decryptedData, err := aesgcm.Open(nil, nonce, encryptedData, nil)
  check(err, "Unable to decrypt data")

  nixIfExists(filePath + ".decrypted")
  outFile, err := os.Create(filePath + ".decrypted")
  check(err, "Unable to create output file")
  defer outFile.Close()
  w := bufio.NewWriter(outFile)

  w.Write(decryptedData)
  w.Flush()
}

func encryptRSA(filePath string) (error) {
  outFilePath := filePath + ".encrypted"
  nixIfExists(outFilePath)

  // Create output file, and Writer
  outFile, err := os.Create(outFilePath)
  defer outFile.Close()
  w := bufio.NewWriter(outFile)

  // Generate AES Session Key; to be RSA encrypted, and used to
  // encrypt input file
  sessionKey, err := generateRandomBytes(32)
  check(err, "Unable to generate sessionKey")

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

  encryptedSessionKey, err := rsa.EncryptOAEP(hash, rng, publicKey.(*rsa.PublicKey), sessionKey, []byte(""))

  encryptedSessionKey64 := base64.StdEncoding.EncodeToString(encryptedSessionKey)

  w.Write([]byte(encryptedSessionKey64))

  // Create new AES block
  aesCipher, err := aes.NewCipher(sessionKey)
  check(err, "Unable to create AES cipher")

  // Generate nonce
  nonce := make([]byte, lenAESNonce)
  if _, err := io.ReadFull(rng, nonce); err != nil {
    panic(err.Error())
  }

  nonce64 := base64.StdEncoding.EncodeToString(nonce)

  w.Write([]byte(nonce64))

  // Slurp file to be encrypted
  fSlurp, err := ioutil.ReadFile(filePath)
  check(err, "Unable to read file to be secured.")

  // Do the deed.
  aesgcm, err := cipher.NewGCM(aesCipher)
  check(err, "Unable to create new AES cipher")

  encryptedBin := aesgcm.Seal(nil, nonce, fSlurp, nil)

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

