package main

import (
  "errors"
  "crypto/rand"
  "crypto/x509"
  "io/ioutil"
  "crypto/rsa"
  "encoding/pem"
  "golang.org/x/crypto/sha3"
  "bufio"
  "os"
  "crypto/aes"
  "io"
  "crypto/cipher"
)

//func decryptRSA(secret string) {
//  hash := sha3.New512()
//}

func encryptRSA(filePath string) (error) {
  // Create output file, and Writer
  of, err := os.Create(filePath + ".encrypted")
  defer of.Close()
  w := bufio.NewWriter(of)

  // Generate AES Session Key; to be RSA encrypted, and used to
  // encrypt input file
  sessionKey, err := generateRandomBytes(32)
  check(err, "Unable to generate sessionKey")

  // Slurp and parse public key to encrypt AES Session Key
  keySlurp, err := ioutil.ReadFile("./id_rsa.pub")
  publicBlock, _ := pem.Decode(keySlurp)
  if publicBlock == nil || publicBlock.Type != "PUBLIC KEY" {
    check(errors.New("Failed to decode PEM block containing public key"), "Failed to decode PEM block containing public key")
  }

  publicKey, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
  check(err, "Unable to parse public key")

  hash := sha3.New512()
  rng := rand.Reader

  encryptedSessionKey, err := rsa.EncryptOAEP(hash, rng, publicKey.(*rsa.PublicKey), sessionKey, []byte(""))

  w.Write(encryptedSessionKey)

  // Create new AES block
  aesCipher, err := aes.NewCipher(sessionKey)
  check(err, "Unable to create AES cipher")

  // Generate nonce
  nonce := make([]byte, 12)
  if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
    panic(err.Error())
  }

  // Slurp file to be encrypted
  fSlurp, err := ioutil.ReadFile(filePath)
  check(err, "Unable to read file to be secured.")

  // Do the deed.
  aesgcm, err := cipher.NewGCM(aesCipher)
  check(err, "Unable to create new AES cipher")

  encrypted := aesgcm.Seal(fSlurp, nonce, fSlurp, nil)

  w.Write(encrypted)
  w.Flush()

  return nil
}

func generateRSA4096(secret string) (privateKey []byte, publicKey []byte, err error) {
  rng := rand.Reader

  // private1 *rsa.PrivateKey;
  // err error;
  private1, err := rsa.GenerateKey(rand.Reader, 4096)

  check(err, "Failed to generate private key.")

  err = private1.Validate()
  check(err, "Validation failed.")

  // Get der format. priv_der []byte
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
  check(err, "Failed to encrupt PEM private block")

  // Resultant private key in PEM format.
  // priv_pem string
  privatePem := pem.EncodeToMemory(private3)

  // Public Key generation
  public1 := private1.PublicKey

  publicDer, err := x509.MarshalPKIXPublicKey(&public1)
  check(err, "Failed to get der format for PublicKey.")

  public2 := pem.Block{
    Type:    "PUBLIC KEY",
    Headers: nil,
    Bytes:   publicDer,
  }
  publicPem := pem.EncodeToMemory(&public2)

  return privatePem, publicPem, nil
}

