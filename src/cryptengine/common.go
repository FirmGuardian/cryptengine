/*
 * TODO: Need to improve consistency
 */

package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"strings"
)

type ErrType struct {
	Code uint16
	Msg  string
}

type OutPathMode uint8

func (t ErrType) code() int {
	return int(t.Code)
}

// Universal constants
const (
	constPassphrase         string      = ""
	legalCryptFileExtension string      = ".lcsf"
	maxInputFileSize        int64       = 1024 * 1024 * 512 // 512MB; uint64 to support >= 4GB
	outPathFile             OutPathMode = 0
	outPathDir              OutPathMode = 1
)

// Error Definitions, all in one spot
var errs = map[string]ErrType{
	"fsCantOpenFile":                    ErrType{200, "Unable to access file"},
	"fsCantCreateFile":                  ErrType{201, "Unable to create file"},
	"fsCantDeleteFile":                  ErrType{202, "Unable to remove existing file"},
	"ioCantReadFromFile":                ErrType{300, "Unable to read from file"},
	"memFileTooBig":                     ErrType{400, "Input file exceeds maximum"},
	"cryptCantReadEncryptedBlock":       ErrType{500, "Error reading encryptedData"},
	"cryptCantDeserializeEncryptedData": ErrType{501, "Unable to deserialize encrypted data"},
	"cryptCantDecodePrivatePEM":         ErrType{501, "Failed to decode PEM block of private key"},
	"cryptCantDecodePublicPEM":          ErrType{502, "Failed to decode PEM block of public key"},
	"cryptCantDecryptPrivateBlock":      ErrType{503, "Unable to decrypt private block"},
	"cryptCantParsePrivateKey":          ErrType{504, "Unable to parse decrypted private key"},
	"cryptCantParsePublicKey":           ErrType{505, "Unable to parse public key"},
	"cryptCantEncryptZip":               ErrType{510, "Unable to encrypt generated zip archive"},
	"cryptCantEncryptOrWrite":           ErrType{511, "Unable to encrypt data, or write encrypted file!"},
	"cryptCantDecryptCipher":            ErrType{520, "Unable to decrypt cipher"},
	"cryptCantDecryptFile":              ErrType{521, "Unable to decrypt file"},
	"cryptAESCantCreateCipher":          ErrType{550, "Unable to create AES cipher"},
	"cryptAESCantCreateGCMBlock":        ErrType{551, "Unable to create GCM Block"},
	"cryptAESCantDecrypt":               ErrType{552, "Unable to decrypt data"},
	"cryptAESCantGenerateSessionKey":    ErrType{553, "Unable to generate sessionKey"},
	"keypairCantReadPublicKey":          ErrType{600, "Error reading public key"},
	"keypairCantReadPrivateKey":         ErrType{601, "Error reading private key"},
	"keypairCantGeneratePrivateKey":     ErrType{602, "Failed to generate private key"},
	"keypairCantValidatePrivateKey":     ErrType{603, "Failed to validate private key"},
	"keypairCantEncryptPrivatePEM":      ErrType{604, "Failed to encrypt private PEM"},
	"keypairCantMarshalPublicKey":       ErrType{605, "Failed to marshal public key block"},
	"panicBadErrType":                   ErrType{1000, "Bad ErrType"},
	"panicWTF":                          ErrType{1001, "Something has gone horribly wrong."},
}

// We use this to fail hard, which is a good thing
func check(err error, errtype ErrType) {
	if errtype.Msg == "" {
		fmt.Println(err)
		panic(errors.New(errs["panicBadErrType"].Msg))
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERR::"+errtype.Msg)
		os.Exit(errtype.code())
	}
}

// Only check if a file exists
func fileExists(filePath string) (bool, os.FileInfo) {
	if stats, err := os.Stat(filePath); err == nil {
		return true, stats
	}

	return false, nil
}

// Used to generate numBytes number of cryptographically random bytes
func generateRandomBytes(numBytes int) ([]byte, error) {
	b := make([]byte, numBytes)
	n, err := rand.Read(b)
	if n != len(b) || err != nil {
		return nil, fmt.Errorf("ERR::Unable to successfully read from the system CSPRNG (%v)", err)
	}

	return b, nil
}

// Removes the legalCryptFileExtension from a filename, if it exists
func getDecryptedFilename(fname string, outpath string) (string, error) {
	opathInfo, _ := os.Stat(outpath)
	opathMode := opathInfo.Mode()
	opathIsDirectory := opathMode.IsDir()

	if strings.LastIndex(fname, legalCryptFileExtension) < 0 {
		return "", errors.New(fname + " does not appear to be a valid LegalCrypt Protected File")
	}
	return strings.Replace(fname, legalCryptFileExtension, "", -1), nil
}

// Adds the legalCryptFileExtension to a filename
func getEncryptedFilename(fname string, outpath string) string {
	return fname + legalCryptFileExtension
}

// If the file exists, delete it
func nixIfExists(filePath string) {
	if exists, _ := fileExists(filePath); exists {
		check(os.Remove(filePath), errs["fsCantDeleteFile"])
	}
}

// Test filename if it has LegalCrypt File Extension
//func isLegalCryptFileName

// Test the outpath
func testOutPath(outpath string) (string, OutPathMode, error) {
	pathExists, opathInfo := fileExists(outpath)

	if pathExists {
		mode := opathInfo.Mode()
		isDir := mode.IsDir()

		if isDir {
			if strings.LastIndex(outpath, string(os.PathSeparator)) < 0 {
				return outpath + "/", outPathDir, nil
			}

			return outpath, outPathDir, nil
		}

		return outpath, outPathFile, nil
	}

	return "", outPathFile, errors.New("Path does not exist.")
}
