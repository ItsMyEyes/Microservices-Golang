package configs

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	connectMongo()
}

var (
	Client *mongo.Client
	Users  *mongo.Collection
	DB     *mongo.Database
	Ctx    = context.TODO()
)

func connectMongo() {
	var err error
	Client, err = mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URL")))
	if err != nil {
		panic(err)
	}
	err = Client.Connect(context.TODO())
	if err != nil {
		panic(err)
	}

	err = Client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	DB = Client.Database(os.Getenv("MONGO_DB"))
	Users = DB.Collection("users")

	fmt.Println("Connected to MongoDB!")
}
