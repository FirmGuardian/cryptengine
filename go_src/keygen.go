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
	// private1 *rsa.PrivateKey;
	// err error;
	private1, err := rsa.GenerateKey(rand.Reader, 4096)

	if err != nil {
		fmt.Println("ERR::Failed to generate private key.")
		return "", "", err
	}

	err = private1.Validate()
	if err != nil {
		fmt.Println("ERR::Validation failed.")
		return "", "", err
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
		return "", "", err
	}

	// Resultant private key in PEM format.
	// priv_pem string
	privatePem := string(pem.EncodeToMemory(private3))

	// Public Key generation
	public1 := private1.PublicKey

	publicDer, err := x509.MarshalPKIXPublicKey(&public1)
	if err != nil {
		fmt.Println("ERR::Failed to get der format for PublicKey.")
		return "", "", err
	}

	public2 := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   publicDer,
	}
	publicPem := string(pem.EncodeToMemory(&public2))

	return privatePem, publicPem, nil
}

func writeKeyPair(privatePem string, publicPem string, encType string) (nPrivate int, nPublic int) {
	privateFilename := "./id_" + strings.ToLower(encType) + ".private.pem"
	publicFilename := "./id_" + strings.ToLower(encType) + ".public.pem"

	fPrivate, err := os.Create(privateFilename)
	check(err)
	defer fPrivate.Close()

	fPublic, err := os.Create(publicFilename)
	check(err)
	defer fPublic.Close()

	nPrivate, err = fPrivate.WriteString(privatePem)
	if err != nil {
		fmt.Println("ERR::Failed to write private key file")
		return 0, 0
	}
	nPublic, err = fPublic.WriteString(publicPem)
	if err != nil {
		fmt.Println("ERR::Failed to write public key file")
		return nPrivate, 0
	}

	fPrivate.Sync()
	fPublic.Sync()

	return nPrivate, nPublic
}

func main() {
	passphrase := "tina_belcher_loves_butts"

	privatePem, publicPem, err := generateRSA4096(passphrase)
	if err != nil {
		fmt.Println("ERR::Something has gone awry.", err)
		return
	}

	szPriv, szPub := writeKeyPair(privatePem, publicPem, "rsa")

	// Parsable output <STATUS>::<SZ_PRIV_KEY>::<SZ_PUB_KEY>
	fmt.Printf("OK::%d::%d\n", szPriv, szPub)
}
