package main

import (
  "errors"
  "crypto/rand"
  "crypto/x509"
  "io/ioutil"
  "crypto/rsa"
  "encoding/pem"
  "golang.org/x/crypto/sha3"
)

//func decryptRSA(secret string) {
//  hash := sha3.New512()
//}

func encryptRSA(filePath string) ([]byte) {
  hash := sha3.New512()
  rng := rand.Reader

  //sessionKey, err := generateRandomBytes(32)
  //check(err, "Unable to generate sessionKey")

  keySlurp, err := ioutil.ReadFile("./id_rsa.pub")
  publicBlock, _ := pem.Decode(keySlurp)
  if publicBlock == nil || publicBlock.Type != "PUBLIC KEY" {
    check(errors.New("Failed to decode PEM block containing public key"), "Failed to decode PEM block containing public key")
  }

  publicKey, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
  check(err, "Unable to parse public key")

  fSlurp, err := ioutil.ReadFile(filePath)
  check(err, "Unable to read file to be secured.")

  encrypted, err := rsa.EncryptOAEP(hash, rng, publicKey.(*rsa.PublicKey), fSlurp, []byte(""))
  check(err, "Failed to encrypt file.")

  return encrypted
}

func writeEncryptedFile(filename string, data []byte) (error) {
  return ioutil.WriteFile(filename, data, 0400)
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

