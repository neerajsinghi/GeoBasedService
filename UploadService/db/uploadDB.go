package db

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

// UploadToS3 will upload a single file to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func UploadToS3(fileDir, filename, userid string) (string, error) {
	e := godotenv.Load() //Load .env file
	if e != nil {
		log.Print(e)
	}
	aws_access_key_id := os.Getenv("access_key_id")
	aws_secret_access_key := os.Getenv("secret_access_key")
	token := ""
	creds := credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, token)
	_, err := creds.Get()
	if err != nil {
		// handle error
	}
	cfg := aws.NewConfig().WithRegion("ap-south-1").WithCredentials(creds)
	s := s3.New(session.New(), cfg)
	// Open the file for use
	file, err := os.Open(fileDir)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)
	result, err := s.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(strings.ToLower(userid))})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				log.Println(s3.ErrCodeBucketAlreadyExists, aerr.Error())
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				log.Println(s3.ErrCodeBucketAlreadyOwnedByYou, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}

	}

	log.Println(result)
	result2, err := s.PutBucketVersioning(&s3.PutBucketVersioningInput{Bucket: aws.String(strings.ToLower(userid)), VersioningConfiguration: &s3.VersioningConfiguration{MFADelete: aws.String("Disabled"), Status: aws.String("Enabled")}})
	if err != nil {

		log.Println(err.Error())

	}

	log.Println(result2)
	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	result3, err := s.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(strings.ToLower(userid)),
		Key:                  aws.String(filename),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		log.Println(err)
	}

	return result3.GoString(), err
}

func DownloadFromS3(fileName, userid, versionid string, w http.ResponseWriter) {
	e := godotenv.Load() //Load .env file
	if e != nil {
		log.Print(e)
	}
	aws_access_key_id := os.Getenv("access_key_id")
	aws_secret_access_key := os.Getenv("secret_access_key")
	token := ""

	file, err := os.Create(fileName + versionid)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
		return
	}

	defer file.Close()

	s, _ := session.NewSession(&aws.Config{
		Region:      aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, token),
	})
	downloader := s3manager.NewDownloader(s)

	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(userid),
			Key:    aws.String(fileName),
		})
	if err != nil {
		exitErrorf("Unable to download item %q, %v", fileName, err)
		return
	}
	FileHeader := make([]byte, 512)
	file.Read(FileHeader)
	FileContentType := http.DetectContentType(FileHeader)
	FileStat, _ := file.Stat()                         //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	//Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)
	//Send the file
	//We read 512 bytes from the file already, so we reset the offset back to 0
	file.Seek(0, 0)
	io.Copy(w, file) //'Copy' the file to the client
	err = os.Remove(fileName + versionid)
	if err != nil {
		log.Println("Failed Toi delete File")
	}

	log.Println("==> done deleting file")

}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
}
