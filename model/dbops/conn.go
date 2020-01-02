package dbops

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

var (
	dbClient *mongo.Client
	err      error
)

func init() {
	log.Println("dbconn init")
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	//TODO: set auth mode on into the mongo db server and client
	dbClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	collection := dbClient.Database("vocaloidDB").Collection("youtubeVideos")

	collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bsonx.Doc{{"id", bsonx.String("text")}},
		Options: options.Index().
			SetBackground(false).
			SetExpireAfterSeconds(10).
			SetName("a").
			SetSparse(false).
			SetUnique(true).
			SetVersion(1).
			SetTextVersion(1).
			SetWeights(bsonx.Doc{}).
			SetSphereVersion(1).
			SetBits(2).
			SetMax(10).
			SetMin(1).
			SetBucketSize(1),
	})

}
