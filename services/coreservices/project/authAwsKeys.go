package project

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)
var AWS_REGION = "us-west-2"

func CheckAWSKeys(accessKey, secretKey string) (bool, error) {
	awsSession := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Region:      aws.String(AWS_REGION),
	}))

	svc := sts.New(awsSession)

	_, err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		fmt.Println("error in authenticating aws access keys",err)
		return false, err
	}

	return true, nil
}
