package san2s3

import (
	"fmt"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/joho/godotenv"
)

func TestAddFileToS3(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Errorf("Error loading .env file")
	}
	ses, err := GetAWSSession()
	if err != nil {
		log.Fatal(err)
	} else {
		err, resp := AddFileToS3(ses, "test.mp4")
		if err != nil {
			t.Errorf("Error while uploading file: %v", err)
		} else {
			fmt.Printf("response %s", awsutil.StringValue(resp))
		}
	}
}
