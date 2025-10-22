package atoms

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPassAtom(pass string) string {
	hash := sha256.Sum256([]byte(pass))
	return hex.EncodeToString(hash[:])
}
