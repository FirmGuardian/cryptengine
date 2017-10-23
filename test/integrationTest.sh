#!/bin/bash
# Integration tests for cryptengine
# Written by Kevin Cisler

echo BEGIN CRYPTENGINE INTEGRATION TESTS
# TODO: check for test files, generate them if not there
echo checking for random test data
if [ ! $(ls -1 *.rando | wc -l)  -eq '8' ] || [ ! $(ls -1 CHECKSUM* | wc -l)  -eq '1' ]
then
  echo missing one or more test files, regenerating...
  rm CHECKSUM*
  rm *.rando
  ./randos <<< 'y'
  echo random data regen complete!
fi
echo beginning tests...
# TODO: run acceptable file size tests

# TODO: run 512 size test

# TODO: test oversized file

# TODO: test multi-file encryption

# TODO: test oversized multi-file encryption
# need to clarify how to determine file size in this scenario