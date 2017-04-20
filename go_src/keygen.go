/*
 * Generates a private/public key pair in PEM format (not Certificate)
 *
 * The generated private key can be parsed with openssl as follows:
 * > openssl rsa -in key.pem -text
 *
 * The generated public key can be parsed as follows:
 * > openssl rsa -pubin -in pub.pem -text
 *
 * TODO: Clean this up, it's gross.
 * TODO: Add support for additional encryption methods
 */
package main

import (
  "crypto/rsa"
  "crypto/rand"
  "crypto/x509"
  "encoding/pem"
  "fmt"
  "os"
)

func check(e error) {
  if e != nil {
    panic(e)
  }
}

func main() {
  passphrase := "tinabelcherlovesbutts"

  // priv *rsa.PrivateKey;
  // err error;
  priv, err := rsa.GenerateKey(rand.Reader, 4096)
  if err != nil {
    fmt.Println(err)
    return
  }
  err = priv.Validate()
  if err != nil {
    fmt.Println("Validation failed.", err)
  }

  // Get der format. priv_der []byte
  priv_der := x509.MarshalPKCS1PrivateKey(priv)

  // pem.Block
  // blk pem.Block
  priv_blk := pem.Block {
    Type: "RSA PRIVATE KEY",
    Headers: nil,
    Bytes: priv_der,
  }

  // Encrypt the pem
  crypted_blk, err := x509.EncryptPEMBlock(rand.Reader, priv_blk.Type, priv_blk.Bytes, []byte(passphrase), x509.PEMCipherAES256)
  if err != nil {
    fmt.Println(err)
    return
  }

  // Resultant private key in PEM format.
  // priv_pem string
  priv_pem := string(pem.EncodeToMemory(crypted_blk))

  f_priv, err := os.Create("./id_rsa.private.pem")
  check(err)
  defer f_priv.Close()

  n_priv, err := f_priv.WriteString(priv_pem)

  f_priv.Sync()

  // Public Key generation
  pub := priv.PublicKey

  pub_der, err := x509.MarshalPKIXPublicKey(&pub)
  if err != nil {
    fmt.Println("Failed to get der format for PublicKey.", err)
    return
  }

  pub_blk := pem.Block {
    Type: "PUBLIC KEY",
    Headers: nil,
    Bytes: pub_der,
  }
  pub_pem := string(pem.EncodeToMemory(&pub_blk))

  f_pub, err := os.Create("./id_rsa.public.pem")
  check(err)
  defer f_pub.Close()

  n_pub, err := f_pub.WriteString(pub_pem)

  f_pub.Sync()

  // Parsable output <STATUS>::<SZ_PRIV_KEY>::<SZ_PUB_KEY>
  fmt.Printf("OK::%d::%d\n", n_priv, n_pub)
}
