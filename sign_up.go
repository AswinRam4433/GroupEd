package main

import (
	// "./new_back.go"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
)

// Full package path

var TWILIO_ACCOUNT_SID string
var TWILIO_AUTH_TOKEN string
var VERIFY_SERVICE_SID string

var client *twilio.RestClient

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	// Initialize environment variables
	TWILIO_ACCOUNT_SID = os.Getenv("TWILIO_ACCOUNT_SID")
	TWILIO_AUTH_TOKEN = os.Getenv("TWILIO_AUTH_TOKEN")
	VERIFY_SERVICE_SID = os.Getenv("VERIFY_SERVICE_SID")

	// Initialize Twilio client
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: TWILIO_ACCOUNT_SID,
		Password: TWILIO_AUTH_TOKEN,
	})

	// fmt.Println((client))
}

func sendOtp(to string) {
	fmt.Println("Sending OTP")

	params := &openapi.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(VERIFY_SERVICE_SID, params)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Sent verification '%s'\n", *resp.Sid)
	}
}

func checkOtp(to string, code string) bool {
	// var code string
	// fmt.Println("Please check your phone and enter the code:")
	// fmt.Scanln(&code)

	params := &openapi.CreateVerificationCheckParams{}
	params.SetTo(to)
	params.SetCode(code)

	resp, err := client.VerifyV2.CreateVerificationCheck(VERIFY_SERVICE_SID, params)

	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if *resp.Status == "approved" {
		fmt.Println("Correct!")
		return true
	} else {
		fmt.Println("Incorrect!")
		return false
	}
}

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
type OtpFormData struct {
	Otp string `json:"otp"`
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
		cookie := http.Cookie{Name: "myCookiePhNo", Value: formData.Phno}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/otp-verification", http.StatusSeeOther)
	} else {
		// Handle the case when 'answer' is false
		// For example, you could display an error message or redirect to another page
	}

	// sendOtp(string(formData.Phno))

}

func otpVerificationPage(w http.ResponseWriter, r *http.Request) {
	// Serve your OTP verification HTML page here

	fmt.Println("Triggered otp verification page function")
	fmt.Println("Cookies in API Call:")

	tokenCookie, err := r.Cookie("myCookiePhNo")
	if err != nil {
		log.Fatalf("Error occured while reading cookie")
	}
	fmt.Println("\nPrinting cookie with phone number as token")
	fmt.Println(tokenCookie)
	fmt.Println("The key is ", tokenCookie.Name)
	fmt.Println("The value is ", tokenCookie.Value)

	sendOtp(tokenCookie.Value)
	fmt.Println("\nSending the OTP")
	// for _, c := range r.Cookies() {
	// 	fmt.Println(c)
	// }

	fmt.Println()
	http.ServeFile(w, r, `static\otp.html`)

	var otpData OtpFormData
	err1 := r.ParseForm()
	if err1 != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	otpData.Otp = r.FormValue("otp")
	fmt.Println(otpData.Otp)

	checkOtp(tokenCookie.Value, otpData.Otp)

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
