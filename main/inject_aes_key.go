package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func TestInjectAesKey() {
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
