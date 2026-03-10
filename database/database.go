package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var DB *gorm.DB
var MongoClient *mongo.Client
var MongoDB *mongo.Database

func DatabaseInit() {

	var err error

	// Ignore error if .env file is not found
	_ = godotenv.Load()

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Create DSN for PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Database connected")
}

func DatabaseInitMongo() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatal(err)
	}

	var result bson.M
	if err := client.Database("n8n").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		log.Fatal(err)
	}

	MongoClient = client
	MongoDB = client.Database("n8n")

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}

func DatabaseCloseMongo() {
	if MongoClient == nil {
		return
	}

	if err := MongoClient.Disconnect(context.TODO()); err != nil {
		log.Println("failed to disconnect mongo:", err)
		return
	}

	log.Println("mongo disconnected")
}
