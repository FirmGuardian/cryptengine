#!/bin/bash
# Integration tests for cryptengine
# Written by Kevin Cisler

testFiles='benchmark-file-megs1.rando
           benchmark-file-megs2.rando
           benchmark-file-megs15.rando
           benchmark-file-megs60.rando
           benchmark-file-megs120.rando
           benchmark-file-megs240.rando
           benchmark-file-megs512.rando
           benchmark-file-megs740.rando'
validSizedFiles='benchmark-file-megs1.rando
                 benchmark-file-megs2.rando
                 benchmark-file-megs15.rando
                 benchmark-file-megs60.rando
                 benchmark-file-megs120.rando
                 benchmark-file-megs240.rando'
oversizedFile='benchmark-file-megs740.rando'
checksumFile='CHECKSUM.SHA512-benchmark-file'
needsFileGen=false

email=foo.bar@bar.foo
password=butts

passedCount=0
failedCount=0
runCount=0
totalTests=1

# reports results of test
reportResults () {
	echo ********** END TESTS **********
	echo $totalTests tests available.
	echo $runCount tests ran.
	echo $passedCount tests passed.
	echo $failedCount tests failed.
    echo *******************************
}

# encrypts a given file, checks that .lcsf file exists
encrypt() {
echo encrypting...
./cryptengine-darwin-amd64 -e -t rsa ./testFiles/$1
if [ ! -e $1.lcsf ] || [ ! $? -eq 0]
then
	return 1
fi
return 0
}

# decrypts a given file, checks that the unencrypted version exists
decrypt() {
echo decrypting...
./cryptengine-darwin-amd64 -eml $email -t rsa -d $1.lcsf
if [ ! -e $1 ] || [ ! $? -eq 0 ]
	return 1
fi
return 0
}

#function used to test valid single file encryption
testED () {
let runCount++
let failedCount++
echo [Test $runCount] Testing encryption/decryption of $1...
#encrypt
encrypt $1
if [ ! $? -eq 0 ]
then
    echo ERROR: failed to encrypt $1
    return 1
fi
#decrypt
decrypt $1.lcsf
if [ ! $? -eq 0 ]
then
    echo ERROR: failed to decrypt $1.lcsf
    return 2
fi

#checksum verification
#randoschecksum= less CHECKSUM.SHA512-benchmark-file | grep benchmark-file-megs1.rando | cut -f2 -d'=' | tr -d '[:space:]'
#testchecksum= openssl dgst -sha512 benchmark-file-megs1.rando
#rm benchmark-file-megs1.rando
#mv benchmark-file-megs1.rando.old benchmark-file-megs1.rando
#if [ ! $randoschecksum -eq $testchecsum ]
#	echo ERROR: checksum mismatch, decrypted file is corrupted
#	exit 1
#fi
let failedCount--
let runCount++
return 0
}

echo ********* BEGIN CRYPTENGINE INTEGRATION TESTS ********
echo setting up test envrionment...
# remove old test files
echo cleaning up directory...
{
rm -f benchmark*
rm -f CHECKSUM*
rm -f id_rsa*
} > /dev/null
echo cleanup done.

# Check for test files
echo verifying test files...
if [ ! -e ./testFiles/$checksumFile ]
then
	needsFileGen=true
fi

for file in $testFiles
do
	if [ ! -e ./testFiles/$file ]
	then
		needsFileGen=true
	fi
done

if $needsFileGen
then
	echo one or more files missing, regenerating...
	{
    rm -rf testFiles
    mkdir testFiles
	} > /dev/null
	./randos <<< 'y'
	echo moving files...
	{
	mv benchmark* testFiles
	mv CHECKSUM* testFiles
	} > /dev/null
	echo test files regen complete.
fi
echo test environement setup complete.


echo beginning tests...

# test keygen
echo [TEST 1] Testing keygen...
let runCount++
let failedCount++
./cryptengine-darwin-amd64 -gen -t rsa -p $password -eml $email
if [ -e id_rsa ] && [ -e id_rsa.pub ] && [ $? -eq 0 ]
then
	echo keypairs successfully generated!
	let passedCount++
	let failedCount--
else
	echo ERROR: failed to gen keys, exiting
	reportResults
	exit 1
fi

# test valid single file encrypt/decrypt
echo testing valid single file encryption...
for file in $testFiles
do
	testED $file
done

'
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
'