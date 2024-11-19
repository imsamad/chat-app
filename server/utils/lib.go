package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), 14)

	return string(bytes), err
}

func CheckPasswordHash(hash string, pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd)) == nil
}
