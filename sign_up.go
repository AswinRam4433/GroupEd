package main

import (
	// "./new_back.go"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"
)

// Full package path

type FormData struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Discord   string `json:"discord"`
	Phno      string `json:"phno"`
	Subscribe string `json:"subscribe"`
	College   string `json:"college"`
	Gender    string `json:"gender"`
	Passw     string `json:"passw"`
}

func submitForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var formData FormData
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	formData.Name = r.FormValue("name")
	formData.Email = r.FormValue("email")
	formData.Discord = r.FormValue("discord")
	formData.Phno = r.FormValue("phno")
	formData.Subscribe = r.FormValue("subscribe")
	formData.College = r.FormValue("college")
	formData.Gender = r.FormValue("gender")
	formData.Passw = r.FormValue("passw")

	// Process the form data as needed
	fmt.Printf("Received data: %+v\n", formData)
	fmt.Println("The discord ID is ", formData.Discord)
	enc_disc, err := GetAESEncrypted(formData.Discord)
	fmt.Println(enc_disc)
	// w.WriteHeader(http.StatusOK)
	// io.WriteString(w, "Form data received")
	answer := valFormData(formData)
	print("Result of answer is ", answer)

	if answer == true {
		cookie := http.Cookie{Name: "myCookie", Value: "someValue"}
		http.SetCookie(w, &cookie)

		// cookie1, err := r.Cookie("myCookie")
		// if err != nil {
		// 	// Handle error
		// 	value := cookie1.Value
		// 	print(value)
		// }

		// http.Redirect(w, r, "/otp.html", http.StatusSeeOther)
		otpVerificationPage(w, r)

	} else {
	}

	// sendOtp(string(formData.Phno))

}

func otpVerificationPage(w http.ResponseWriter, r *http.Request) {
	// Serve your OTP verification HTML page here
	fmt.Println("Triggered otp verification page function")
	http.ServeFile(w, r, `static\otp.html`)
}

func valFormData(f FormData) bool {
	if f.Name != "" && f.Discord != "" && f.Email != "" && f.Passw != "" && f.Gender != "" && f.Phno != "" && valPassword(f.Passw) {
		return true
	}
	return false
}

func valPassword(p string) bool {
	pattern := `(?i)[A-Z]+.*\d+|(?i)\d+.*[A-Z]+`
	re := regexp.MustCompile(pattern)
	fmt.Println("The password is ", p)

	if len(p) >= 8 && re.MatchString(p) {
		fmt.Println("Password passed")
		return true

	}
	fmt.Println("Password failed")
	return false

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
		return nil, fmt.Errorf("block size can't be zero")
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

// func main() {
// 	r := http.NewServeMux()
// 	// http.Handle("/", http.FileServer(http.Dir("static")))
// 	r.HandleFunc("/submit-form", submitForm)

// 	http.Handle("/", r)
// 	http.ListenAndServe(":8080", nil)
// }

func main() {
	// Serve the HTML file
	http.Handle("/", http.FileServer(http.Dir("static")))

	// Handle form submissions
	http.HandleFunc("/submit-form", submitForm)

	http.HandleFunc("/otp-verification", otpVerificationPage)

	fmt.Println("Server is running on :8080...")
	http.ListenAndServe(":8080", nil)
}
