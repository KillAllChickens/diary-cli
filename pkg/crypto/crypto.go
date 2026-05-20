package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

const (
	saltSize   = 16
	nonceSize  = 12
	keyLen     = 32
	scryptN    = 131072
	oldScryptN = 32768
	scryptR    = 8
	scryptP    = 1
)

func HashPassword(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyPassword(password []byte, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), password) == nil
}

func deriveKey(password, salt []byte, n int) ([]byte, error) {
	return scrypt.Key(password, salt, n, scryptR, scryptP, keyLen)
}

func Encrypt(data, password []byte) ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	key, err := deriveKey(password, salt, scryptN)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, data, nil)

	result := make([]byte, 0, saltSize+nonceSize+len(ciphertext))
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result, nil
}

func attemptDecrypt(salt, nonce, ciphertext, password []byte, n int) ([]byte, error) {
	key, err := deriveKey(password, salt, n)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, nonce, ciphertext, nil)
}

func Decrypt(data, password []byte) ([]byte, error) {
	if len(data) < saltSize+nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	salt := data[:saltSize]
	nonce := data[saltSize : saltSize+nonceSize]
	ciphertext := data[saltSize+nonceSize:]

	plaintext, err := attemptDecrypt(salt, nonce, ciphertext, password, scryptN)
	if err == nil {
		return plaintext, nil
	}

	plaintext, fallbackErr := attemptDecrypt(salt, nonce, ciphertext, password, oldScryptN)
	if fallbackErr == nil {
		return plaintext, nil
	}

	return nil, errors.New("authentication failed (incorrect password or corrupted data)")
}
