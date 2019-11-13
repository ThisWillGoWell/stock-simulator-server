package user

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func hashAndSalt(password string) string {
	pwd := []byte(password)
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	bytePW := []byte(plainPwd)
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePW)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
