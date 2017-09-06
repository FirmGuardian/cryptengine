package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
	"io/ioutil"
	"os"
)

func decryptRSA(filePath string, secret string, email string) {
	fileInfo, _ := os.Stat(filePath)
	fileSize := fileInfo.Size()

	inFile, err := os.Open(filePath)
	check(err, errs["fsCantOpenFile"])
	r := bufio.NewReaderSize(inFile, int(fileSize))

	publicKeyHash64 := make([]byte, 88)
	_, err = r.Read(publicKeyHash64)
	check(err, errs["keypairCantReadPublicKey"])

	encryptedKey64 := make([]byte, 684)
	_, err = r.Read(encryptedKey64)
	check(err, errs["keypairCantReadPrivateKey"])
	encryptedKey, err := base64.StdEncoding.DecodeString(string(encryptedKey64))

	nonce64 := make([]byte, 16)
	_, err = r.Read(nonce64)
	nonce, _ := base64.StdEncoding.DecodeString(string(nonce64))

	szEncryptedData := uint64(r.Buffered())

	if szEncryptedData > maxInputFileSize+4096 { // pad max filesize by arbitrary 4k to account for our dick meta
		check(errors.New(errs["memFileTooBig"].Msg), errs["memFileTooBig"])
	}

	encryptedData64 := make([]byte, szEncryptedData) // max encryptable filesize + pad
	_, err = r.Read(encryptedData64[:cap(encryptedData64)])
	inFile.Close()
	check(err, errs["cryptCantReadEncryptedBlock"])

	encryptedData, err := base64.StdEncoding.DecodeString(string(encryptedData64))
	check(err, errs["cryptCantDeserializeEncryptedData"])

	keySlurp, err := ioutil.ReadFile("./id_rsa")
	check(err, errs["keypairCantReadPrivateKey"])
	privateBlock, _ := pem.Decode(keySlurp)
	if privateBlock == nil || privateBlock.Type != "RSA PRIVATE KEY" {
		check(errors.New(errs["cryptCantDecodePrivatePEM"].Msg), errs["cryptCantDecodePrivatePEM"])
	}

	der, err := x509.DecryptPEMBlock(privateBlock, scryptify(secret, email, 64))
	check(err, errs["cryptCantDecryptPrivateBlock"])

	privatePKCS, err := x509.ParsePKCS1PrivateKey(der)
	check(err, errs["cryptCantParsePrivateKey"])

	hash := sha3.New512()
	rng := rand.Reader
	sessionKey, err := rsa.DecryptOAEP(hash, rng, privatePKCS, encryptedKey, []byte(""))
	check(err, errs["cryptCantDecryptCipher"])

	// BEGIN AES DECRYPT (sessionKey, nonce, encryptedData)
	decryptedData, err := decryptAES(sessionKey, nonce, encryptedData)
	check(err, errs["cryptCantDecryptFile"])
	// END AES DECRYPT

	outFilePath, err := getDecryptedFilename(filePath)
	fmt.Println("FILE::" + outFilePath)
	nixIfExists(outFilePath)
	outFile, err := os.Create(outFilePath)
	check(err, errs["fsCantCreateFile"])
	defer outFile.Close()
	w := bufio.NewWriter(outFile)

	w.Write(decryptedData)
	w.Flush()
}

func encryptRSA(filePath string) error {
	outFilePath := getEncryptedFilename(filePath)
	nixIfExists(outFilePath)

	// Create output file, and Writer
	outFile, err := os.Create(outFilePath)
	defer outFile.Close()
	w := bufio.NewWriter(outFile)

	// Slurp and parse public key to encrypt AES Session Key
	keySlurp, err := ioutil.ReadFile("./id_rsa.pub")
	check(err, errs["keypairCantReadPublicKey"])
	publicBlock, _ := pem.Decode(keySlurp)
	if publicBlock == nil || publicBlock.Type != "PUBLIC KEY" {
		check(errors.New(errs["cryptCantDecodePublicPEM"].Msg), errs["cryptCantDecodePublicPEM"])
	}

	publicKey, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	check(err, errs["cryptCantParsePublicKey"])

	publicKeyHash := sha3.New512()
	publicKeyHash.Write(publicBlock.Bytes)

	w.WriteString(base64.StdEncoding.EncodeToString(publicKeyHash.Sum(nil)))

	hash := sha3.New512()
	rng := rand.Reader

	// Slurp file to be encrypted
	fSlurp, err := ioutil.ReadFile(filePath)
	check(err, errs["fsCantOpenFile"])

	// BEGIN AES ENCRYPTION
	encryptedBin, nonce, sessionKey, err := encryptAES(fSlurp)

	encryptedSessionKey, err := rsa.EncryptOAEP(hash, rng, publicKey.(*rsa.PublicKey), sessionKey, []byte(""))
	encryptedSessionKey64 := base64.StdEncoding.EncodeToString(encryptedSessionKey)

	w.Write([]byte(encryptedSessionKey64))

	nonce64 := base64.StdEncoding.EncodeToString(nonce)
	w.Write([]byte(nonce64))

	encrypted64 := base64.StdEncoding.EncodeToString(encryptedBin)
	w.Write([]byte(encrypted64))

	w.Flush()

	return nil
}

func generateRSA4096(secret []byte) {
	privateFilename := "./id_rsa"
	publicFilename := privateFilename + ".pub"
	if fileExists(privateFilename) && fileExists(publicFilename) {
		return
	}

	rng := rand.Reader

	// private1 *rsa.PrivateKey;
	// err error;
	private1, err := rsa.GenerateKey(rng, 4096)
	check(err, errs["keypairCantGeneratePrivateKey"])

	err = private1.Validate()
	check(err, errs["keypairCantValidatePrivateKey"])

	privateDer := x509.MarshalPKCS1PrivateKey(private1)

	// pem.Block
	// blk pem.Block
	private2 := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateDer,
	}

	// Encrypt the pem
	private3, err := x509.EncryptPEMBlock(rng, private2.Type, private2.Bytes, secret, x509.PEMCipherAES256)
	check(err, errs["keypairCantEncryptPrivatePEM"])

	// Resultant private key in PEM format.
	// priv_pem string
	privatePem := pem.EncodeToMemory(private3)

	// Public Key generation
	publicDer, err := x509.MarshalPKIXPublicKey(&private1.PublicKey)
	check(err, errs["keypairCantMarshalPublicKey"])

	public2 := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   publicDer,
	}
	publicPem := pem.EncodeToMemory(&public2)

	_ = ioutil.WriteFile(privateFilename, privatePem, 0400)
	_ = ioutil.WriteFile(publicFilename, publicPem, 0644)
}
