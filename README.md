# Getting Started

## Tooling

1) Golang: This version of `cryptengine` uses Go v1.9.x. You can find that [here](https://golang.org/dl/).
2) GB: This project has been scaffolded for GB. You can find that [here](https://getgb.io/docs/install/).
3) Godo: This project uses a Golang task runner that's in the same vein as Gulp or Rake. From the command line, run `go get -u gopkg.in/godo.v2/cmd/godo`
4) [Optional] It is highly recommended to utilize Gometalinter, a wrapper around a set of static analyzers and linting tools. You can find it [here](https://github.com/alecthomas/gometalinter).

## Building the Project

1) In the root of the project, type `godo` to build the project for macOS and Windows.
2) High-five yourself for winning.

## Updating Vendor Libraries

Vendor libraries should only be updated as part of the development of the application (e.g. golang version upgrades, bugfixes to features we use, performance enhancements, etc.) . You should not have to update vendor libraries as part of your dev environment setup, and you should never do a blanket update.

To update a library:
In the project root, type `gb vendor update <library_name>`.

To blindly update all vendor libraries (again, don't do this), use `-all` instead of the name of a library.

# Command: `cryptengine`

## Usage

### Basic Usage
Basic usage will be removed from the 
```bash
cryptengine <options> [file1 file2...]
 -e      Encrypt a file
 -d      Decrypt a file
 -gen    Generate keypair
 
 -t      Type of encryption/keys, defaults to "rsa"
 -dt     Decrypt token; currently does nothing
 
 -p      Password for keygen
 -eml    Optional: email address, used in keygen
```

### Generate keypair
```bash
cryptengine -gen -t rsa -p <password> -eml <email>
```

**NOTE:** For reasons yet unknown to me, please keep the arguments in the stated order. *I know, I know.* I have to look into this.

### Encrypt a File
```bash
cryptengine -e -t rsa file1 file2...
```

* Encrypts one or more files
* Supports directories
* Encrypting multiple files first creates a zip file containing the given files, then encrypts the zip file
* Skips non-standard files/directories (e.g. `/dev/null`, `/dev/ttyS2`, `/dev/hda0s13`, etc)
* Reports skipped files, but doesn't fail

### Decrypt a File
```bash
cryptengine -p <password> -eml <email> -t rsa -d filename
```
**NOTE:** For reasons yet unknown to me, please keep the arguments in the stated order. *I know, I know.* I have to look into this.

**NOTE:** The `-t` flag will be deprecated upon autoselection of correct private key based on public key hash, at which point it will be ignored

### Test Scrypt Pwd-based KDF
```bash
cryptengine -scrypt passphrase
```
This will take a few seconds, then output the base64-encoded derived key to the console. This is more PoC, than anything else, but will be written into the libraries used by both front- and back ends.

## Supported Features

### Encryption Methods

* RSA-4096 + AES-256

### Password-Based Key Derivation

* Scrypt
