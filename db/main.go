package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func dbConnect(user, password string) *mongo.Client {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://root:root@clusterexp-a6bbr.mongodb.net/test?retryWrites=true&w=majority&authSource=admin"))
	check(err)
	err = client.Ping(ctx, readpref.Primary())
	check(err)
	return client
}

func dbInsertOne(collection *mongo.Collection, doc interface{}) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := collection.InsertOne(ctx, doc)
	check(err)
}

func dbList(collection *mongo.Collection) {
	var result bson.M

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.D{})
	check(err)
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		err = cursor.Decode(&result)
		check(err)
		fmt.Println(result)
	}
}

func dbDeleteFiltered(collection *mongo.Collection, filter interface{}) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := collection.DeleteOne(ctx, filter)
	check(err)
}

func main() {
	atlasUser := "root"
	atlasPassword := "root"
	dbName := "test"
	collName := "numbers"
	doc := bson.M{"name": "pi", "value": 3.14159}
	// filter := bson.D{{
	// "name", bson.D{{
	// "$in",
	// bson.A{"pi"},
	// }},
	// }}

	// c := color.New(color.FgRed)

	client := dbConnect(atlasUser, atlasPassword)
	collection := client.Database(dbName).Collection(collName)
	dbInsertOne(collection, doc)
	dbList(collection)
	// dbDeleteFiltered(collection, filter)
	// dbDeleteFiltered(collection, filter)
	// dbDeleteFiltered(collection, filter)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
