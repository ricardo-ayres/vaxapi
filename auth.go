package main

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
)

const Size = sha256.Size

type Credentials struct {
	hash [Size]byte
	salt [Size]byte
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

	creds.salt, err = genSalt()
	creds.hash = passwordHash(pwd, creds.salt)
	return creds, err
}

func CheckPassword(pwd string, creds Credentials) bool {
	newhash := passwordHash(pwd, creds.salt)
	return newhash == creds.hash
}
