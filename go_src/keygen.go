/*
 * Generates a private/public key pair in PEM format (not Certificate)
 *
 * The generated private key can be parsed with openssl as follows:
 * > openssl rsa -in key.pem -text
 *
 * The generated public key can be parsed as follows:
 * > openssl rsa -pubin -in pub.pem -text
 *
 * TODO: Add support for additional encryption methods -Liam
 */
package main

import (
  "crypto/rsa"
  "crypto/rand"
  "crypto/x509"
  "encoding/pem"
  "fmt"
  "os"
  "strings"
)

func check(e error) {
  if e != nil {
    fmt.Println("ERR::PANIC")
    panic(e)
  }
}

func generateRSA4096(secret string) (privateKey string, publicKey string, err error) {
  // private_1 *rsa.PrivateKey;
  // err error;
  private_1, err := rsa.GenerateKey(rand.Reader, 4096)

  if err != nil {
    fmt.Println("ERR::Failed to generate private key.")
    return "", "", err
  }

  err = private_1.Validate()
  if err != nil {
    fmt.Println("ERR::Validation failed.")
    return "", "", err
  }

  // Get der format. priv_der []byte
  private_der := x509.MarshalPKCS1PrivateKey(private_1)

  // pem.Block
  // blk pem.Block
  private_2 := pem.Block {
    Type: "RSA PRIVATE KEY",
    Headers: nil,
    Bytes: private_der,
  }

  // Encrypt the pem
  private_3, err := x509.EncryptPEMBlock(rand.Reader, private_2.Type, private_2.Bytes, []byte(secret), x509.PEMCipherAES256)
  if err != nil {
    fmt.Println("ERR::Failed to encrupt PEM private block")
    return "", "", err
  }

  // Resultant private key in PEM format.
  // priv_pem string
  private_pem := string(pem.EncodeToMemory(private_3))

  // Public Key generation
  public_1 := private_1.PublicKey

  public_der, err := x509.MarshalPKIXPublicKey(&public_1)
  if err != nil {
    fmt.Println("ERR::Failed to get der format for PublicKey.")
    return "", "", err
  }

  public_2 := pem.Block {
    Type: "PUBLIC KEY",
    Headers: nil,
    Bytes: public_der,
  }
  public_pem := string(pem.EncodeToMemory(&public_2))

  return private_pem, public_pem, nil
}

func writeKeyPair(privatePem string, publicPem string, encType string) (n_private int, n_public int) {
  privateFilename := "./id_" + strings.ToLower(encType) + ".private.pem"
  publicFilename  := "./id_" + strings.ToLower(encType) + ".public.pem"

  f_private, err := os.Create(privateFilename)
  check(err)
  defer f_private.Close()

  f_public, err := os.Create(publicFilename)
  check(err)
  defer f_public.Close()

  n_private, err = f_private.WriteString(privatePem)
  if err != nil {
    fmt.Println("ERR::Failed to write private key file")
    return 0, 0
  }
  n_public, err  = f_public.WriteString(publicPem)
  if err != nil {
    fmt.Println("ERR::Failed to write public key file")
    return n_private, 0
  }

  f_private.Sync()
  f_public.Sync()

  return n_private, n_public
}

func main() {
  passphrase := "tina_belcher_loves_butts"

  privatePem, publicPem, err := generateRSA4096(passphrase)
  if err != nil {
    fmt.Println("ERR::Something has gone awry.", err)
    return
  }


  n_priv, n_pub := writeKeyPair(privatePem, publicPem, "rsa")

  // Parsable output <STATUS>::<SZ_PRIV_KEY>::<SZ_PUB_KEY>
  fmt.Printf("OK::%d::%d\n", n_priv, n_pub)
}
