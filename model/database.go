package model

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongodatabase struct {
	Mclient   *mongo.Client
	Mdatabase *mongo.Database
}

var MongoInstance Mongodatabase
var (
	dbName = "KoulNetwork"
)

func ConnectMongo() error {
	godotenv.Load()

	connectionstring := os.Getenv("ConnectionString")

	ClientOption := options.Client().ApplyURI(connectionstring)
	Client, err := mongo.Connect(context.TODO(), ClientOption)
	if err != nil {
		return err
	}
	db := Client.Database(dbName)

	MongoInstance = Mongodatabase{Mclient: Client, Mdatabase: db}
	return nil
}
