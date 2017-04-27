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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
  "io/ioutil"
	"strings"
  //"golang.org/x/crypto/ed25519"
  //"github.com/mikesmitty/edkey"
  "golang.org/x/crypto/sha3"
  "crypto"
)

func check(e error, msg string) {
	if e != nil {
		fmt.Println("ERR::PANIC")
		panic(e)
	}
}

func decryptRSA(secret string) {
  hash := sha3.New512()
}

func encryptRSA(secret string, filePath string) {
  hash := sha3.New512()
  rng := rand.Reader

  keySlurp, err := ioutil.ReadFile("./id_rsa.pub")
  publicBlock, _ := pem.Decode(keySlurp)
  var publicCert* x509.Certificate
  publicCert, _ = x509.ParseCertificate(publicCert.Bytes)
  publicKey := publicCert.PublicKey.(*rsa.PublicKey)

  fSlurp, err := ioutil.ReadFile(filePath)
  check(err, "ERR::Unable to read file to be secured.")
}

func generateRSA4096(secret string) (privateKey []byte, publicKey []byte, err error) {
  rng := rand.Reader

	// private1 *rsa.PrivateKey;
	// err error;
	private1, err := rsa.GenerateKey(rand.Reader, 4096)

	check(err, "ERR::Failed to generate private key.")

	err = private1.Validate()
	check(err, "ERR::Validation failed.")

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
	check(err, "ERR::Failed to encrupt PEM private block")

	// Resultant private key in PEM format.
	// priv_pem string
	privatePem := pem.EncodeToMemory(private3)

	// Public Key generation
	public1 := private1.PublicKey

	publicDer, err := x509.MarshalPKIXPublicKey(&public1)
	check(err, "ERR::Failed to get der format for PublicKey.")

	public2 := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   publicDer,
	}
	publicPem := pem.EncodeToMemory(&public2)

	return privatePem, publicPem, nil
}

func writeKeyPair(privatePem []byte, publicPem []byte, encType string) {
	privateFilename := "./id_" + strings.ToLower(encType)
	publicFilename := "./id_" + strings.ToLower(encType) + ".pub"

	_ = ioutil.WriteFile(privateFilename, privatePem, 0400)
  _ = ioutil.WriteFile(publicFilename,  publicPem, 0644)
}

func main() {
	passphrase := "t1n@ b3LcHeR_lov3s!bUtts+"

	privatePem, publicPem, err := generateRSA4096(passphrase)
	check(err, "ERR::Something has gone awry.")

	writeKeyPair(privatePem, publicPem, "rsa")

	// Parsable output <STATUS>::<SZ_PRIV_KEY>::<SZ_PUB_KEY>
	fmt.Println("OK")
}
