package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func encrypt(text string, key []byte) (string, error) {
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a nonce. Nonce should be unique for each encryption.
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)

	// Combine the nonce and ciphertext for storage or transmission
	result := append(nonce, ciphertext...)

	return base64.URLEncoding.EncodeToString(result), nil
}

func decrypt(ciphertext string, key []byte) (string, error) {
	ciphertextBytes, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertextBytes) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce := ciphertextBytes[:nonceSize]
	ciphertextBytes = ciphertextBytes[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func main() {
	// Replace this key with a secure key of 16, 24, or 32 bytes for AES-128, AES-192, or AES-256.
	key := []byte("thisisasecretkey")

	// Text to be encrypted
	originalText := "Hello, Golang!Hello, Golang!Hello, Golang!Hello, Golang!Hello, Golang!Hello, Golang!"

	// Encrypt
	encryptedText, err := encrypt(originalText, key)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return
	}
	fmt.Println("Encrypted:", encryptedText)

	// Decrypt
	decryptedText, err := decrypt(encryptedText, key)
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}
	fmt.Println("Decrypted:", decryptedText)
}
