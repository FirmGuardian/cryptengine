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
 *
 * TODO: Improve os.Exit calls, and associate error messages to error codes
 */
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
  "runtime"
)

func generateKeypairs(passphrase string, email string) {
	fmt.Println(";;Generating keypair")
	generateRSA4096(scryptify(passphrase, email, 64))
}

func main() {
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile := flag.String("memprofile", "", "write memory profile to file")

	decryptPtr := flag.Bool("d", false, "Decrypt the given file")
	encryptPtr := flag.Bool("e", false, "Encrypt the given file")
	keygenPtr := flag.Bool("gen", false, "Generates a new key pair")

	methodPtr := flag.String("t", "rsa", "Declares method of encryption/keygen")
	decryptToken := flag.String("dt", constPassphrase, "Decrypt token provided by server")
	passPtr := flag.String("p", constPassphrase, "User passphrase")
	emailPtr := flag.String("eml", "", "User email")

	flag.Parse()

	tail := flag.Args()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		//defer pprof.StopCPUProfile()
		defer pprof.StopCPUProfile()
	}

	fmt.Println(";;Email: " + *emailPtr)
	fmt.Println(";;Pass: " + *passPtr)

	numFiles := len(tail)
	fmt.Printf(";;Tail Size %d\n", numFiles)

	// Generate Keypair
	if *keygenPtr {
		if *emailPtr != "" {
			generateKeypairs(*passPtr, *emailPtr)
		} else {
			fmt.Println("ERR::An email is necessary to generate keypairs.")
			os.Exit(1000)
		}
	} else if numFiles > 0 {

		// Decrypt a file
		if *decryptPtr {

			// We need an email to perform decryption
			if *emailPtr != "" {
				fmt.Println(";;Decrypting file")
				switch *methodPtr {
				default:
					fmt.Println("ERR::Unknown decryption method")
					os.Exit(1000)
				case "rsa":
					fmt.Printf(";;DecryptToken = %s\n", *decryptToken)

					decryptRSA(tail[0], *passPtr, *emailPtr)
				}
			} else {
				fmt.Println("ERR::Flag eml is required when decrypting")
				os.Exit(1000)
			}

			// Encrypt shit
		} else if *encryptPtr {
			fmt.Println(";;Encrypting file(s)")

			// File checks on the first file to be encrypted
			f0info, err := os.Stat(tail[0])
			check(err, errs["fsCantOpenFile"])
			f0mode := f0info.Mode()
			f0isRegular := f0mode.IsRegular()
			f0isDirectory := f0mode.IsDir()

			// Multiple files, or a directory
			if numFiles > 1 || (numFiles == 1 && f0isDirectory) {
				zipPath := archiveFiles(tail)
				err := encryptRSA(zipPath)
				check(err, errs["cryptCantEncryptZip"])

				os.Remove(zipPath)

				fmt.Println("FILE::" + zipPath + legalCryptFileExtension)

				// Just one file, and it's normal (e.g. not /dev/null)
			} else if numFiles == 1 && f0isRegular == true {
				switch *methodPtr {
				default:
					fmt.Println("ERR::Unknown encryption method")
					os.Exit(1000)
				case "rsa":
					err := encryptRSA(tail[0])
					check(err, errs["cryptCantEncryptOrWrite"])
					fmt.Println("FILE::" + tail[0] + legalCryptFileExtension)
				}

				// Something really bizarre has happened
			} else {
				check(errors.New(errs["panicWTF"].Msg), errs["panicWTF"])
			}
		}

		// Usage Error
	} else {
		fmt.Println("ERR::Usage: cryptengine [options] file1 (file2 file3...)")

		// TODO: Remove this soon. We won't want to display usage on the cli
		flag.PrintDefaults()
		os.Exit(1000)
	}

	// Parsable output <STATUS>
	fmt.Println("OK")

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		runtime.GC()
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}
}
