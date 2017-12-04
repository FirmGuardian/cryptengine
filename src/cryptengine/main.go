// Cryptengine is the cryptographic core used by the electron client app.
//
// Its function is to facilitate keypair generation, as well as
// file de/encryption. The code may seem simple and mundane, but I
// followed one important rule throughout the application's source:
//
// "Don't be clever."
//
// WRT the mantra of this codebase, I was sure to ask myself:
// 1) Am I a cryptologist?
// 2) Am I a wizbang mathemetician?
// 3) Am I a contributor to a FIPS-certified cryptographic process?
//
// "Don't be clever."
//
// TODO: Improve os.Exit calls, and associate error messages to error codes

package main

import (
	"errors"
	"flag"
	"fmt"
	//"log"
	"os"
	//"runtime"
	//"runtime/pprof"

	"github.com/FirmGuardian/legalcrypt-protos/messages"
)

func generateKeypairs(passphrase string, email string) {
	fmt.Println(";;Generating keypair")
	generateRSA4096(deriveKey(passphrase, email, 64))
}

// TODO: Abstract the main method's logic. It's getting all spaghetti in there.
func main() {
	// First things first: ensure we have the directory scaffold we desire.
	scaffoldAppDirs()
	// TODO: Remove these at some point
	// Please keep for now
	//cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	//memprofile := flag.String("memprofile", "", "write memory profile to file")

	decryptPtr := flag.Bool("d", false, "Decrypt the given file")
	encryptPtr := flag.Bool("e", false, "Encrypt the given file")
	keygenPtr := flag.Bool("gen", false, "Generates a new key pair")

	outpathPtr := flag.String("o", "", "Output filename or path")
	methodPtr := flag.String("t", "rsa", "Declares method of encryption/keygen")
	decryptToken := flag.String("dt", constPassphrase, "Decrypt token provided by server")
	passPtr := flag.String("p", constPassphrase, "User passphrase")
	emailPtr := flag.String("eml", "", "User email")

	flag.Parse()

	tail := flag.Args()

	// TODO: Remove this, at some point
	// Please keep for now...
	//if *cpuprofile != "" {
	//	f, err := os.Create(*cpuprofile)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	pprof.StartCPUProfile(f)
	//	defer pprof.StopCPUProfile()
	//}

	fmt.Println(";;Email: " + *emailPtr)
	fmt.Println(";;Pass: " + *passPtr)

	numFiles := len(tail)
	fmt.Printf(";;Tail Size %d\n", numFiles)

	// Generate Keypair
	if *keygenPtr {
		generateKeyPair(*passPtr, *emailPtr)
	} else if numFiles > 0 {
		// Decrypt a file
		if *decryptPtr {
			decryptFiles(tail[0], *methodPtr, *passPtr, *emailPtr, *decryptToken, *outpathPtr)
			// Encrypt file(s)
		} else if *encryptPtr {
			encryptFiles(tail, *methodPtr, *outpathPtr)
		}

		// Usage Error
	} else {
		// TODO: Remove this soon. We won't want to display usage on the cli
		fmt.Println("ERR::Usage: cryptengine [options] file1 (file2 file3...)")
		flag.PrintDefaults()
		os.Exit(1000)
	}

	// Parsable output <STATUS>
	fmt.Println("OK")

	// TODO: remove this, at some point.
	// Please keep for now...
	//if *memprofile != "" {
	//	f, err := os.Create(*memprofile)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	runtime.GC()
	//	pprof.WriteHeapProfile(f)
	//	f.Close()
	//	return
	//}
}

func generateKeyPair(passwd string, email string) {
	if passwd != "" && email != "" {
		generateKeypairs(passwd, email)
	} else {
		// TODO: Make an error out of this...
		fmt.Println("ERR::An email is necessary to generate keypairs.")
		os.Exit(1000)
	}
}

func decryptFiles(path string, method string, passwd string, email string, decryptToken string, outPath string) {
	// We need an email to perform decryption
	if outPath == "" {
		outPath = outDirDec() // Get calculated default outdir
	}
	if email != "" {
		fmt.Println(";;Decrypting file")
		switch method {
		default:
			// TODO: Make an error out of this...
			fmt.Println("ERR::Unknown decryption method")
			os.Exit(1000)
		case "rsa":
			fmt.Printf(";;DecryptToken = %s\n", decryptToken)

			decryptRSA(path, passwd, email, outPath)
		}
	} else {
		// TODO: Make an error out of this...
		fmt.Println("ERR::Flag eml is required when decrypting")
		os.Exit(1000)
	}
}

func encryptFiles(files []string, method string, outPath string) {
	if outPath == "" {
		outPath = outDirEnc() // Get calculated default outdir
	}
	fmt.Println(";;Encrypting file(s)")
	numFiles := len(files)
	// File checks on the first file to be encrypted
	f0 := pathInfo(files[0])
	if !f0.Exists {
		check(errors.New(errs["fsCantOpenFile"].Msg), errs["fsCantOpenFile"])
	}

	// Multiple files, or a directory
	if numFiles > 1 || (numFiles == 1 && f0.IsDir) {
		zipPath := archiveFiles(files)
		err := encryptRSA(zipPath, outPath, messages.MType_LCSZ)
		check(err, errs["cryptCantEncryptZip"])

		os.Remove(zipPath)

		fmt.Println("FILE::" + zipPath + legalCryptFileExtension)

		// Just one file, and it's normal (e.g. not /dev/null)
	} else if numFiles == 1 && f0.IsReg {
		switch method {
		default:
			// TODO: Make an error out of this...
			fmt.Println("ERR::Unknown encryption method")
			os.Exit(1000)
		case "rsa":
			err := encryptRSA(files[0], outPath, messages.MType_LCSF)
			check(err, errs["cryptCantEncryptOrWrite"])
			fmt.Println("FILE::" + files[0] + legalCryptFileExtension)
		}

		// Something really bizarre has happened
	} else {
		check(errors.New(errs["panicWTF"].Msg), errs["panicWTF"])
	}
}
