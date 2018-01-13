# Cryptengine

## Build
`go build -race -a -o ./bin/cryptengine`

### Tooling Dependencies

#### Go
Go is the language used to build the cryptengine. You'll need to download version 1.9.2 or newer from their site [here](https://golang.org/).

#### Glide
Glide is our vendor package manager. You can find it [here](https://glide.sh/). It is also recommended to install the glide plugin `glide-pin`, which can be installed by running the following command after installing glide: `go get github.com/multiplay/glide-pin`

#### Putting it Together
1. Clone this repo into your `$GOPATH/src` directory. If your `$GOPATH` isn't set, it's best to set it to `~/go`, so that your cloned repo is in `~/go/src/cryptengine`.
2. In the repo's root directory, run `glide install` to install the vendor dependencies.
3. You should be good to go. (See what I did there?) Try a build!

## Usage
Command: `cryptengine`

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
