#!/bin/sh
{
rm -rf testFiles
mkdir testFiles
} > /dev/null
./randos <<< 'y'
{
mv benchmark* testFiles
mv CHECKSUM* testFiles
} > /dev/null