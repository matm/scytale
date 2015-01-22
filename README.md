[![GoDoc](https://godoc.org/github.com/matm/scytale?status.svg)](https://godoc.org/github.com/matm/scytale)

**Scytale** is a simple wrapper library to make use of encryption with Go
fast and easy.

Most of the credits goes to the wonderful [go.crypto](https://code.google.com/p/go.crypto/) library.

![logo](http://go-tsunami.com/assets/images/scytaleLogo.png)

## Installation

Use `go get` to install the package:

    $ go get github.com/matm/scytale

## Tools

The `bin/aesenc` and `bin/aeszip` CLI tools allow file encryption using
[AES](http://en.wikipedia.org/wiki/Advanced_Encryption_Standard)-256 operating in [CBC mode](http://en.wikipedia.org/wiki/Block_cipher_mode_of_operation),
[password-based encryption](http://en.wikipedia.org/wiki/Password-based_cryptography) (PBE)
and a [PBKDF2](http://en.wikipedia.org/wiki/PBKDF2) password-based key derivation function.

The former can be used to encrypt/decrypt a single file:

    $ go install github.com/matm/scytale/bin/aesenc
    $ aesenc -o out.enc myfile.pdf
    $ aesenc -o myfile.pdf -d out.enc

The latter encrypts a bunch of files into a standard [ZIP](http://en.wikipedia.org/wiki/Zip_%28file_format%29)
file:

    $ go install github.com/matm/scytale/bin/aeszip
    $ aeszip -o secure.zip *.pdf

Both `aeszip -l` and `unzip -l` can be used to list the content of the archive. All files within
the archive can be decrypted and extracted with

    $ aeszip -x -o . secure.zip

The second file in the archive (position 1) can be decrypted and extracted with

    $ aeszip -x -o . -n 1 secure.zip

See `aeszip -h` output for usage.

## API Docs

The API doc is available on [GoDoc](https://godoc.org/github.com/matm/scytale).

## License

This is free software, see LICENSE.
