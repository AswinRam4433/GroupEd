package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	// "gorila/mux"
	"net/http"
)

type FormData struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Discord string `json:"discord"`
}

func submitForm(w http.ResponseWriter, r *http.Request) {
	var formData FormData
	err := json.NewDecoder(r.Body).Decode(&formData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Access the Name and Email fields from the formData struct
	name := formData.Name
	email := formData.Email
	discord := formData.Discord

	// You can now use 'name' and 'email' as needed for processing
	// For example, you can save them to a database or send emails

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Data received: Name = "+name+", Email = "+email+"Discord ID ="+discord)
	dummy, err := GetAESEncrypted(discord)
	if err != nil {
		fmt.Println("Error has occurred")
	} else {
		fmt.Println(dummy)
	}

	dummy1, err := GetAESDecrypted(dummy)
	if err != nil {
		fmt.Println("Error has occurred")
	} else {
		fmt.Println(string(dummy1))
	}
}

func GetAESDecrypted(encrypted string) ([]byte, error) {
	key := "my32digitkey12345678901234567890"
	iv := "my16digitIvKey12"

	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)

	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return nil, err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("block size cant be zero")
	}

	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = PKCS5UnPadding(ciphertext)

	return ciphertext, nil
}

// PKCS5UnPadding  pads a certain blob of data with necessary data to be used in AES block cipher
func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])

	return src[:(length - unpadding)]
}

// GetAESEncrypted encrypts given text in AES 256 CBC
func GetAESEncrypted(plaintext string) (string, error) {
	key := "my32digitkey12345678901234567890"
	iv := "my16digitIvKey12"

	var plainTextBlock []byte
	length := len(plaintext)

	if length%16 != 0 {
		extendBlock := 16 - (length % 16)
		plainTextBlock = make([]byte, length+extendBlock)
		copy(plainTextBlock[length:], bytes.Repeat([]byte{uint8(extendBlock)}, extendBlock))
	} else {
		plainTextBlock = make([]byte, length)
	}

	copy(plainTextBlock, plaintext)
	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(plainTextBlock))
	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, plainTextBlock)

	str := base64.StdEncoding.EncodeToString(ciphertext)

	return str, nil
}

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/submit-form", submitForm)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
