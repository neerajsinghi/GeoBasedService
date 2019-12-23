package db

import (
	model "GeoLoginSer/models"
	"fmt"
	"log"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Insert allows populating database
func Register(user model.User) (res model.ResponseResult) {

	var result model.User
	res.Success = false

	err := model.FindOne(bson.M{"$or": bson.A{bson.M{"phoneno": user.Email}, bson.M{"email": user.Email}}}, bson.M{}).Decode(&result)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

			if err != nil {
				res.Error = "Error While Hashing Password, Try Again"
				return
			}
			user.Password = string(hash)
			user.ID = primitive.NewObjectID()
			user.Location.Type = "Point"
			var long, lat float64
			if long, err = strconv.ParseFloat(user.Long, 64); err != nil {
				log.Println(err.Error()) // 3.1415927410125732
			}
			if lat, err = strconv.ParseFloat(user.Lat, 64); err != nil {
				log.Println(err.Error()) // 3.1415927410125732
			}
			user.Location.Coordinates = append(user.Location.Coordinates, long, lat)
			user.Long = ""
			user.Lat = ""
			if user.EmergencyContacts == nil {
				arr := make([]model.EmergencyContact, 0)
				user.EmergencyContacts = arr
			}
			if user.Devid != "" {
				user.Devids = append(user.Devids, user.Devid)
			} else {
				arr := make([]string, 0)
				user.Devids = arr
			}
			_, err = model.InsertOne(user)
			if err != nil {
				res.Error = "Error While Creating User, Try Again"
				return
			}
			res.Success = true

			res.Result = "Registration Successful"
			return
		}

		res.Error = err.Error()
		return
	}

	res.Result = "Username already Exists!!"
	return
}
func CheckUser(user model.User) (res model.ResponseResult) {
	var result model.User

	filter := bson.M{"$or": bson.A{bson.M{"phoneno": user.PhoneNo}, bson.M{"email": user.PhoneNo}}}
	err := model.FindOne(filter, bson.M{}).Decode(&result)
	res.Success = false

	if err != nil {
		res.Success = false
		res.Error = "Username doesn't exist"
		return
	} else {
		res.Error = "Username already exist"
		res.Success = true
	}
	return
}

func Login(user model.User) (res model.ResponseResult, result model.User) {

	filter := bson.M{"$or": bson.A{bson.M{"phoneno": user.PhoneNo}, bson.M{"email": user.PhoneNo}}}
	err := model.FindOne(filter, bson.M{}).Decode(&result)
	res.Success = false

	if err != nil {
		res.Error = "Invalid username"
		return
	}
	var long, lat float64

	if long, err = strconv.ParseFloat(user.Long, 64); err != nil {
		log.Println(err.Error()) // 3.1415927410125732
	}
	if lat, err = strconv.ParseFloat(user.Lat, 64); err != nil {
		log.Println(err.Error()) // 3.1415927410125732
	}
	if user.Devid != "" {
		set := bson.M{"$addToSet": bson.D{{"devids", user.Devid}}, "$set": bson.D{{"devid", user.Devid}}}
		_, err = model.UpdateOne(filter, set)
	}
	if long != 0 && lat != 0 {
		user.Location.Coordinates = append(user.Location.Coordinates, long, lat)

		_, err = model.UpdateOne(filter, bson.M{"$set": bson.M{"location": bson.M{"type": "Point", "coordinates": user.Location.Coordinates}}})
		if err != nil {
			log.Println(err.Error())
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))

	if err != nil {
		res.Error = "Invalid password"
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phoneno":   result.PhoneNo,
		"firstname": result.FirstName,
		"lastname":  result.LastName,
		"devid":     user.Devid,
	})
	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		res.Error = "Error while generating token,Try again"
		return
	}
	err = model.FindOne(filter, bson.M{}).Decode(&result)

	result.Token = tokenString
	res.Success = true
	result.Password = ""
	return
}

func ForgotPassword(user model.User) (res model.ResponseResult) {
	var result model.User
	res.Success = false

	filter := bson.M{"$or": bson.A{bson.M{"phoneno": user.Email}, bson.M{"email": user.Email}}}
	err := model.FindOne(filter, bson.M{}).Decode(&result)

	if err != nil {
		res.Success = false
		res.Error = "Invalid username"
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

	if err != nil {
		res.Error = "Error While Hashing Password, Try Again"
		return
	}
	user.Password = string(hash)
	_, err = model.UpdateOne(filter, bson.M{"$set": bson.M{"password": user.Password}})
	res.Success = true
	result.Password = ""
	res.Result = "Password changed successfully"
	return
}

func ChangePassword(user model.User) (res model.ResponseResult) {
	var result model.User
	res.Success = false

	id, _ := primitive.ObjectIDFromHex(user.UID)
	filter := bson.M{"_id": id}
	err := model.FindOne(filter, bson.M{}).Decode(&result)

	if err != nil {
		res.Error = "Invalid user"
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Oldpassword))
	if err != nil {
		res.Error = "Invalid password"
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)

	if err != nil {
		res.Error = "Error While Hashing Password, Try Again"
		return
	}
	user.Password = string(hash)
	user.ID = primitive.NewObjectID()
	_, err = model.UpdateOne(filter, bson.M{"$set": bson.M{"password": user.Password}})
	res.Success = true

	result.Password = ""
	res.Result = "Password changed successfully"
	return
}

func Logout(tokenString string) (res model.ResponseResult) {

	var result model.User
	res.Success = false

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result.PhoneNo = claims["phoneno"].(string)
		result.Devid = claims["devid"].(string)
	}
	resul, err := model.UpdateOne(bson.M{"phoneno": result.PhoneNo}, bson.M{"$pull": bson.M{"devids": result.Devid}})
	if err != nil {
		res.Error = err.Error()
		return
	}
	if resul.ModifiedCount == 0 {
		res.Result = "Nothing Changed"
		return
	}
	res.Success = true

	res.Result = "Logged out successfully!!"
	return
}
