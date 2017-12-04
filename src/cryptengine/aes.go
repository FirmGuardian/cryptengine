package main

// TODO: We can add senderIDs to the additional data of the gcm methods. This data would be plaintext, but authenticated

import (
	"crypto/aes"
	"crypto/cipher"
)

const (
	lenAESNonce = 12
	lenAESKey   = 32 // 32bytes * 8bits = 256bits
)

func decryptAES(key []byte, nonce []byte, encryptedData []byte) ([]byte, error) {
	aesCipher, err := aes.NewCipher(key)
	check(err, errs["cryptAESCantCreateCipher"])

	gcm, err := cipher.NewGCM(aesCipher)
	check(err, errs["cryptAESCantCreateGCMBlock"])

	decryptedData, err := gcm.Open(nil, nonce, encryptedData, nil)
	check(err, errs["cryptAESCantDecrypt"])

	return decryptedData, err
}

func encryptAES(unencryptedData []byte) ([]byte, []byte, []byte, error) {
	// Generate AES Session Key; to be RSA encrypted, and used to
	// encrypt input file
	key, err := generateRandomBytes(lenAESKey)
	check(err, errs["cryptAESCantGenerateSessionKey"])

	aesCipher, err := aes.NewCipher(key)
	check(err, errs["cryptAESCantCreateCipher"])

	// Generate nonce
	nonce, _ := generateRandomBytes(lenAESNonce)

	gcm, err := cipher.NewGCM(aesCipher)
	check(err, errs["cryptAESCantCreateGCMBlock"])

	encryptedBin := gcm.Seal(nil, nonce, unencryptedData, nil)

	return encryptedBin, nonce, key, err
}
