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
)

func check(e error) {
	if e != nil {
		fmt.Println("ERR::PANIC")
		panic(e)
	}
}

func generateRSA4096(secret string) (privateKey []byte, publicKey []byte, err error) {
	// private1 *rsa.PrivateKey;
	// err error;
	private1, err := rsa.GenerateKey(rand.Reader, 4096)

	if err != nil {
		fmt.Println("ERR::Failed to generate private key.")
		return nil, nil, err
	}

	err = private1.Validate()
	if err != nil {
		fmt.Println("ERR::Validation failed.")
		return nil, nil, err
	}

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
	private3, err := x509.EncryptPEMBlock(rand.Reader, private2.Type, private2.Bytes, []byte(secret), x509.PEMCipherAES256)
	if err != nil {
		fmt.Println("ERR::Failed to encrupt PEM private block")
		return nil, nil, err
	}

	// Resultant private key in PEM format.
	// priv_pem string
	privatePem := pem.EncodeToMemory(private3)

	// Public Key generation
	public1 := private1.PublicKey

	publicDer, err := x509.MarshalPKIXPublicKey(&public1)
	if err != nil {
		fmt.Println("ERR::Failed to get der format for PublicKey.")
		return nil, nil, err
	}

	public2 := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   publicDer,
	}
	publicPem := pem.EncodeToMemory(&public2)

	return privatePem, publicPem, nil
}

func writeKeyPair(privatePem []byte, publicPem []byte, encType string) {
	privateFilename := "./id_" + strings.ToLower(encType) + ".private.pem"
	publicFilename := "./id_" + strings.ToLower(encType) + ".public.pem"

	_ = ioutil.WriteFile(privateFilename, privatePem, 0600)
  _ = ioutil.WriteFile(publicFilename,  publicPem, 0644)
}

func main() {
	passphrase := "tina_belcher_loves_butts"

	privatePem, publicPem, err := generateRSA4096(passphrase)
	if err != nil {
		fmt.Println("ERR::Something has gone awry.", err)
		return
	}

	writeKeyPair(privatePem, publicPem, "rsa")

	// Parsable output <STATUS>::<SZ_PRIV_KEY>::<SZ_PUB_KEY>
	fmt.Println("OK")
}
