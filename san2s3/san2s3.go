package san2s3

import (
	"bytes"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

//its s3 bucket name
var S3Bucket = "mangospring"

//its s3 region name
var S3Region = "us-east-1"

// Create a single AWS session (we can re use this if we're uploading many files)
func GetAWSSession() (*session.Session, error) {
	creds := credentials.NewEnvCredentials()
	s, err := session.NewSession(&aws.Config{Region: aws.String(S3Region)}, &aws.Config{Credentials: creds})
	return s, err
}

// AddFileToS3 will upload a single file to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func AddFileToS3(s *session.Session, fileDir string) (error, *s3.PutObjectOutput) {

	// Open the file for use
	file, err := os.Open(fileDir)
	if err != nil {
		return err, nil
	}
	defer file.Close()
	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)
	contentType := http.DetectContentType(buffer)
	perm := "private"
	if strings.Contains(contentType, "video") {
		perm = "public-read"
	}
	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	resp, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(S3Bucket),
		Key:                  aws.String(fileDir),
		ACL:                  aws.String(perm),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(contentType),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err, resp
}

func desanitizeString(str string) string {
	str = strings.Replace(str, "&amp;", "&", -1)
	str = strings.Replace(str, "&#39;", "\\'", -1)
	str = strings.Replace(str, "&#039;", "\\'", -1)
	str = strings.Replace(str, "&#34;", "\"", -1)
	str = strings.Replace(str, "&#47;", "\\/", -1)
	str = strings.Replace(str, "&#92;", "\\\\", -1)
	str = strings.Replace(str, "&lt;", "\\<", -1)
	str = strings.Replace(str, "&gt;", "\\>", -1)
	str = strings.Replace(str, "&rsquo;", "\\'", -1)
	str = strings.Replace(str, "&quot;", "\"", -1)
	return str
}
