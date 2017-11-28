package main

import (
	"errors"
	"fmt"
	"os"
)

// ErrType is an exported type used throughout the application to map simple, semantic string keys
// to structs that provide both an error code for exit values, as well as a more human readable
// error message.
type ErrType struct {
	Code uint16
	Msg  string
}

func (t ErrType) code() int {
	return int(t.Code)
}

// Error Definitions, all in one spot
var errs = map[string]ErrType{
	"fsCantOpenFile":                    {200, "Unable to access file"},
	"fsCantCreateFile":                  {201, "Unable to create file"},
	"fsCantDeleteFile":                  {202, "Unable to remove existing file"},
	"ioCantReadFromFile":                {300, "Unable to read from file"},
	"memFileTooBig":                     {400, "Input file exceeds maximum"},
	"cryptCantReadEncryptedBlock":       {500, "Error reading encryptedData"},
	"cryptCantDeserializeEncryptedData": {501, "Unable to deserialize encrypted data"},
	"cryptCantDecodePrivatePEM":         {501, "Failed to decode PEM block of private key"},
	"cryptCantDecodePublicPEM":          {502, "Failed to decode PEM block of public key"},
	"cryptCantDecryptPrivateBlock":      {503, "Unable to decrypt private block"},
	"cryptCantParsePrivateKey":          {504, "Unable to parse decrypted private key"},
	"cryptCantParsePublicKey":           {505, "Unable to parse public key"},
	"cryptCantEncryptZip":               {510, "Unable to encrypt generated zip archive"},
	"cryptCantEncryptOrWrite":           {511, "Unable to encrypt data, or write encrypted file!"},
	"cryptCantDecryptCipher":            {520, "Unable to decrypt cipher"},
	"cryptCantDecryptFile":              {521, "Unable to decrypt file"},
	"cryptAESCantCreateCipher":          {550, "Unable to create AES cipher"},
	"cryptAESCantCreateGCMBlock":        {551, "Unable to create GCM Block"},
	"cryptAESCantDecrypt":               {552, "Unable to decrypt data"},
	"cryptAESCantGenerateSessionKey":    {553, "Unable to generate sessionKey"},
	"keypairCantReadPublicKey":          {600, "Error reading public key"},
	"keypairCantReadPrivateKey":         {601, "Error reading private key"},
	"keypairCantGeneratePrivateKey":     {602, "Failed to generate private key"},
	"keypairCantValidatePrivateKey":     {603, "Failed to validate private key"},
	"keypairCantEncryptPrivatePEM":      {604, "Failed to encrypt private PEM"},
	"keypairCantMarshalPublicKey":       {605, "Failed to marshal public key block"},
	"panicBadErrType":                   {1000, "Bad ErrType"},
	"panicWTF":                          {1001, "Something has gone horribly wrong."},
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
