# Getting Started

## Building the Project

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
