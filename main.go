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
  "flag"
  "os"
)

func main() {
  decryptPtr := flag.Bool("decrypt", false, "Decrypt the given file")
  encryptPtr := flag.Bool("encrypt", false, "Encrypt the given file")
  keygenPtr := flag.Bool("keygen", false, "Generates a new key pair")

  methodPtr := flag.String("method", "", "Declares method of encryption/keygen")

  flag.Parse()

  if *methodPtr == "" {
    fmt.Fprintf(os.Stderr, "You must provide a cryptographic method.\n", os.Args[0])
    fmt.Fprintln(os.Stderr, "")
    flag.PrintDefaults()
    os.Exit(0)
  } else {
    if *keygenPtr {
      passphrase := "t1n@ b3LcHeR_lov3s!bUtts+"

      fmt.Println(";;Generating keypair")
      generateRSA4096(passphrase)
    } else if *decryptPtr {
      fmt.Println("DO DECRYPTION")
    } else if *encryptPtr {
      fmt.Println(";;Encrypting file")
      err := encryptRSA("./mysecretdata.txt")
      check(err, "Could not encrypt data, or write encrypted file!")
    }
  }

	// Parsable output <STATUS>::<SZ_PRIV_KEY>::<SZ_PUB_KEY>
	fmt.Println("OK")
}
