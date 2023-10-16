package main

import (
	// "./new_back.go"
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		// http.Redirect(w, r, "/your-next-page", http.StatusSeeOther)
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

type Exam struct {
	Name        string
	DueDate     int
	PostedBy    string
	Description string
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

	formData.Passw, _ = GetAESEncrypted(formData.Passw)

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
		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer client.Disconnect(context.Background())

		// Access the "mydb" database and the "formData" collection
		db := client.Database("mydb")
		collection := db.Collection("formData")
		_, err = collection.InsertOne(context.Background(), formData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Wrote to MongoDB")
		cookie := http.Cookie{Name: "myCookiePhNo", Value: formData.Phno}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/otp-verification", http.StatusSeeOther)
	} else {
		// Handle the case when 'answer' is false
		// For example, you could display an error message or redirect to another page
	}

	// sendOtp(string(formData.Phno))

}

type signInStruct struct {
	Discord string `json:"discord-id"`
	Passw   string `json:"password"`
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }
	fmt.Println("The request method is ", r.Method)

	var formData signInStruct
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	formData.Discord = r.FormValue("discord-id")
	formData.Passw = r.FormValue("password")

	// Process the form data as needed
	fmt.Printf("Received data: %+v\n", formData)
	fmt.Println("The discord ID is ", formData.Discord)
	fmt.Println("The password is ", formData.Passw)
}

func nextPageSenderHandler(w http.ResponseWriter, r *http.Request) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Define the database and collection
	db := client.Database("mydb")
	collection := db.Collection("alerts")
	var exams []Exam

	// Filter, if needed
	// filter := bson.M{"Name": "Udemy Go Lang"}

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var exam Exam
		if err := cursor.Decode(&exam); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		exams = append(exams, exam)
		// Print the retrieved document to the console
		fmt.Println("Retrieved Document:", exam)
	}

	// Render the HTML template
	tmpl := `
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #ebe9a1;
				}
	
				.container {
					max-width: 800px;
					margin: 0 auto;
					padding: 20px;
					box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.2);
					border-radius: 10px;
					background-color: #f2f2f2;
				}
	
				h1 {
					text-align: center;
				}
	
				table {
					width: 100%;
					border-collapse: collapse;
					margin-top: 20px;
				}
	
				th, td {
					border: 1px solid #dddddd;
					text-align: left;
					padding: 8px;
				}
	
				th {
					background-color: #f2f2f2;
					text-align: center;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Your Schedule</h1>
				<table>
					<thead>
						<tr>
							<th>Exam Name</th>
							<th>Due Date</th>
							<th>Posted By</th>
							<th>Description</th>
						</tr>
					</thead>
					<tbody>
						{{range .}}
						<tr>
							<td>{{.Name}}</td>
							<td>{{.DueDate}}</td>
							<td>{{.PostedBy}}</td>
							<td>{{.Description}}</td>
						</tr>
						{{end}}
					</tbody>
				</table>
			</div>
		</body>
		</html>
		`
	t := template.Must(template.New("exam").Parse(tmpl))
	if err := t.Execute(w, exams); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func nextPageSender(w http.ResponseWriter, r *http.Request) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Define the database and collection
	db := client.Database("mydb")
	collection := db.Collection("alerts")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var exams []Exam

		// filter := bson.M{"Name": "Udemy Go Lang"}

		// Find the document with the specified filter
		// cursor, err := collection.Find(context.TODO(), filter)
		cursor, err := collection.Find(context.TODO(), bson.D{})

		if err != nil {
			log.Fatal(err)
		}

		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var exam Exam
			if err := cursor.Decode(&exam); err != nil {
				fmt.Println("Each iteration of cursor.Next")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			exams = append(exams, exam)

			// Print the retrieved document to the console
			fmt.Println("Retrieved Document:", exam)
		}

		// Render the HTML template
		tmpl := `
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body {
					font-family: Arial, sans-serif;
					background-color: #ebe9a1;
				}
		
				.container {
					max-width: 800px;
					margin: 0 auto;
					padding: 20px;
					box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.2);
					border-radius: 10px;
					background-color: #f2f2f2;

				}
		
				h1 {
					text-align: center;
				}
		
				table {
					width: 100%;
					border-collapse: collapse;
					margin-top: 20px;
				}
		
				th, td {
					border: 1px solid #dddddd;
					text-align: left;
					padding: 8px;
				}
		
				th {
					background-color: #f2f2f2;
					text-align: center;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Your Schedule</h1>
				<table>
					<thead>
						<tr>
							<th>Exam Name</th>
							<th>Due Date</th>
							<th>Posted By</th>
							<th>Description</th>
						</tr>
					</thead>
					<tbody>
						{{range .}}
						<tr>
							<td>{{.Name}}</td>
							<td>{{.DueDate}}</td>
							<td>{{.PostedBy}}</td>
							<td>{{.Description}}</td>
						</tr>
						{{end}}
					</tbody>
				</table>
			</div>
		</body>
		</html>
		
`
		t := template.Must(template.New("exam").Parse(tmpl))
		if err := t.Execute(w, exams); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func otpVerificationPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Triggered otp verification page function")
	fmt.Println("Cookies in API Call:")

	tokenCookie, err := r.Cookie("myCookiePhNo")
	if err != nil {
		log.Fatalf("Error occurred while reading cookie")
	}
	fmt.Println("\nPrinting cookie with phone number as token")
	fmt.Println(tokenCookie)
	fmt.Println("The key is ", tokenCookie.Name)
	fmt.Println("The value is ", tokenCookie.Value)

	sendOtp(tokenCookie.Value)
	fmt.Println("\nSending the OTP")

	var otpData OtpFormData
	err1 := r.ParseForm()
	if err1 != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	otpData.Otp = r.FormValue("otp")
	fmt.Println(otpData.Otp)

	otpStatus := checkOtp(tokenCookie.Value, otpData.Otp)
	if otpStatus == true {
		http.Redirect(w, r, "/your-next-page", http.StatusSeeOther)
		return // Make sure to return after the redirection
	}

	// If OTP is not successfully verified, serve the OTP verification HTML page
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

func main() {

	// http.HandleFunc("/signin", signInHandler)
	http.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/sign_in.html")
	})

	// Serve the HTML file
	http.Handle("/", http.FileServer(http.Dir("static")))

	// Handle form submissions
	// http.HandleFunc("/signin", signInHandler)
	http.HandleFunc("/submit-form", submitForm)
	http.HandleFunc("/sign-in-val", signInHandler)
	http.HandleFunc("/otp-verification", otpVerificationPage)

	http.HandleFunc("/your-next-page", nextPageSenderHandler)

	fmt.Println("Server is running on :8080...")
	http.ListenAndServe(":8080", nil)
}
