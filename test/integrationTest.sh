#!/bin/bash
# Integration tests for cryptengine
# Written by Kevin Cisler

echo BEGIN CRYPTENGINE INTEGRATION TESTS
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
# test keypair generation
echo generating keypair...
rm id_rsa* <<< 'y'
./cryptengine-darwin-amd64 -gen -t rsa -p testPassword -eml alfonz.gangwar@gmail.com
if [ -e id_rsa ] && [ -e id_rsa.pub ]
then
	echo keypairs successfully generated!
else
	echo ERROR: failed to gen keys, exiting
	exit 1
fi


# TODO: run acceptable file size tests

# TODO: run 512 size test

# TODO: test oversized file

# TODO: test multi-file encryption

# TODO: test oversized multi-file encryption
# need to clarify how to determine file size in this scenario