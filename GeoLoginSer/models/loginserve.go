package model

import (
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID                primitive.ObjectID `bson:"_id" json:"userid,omitempty"`
	UID               string             `json:"uid,omitempty"`
	Username          string             `json:"username,omitempty"`
	Email             string             `json:"email,omitempty"`
	PhoneNo           string             `json:"phoneno,omitempty"`
	FirstName         string             `json:"firstname,omitempty"`
	LastName          string             `json:"lastname,omitempty"`
	Oldpassword       string             `json:"oldpassword,omitempty"`
	Password          string             `json:"password,omitempty"`
	Devid             string             `json:"devid,omitempty"`
	Devids            []string           `json:"devids,omitempty"`
	Token             string             `json:"token,omitempty"`
	Lat               string             `json:"lat,omitempty"`
	Long              string             `json:"long,omitempty"`
	Location          LocStruct          `json:"location,omitempty"`
	EmergencyContacts []EmergencyContact `bson:"emergency_contacts" json:"emergency_contacts,omitempty"`
}

type EmergencyContact struct {
	Name        string `json:"name,omitempty"`
	PhoneNumber string `bson:"phone_number" json:"phone_number,omitempty"`
}

type LocStruct struct {
	Type        string    `bson:"type" json:"type,omitempty"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates,omitempty"`
}
type ResponseResult struct {
	Error   string `json:"error,omitempty"`
	Result  string `json:"result,omitempty"`
	Success bool   `json:"success"`
}

var collections *mongo.Collection

func init() {
	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}
	var err error
	dbPort := os.Getenv("db_port")
	dbHost := os.Getenv("db_host")
	db := os.Getenv("db")
	dbName := os.Getenv("db_name")
	collName := os.Getenv("collection_name")
	uri := db + "://" + dbHost + ":" + dbPort
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	client.Connect(nil)
	collections = client.Database(dbName).Collection(collName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("MongoDB Connected")
}

func MongoGoCollection() *mongo.Collection {
	return collections
}

func FindOne(filter, projection bson.M) *mongo.SingleResult {
	return collections.FindOne(nil, filter, options.FindOne().SetProjection(projection))
}
func Find(filter, projection bson.M) (*mongo.Cursor, error) {
	return collections.Find(nil, filter, options.Find().SetProjection(projection))
}
func InsertMany(document []interface{}) (*mongo.InsertManyResult, error) {
	return collections.InsertMany(nil, document)
}
func InsertOne(document interface{}) (*mongo.InsertOneResult, error) {
	return collections.InsertOne(nil, document)
}
func UpdateOne(filter, update interface{}) (*mongo.UpdateResult, error) {
	return collections.UpdateOne(nil, filter, update)
}
