package handler

import (
	db "GeoProfileService/db"
	model "GeoProfileService/models"
	util "GeoProfileService/utils"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

// GetPeople is an httpHandler for route POST /updatepost
// This is the api used for updating data in post
func GetPeople(w http.ResponseWriter, r *http.Request) {
	data := model.SearchStruct{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	data.Limit = 100
	if data.Page == 0 {
		data.Page = 1
	}
	data.Skip = (data.Page - 1) * 100
	people, err2 := db.GetPeople(data)
	if err2 != "" {
		json.NewEncoder(w).Encode(bson.M{"success": false, "error": err2})
		return
	}
	counts, err := db.GetNumOfUser()
	pages := int64(counts/100) + 1
	if len(people) == 0 {
		var arr = make([]string, 0)
		json.NewEncoder(w).Encode(bson.M{"success": true, "error": err2, "page": data.Page, "total_pages": pages, "data": arr})

	} else {
		json.NewEncoder(w).Encode(bson.M{"success": true, "error": err2, "page": data.Page, "total_pages": pages, "data": people})

	}
}

func SearchPeople(w http.ResponseWriter, r *http.Request) {
	data := model.SearchStruct{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	people, err2 := db.SearchPeople(data)
	if err2 != "" {
		json.NewEncoder(w).Encode(bson.M{"success": false, "error": err2})
		return
	}
	json.NewEncoder(w).Encode(bson.M{"success": true, "error": err2, "data": people})
}

func UpdateLocation(w http.ResponseWriter, r *http.Request) {
	data := model.UserData{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(db.UpdateLocation(data))

}
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	data := model.UserData{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(db.Update(data))

}
func UploadEmergency(w http.ResponseWriter, r *http.Request) {
	data := model.UserData{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
	}
	if data.Ename != "" && data.Ephonenumber != "" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(db.Update(data))
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"success": false, "Error": "Missing Parameters"})
	}

}
func UploadImage(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	data := model.UserData{}
	data.UID = r.FormValue("UID")
	data.ImageType = r.FormValue("Type")
	if userImage, userImageHeader, err := r.FormFile("UserImage"); err == nil {
		defer userImage.Close()
		con, err := os.OpenFile(userImageHeader.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
			return
		}
		defer con.Close()
		io.Copy(con, userImage)
		log.Printf(userImageHeader.Filename)
		version, err2 := util.PostFile(userImageHeader.Filename, userImageHeader.Filename, data.UID)
		if err2 != nil {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(bson.M{"success": false, "Error": "Not Able to Upload"})
		}
		data.Userimageversion = version
		data.Userimage = userImageHeader.Filename
		err = os.Remove(userImageHeader.Filename)
	} else if err != nil {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"success": false, "Error": "image not present please varify the fieldname is UserImage"})
	}
	if data.Userimage == "" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bson.M{"success": false, "Error": "image not Uploaded"})
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(db.UploadImage(data))
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	data := model.SearchStruct{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	people, err2 := db.GetProfile(data)
	if err2 != "" {
		json.NewEncoder(w).Encode(bson.M{"success": false, "error": err2})
		return
	}
	json.NewEncoder(w).Encode(bson.M{"success": true, "error": err2, "data": people})
}

func GetUserData(w http.ResponseWriter, r *http.Request) {
	data := model.SearchStruct{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	people, err2 := db.GetProfile(data)
	if err2 != "" {
		json.NewEncoder(w).Encode(bson.M{"success": false, "error": err2})
		return
	}

	json.NewEncoder(w).Encode(bson.M{"success": true, "error": err2, "data": people})
}
