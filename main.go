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
	"fmt"
)

func main() {
	passphrase := "t1n@ b3LcHeR_lov3s!bUtts+"

  fmt.Println(";;Generating keypair")
	privatePem, publicPem, err := generateRSA4096(passphrase)
	check(err, "Something has gone awry.")

  fmt.Println(";;Writing keypair")
	writeKeyPair(privatePem, publicPem, "rsa")

  fmt.Println(";;Encrypting file")
  encrypted := encryptRSA("./trump.gif")

  fmt.Println(";;Writing encrypted file")
  err = writeEncryptedFile("./trump.legalcrypt", encrypted)
  check(err, "Could not write encrypted file!")

	// Parsable output <STATUS>::<SZ_PRIV_KEY>::<SZ_PUB_KEY>
	fmt.Println("OK")
}
