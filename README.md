# Getting Started

## Building the Project

1) Golang: This version of `cryptengine` uses Go v1.8.3. You can find that [here](https://golang.org/dl/).
2) GB: This project has been scaffolded for GB. You can find that [here](https://getgb.io/docs/install/).
3) In the root of the project, type `gb build all`.

## Updating Vendor Libraries

In the project root, type `gb vendor update <library_name>`.

To blindly update all vendor libraries, use `-all` instead of the name of a library.

# Command: `cryptengine`

## Usage

### Basic Usage
```
cryptengine <options> [file1 file2...]
 -e     Encrypt a file
 -d     Decrypt a file
 -gen   Generate keypair
 
 -t     Type of encryption/keys, defaults to "rsa"
 -dt    Decrypt token; currently does nothing
```

### Generate keypair
```
cryptengine -gen -t rsa
```

### Encrypt a File
```
cryptengine -e -t rsa file1 file2...
```

* Encrypts one or more files
* Supports directories
* Encrypting multiple files first creates a zip file containing the given files, then encrypts the zip file
* Skips non-standard files/directories (e.g. `/dev/null`, `/dev/ttyS2`, `/dev/hda0s13`, etc)
* Reports skipped files, but doesn't fail

### Decrypt a File
```
cryptengine -d -t rsa filename
```

NOTE: The `-t` flag will be deprecated upon autoselection of correct private key based on public key hash, at which point it will be ignored

## Supported Features

### Encryption Methods

* RSA-4096 + AES-256
