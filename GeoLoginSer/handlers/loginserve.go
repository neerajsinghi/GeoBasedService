package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	db "GeoLoginSer/db"
	model "GeoLoginSer/models"

	"go.mongodb.org/mongo-driver/bson"
)

// GetPeople is an httpHandler for route GET /people
func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)

	res := db.Register(user)
	if err != nil {
		res.Error = err.Error()
	}
	json.NewEncoder(w).Encode(res)
}

// GetPerson is an httpHandler for route GET /people/{id}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}
	res, result := db.Login(user)
	if res.Error == "" {
		json.NewEncoder(w).Encode(bson.M{"success": true, "data": result})
		return
	}
	res.Success = false
	json.NewEncoder(w).Encode(res)
}

func CheckUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}
	res := db.CheckUser(user)

	json.NewEncoder(w).Encode(res)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	tokenString := r.Header.Get("Authorization")
	token := strings.Split(tokenString, " ")
	json.NewEncoder(w).Encode(db.Logout(token[1]))
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(db.ForgotPassword(user))

}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(db.ChangePassword(user))

}
