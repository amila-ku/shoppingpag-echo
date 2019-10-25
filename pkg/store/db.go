package store

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/amila-ku/shoppingpal/pkg/entity"
)

var region = "eu-central-1"
var endpoint = "http://localhost:8000"

//DatabaseService defines a general database provider type
type DatabaseService struct {
	itemTable *dynamodb.DynamoDB
	tableName string
}

//NewTable Creates a pointer to new new DynamoDB type
func NewTable(tableName string) (*DatabaseService, error) {

	config := &aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
	}

	sess := session.Must(session.NewSession(config))

	svc := dynamodb.New(sess)

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("name"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Id"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("name"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		log.Println("Got error calling CreateTable:")
		log.Fatalln(err.Error())
		os.Exit(1)
	}

	return &DatabaseService{
		itemTable: svc,
		tableName: tableName,
	}, nil
}

//CreateItem is a receiver function for DatabaseService, this creates new item in database
func (i *DatabaseService) CreateItem(item entity.Item) error {
	//now := time.Now()
	//item.CreatedAt = now
	//item.UpdatedAt = now
	//item.Id = xid.New().String()
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Println("Got error marshalling new item:")
		log.Println(err.Error())
		os.Exit(1)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(i.tableName),
	}

	_, err = i.itemTable.PutItem(input)
	if err != nil {
		log.Println("Got error calling PutItem:")
		log.Println(err.Error())
		os.Exit(1)
	}

	return err
}
