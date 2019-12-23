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

//FeedStruct Data to be added to the database
type SearchStruct struct {
	Key      string   `json:"key,omitempty"`
	UID      string   `json:"uid,omitempty"`
	Userid   string   `json:"userid,omitempty"`
	Userids  []string `json:"userids,omitempty"`
	Limit    int64    `json:"limit,omitempty"`
	Skip     int64    `json:"skip,omitempty"`
	Distance float64  `json:"distance,omitempty"`
	Page     int64    `json:"page,omitempty"`
	Lat      string   `json:"lat,omitempty"`
	Long     string   `json:"long,omitempty"`
}

type UserData struct {
	ID                primitive.ObjectID `bson:"_id" json:"userid,omitempty"`
	UID               string             `json:"uid,omitempty"`
	ImageType         string             `json:"image_type,omitempty"`
	Username          string             `json:"username,omitempty"`
	AboutMe           string             `json:"about_me,omitempty"`
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
	DateOfBirth       string             `json:"date_of_birth,omitempty"`
	Userimage         string             `json:"userimage,omitempty"`
	Userimageversion  string             `json:"userimageversion,omitempty"`
	Userimageurl      string             `json:"userimageurl,omitempty"`
	Adhaarurl         string             `json:"adhaarurl,omitempty"`
	Dist              float64            `json:"dist"`
	LocationUpdated   int64              `bson:"location_updated" json:"location_updated,omitempty"`
	EmergencyContacts []EmergencyContact `bson:"emergency_contacts" json:"emergency_contacts,omitempty"`
	Ename             string             `json:"ename,omitempty"`
	Ephonenumber      string             `json:"ephonenumber,omitempty"`
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
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Result  string `json:"result,omitempty"`
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
}

//FindPost find list of documents from db
func Find(filter, projection, sort bson.M, limit, skip int64) (*mongo.Cursor, error) {
	return collections.Find(nil, filter, options.Find().SetProjection(projection), options.Find().SetSort(sort), options.Find().SetSkip(skip), options.Find().SetLimit(limit))
}

//FindPostOne find one  document from db
func FindOne(filter, projection bson.M) *mongo.SingleResult {
	return collections.FindOne(nil, filter, options.FindOne().SetProjection(projection))
}

//InsertMany Insert many documents to db
func InsertMany(document []interface{}) (*mongo.InsertManyResult, error) {
	return collections.InsertMany(nil, document)
}

//InsertOne Insert one document to db
func InsertOne(document bson.M) (*mongo.InsertOneResult, error) {
	return collections.InsertOne(nil, document)
}

//UpdateOne update one document
func UpdateOne(filter, set bson.M) (*mongo.UpdateResult, error) {
	return collections.UpdateOne(nil, filter, set)
}
func Aggregate(pipeline bson.A) (*mongo.Cursor, error) {
	return collections.Aggregate(nil, pipeline)
}

func Count() (int64, error) {
	return collections.CountDocuments(nil, bson.M{})
}
