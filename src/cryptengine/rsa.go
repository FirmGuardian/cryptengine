package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/FirmGuardian/legalcrypt-protos/messages"
	"github.com/golang/protobuf/proto"
	"golang.org/x/crypto/sha3"
	"io/ioutil"
	"os"
)

func decryptRSA(filePath string, secret string, email string, outpath string) {
	fileInfo, _ := os.Stat(filePath)
	fileSize := fileInfo.Size()

	// 1) Open the file
	inFile, err := os.Open(filePath)
	check(err, errs["fsCantOpenFile"])
	r := bufio.NewReaderSize(inFile, int(fileSize))
	buf := make([]byte, int(fileSize))
	_, err = r.Read(buf)
	check(err, errs["ioCantReadFromFile"])
	inFile.Close()

	decryptFile := &messages.EncryptedFile{}
	proto.Unmarshal(buf, decryptFile)

	// 2) Read the recipient hashes
	_ = decryptFile.GetRecipientHashes()

	// 3) Read the encrypted aes key
	encryptedKey := decryptFile.GetCipherKey()

	// 4) Read nonce
	nonce := decryptFile.GetCipherNonce()

	// 5) Read encrypted data
	encryptedData := decryptFile.GetEncryptedData()

	// 6) Read the private key, and pem decode it
	keySlurp, err := ioutil.ReadFile("./id_rsa")
	check(err, errs["keypairCantReadPrivateKey"])
	privateBlock, _ := pem.Decode(keySlurp)
	if privateBlock == nil || privateBlock.Type != "RSA PRIVATE KEY" {
		check(errors.New(errs["cryptCantDecodePrivatePEM"].Msg), errs["cryptCantDecodePrivatePEM"])
	}

	// 7) Decrypt the private key using the password and email
	der, err := x509.DecryptPEMBlock(privateBlock, scryptify(secret, email, 64))
	check(err, errs["cryptCantDecryptPrivateBlock"])

	// 8) Unmarshal the private key
	privatePKCS, err := x509.ParsePKCS1PrivateKey(der)
	check(err, errs["cryptCantParsePrivateKey"])

	hash := sha3.New512()
	sessionKey, err := rsa.DecryptOAEP(hash, rand.Reader, privatePKCS, encryptedKey, []byte(""))
	check(err, errs["cryptCantDecryptCipher"])

	// BEGIN AES DECRYPT (sessionKey, nonce, encryptedData)
	decryptedData, err := decryptAES(sessionKey, nonce, encryptedData)
	check(err, errs["cryptCantDecryptFile"])
	// END AES DECRYPT

	// TODO: _ is an err; write a check for it.
	outFilePath, _ := getDecryptedFilename(filePath, outpath)
	fmt.Println("FILE::" + outFilePath)
	nixIfExists(outFilePath)
	outFile, err := os.Create(outFilePath)
	check(err, errs["fsCantCreateFile"])
	defer outFile.Close()
	w := bufio.NewWriter(outFile)

	w.Write(decryptedData)
	w.Flush()
}

func encryptRSA(filePath string, outpath string) error {
	fileInfo, _ := os.Stat(filePath)
	fileSize := fileInfo.Size()
	if fileSize > maxInputFileSize {
		check(errors.New(errs["memFileTooBig"].Msg), errs["memFileTooBig"])
	}

	outFilePath := getEncryptedFilename(filePath, outpath)
	nixIfExists(outFilePath)

	// Create output file, and Writer; TODO: _ is an err, write a check for it
	outFile, _ := os.Create(outFilePath)
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

	// This supports multiple recipient hashes. we're still rolling with one, for now.
	publicKeyHash := sha3.New512()
	publicKeyHash.Write(publicBlock.Bytes)

	publicKeyHashes := make([][]byte, 1)
	publicKeyHashes[0] = publicKeyHash.Sum(nil)

	hash := sha3.New512()
	rng := rand.Reader

	// Slurp file to be encrypted
	fSlurp, err := ioutil.ReadFile(filePath)
	check(err, errs["fsCantOpenFile"])

	// BEGIN AES ENCRYPTION
	// TODO: _ is an err, write a check
	encryptedBin, nonce, sessionKey, _ := encryptAES(fSlurp)

	// CipherKey
	// TODO: _ is an err, write a check for it
	encryptedSessionKey, _ := rsa.EncryptOAEP(hash, rng, publicKey.(*rsa.PublicKey), sessionKey, []byte(""))

	encryptedFileProto := &messages.EncryptedFile{
		Mtype:           messages.MType_LCSF, // It's all LCSF, atm. this should be passed in from main
		RecipientHashes: publicKeyHashes,
		CipherKey:       encryptedSessionKey,
		CipherNonce:     nonce,
		EncryptedData:   encryptedBin,
	}

	// TODO: _ is an err, write a check for it
	encryptedFile, _ := proto.Marshal(encryptedFileProto)

	w.Write(encryptedFile)
	w.Flush()

	return nil
}

func generateRSA4096(secret []byte) {
	privateFilename := "./id_rsa"
	publicFilename := privateFilename + ".pub"

	existsPriv, _ := fileExists(privateFilename)
	existsPub, _ := fileExists(publicFilename)
	if existsPriv && existsPub {
		return
	}

	// private1 *rsa.PrivateKey;
	// err error;
	cipherKey, err := rsa.GenerateKey(rand.Reader, 4096)
	check(err, errs["keypairCantGeneratePrivateKey"])

	check(cipherKey.Validate(), errs["keypairCantValidatePrivateKey"])

	_ = ioutil.WriteFile(privateFilename, derivePrivatePem(cipherKey, secret), 0400)
	_ = ioutil.WriteFile(publicFilename, derivePublicPem(cipherKey), 0644)
}

func derivePrivatePem(cipherKey *rsa.PrivateKey, secret []byte) []byte {
	// pem.Block
	// blk pem.Block
	privateBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(cipherKey),
	}

	// Encrypt the pem
	priv509, err := x509.EncryptPEMBlock(rand.Reader, privateBlock.Type, privateBlock.Bytes, secret, x509.PEMCipherAES256)
	check(err, errs["keypairCantEncryptPrivatePEM"])

	// Resulting private key in PEM format.
	// priv_pem string
	return pem.EncodeToMemory(priv509) // []byte
}

func derivePublicPem(cipherKey *rsa.PrivateKey) []byte {
	// Public Key generation
	publicDer, err := x509.MarshalPKIXPublicKey(&cipherKey.PublicKey)
	check(err, errs["keypairCantMarshalPublicKey"])

	publicBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   publicDer,
	}
	return pem.EncodeToMemory(&publicBlock)
}
