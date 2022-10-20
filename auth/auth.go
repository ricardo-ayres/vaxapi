package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
)

const Size = sha256.Size

type Credentials struct {
	Hash [Size]byte
	Salt [Size]byte
}

func genSalt() ([Size]byte, error) {
	var salt [Size]byte
	var err error

	n, err := rand.Read(salt[:])
	if err != nil {
		return salt, err
	}

	if n != Size {
		err = errors.New("Salt generated with wrong size.")
	}
	return salt, err
}

func passwordHash(pwd string, salt [Size]byte) [Size]byte {
	salted := append([]byte(pwd), salt[:]...)
	return sha256.Sum256(salted)
}

func NewCredentials(pwd string) (Credentials, error) {
	var creds Credentials
	var err error

	creds.Salt, err = genSalt()
	creds.Hash = passwordHash(pwd, creds.Salt)
	return creds, err
}

func CheckPassword(pwd string, creds Credentials) bool {
	newhash := passwordHash(pwd, creds.Salt)
	return newhash == creds.Hash
}
