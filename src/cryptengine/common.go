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

// Universal constants
const (
	constPassphrase         string = ""
	legalCryptFileExtension string = ".lcsf"
	maxInputFileSize        uint64 = 1024 * 1024 * 512 // 512MB; uint64 to support >= 4GB
)

type ErrType struct {
	Code uint8
	Msg  string
}

// We use this to fail hard, which is a good thing
func check(e error, msg string) {
	if e != nil {
		fmt.Fprintln(os.Stderr, "ERR::"+msg)
		panic(e)
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
func getDecryptedFilename(fname string) (string, error) {
	if strings.LastIndex(fname, legalCryptFileExtension) < 0 {
		return "", errors.New(fname + " does not appear to be a valid LegalCrypt Protected File")
	}
	return strings.Replace(fname, legalCryptFileExtension, "", -1), nil
}

// Adds the legalCryptFileExtension to a filename
func getEncryptedFilename(fname string) string {
	return fname + legalCryptFileExtension
}

// If the file exists, delete it
func nixIfExists(filePath string) {
	if _, err := os.Stat(filePath); err == nil {
		check(os.Remove(filePath), "Unable to remove existing file")
	} else {
		check(err, "Unable to remove "+filePath)
	}
}
