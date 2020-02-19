package crypto

const TestString = "encryptme"

// func TestDecrypt(t *testing.T) {
// 	t.Run("success decryption", func(t *testing.T) {
// 		enryptedMess, err := Encrypt(TestString)
// 		assert.Nil(t, err)
// 		actual, err := Decrypt(enryptedMess)
// 		assert.Nil(t, err)
// 		assert.Equal(t, TestString, actual)
// 	})
// 	t.Run("failure decryption", func(t *testing.T) {
// 		_, err := Decrypt(TestString)
// 		assert.NotNil(t, err)
// 	})
// }
