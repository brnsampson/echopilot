package signing

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/fs"
	"os"
)

const (
	KeyFilePerms    fs.FileMode = 0600
	PubKeyFilePerms fs.FileMode = 0644
)

func encode(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) ([]byte, []byte, error) {
	x509Encoded, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, nil, err
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	x509EncodedPub, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, nil, err
	}
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return pemEncoded, pemEncodedPub, nil
}

func encodePub(publicKey *ecdsa.PublicKey) ([]byte, error) {
	x509EncodedPub, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return pemEncodedPub, nil
}

func decode(pemEncoded []byte, pemEncodedPub []byte) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	// BUG WATCH
	// This only processes the first PEM formatted block, so if this were a cert chain or something
	// we could not throw away the second return value (the remaining input)
	block, _ := pem.Decode(pemEncoded)
	x509Encoded := block.Bytes
	privateKey, err := x509.ParseECPrivateKey(x509Encoded)
	if err != nil {
		return nil, nil, err
	}

	// BUG WATCH - See above
	blockPub, _ := pem.Decode(pemEncodedPub)

	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return privateKey, publicKey, nil
}

func decodePub(pemEncodedPub []byte) (*ecdsa.PublicKey, error) {
	// BUG WATCH
	// This only processes the first PEM formatted block, so if this were a cert chain or something
	// we could not throw away the second return value (the remaining input)
	blockPub, _ := pem.Decode(pemEncodedPub)

	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return publicKey, nil
}

func GenerateECDSA() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	pub := &key.PublicKey
	return key, pub, nil
}

func KeyFilePermsValid(keyFile, pubFile string) (keyFileValid, pubFileValid bool, err error) {
	// SECURITY - key files should be file permissions 600.
	// Ideally, pub key files would be 644 but as long as group and other
	// users can't write it we should be fine. For simplicity, we're just
	// going to require 644

	keyFileValid = false
	pubFileValid = false
	keyStat, keyErr := os.Stat(keyFile)
	pubStat, pubErr := os.Stat(pubFile)

	if pubErr == nil {
		pubMode := pubStat.Mode()
		if pubMode^PubKeyFilePerms == 0 {
			// pubMode exactly matches binary 110 100 100 -> 644 file permission mask
			// This allows other users/groups to read the public key if they need to
			// verify signatures and such.
			pubFileValid = true
		}
	} else {
		err = pubErr
	}

	if keyErr == nil {
		keyMode := keyStat.Mode()
		if keyMode^KeyFilePerms == 0 {
			// keyMode exactly matches binary 110 000 000 -> 600 file permission mask
			keyFileValid = true
		}
	} else {
		err = keyErr
	}

	return
}

func PubKeyFilePermsValid(pubFile string) (pubFileValid bool, err error) {
	// SECURITY - Ideally, pub key files would be 644 but as long as group and other
	// users can't write it we should be fine. For simplicity, we're just
	// going to require 644

	pubFileValid = false
	pubStat, err := os.Stat(pubFile)

	if err == nil {
		pubMode := pubStat.Mode()
		if pubMode^PubKeyFilePerms == 0 {
			// pubMode exactly matches binary 110 100 100 -> 644 file permission mask
			// This allows other users/groups to read the public key if they need to
			// verify signatures and such.
			pubFileValid = true
		}
	} else {
		return false, err
	}
	return
}

func SetKeyFilePerms(keyFile, pubFile string) error {
	// SECURITY - key files should be file permissions 600.
	// Ideally, pub key files would be 644 but as long as group and other
	// users can't write it we should be fine. For simplicity, we're just
	// going to require 644

	err := os.Chmod(keyFile, KeyFilePerms)
	if err != nil {
		return err
	}

	err = os.Chmod(pubFile, PubKeyFilePerms)
	if err != nil {
		return err
	}

	return nil
}

func ReadECDSAFiles(keyFile, pubFile string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	// BUG WATCH - this function is untested despite being part of the public interface. This is
	// because it only interacts with tested interfaces, external resources, and error propogation.
	keyValid, pubValid, err := KeyFilePermsValid(keyFile, pubFile)
	if err != nil {
		return nil, nil, fmt.Errorf("Error while checking key/pubkey permissions: %w", err)
	}

	if keyValid != true || pubValid != true {
		// If we were writing the key/pubkey then we would set the perms here, but for reading just bail.
		return nil, nil, fmt.Errorf("File permisisons for signing keys unacceptable. %s should be 0600 and %s should be 0644.", keyFile, pubFile)
	}

	keyReader, err := os.Open(keyFile)
	if err != nil {
		return nil, nil, err
	}
	defer keyReader.Close()

	pubReader, err := os.Open(pubFile)
	if err != nil {
		return nil, nil, err
	}
	defer pubReader.Close()

	return ReadECDSA(keyReader, pubReader)
}

func ReadPubECDSAFile(pubFile string) (*ecdsa.PublicKey, error) {
	// BUG WATCH - this function is untested despite being part of the public interface. This is
	// because it only interacts with tested interfaces, external resources, and error propogation.
	pubValid, err := PubKeyFilePermsValid(pubFile)
	if err != nil {
		return nil, fmt.Errorf("Error while checking key/pubkey permissions: %w", err)
	}

	if pubValid != true {
		// If we were writing the key/pubkey then we would set the perms here, but for reading just bail.
		return nil, fmt.Errorf("File permisisons for signing keys unacceptable. %s should be 0644.", pubFile)
	}

	pubReader, err := os.Open(pubFile)
	if err != nil {
		return nil, err
	}
	defer pubReader.Close()

	return ReadPubECDSA(pubReader)
}

func ReadECDSA(keyReader, pubReader io.Reader) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	// BUG WATCH a key should be small so it should be fine to just use ReadAll,
	// but what if someone passed us a massive file instead?
	keyEncoded, err := io.ReadAll(keyReader)
	if err != nil {
		return nil, nil, err
	}

	pubEncoded, err := io.ReadAll(pubReader)
	if err != nil {
		return nil, nil, err
	}

	return decode(keyEncoded, pubEncoded)
}

func ReadPubECDSA(pubReader io.Reader) (*ecdsa.PublicKey, error) {
	// BUG WATCH a key should be small so it should be fine to just use ReadAll,
	// but what if someone passed us a massive file instead?
	pubEncoded, err := io.ReadAll(pubReader)
	if err != nil {
		return nil, err
	}

	return decodePub(pubEncoded)
}

func WriteECDSAFiles(key *ecdsa.PrivateKey, keyFile string, pub *ecdsa.PublicKey, pubFile string) error {
	// BUG WATCH - this function is untested despite being part of the public interface.
	// This is because this bunction just interacts with external resources and calling other
	// tested functions.

	// There is an implicit contract to the user that this function will create the keyfiles
	// if they do not exist.

	// SECURITY - always write keys with restrictive permissions of 0600
	keyWriter, err := os.Open(keyFile)
	if err != nil {
		return err
	}
	defer keyWriter.Close()

	// SECURITY - public keys should be readable by anyone but only writable by the service user: 0644
	pubWriter, err := os.Open(pubFile)
	if err != nil {
		return err
	}
	defer pubWriter.Close()

	keyValid, pubValid, err := KeyFilePermsValid(keyFile, pubFile)
	if err != nil {
		return fmt.Errorf("Error while checking key/pubkey permissions: %w", err)
	}

	if keyValid != true || pubValid != true {
		err := SetKeyFilePerms(keyFile, pubFile)
		if err != nil {
			return fmt.Errorf("Error while attempting to set key/pubkey permissions: %w", err)
		}
	}

	return WriteECDSA(key, keyWriter, pub, pubWriter)
}

func WriteECDSA(key *ecdsa.PrivateKey, keyWriter io.Writer, pub *ecdsa.PublicKey, pubWriter io.Writer) error {
	keyEncoded, pubEncoded, err := encode(key, pub)
	if err != nil {
		return err
	}

	keyLen := len(keyEncoded)
	keyWritten := 0
	for keyWritten < keyLen {
		written, err := keyWriter.Write(keyEncoded[keyWritten:])
		if err != nil && written != (keyLen-keyWritten) {
			// if we did not write the entire remaining buffer we expect an error
			if err != io.ErrShortWrite {
				return fmt.Errorf("Error writing private key: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("Error writing private key: %w", err)
		}
		keyWritten += written
	}

	pubLen := len(pubEncoded)
	pubWritten := 0
	for pubWritten < pubLen {
		written, err := pubWriter.Write(pubEncoded[pubWritten:])
		if err != nil && written != (pubLen-pubWritten) {
			// if we did not write the entire remaining buffer we expect an error
			if err != io.ErrShortWrite {
				return fmt.Errorf("Error writing public key: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("Error writing public key: %w", err)
		}
		pubWritten += written
	}

	return nil
}

func EncodedECDSAReaders(key *ecdsa.PrivateKey, pub *ecdsa.PublicKey) (keyReader, pubReader io.Reader, err error) {
	// BUG WATCH - this function is untested despite being part of the public interface. This is because it
	// doesn't really do much...
	keyEncoded, pubEncoded, err := encode(key, pub)
	if err != nil {
		return nil, nil, err
	}

	return bytes.NewReader(keyEncoded), bytes.NewReader(pubEncoded), nil
}

func EncodedPubECDSAReader(pub *ecdsa.PublicKey) (pubReader io.Reader, err error) {
	// BUG WATCH - this function is untested despite being part of the public interface. This is because it
	// doesn't really do much...
	pubEncoded, err := encodePub(pub)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(pubEncoded), nil
}
