package box

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/xtls/xray-core/constant"
	"io/ioutil"
	"os"
	"strings"
)

func DecryptAES(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(constant.ENCRYPT_KEY))
	mode := cipher.NewCBCDecrypter(block, []byte(constant.ENCRYPT_KEY_IV))
	plaintext := make([]byte, len(ciphertext))
	ciphertext, err = base64.StdEncoding.DecodeString(string(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)
	padSize := int(plaintext[len(plaintext)-1])
	return plaintext[:len(plaintext)-padSize], err
}

func InjectAesKey() {
	content := fmt.Sprintf(`package constant

var (
	ENCRYPTED_CONFIG = false
	ENCRYPT_KEY      = "%s"
	ENCRYPT_KEY_IV   = "%s"
)
`, getSecret("ENCRYPT_KEY"),
		getSecret("ENCRYPT_KEY_IV"))

	err := ioutil.WriteFile("./constant/constants.go", []byte(content), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func getSecret(secretName string) string {
	secret, exists := GetEnv(secretName)
	if !exists {
		fmt.Printf("Error: Secret %s not found\n", secretName)
	}
	return secret
}

func GetEnv(key string) (string, bool) {
	value, exists := LookupEnv(key)
	if !exists {
		// Fallback to uppercase key if not found
		value, exists = LookupEnv(ToUpperCase(key))
	}
	return value, exists
}

func LookupEnv(key string) (string, bool) {
	val, exists := GetEnvVar(key)
	return val, exists
}

func GetEnvVar(key string) (string, bool) {
	val, exists := LookupEnvVar(key)
	return val, exists
}

func LookupEnvVar(key string) (string, bool) {
	val, exists := lookupEnv(key)
	return val, exists
}

// ToUpperCase converts a string to uppercase.
func ToUpperCase(s string) string {
	return strings.ToUpper(s)
}

func lookupEnv(key string) (string, bool) {
	value, exists := os.LookupEnv(key)
	return value, exists
}
