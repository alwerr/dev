package dev

import (
	"encoding/hex"

	"golang.org/x/crypto/blake2b"
)

func PassHash(password string) string {
	salt := "./data/db.db?mode=rwc"

	hasher, err := blake2b.New(64, []byte(salt))
	if err != nil {
		panic(err)
	}
	hasher.Write([]byte(password))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

func PassCheck(password, hash string) bool {
	return password == PassHash(hash)
}
