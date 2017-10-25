#!/bin/bash
# Integration tests for cryptengine
# Written by Kevin Cisler


testCount=3
passedCount=0
failedCount=testCount

echo ********* BEGIN CRYPTENGINE INTEGRATION TESTS ********
# Check for test data, generate if needed
echo checking for random test data and keys
if [ ! $(ls -1 *.rando | wc -l)  -eq '8' ] || [ ! $(ls -1 CHECKSUM* | wc -l)  -eq '1' ]
then
  echo missing one or more test files, regenerating...
  rm CHECKSUM*
  rm *.rando
  ./randos <<< 'y'
  echo random data regen complete!
fi

echo beginning tests...

# test keypair generation w/o password
echo [TEST 1] Testing keygen without password...
rm id_rsa* <<< 'y'
./cryptengine-darwin-amd64 -gen -t rsa -eml alfonz.gangwar@gmail.com
if [ -e id_rsa ] && [ -e id_rsa.pub ]
then
	echo keypairs successfully generated!
	let passedCount++
	let failedCount--
else
	echo ERROR: failed to gen keys without password, exiting
	exit 1
fi

# test encrypt/decrypt w/o password
#encrypt
echo [TEST 2] Testing encryption/decryption w/o password...
./cryptengine-darwin-amd64 -e -t rsa benchmark-file-megs1.rando
if [ !-e benchmark-file-megs1.rando.lcsf ]
then
	echo ERROR: failed to encrypt file, exiting
	exit 1
fi
#decrypt
echo encryption successful. Descrypting...
mv benchmark-file-megs1.rando benchmark-file-megs1.rando.old
./cryptengine-darwin-amd64 -eml alfonz.gangwar@gmail.com -t rsa -d benchmark-file-megs1.rando.lcsf
if [ ! -e benchmark-file-megs1.rando ]
	echo ERROR: failed to decrypt file, exiting
	exit 1
fi
#checksum verification
randoschecksum= less CHECKSUM.SHA512-benchmark-file | grep benchmark-file-megs1.rando | cut -f2 -d'=' | tr -d '[:space:]'
testchecksum= openssl dgst -sha512 benchmark-file-megs1.rando
rm benchmark-file-megs1.rando
mv benchmark-file-megs1.rando.old benchmark-file-megs1.rando
if [ ! $randoschecksum -eq $testchecsum ]
	echo ERROR: checksum mismatch, decrypted file is corrupted
	exit 1
fi
echo encrypt/decrypt w/o password succeeded!
let passedCount++
let failedCount--


# test keypair generation w/ password
echo [TEST 3] Testing keygen with password...
rm id_rsa* <<< 'y'
./cryptengine-darwin-amd64 -gen -t rsa -p testPassword -eml alfonz.gangwar@gmail.com
if [ -e id_rsa ] && [ -e id_rsa.pub ]
then
	echo keypairs successfully generated!
	let passedCount++
	let failedCount--
else
	echo ERROR: failed to gen keys with password, exiting
	exit 1
fi

# TODO: run acceptable file size tests

# TODO: run 512 size test

# TODO: test oversized file

# TODO: test multi-file encryption

# TODO: test oversized multi-file encryption
# need to clarify how to determine file size in this scenario
echo $testCount tests total
echo $passedCount tests passed
echo $failedCount tests failed