// TODO: Need to improve consistency

package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Universal constants
const (
	constPassphrase               = ""
	legalCryptFileExtension       = ".lcsf"
	maxInputFileSize        int64 = 1024 * 1024 * 512 // 512MB; uint64 to support >= 4GB
)

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

	if !pathEndsWithLCSF(strings.ToLower(fname)) {
		return "", errors.New(fname + " does not appear to be a valid LegalCrypt Protected File")
	}

	if opathIsDirectory {
		// TODO: write this as an LC Error
		return "", errors.New("output path exists and is a directory")
	}

	return strings.Replace(outpath, legalCryptFileExtension, "", -1), nil
}

// Adds the legalCryptFileExtension to a filename
func getEncryptedFilename(fname string) string {
	return fname + legalCryptFileExtension
}
