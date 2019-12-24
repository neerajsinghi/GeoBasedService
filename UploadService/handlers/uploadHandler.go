package handler

import (
	"UploadService/db"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Param is used to identify a person
//
// swagger:parameters listPerson
type DownloadCreds struct {
	Filename  string
	UserID    string
	Versionid string
}

func Upload(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)
	UserID := r.FormValue("userID")
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	//log.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	version, err := db.UploadToS3(handler.Filename, handler.Filename, UserID)
	versionid := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(strings.Split(strings.Split(version, ",")[2], ":")[1], "\"", ""), "}", ""))
	err = os.Remove(handler.Filename)
	if err != nil {
		log.Println("Failed Toi delete File")
	}

	log.Println("==> done deleting file")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Content-Type", "application/json;")
	w.WriteHeader(http.StatusOK)
	data := make(map[string]interface{})
	if err != nil {
		data["success"] = "false"
		data["err"] = err.Error()
		json.NewEncoder(w).Encode(data)
	}
	data["success"] = "true"
	data["err"] = ""
	data["versionid"] = versionid
	json.NewEncoder(w).Encode(data)

}

func Download(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var dCreds DownloadCreds

	dCreds.Filename = r.URL.Query().Get("name")
	dCreds.UserID = r.URL.Query().Get("user")
	dCreds.Versionid = r.URL.Query().Get("ver")

	db.DownloadFromS3(dCreds.Filename, dCreds.UserID, dCreds.Versionid, w)

}
