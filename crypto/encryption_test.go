package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	t.Run("success encryption", func(t *testing.T) {
		enryptedMess, err := Encrypt(TestString)
		assert.Nil(t, err)
		_, err = Decrypt(enryptedMess)
		assert.Nil(t, err)
	})
}
