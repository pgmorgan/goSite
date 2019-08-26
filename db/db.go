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

//Book holds the Author and Title strings
type Book struct {
	Title string
}

func dbConnect(user, password, dbName string) (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// "mongodb://localhost:27017"
	uri := "mongodb+srv://" + user + ":" + password + "@clusterexp-a6bbr.mongodb.net/" +
		dbName + "?retryWrites=true&w=majority&authSource=admin"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func DBinsertOne(book Book) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	doc := bson.M{}
	doc["Title"] = book.Title
	_, err := globCollection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	return nil
}

func DBlist() ([]Book, error) {
	var result bson.M
	var list []Book
	var tmp Book
	// var m map[string]string
	// m := make(map[string]string)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// err := globClient.Ping(ctx, readpref.Primary())
	// if err != nil {
	// 	fmt.Println("reached here too!")
	// 	return nil, nil
	// }

	cursor, err := globCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		// fmt.Println(result)
		if result["Title"] == nil {
			fmt.Println("reached here!")

			return nil, nil
		}
		// fmt.Println(result, m)

		tmp = Book{Title: result["Title"].(string)}
		// fmt.Println(m)
		list = append(list, tmp)
	}
	// return listParse(result)
	fmt.Println(list)
	return list, nil
}

func DBdeleteFiltered(title string) error {
	filter := bson.D{{
		"Title", bson.D{{
			"$in",
			bson.A{title},
		}},
	}}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := globCollection.DeleteOne(ctx, filter)
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
// 	globCollection := client.Database(dbName).Collection(collName)
// 	err = dbInsertOne(globCollection, doc)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = dbDeleteFiltered(globCollection, filter)
// 	if err != nil {
// 		return nil, err
// 	}
// 	list, err := dbList(globCollection)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Println(list)
// 	return list, nil
// }

var globClient *mongo.Client
var globCollection *mongo.Collection

func init() {
	atlasUser := "root"
	atlasPassword := "root"
	dbName := "bookshelf"
	collName := "books"

	globClient, err := dbConnect(atlasUser, atlasPassword, dbName)
	if err != nil {
		log.Fatal(err)
	}
	globCollection = globClient.Database(dbName).Collection(collName)
}
