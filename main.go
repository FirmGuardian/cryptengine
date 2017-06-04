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

  methodPtr  := flag.String("t", "rsa", "Declares method of encryption/keygen")
  //filePtr    := flag.String("f", "", "File to de/encrypt")
  decryptToken := flag.String("dt", "", "Decrypt token provided by server")

  flag.Parse()

  tail := flag.Args()

  numFiles := len(tail)
  fmt.Printf(";;Tail Size %d\n", numFiles)

  if *keygenPtr {
    genThoseKeys()
  } else if numFiles > 0 {
    if *decryptPtr {
      fmt.Println(";;Decrypting file")
      switch *methodPtr {
      default:
        fmt.Println("ERR::Unknown decryption method")
        os.Exit(2)
      case "rsa":
        fmt.Printf(";;DecryptToken = %s\n", *decryptToken)

        decryptRSA(tail[0])
      }
    } else if *encryptPtr {
      fmt.Println(";;Encrypting file(s)")

      // the following is temporary until multiple files/dirs are supported
      f0info, err := os.Stat(tail[0])
      check(err, "Something is fucky with " + tail[0])
      f0mode := f0info.Mode()
      isRegular := f0mode.IsRegular()
      isDirectory := f0mode.IsDir()

      if numFiles > 1 || (numFiles == 1 && isDirectory) {
        archiveFiles(tail)
      }

      if isRegular == true {
        switch *methodPtr {
        default:
          fmt.Println("ERR::Unknown encryption method")
          os.Exit(2)
        case "rsa":
          err := encryptRSA(tail[0])
          check(err, "Could not encrypt data, or write encrypted file!")
        }
      }
    }
  } else {
    fmt.Println("ERR::Usage: cryptengine [options] file1 (file2 file3...)")
    flag.PrintDefaults()
    os.Exit(1)
  }

	// Parsable output <STATUS>
	fmt.Println("OK")
}
