package db

import (
	model "GeoProfileService/models"
	"context"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

func GetPeople(user model.SearchStruct) (result []model.UserData, erro string) {
	var long, lat float64
	var err error
	if long, err = strconv.ParseFloat(user.Long, 64); err != nil {
		log.Println(err.Error()) // 3.1415927410125732
	}
	if lat, err = strconv.ParseFloat(user.Lat, 64); err != nil {
		log.Println(err.Error()) // 3.1415927410125732
	}
	id, err := primitive.ObjectIDFromHex(user.UID)
	//add timecheck when required
	//timeCheck := time.Now().Unix() - 30*60
	//, {"location_updated", bson.D{{"$gte": timeCheck}}
	pipeline := bson.A{
		bson.M{
			"$geoNear": bson.D{
				{"near", bson.M{"type": "Point", "coordinates": bson.A{long, lat}}},
				{"key", "location"},
				{"distanceField", "dist"},
				{"spherical", true},
				{"distanceMultiplier", 0.001},
				{"maxDistance", user.Distance * 1000},
			},
		},
		bson.M{"$match": bson.D{{"_id", bson.D{{"$ne", id}}}}},
		bson.M{"$skip": user.Skip},

		bson.M{"$limit": user.Limit},
	}

	cursor, err := model.Aggregate(pipeline)
	if err != nil {
		erro = err.Error()
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var people model.UserData
		err = cursor.Decode(&people)
		if err != nil {
			erro = err.Error()
			return
		}
		people.Password = ""
		result = append(result, people)
	}
	return
}
func SearchPeople(user model.SearchStruct) (result []model.UserData, erro string) {
	project := bson.M{}

	cursor, err := model.Find(bson.M{"username": bson.M{"$regex": user.Key, "$options": "i"}}, project, bson.M{}, user.Limit, user.Skip)
	if err != nil {
		erro = err.Error()
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var people model.UserData
		err = cursor.Decode(&people)
		if err != nil {
			erro = err.Error()
			return
		}
		result = append(result, people)
	}
	return
}

func UpdateLocation(user model.UserData) (res model.ResponseResult) {
	id, err := primitive.ObjectIDFromHex(user.UID)
	filter := bson.M{"_id": id}

	if err != nil {
		res.Error = err.Error()
		return
	}
	var long, lat float64

	if long, err = strconv.ParseFloat(user.Long, 64); err != nil {
		log.Println(err.Error()) // 3.1415927410125732
	}
	if lat, err = strconv.ParseFloat(user.Lat, 64); err != nil {
		log.Println(err.Error()) // 3.1415927410125732
	}
	if long != 0 && lat != 0 {
		user.Location.Coordinates = append(user.Location.Coordinates, long, lat)

		resul, err := model.UpdateOne(filter, bson.M{"$set": bson.M{"location": bson.M{"type": "Point", "coordinates": user.Location.Coordinates}, "location_updated": time.Now().Unix()}})
		if err != nil {
			res.Error = err.Error()
			return
		}

		if resul.ModifiedCount == 0 {
			res.Result = "Nothing Changed"
			return
		}
	}

	res.Success = true
	res.Result = "Data Updated successfully!!"
	return
}
func UploadImage(user model.UserData) (res model.ResponseResult) {
	id, _ := primitive.ObjectIDFromHex(user.UID)
	userimageurl := "?user=" + user.UID + "&ver=" + user.Userimageversion + "&name=" + user.Userimage
	var imageType = ""
	if user.ImageType == "profile" {
		imageType = "userimageurl"
	} else {
		imageType = "adhaarurl"
	}
	resul, err := model.UpdateOne(bson.M{"_id": id}, bson.M{"$set": bson.M{imageType: userimageurl}})
	if err != nil {
		res.Error = err.Error()
		return
	}
	if resul.ModifiedCount == 0 {
		res.Result = "Nothing Changed"
		return
	}
	res.Success = true
	res.Result = "Image Uploaded successfully!!"
	return
}

func GetProfile(user model.SearchStruct) (result map[string]interface{}, erro string) {
	id, _ := primitive.ObjectIDFromHex(user.Userid)

	result = make(map[string]interface{})
	project := bson.M{}

	err := model.FindOne(bson.M{"_id": id}, project).Decode(&result)
	if err != nil {
		erro = err.Error()
		return
	}
	for key, val := range result {
		if val == "" {
			delete(result, key)
		}
	}
	if result["emergency_contacts"] == nil {
		var arr = make([]model.EmergencyContact, 0)
		result["emergency_contacts"] = arr
	}
	if result["devids"] == nil {
		var arr = make([]string, 0)
		result["devids"] = arr
	}
	delete(result, "password")
	return
}
func GetDevid(user model.SearchStruct) (result []string, erro string) {
	for _, uid := range user.Userids {
		res := make(map[string][]string)
		id, _ := primitive.ObjectIDFromHex(uid)
		err := model.FindOne(bson.M{"_id": id}, bson.M{"_id": 0, "devids": 1}).Decode(&res)
		if err != nil {
			erro = err.Error()
			return
		}

		result = append(result, res["devids"]...)

	}

	return
}

func Update(user model.UserData) (res model.ResponseResult) {
	id, _ := primitive.ObjectIDFromHex(user.UID)
	filter := bson.M{"_id": id}

	set := bson.M{}
	if user.Ename != "" && user.Ephonenumber != "" {
		emergencySet := bson.M{"emergency_contacts": bson.M{"name": user.Ename, "phone_number": user.Ephonenumber}}
		set["$addToSet"] = emergencySet
	}
	if user.Username != "" {
		set["$set"] = bson.M{"username": user.Username}
	}
	if user.Email != "" {
		set["$set"] = bson.M{"email": user.Email}
	}
	resul, err := model.UpdateOne(filter, set)
	if err != nil {
		res.Error = err.Error()
		return
	}
	if resul.ModifiedCount == 0 {
		res.Result = "Nothing Changed"
		return
	}
	res.Success = true
	res.Result = "Uploaded successfully!!"
	return
}

func GetNumOfUser() (int64, error) {
	return model.Count()
}
