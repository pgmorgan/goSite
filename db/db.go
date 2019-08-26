package db

import (
	"context"
	"fmt"
	"time"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//Book holds the Author and Title strings
type Book struct {
	Title	string
}

func dbConnect(user, password string) (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// Client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://root:root@clusterexp-a6bbr.mongodb.net/test?retryWrites=true&w=majority&authSource=admin"))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func dbInsertOne(collection *mongo.Collection, doc interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := collection.InsertOne(ctx, doc)
	fmt.Println(doc, err)
	if err != nil {
		return err
	}
	return nil
}

func dbList(collection *mongo.Collection) ([]map[string]string, error) {
	var result bson.M
	var list []map[string]string
	m := make(map[string]string)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		// fmt.Println(result["title"])
		if result["title"] == nil {
			return nil, nil
		}
		m["title"] = result["title"].(string)
		// fmt.Println(m)
		list = append(list, m)
	}
	// return listParse(result)
	return list, nil
}

func dbDeleteFiltered(collection *mongo.Collection, filter interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

// func Launch() ([]map[string]string, error) {
// 	// doc := bson.M{"name": "pi", "value": 3.2}
// 	// filter := bson.D{{
// 	// 	"name", bson.D{{
// 	// 		"$in",
// 	// 		bson.A{"pi"},
// 	// 	}},
// 	// }}

// 	// c := color.New(color.FgRed)

// 	client, err := dbConnect(atlasUser, atlasPassword)
// 	if err != nil {
// 		return nil, err
// 	}
// 	collection := client.Database(dbName).Collection(collName)
// 	err = dbInsertOne(collection, doc)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = dbDeleteFiltered(collection, filter)
// 	if err != nil {
// 		return nil, err
// 	}
// 	list, err := dbList(collection)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Println(list)
// 	return list, nil
// }

func PublicList() ([]map[string]string, error) {
	return dbList(globCollection)
}

func PublicInsertOne(book Book) error {
	doc := bson.M{}
	doc["Title"] = book.Title
	// fmt.Println(doc)
	err := dbInsertOne(globCollection, doc)
	return err
}

var globClient *mongo.Client
var globCollection *mongo.Collection

func init() {
	atlasUser := "root"
	atlasPassword := "root"
	dbName := "bookshelf"
	collName := "books"

	globClient, err := dbConnect(atlasUser, atlasPassword)
	if err != nil {
		log.Fatal(err)
	}
	globCollection = globClient.Database(dbName).Collection(collName)
}
