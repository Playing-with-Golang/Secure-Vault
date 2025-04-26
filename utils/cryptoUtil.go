package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
)

// EncryptFile encrypts the input file and writes it to output file, returns the IV
func EncryptFile(inputPath, outputPath string, key []byte) ([]byte, error) {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return nil, err
	}
	defer outFile.Close()

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)
	writer := &cipher.StreamWriter{S: stream, W: outFile}

	if _, err := io.Copy(writer, inFile); err != nil {
		return nil, err
	}

	return iv, nil
}

// DecryptFile decrypts the input file and writes it to output file
func DecryptFile(inputPath, outputPath string, key, iv []byte) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	stream := cipher.NewCTR(block, iv)
	reader := &cipher.StreamReader{S: stream, R: inFile}

	if _, err := io.Copy(outFile, reader); err != nil {
		return err
	}

	return nil
}
