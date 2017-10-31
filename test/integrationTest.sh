#!/bin/bash
# Integration tests for cryptengine
# Written by Kevin Cisler
dickbutt='benchmark-file-megs1.rando'
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
                 benchmark-file-megs240.rando
                 benchmark-file-megs512.rando'
oversizedFile='benchmark-file-megs740.rando'
checksumFile='CHECKSUM.SHA512-benchmark-file'
needsFileGen=false

email=foo.bar@bar.foo
password=butts

passedCount=0
failedCount=0
runCount=0
totalTests=8

# reports results of test
reportResults () {
	echo "********** END TESTS **********"
    echo "$totalTests tests available."
    echo "$runCount ran"
    echo "$passedCount passed."
    echo "$failedCount failed."
    echo "*******************************"
}

# encrypts a given file, checks that .lcsf file exists
encrypt() {
    echo "encrypting..."
    cp ./testFiles/$1 ./
    ./cryptengine-darwin-amd64 -e -t rsa $1
    exitcode=$?
    rm $1
    if [ ! -e $2 ] || [ ! $exitcode -eq 0 ]
    then
    	return 1
    fi
    return 0
}

# decrypts a given file, checks that the unencrypted version exists
decrypt() {
    echo "decrypting..."
    ./cryptengine-darwin-amd64 -d -p $password -eml $email -t rsa -d $2
    exitcode=$?
    if [ ! -e $1 ] || [ ! $exitcode -eq 0 ]
    then
    	return 1
    fi
    return 0
}
# verifies checksum of given .rando file
checksumCheck() {
    echo "comparing checksums..."
    randoschecksum=$(less ./testFiles/$checksumFile | grep $1 | cut -f2 -d'=' | tr -d '[:space:]')
    testchecksum=$(openssl dgst -sha512 $1 | cut -f2 -d'=' | tr -d '[:space:]')
    #echo "randoschecksum equals $randoschecksum"
    #echo "testchecksum equals   $testchecksum"
    if [ ! $randoschecksum == $testchecksum ]
    then
        return 1
    fi
    echo "checksums match"
    return 0
}

#function used to test valid single file encryption
testED () {
    decryptfile=$1.lcsf
    let runCount++
    let failedCount++
    echo "[Test $runCount] Testing encryption/decryption of $1..."
    #encrypt
    encrypt $1 $decryptfile
    #echo "return from encrypt function is: $?"
    if [ ! $? -eq 0 ]
    then
        echo "ERROR: failed to encrypt $1"
        return 1
    fi
    #decrypt
    decrypt $1 $decryptfile
    if [ ! $? -eq 0 ]
    then
        echo "ERROR: failed to decrypt $1.lcsf"
        return 2
    fi
    #checksum verification
    checksumCheck $1
    if [ ! $? -eq 0 ]
    then
        echo "ERROR: checksum mismatch for decrypted version of $1"
        return 3
    fi
    let failedCount--
    let passedCount++
    return 0
}

echo "********* BEGIN CRYPTENGINE INTEGRATION TESTS ********"
echo "setting up test envrionment..."
# remove old test files
echo "cleaning up directory..."
{
rm -f benchmark*
rm -f CHECKSUM*
rm -f id_rsa*
} > /dev/null

# Check for test files
echo "checking test files..."
{
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
} > /dev/null

if $needsFileGen
then
	echo "one or more files missing, regenerating..."
	{
    rm -rf testFiles
    mkdir testFiles
	} > /dev/null
	./randos <<< 'y'
	{
	mv benchmark* testFiles
	mv CHECKSUM* testFiles
	} > /dev/null
	echo "test files regen complete."
fi
echo "test environment setup complete."


echo "beginning tests..."

# test keygen
echo "[TEST 1] Testing keygen..."
let runCount++
let failedCount++
./cryptengine-darwin-amd64 -gen -t rsa -p $password -eml $email
if [ -e id_rsa ] && [ -e id_rsa.pub ] && [ $? -eq 0 ]
then
	echo "keypairs successfully generated!"
	let passedCount++
	let failedCount--
else
	echo "ERROR: failed to gen keys, exiting"
	reportResults
	exit 1
fi

echo "testing valid single file encryption..."
for file in $validSizedFiles
do
	testED $file
done

#cleanup and report
{
rm -f benchmark*
rm -f CHECKSUM*
} > /dev/null
reportResults