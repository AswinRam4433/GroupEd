// package main

// import (
// 	"context"
// 	"fmt"
// 	"html/template"
// 	"log"
// 	"net/http"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// type Exam struct {
// 	Name        string
// 	DueDate     string
// 	PostedBy    string
// 	Description string
// }

// func main() {
// 	// Set up MongoDB connection
// 	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
// 	client, err := mongo.NewClient(clientOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	err = client.Connect(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer client.Disconnect(ctx)

// 	// Define the database and collection
// 	db := client.Database("mydb")
// 	collection := db.Collection("alerts")

// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		var exams []Exam

// 		filter := bson.M{"name": "Alice"}

// 		// Find the document with the specified filter
// 		cursor, err := collection.Find(context.TODO(), filter)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		defer cursor.Close(ctx)

// 		for cursor.Next(ctx) {
// 			var exam Exam
// 			if err := cursor.Decode(&exam); err != nil {
// 				fmt.Println("Each iteration of cursor.Next")
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 				return
// 			}
// 			exams = append(exams, exam)

// 			// Print the retrieved document to the console
// 			fmt.Println("Retrieved Document:", exam)
// 		}

// 		// Render the HTML template
// 		tmpl := `
// <!DOCTYPE html>
// <html>
// <head>
//     <style>
//         /* Add your CSS styles here */
//     </style>
// </head>
// <body>
//     <div class="container">
//         <table>
//             <thead>
//                 <tr>
//                     <th>Exam Name</th>
//                     <th>Due Date</th>
//                     <th>Posted By</th>
//                     <th>Description</th>
//                 </tr>
//             </thead>
//             <tbody>
//                 {{range .}}
//                 <tr>
//                     <td>{{.Name}}</td>
//                     <td>{{.DueDate}}</td>
//                     <td>{{.PostedBy}}</td>
//                     <td>{{.Description}}</td>
//                 </tr>
//                 {{end}}
//             </tbody>
//         </table>
//     </div>
// </body>
// </html>
// `
// 		t := template.Must(template.New("exam").Parse(tmpl))
// 		if err := t.Execute(w, exams); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	})

// 	fmt.Println("Server is running on :8080...")
// 	http.ListenAndServe(":8080", nil)
// }

package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Exam struct {
	Name        string
	DueDate     int
	PostedBy    string
	Description string
}

func main() {
	// Set up MongoDB connection
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

		filter := bson.M{"Name": "Udemy Go Lang"}

		// Find the document with the specified filter
		cursor, err := collection.Find(context.TODO(), filter)
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
				<h1>Schedule</h1>
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

	fmt.Println("Server is running on :8080...")
	http.ListenAndServe(":8080", nil)
}
