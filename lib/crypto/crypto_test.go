package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrypto(t *testing.T) {
	key := "test"
	copy(secretKey[:], []byte(key))

	msg := "foobar"
	box, err := Encrypt(msg)
	assert.NoError(t, err)

	secretMsg := Decrypt(box)
	assert.Equal(t, msg, secretMsg)
}
