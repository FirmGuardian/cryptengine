package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
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
		// TODO: Move to errors.go
		return nil, fmt.Errorf("ERR::Unable to successfully read from the system CSPRNG (%v)", err)
	}

	return b, nil
}

// Removes the legalCryptFileExtension from a filename, if it exists
func getDecryptedFilename(fname string, outpath string) (string, error) {
	fInfo := pathInfo(fname)
	outInfo := pathInfo(outpath)

	var decryptPath string

	if !pathEndsWithLCSF(fname) {
		// TODO: Make a proper error out of this in errors.go
		return "", errors.New(fname + " does not appear to be a valid LegalCrypt Protected File")
	}

	if outpath == "" {
		// cut out early, if empty outpath
		return stripTrailingLCExt(outpath), nil
	}

	if outInfo.Exists {
		if outInfo.IsDir {
			decryptPath = path.Join(outInfo.Clean, fInfo.File)
		} else if outInfo.IsReg {
			// looks like it exists as a regular file. Nix it
			decryptPath = stripTrailingLCExt(outInfo.Clean)
			nixIfExists(decryptPath)
		} else {
			// This is not a file we should be touching
			return "", fmt.Errorf("ERR::Illegal attempt to overwrite special file (%v)", outInfo.Clean)
		}
	} else {
		// Treat outpath as directory
		if outInfo.Ext == "" {
			// Make the directory, since it doesn't exist
			err := os.MkdirAll(outInfo.Clean, 0600)
			if err != nil {
				return "", fmt.Errorf("ERR::Unable to create directory (%v)", outInfo.Clean)
			}
			decryptPath = path.Join(outInfo.Clean, fInfo.File)
		} else {
			// Looks like a file. Strip any LC extension, and return
			decryptPath = stripTrailingLCExt(outInfo.Clean)
			dpInfo := pathInfo(decryptPath)
			os.MkdirAll(dpInfo.Dir, 0700)
		}
	}

	return decryptPath, nil
}

// Adds the legalCryptFileExtension to a filename
func getEncryptedFilename(f string, o string) string {
	// TODO: Check if exists, and if not, ensure dirs-to-file exist w/ MkDirAll()
	fInfo := pathInfo(f)
	oInfo := pathInfo(o)

	if fInfo.File == "" {
		// TODO: Make this a proper error
		log.Fatalln("ERR::fInfo is not a file")
	}
	var rPath string
	if oInfo.Exists && oInfo.IsDir {
		rPath = path.Join(oInfo.Clean, fInfo.File)
	} else if !oInfo.Exists && oInfo.Ext == "" {
		os.MkdirAll(oInfo.Clean, 0700)
		rPath = path.Join(oInfo.Clean, fInfo.File)
	} else {
		os.MkdirAll(fInfo.Dir, 0700)
		rPath = f
	}

	return appendTrailingLCExt(rPath)
}
