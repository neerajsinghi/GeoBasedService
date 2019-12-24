package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func PostFile(filedir, filename, userID string) (version string, err error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		log.Println("error writing to buffer")
		return "", err
	}

	// open file handle
	fh, err := os.Open(filedir)
	if err != nil {
		log.Println("error opening file")
		return "", err
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return "", err
	}
	fw, err := bodyWriter.CreateFormField("userID")
	if err != nil {
		log.Println("error opening file")
		return "", err
	}
	_, err = io.Copy(fw, strings.NewReader(userID))
	if err != nil {
		return "", err
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	targetURL := "http://localhost:4000/v1/upload"
	resp, err := http.Post(targetURL, contentType, bodyBuf)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	log.Println(resp.Status)
	log.Println(string(respBody))
	var response map[string]string
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		panic(err)
	}
	return response["versionid"], err
}
