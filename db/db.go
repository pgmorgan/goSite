package db

import (
	"context"
	"errors"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Book struct {
	Title, URLtitle, Author, ID, Price, BuyLink, ThumbURL string
}

var wait time.Duration = 10
var globClient *mongo.Client
var globCollection *mongo.Collection
var dbName = "bookstore"

func init() {
	var err error
	if globClient, err = dbConnect(); err != nil {
		log.Fatal(err)
	}
}

/*	CONNECTION TO DATABASE - Called by init()	*/
func dbConnect() (*mongo.Client, error) {
	if err := godotenv.Load("./dev.env"); err != nil {
		return nil, errors.New(`No .env file found at root of repository.
		Must set the following environment variables:\n
		MONGODB_URL=<MongoDB Connection String>\n
		GOOGLE_DEV_API_KEY=<API Key from https://console.developers.google.com>\n`)
	}
	uri, exists := os.LookupEnv("MONGODB_URL")
	if !exists {
		return nil, errors.New(`Missing MONGODB_URL environment variable in .env file at root of repository`)
	}
	ctx, _ := context.WithTimeout(context.Background(), wait*time.Second)
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

/*	DATABASE ENTRY CREATION */
func DBinsertOne(book Book, userEmail string) error {
	coll := globClient.Database(dbName).Collection(userEmail)
	ctx, _ := context.WithTimeout(context.Background(), wait*time.Second)
	doc := bson.M{}
	doc["Title"] = book.Title
	doc["Author"] = book.Author
	doc["Price"] = book.Price
	doc["BuyLink"] = book.BuyLink
	doc["ID"] = book.ID
	doc["ThumbURL"] = book.ThumbURL
	_, err := coll.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	return nil
}

/*	DATABASE ENTRY READING */
func DBlist(userEmail string) ([]Book, error) {
	var result bson.M
	var list []Book
	var tmp Book

	coll := globClient.Database(dbName).Collection(userEmail)
	ctx, _ := context.WithTimeout(context.Background(), wait*time.Second)
	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		if result["Title"] == nil {
			return nil, nil
		} else {
			tmp.Title = result["Title"].(string)
			tmp.URLtitle = url.QueryEscape(result["Title"].(string))

		}
		if result["Author"] == nil {
			tmp.Author = ""
		} else {
			tmp.Author = result["Author"].(string)
		}
		if result["Price"] == nil {
			tmp.Price = ""
		} else {
			tmp.Price = result["Price"].(string)
		}
		if result["BuyLink"] == nil {
			tmp.BuyLink = ""
		} else {
			tmp.BuyLink = result["BuyLink"].(string)
		}
		if result["ThumbURL"] == nil {
			tmp.ThumbURL = ""
		} else {
			tmp.ThumbURL = result["ThumbURL"].(string)
		}
		list = append(list, tmp)
	}
	return list, nil
}

/*	DATABASE ENTRY REMOVAL	*/
func DBdeleteOne(title, userEmail string) error {
	coll := globClient.Database(dbName).Collection(userEmail)
	filter := bson.D{{
		"Title", bson.D{{
			"$in",
			bson.A{title},
		}},
	}}
	ctx, _ := context.WithTimeout(context.Background(), wait*time.Second)
	_, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

/*	DATABASE ENTRY DUPLICATE CHECK	*/
func DBidAlreadyListed(id, userEmail string) (bool, error) {
	var result bson.M

	coll := globClient.Database(dbName).Collection(userEmail)
	ctx, _ := context.WithTimeout(context.Background(), wait*time.Second)
	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return true, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return true, err
		}
		if result["ID"] == id {
			return true, nil
		}
	}
	return false, nil
}
