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

func (t ErrType) code() int {
	return int(t.Code)
}

// Universal constants
const (
	constPassphrase         string = ""
	legalCryptFileExtension string = ".lcsf"
	maxInputFileSize        int64  = 1024 * 1024 * 512 // 512MB; uint64 to support >= 4GB
)

// Error Definitions, all in one spot
var errs = map[string]ErrType{
	"fsCantOpenFile":                    ErrType{200, "Unable to access file"},
	"fsCantCreateFile":                  ErrType{201, "Unable to create file"},
	"fsCantDeleteFile":                  ErrType{202, "Unable to remove existing file"},
	"memFileTooBig":                     ErrType{300, "Input file exceeds maximum"},
	"cryptCantReadEncryptedBlock":       ErrType{400, "Error reading encryptedData"},
	"cryptCantDeserializeEncryptedData": ErrType{401, "Unable to deserialize encrypted data"},
	"cryptCantDecodePrivatePEM":         ErrType{401, "Failed to decode PEM block of private key"},
	"cryptCantDecodePublicPEM":          ErrType{402, "Failed to decode PEM block of public key"},
	"cryptCantDecryptPrivateBlock":      ErrType{403, "Unable to decrypt private block"},
	"cryptCantParsePrivateKey":          ErrType{404, "Unable to parse decrypted private key"},
	"cryptCantParsePublicKey":           ErrType{405, "Unable to parse public key"},
	"cryptCantEncryptZip":               ErrType{410, "Unable to encrypt generated zip archive"},
	"cryptCantEncryptOrWrite":           ErrType{411, "Unable to encrypt data, or write encrypted file!"},
	"cryptCantDecryptCipher":            ErrType{420, "Unable to decrypt cipher"},
	"cryptCantDecryptFile":              ErrType{421, "Unable to decrypt file"},
	"cryptAESCantCreateCipher":          ErrType{450, "Unable to create AES cipher"},
	"cryptAESCantCreateGCMBlock":        ErrType{451, "Unable to create GCM Block"},
	"cryptAESCantDecrypt":               ErrType{452, "Unable to decrypt data"},
	"cryptAESCantGenerateSessionKey":    ErrType{453, "Unable to generate sessionKey"},
	"keypairCantReadPublicKey":          ErrType{500, "Error reading public key"},
	"keypairCantReadPrivateKey":         ErrType{501, "Error reading private key"},
	"keypairCantGeneratePrivateKey":     ErrType{502, "Failed to generate private key"},
	"keypairCantValidatePrivateKey":     ErrType{503, "Failed to validate private key"},
	"keypairCantEncryptPrivatePEM":      ErrType{504, "Failed to encrypt private PEM"},
	"keypairCantMarshalPublicKey":       ErrType{505, "Failed to marshal public key block"},
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
func fileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	}

	return false
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
	opath_info, _ := os.Stat(outpath)
	opath_mode := opath_info.Mode()
	opath_isDirectory := opath_mode.IsDir()

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
	if _, err := os.Stat(filePath); err == nil {
		check(os.Remove(filePath), errs["fsCantDeleteFile"])
	}
}
