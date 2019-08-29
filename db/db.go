package db

import (
	"context"
	"log"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//Book holds the Author and Title strings
type Book struct {
	Title, URLtitle, Author, ID, Price, BuyLink string
}

var wait time.Duration

func dbConnect(user, password, clusterName, dbName string) (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), wait*time.Second)
	// uri := "mongodb://localhost:27017"
	uri := "mongodb+srv://" + user + ":" + password + "@" + clusterName +
		"-a6bbr.mongodb.net/" + dbName + "?retryWrites=true&w=majority&authSource=admin"
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
	ctx, _ := context.WithTimeout(context.Background(), wait*time.Second)
	doc := bson.M{}
	doc["Title"] = book.Title
	doc["Author"] = book.Author
	doc["Price"] = book.Price
	// doc["Currency"] = book.Currency
	doc["BuyLink"] = book.BuyLink
	doc["ID"] = book.ID
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

	ctx, _ := context.WithTimeout(context.Background(), wait*time.Second)
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
		// if result["Currency"] == nil {
		// tmp.Currency = ""
		// } else {
		// tmp.Currency = result["Currency"].(string)
		// }
		if result["BuyLink"] == nil {
			tmp.BuyLink = ""
		} else {
			tmp.BuyLink = result["BuyLink"].(string)
		}
		list = append(list, tmp)
	}
	return list, nil
}

func DBdeleteOne(title string) error {
	filter := bson.D{{
		"Title", bson.D{{
			"$in",
			bson.A{title},
		}},
	}}
	ctx, _ := context.WithTimeout(context.Background(), wait*time.Second)
	_, err := globCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

var globClient *mongo.Client
var globCollection *mongo.Collection

func DBidAlreadyListed(id string) (bool, error) {
	var result bson.M

	ctx, _ := context.WithTimeout(context.Background(), wait*time.Second)
	cursor, err := globCollection.Find(ctx, bson.D{})
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

func init() {
	wait = 10
	atlasUser := "root"
	atlasPassword := "root"
	clusterName := "cluster0"
	dbName := "bookstore"
	collName := "books"

	globClient, err := dbConnect(atlasUser, atlasPassword, clusterName, dbName)
	if err != nil {
		log.Fatal(err)
	}
	globCollection = globClient.Database(dbName).Collection(collName)
}
