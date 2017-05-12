/*
 * Cryptengine is the cryptographic core of the electron client app.
 *
 * Its function is to facilitate keypair generation, as well as
 * file de/encryption. The code may seem simple and mundane, but I
 * followed one important rule throughout the application's source:
 *
 * "Don't be crafty."
 *
 * WRT the mantra of this codebase, I was sure to ask myself:
 * 1) Am I a cryptologist?
 * 2) Am I a wizbang mathemetician?
 * 3) Am I a contributor to a FIPS-certified cryptographic process?
 *
 * "Don't be crafty."
 */
package main

import (
	"fmt"
  "flag"
  "os"
)

func genThoseKeys() {
  fmt.Println(";;Generating keypair")
  generateRSA4096(constPassphrase)
}

func main() {
  decryptPtr := flag.Bool("d", false, "Decrypt the given file")
  encryptPtr := flag.Bool("e", false, "Encrypt the given file")
  keygenPtr  := flag.Bool("gen", false, "Generates a new key pair")

  methodPtr  := flag.String("t", "", "Declares method of encryption/keygen")
  filePtr    := flag.String("f", "", "File to de/encrypt")

  flag.Parse()

  if *methodPtr == "" {
    fmt.Fprintf(os.Stderr, "You must provide a cryptographic method.\n", os.Args[0])
    fmt.Fprintln(os.Stderr, "")
    flag.PrintDefaults()
    os.Exit(1)
  } else {
    if *keygenPtr {
      genThoseKeys()
    } else if *filePtr != "" {
      if *decryptPtr {
        fmt.Println(";;Decrypting file")
        switch *methodPtr {
        default:
          fmt.Println("ERR::Unknown decryption method")
          os.Exit(2)
        case "rsa":
          decryptRSA(*filePtr)
        }
      } else if *encryptPtr {
        fmt.Println(";;Encrypting file")
        switch *methodPtr {
        default:
          fmt.Println("ERR::Unknown encryption method")
          os.Exit(2)
        case "rsa":
          err := encryptRSA(*filePtr)
          check(err, "Could not encrypt data, or write encrypted file!")
        }
      }
    } else {
      flag.PrintDefaults()
      os.Exit(1)
    }
  }

	// Parsable output <STATUS>::<SZ_PRIV_KEY>::<SZ_PUB_KEY>
	fmt.Println("OK")
}
