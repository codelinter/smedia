package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"

	"golang.org/x/crypto/nacl/secretbox"
)

func getSecretBytes() ([]byte, error) {
	secretKeyBytes := []byte(appSymPassword)
	if len(appSymPassword) == 0 {
		return nil, fmt.Errorf("Symmetric password not found")
	}
	return secretKeyBytes, nil
}

func encrypt(plaindata []byte) ([]byte, error) {
	var secretKey [32]byte
	secretKeyBytes, err := getSecretBytes()
	if err != nil {
		return nil, err
	}
	copy(secretKey[:], secretKeyBytes)
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}
	return secretbox.Seal(nonce[:], plaindata, &nonce, &secretKey), nil
}

func decrypt(eData []byte) ([]byte, error) {
	var secretKey [32]byte
	secretKeyBytes, err := getSecretBytes()
	if err != nil {
		return nil, err
	}
	copy(secretKey[:], secretKeyBytes)
	var decryptNonce [24]byte
	copy(decryptNonce[:], eData[:24])
	dData, ok := secretbox.Open(nil, eData[24:], &decryptNonce, &secretKey)
	if !ok {
		return nil, fmt.Errorf(`Decryption failed`)
	}
	return dData, nil
}

// encryptAndWriteToFile encrypts and then writes data to file at the provided path
func encryptAndWriteToFile(path string, data []byte) error {
	eData, err := encrypt(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, eData, 0755)
}

func readFromFileAndDecrypt(path string) ([]byte, error) {
	eData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	data, err := decrypt(eData)
	if err != nil {
		return nil, err
	}
	return data, nil
}
