package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Access a collection
	// collection := client.Database("mydb").Collection("mycollection")
	collection := client.Database("mydb").Collection("alerts")

	// // Insert a document
	// insertResult, err := collection.InsertOne(context.TODO(), bson.M{"name": "John Doe", "age": 30})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Inserted document ID: %s\n", insertResult.InsertedID)

	// Find a document
	var result bson.M
	// err = collection.FindOne(context.TODO(), bson.M{"name": "John Doe"}).Decode(&result)
	err = collection.FindOne(context.TODO(), bson.M{"name": "Alice"}).Decode(&result)
	// _, err1 := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Found a document:", result)

	// Close the MongoDB connection when done
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
