package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
	"github.com/labstack/gommon/log"
	"onair/src/config"
)

type DB struct {
	DynamoDB  *dynamo.DB
	CoreTable dynamo.Table
	SDK       *dynamodb.DynamoDB
}

var db *DB

func GetDB() *DB {
	if db != nil {
		return db
	}

	db = new(DB)
	credential := credentials.NewStaticCredentials(config.GetEnv("AWS_ACCESS_ID"), config.GetEnv("AWS_SECRET_KEY"), "")
	awsConfig := &aws.Config{Credentials: credential, Region: aws.String("ap-northeast-2")}

	awsSession, err := session.NewSession(awsConfig)

	if err != nil {
		log.Error(err)
		panic(err)
	}

	sdk := dynamodb.New(awsSession)

	db.DynamoDB = dynamo.New(awsSession)
	db.CoreTable = db.DynamoDB.Table("the-football")
	db.SDK = sdk
	return db
}