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
	"log"
	"os"
	"path"
)

const (
	idRSA  = "id_rsa"
	pubRSA = "id_rsa.pub"
)

func decryptRSA(filePath string, secret string, email string, outpath string) {
	fmt.Printf(";;Unused email: %v\n", email)
	fileInfo := pathInfo(filePath)

	// 1) Open the file
	inFile, err := os.OpenFile(filePath, os.O_RDONLY|os.O_SYNC, 0600)
	check(err, errs["fsCantOpenFile"])
	r := bufio.NewReaderSize(inFile, int(fileInfo.Size))
	buf := make([]byte, int(fileInfo.Size))
	_, err = r.Read(buf)
	check(err, errs["ioCantReadFromFile"])
	inFile.Close()

	decryptFile := &messages.EncryptedFile{}
	proto.Unmarshal(buf, decryptFile)

	mType := decryptFile.GetMtype()

	// 2) Read the recipient hashes
	_ = decryptFile.GetRecipientHashes()

	// 3) Read the encrypted aes key
	encryptedKey := decryptFile.GetCipherKey()

	// 4) Read nonce
	nonce := decryptFile.GetCipherNonce()

	// 5) Read encrypted data
	encryptedData := decryptFile.GetEncryptedData()

	decryptFile = nil

	// 6) Read the private key, and pem decode it
	keySlurp, err := ioutil.ReadFile(path.Join(keyDir(), idRSA))
	check(err, errs["keypairCantReadPrivateKey"])
	privateBlock, _ := pem.Decode(keySlurp)
	if privateBlock == nil || privateBlock.Type != "RSA PRIVATE KEY" {
		check(errors.New(errs["cryptCantDecodePrivatePEM"].Msg), errs["cryptCantDecodePrivatePEM"])
	}

	// 7) Decrypt the private key using the password and email
	// TODO: You know the drill.
	uSalt := make([]byte, 0)
	xSalt, _ := ioutil.ReadFile(path.Join(keyDir(), "nacl"))
	der, err := x509.DecryptPEMBlock(privateBlock, deriveKey(secret, append(uSalt, xSalt...)))
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
	outFile, err := os.OpenFile(outFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC, 0600)
	check(err, errs["fsCantCreateFile"])
	n, err := outFile.Write(decryptedData)
	if err != nil {
		// TODO: You know the drill
		log.Fatalln(err)
	}
	if n != len(decryptedData) {
		// TODO: You now the drill
		log.Fatalln("Bytes written not equal to bytes to write!")
	}

	outFile.Close()

	if mType == messages.MType_LCSZ {
		err := unarchiveFiles(outFilePath)
		if err != nil {
			// TODO: You know what to do here.
			log.Fatalln(err)
		}
	}
}

func encryptRSA(filePath string, outPath string, mType messages.MType) error {
	fileInfo := pathInfo(filePath)
	if fileInfo.Size > maxInputFileSize {
		check(errors.New(errs["memFileTooBig"].Msg), errs["memFileTooBig"])
	}

	outFilePath := getEncryptedFilename(filePath, outPath)
	nixIfExists(outFilePath)

	// Slurp and parse public key to encrypt AES Session Key
	keySlurp, err := ioutil.ReadFile(path.Join(keyDir(), pubRSA))
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
		Mtype:           mType,
		RecipientHashes: publicKeyHashes,
		CipherKey:       encryptedSessionKey,
		CipherNonce:     nonce,
		EncryptedData:   encryptedBin,
	}

	// TODO: _ is an err, write a check for it
	encryptedFile, _ := proto.Marshal(encryptedFileProto)

	// Create output file, and Writer; TODO: _ is an err, write a check for it
	outFile, _ := os.OpenFile(outFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_SYNC, 0600)
	outFile.Write(encryptedFile)
	outFile.Close()

	return nil
}

func generateRSA4096(secret []byte) {
	privateFilename := path.Join(keyDir(), idRSA)
	publicFilename := path.Join(keyDir(), pubRSA)

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

	privateFile, err := os.OpenFile(privateFilename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_SYNC, 0600)
	if err != nil {
		// TODO: Create an error here
		log.Fatalln(err)
	}
	_, err = privateFile.Write(derivePrivatePem(cipherKey, secret))
	privateFile.Close()
	if err != nil {
		// TODO: Create an error here
		log.Fatalln(err)
	}
	existsPriv, _ = fileExists(privateFilename)
	if !existsPriv {
		// TODO: Create an error here
		log.Fatalln("ERR::PrivateKey not created.")
	}
	publicFile, err := os.OpenFile(publicFilename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_SYNC, 0600)
	if err != nil {
		// TODO: Create an error here
		log.Fatalln(err)
	}
	_, err = publicFile.Write(derivePublicPem(cipherKey))
	publicFile.Close()
	if err != nil {
		// TODO: Create an error here
		log.Fatalln(err)
	}
	existsPub, _ = fileExists(publicFilename)
	if !existsPub {
		// TODO: Create an error here
		log.Fatalln("ERR::PublicKey not created.")
	}
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
