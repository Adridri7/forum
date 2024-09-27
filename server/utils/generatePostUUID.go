package posts

import (
	"crypto/rand"
	"fmt"
)

// GenerateUUID génère un UUID aléatoire
func GenerateUUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	b[6] = b[6]&^0xf0 | 0x40 // Version 4
	b[8] = b[8]&^0x3f | 0x80 // Variante RFC 4122
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
