package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"log"

	"golang.org/x/crypto/argon2"
)

func generateSalt(length int) []byte {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		log.Fatalf("error generating salt: %s", err)
	}
	return salt
}

func deriveKey(password string, salt []byte, time, memory uint32, threads uint8, keyLen uint32) []byte {
	return argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)
}

func encryptBytes(key []byte, data []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("aes cipher error: %s", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("error making gcm: %s", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		log.Fatalf("error generating salt: %s", err)
	}

	encryptedData := gcm.Seal(nonce, nonce, data, nil)

	return encryptedData
}

func decryptBytes(key []byte, data []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("aes cipher error: %s", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("error making gcm: %s", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		log.Fatal("ciphertext too short")
	}
	nonce, data := data[:nonceSize], data[nonceSize:]

	decryptedData, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		log.Fatalf("error decrypting data: %s", err)
	}

	return decryptedData
}
