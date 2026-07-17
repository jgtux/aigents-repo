package atoms

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassAtom(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePassAtom(hashedPass, tryPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(tryPass)) == nil
}
