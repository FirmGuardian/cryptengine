package main

import (
	"crypto/aes"
	"crypto/cipher"
)

const lenAESNonce uint8 = 12

func decryptAES(key []byte, nonce []byte, encryptedData []byte) ([]byte, error) {
	aesCipher, err := aes.NewCipher(key)
	check(err, errs["cryptAESCantCreateCipher"])

	aesgcm, err := cipher.NewGCM(aesCipher)
	check(err, errs["cryptAESCantCreateGCMBlock"])

	decryptedData, err := aesgcm.Open(nil, nonce, encryptedData, nil)
	check(err, errs["cryptAESCantDecrypt"])

	return decryptedData, err
}

func encryptAES(unencryptedData []byte) ([]byte, []byte, []byte, error) {
	// Generate AES Session Key; to be RSA encrypted, and used to
	// encrypt input file
	key, err := generateRandomBytes(32) // 32bytes * 8bits = 256bits
	check(err, errs["cryptAESCantGenerateSessionKey"])

	aesCipher, err := aes.NewCipher(key)
	check(err, errs["cryptAESCantCreateCipher"])

	// Generate nonce
	nonce, _ := generateRandomBytes(int(lenAESNonce))

	aesgcm, err := cipher.NewGCM(aesCipher)
	check(err, errs["cryptAESCantCreateGCMBlock"])

	encryptedBin := aesgcm.Seal(nil, nonce, unencryptedData, nil)

	return encryptedBin, nonce, key, err
}
