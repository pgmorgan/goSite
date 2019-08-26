package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//Book holds the author and title strings
type Book struct {
	//	author string
	title string
}

func DBconnect(user, password string) (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
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

func DBinsertOne(collection *mongo.Collection, doc interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	return nil
}

func DBlist(collection *mongo.Collection) ([]map[string]float64, error) {
	var result bson.M
	var list []map[string]float64
	m := make(map[string]float64)

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
		fmt.Println(result["name"], result["value"])
		m[result["name"].(string)] = result["value"].(float64)
		fmt.Println(m)
		list = append(list, m)
	}
	// return listParse(result)
	return list, nil
}

func DBdeleteFiltered(collection *mongo.Collection, filter interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func Launch() ([]map[string]float64, error) {
	atlasUser := "root"
	atlasPassword := "root"
	dbName := "test"
	collName := "numbers"
	doc := bson.M{"name": "pi", "value": 3.2}
	filter := bson.D{{
		"name", bson.D{{
			"$in",
			bson.A{"pi"},
		}},
	}}

	// c := color.New(color.FgRed)

	client, err := DBconnect(atlasUser, atlasPassword)
	if err != nil {
		return nil, err
	}
	collection := client.Database(dbName).Collection(collName)
	err = DBinsertOne(collection, doc)
	if err != nil {
		return nil, err
	}
	err = DBdeleteFiltered(collection, filter)
	if err != nil {
		return nil, err
	}
	list, err := DBlist(collection)
	if err != nil {
		return nil, err
	}
	fmt.Println(list)
	return list, nil
}
