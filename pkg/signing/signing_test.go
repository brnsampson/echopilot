package signing_test

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/brnsampson/echopilot/pkg/signing"
)

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

// equalsOctal fails the test if exp is not equal to act. The difference to equals() is that it only works with
// numbers and formats output as octal.
func equalsOctal(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %O\n\n\tgot: %O\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func testKeys(t *testing.T, key *ecdsa.PrivateKey, pub *ecdsa.PublicKey) (bool, error) {
	// Just sign then verify with the key to make sure it's functional
	msg := "This is a seeeeecret"
	hash := sha256.Sum256([]byte(msg))
	sig, err := ecdsa.SignASN1(rand.Reader, key, hash[:])
	if err != nil {
		return false, err
	}

	valid := ecdsa.VerifyASN1(pub, hash[:], sig)
	return valid, nil
}

func TestGenerateECDSA(t *testing.T) {
	key, pub, err := signing.GenerateECDSA()
	ok(t, err)

	valid, err := testKeys(t, key, pub)
	ok(t, err)
	assert(t, valid, "Could not use generated keys to sign and verify a hash")
}

func TestReadECDSA(t *testing.T) {
	key, pub, err := signing.GenerateECDSA()
	ok(t, err)

	keyReader, pubReader, err := signing.EncodedECDSAReaders(key, pub)
	ok(t, err)

	newKey, newPub, err := signing.ReadECDSA(keyReader, pubReader)
	ok(t, err)

	valid, err := testKeys(t, newKey, newPub)
	ok(t, err)
	assert(t, valid, "Could not use generated keys to sign and verify a hash")
	equals(t, key, newKey)
	equals(t, pub, newPub)
}

func TestReadPubECDSA(t *testing.T) {
	key, pub, err := signing.GenerateECDSA()
	ok(t, err)

	pubReader, err := signing.EncodedPubECDSAReader(pub)
	ok(t, err)

	newPub, err := signing.ReadPubECDSA(pubReader)
	ok(t, err)

	valid, err := testKeys(t, key, newPub)
	ok(t, err)
	assert(t, valid, "Could not use generated keys to sign and verify a hash")
	equals(t, pub, newPub)
}

func TestWriteECDSA(t *testing.T) {
	var keySlice []byte
	var pubSlice []byte
	keyBuf := bytes.NewBuffer(keySlice)
	pubBuf := bytes.NewBuffer(pubSlice)

	key, pub, err := signing.GenerateECDSA()
	ok(t, err)

	signing.WriteECDSA(key, keyBuf, pub, pubBuf)

	newKey, newPub, err := signing.ReadECDSA(keyBuf, pubBuf)
	ok(t, err)

	valid, err := testKeys(t, newKey, newPub)
	ok(t, err)
	assert(t, valid, "Could not use generated keys to sign and verify a hash")
	equals(t, key, newKey)
	equals(t, pub, newPub)
}

func TestKeyFilePermsValid(t *testing.T) {
	const keyPerms fs.FileMode = 0600
	const pubPerms fs.FileMode = 0644

	keyFile, err := os.CreateTemp("", "go_signing_test_*")
	defer keyFile.Close()
	defer os.Remove(keyFile.Name())
	ok(t, err)
	keyFile.Chmod(keyPerms)
	info, err := keyFile.Stat()
	ok(t, err)
	equalsOctal(t, info.Mode(), keyPerms)

	pubFile, err := os.CreateTemp("", "go_signing_test_*")
	defer pubFile.Close()
	defer os.Remove(pubFile.Name())
	ok(t, err)
	pubFile.Chmod(pubPerms)
	info, err = pubFile.Stat()
	equalsOctal(t, info.Mode(), pubPerms)
	ok(t, err)

	keyValid, pubValid, err := signing.KeyFilePermsValid(keyFile.Name(), pubFile.Name())
	ok(t, err)
	assert(t, keyValid, "keyfile permissions 600 was reported invalid, but should have been valid!")
	assert(t, pubValid, "public keyfile permissions 644 was reported invalid, but should have been valid!")

	// Set the wrong perms!
	keyFile.Chmod(pubPerms)
	pubFile.Chmod(keyPerms)
	keyValid, pubValid, err = signing.KeyFilePermsValid(keyFile.Name(), pubFile.Name())
	ok(t, err)
	assert(t, !keyValid, "keyfile permissions 644 was reported valid, but should have been invalid!")
	assert(t, !pubValid, "public keyfile permissions 600 was reported valid, but should have been invalid!")
}

func TestPubKeyFilePermsValid(t *testing.T) {
	const keyPerms fs.FileMode = 0600
	const pubPerms fs.FileMode = 0644

	pubFile, err := os.CreateTemp("", "go_signing_test_*")
	defer pubFile.Close()
	defer os.Remove(pubFile.Name())
	ok(t, err)
	pubFile.Chmod(pubPerms)
	info, err := pubFile.Stat()
	equalsOctal(t, info.Mode(), pubPerms)
	ok(t, err)

	pubValid, err := signing.PubKeyFilePermsValid(pubFile.Name())
	ok(t, err)
	assert(t, pubValid, "public keyfile permissions 644 was reported invalid, but should have been valid!")

	// Set the wrong perms!
	pubFile.Chmod(keyPerms)
	pubValid, err = signing.PubKeyFilePermsValid(pubFile.Name())
	ok(t, err)
	assert(t, !pubValid, "public keyfile permissions 600 was reported valid, but should have been invalid!")
}

func TestSetKeyFilePerms(t *testing.T) {
	const keyPerms fs.FileMode = 0600
	const pubPerms fs.FileMode = 0644

	keyFile, err := os.CreateTemp("", "go_signing_test_*")
	defer keyFile.Close()
	defer os.Remove(keyFile.Name())
	ok(t, err)
	// Intentionally wrong perms
	keyFile.Chmod(pubPerms)

	pubFile, err := os.CreateTemp("", "go_signing_test_*")
	defer pubFile.Close()
	defer os.Remove(pubFile.Name())
	ok(t, err)
	// Intentionally wrong perms
	pubFile.Chmod(keyPerms)

	println(keyFile.Name())
	err = signing.SetKeyFilePerms(keyFile.Name(), pubFile.Name())
	ok(t, err)

	info, err := keyFile.Stat()
	ok(t, err)
	equalsOctal(t, info.Mode(), keyPerms)

	info, err = pubFile.Stat()
	ok(t, err)
	equalsOctal(t, info.Mode(), pubPerms)

	keyValid, pubValid, err := signing.KeyFilePermsValid(keyFile.Name(), pubFile.Name())
	ok(t, err)
	assert(t, keyValid, "keyfile perms were not set correctly by SetKeyFilePerms")
	assert(t, pubValid, "public keyfile permissions were not set correclty by SetKeyFilePerms")
}
