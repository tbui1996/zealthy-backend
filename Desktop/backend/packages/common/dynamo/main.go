package dynamo

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Database interface {
	Get(item *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	Create(item interface{}) error
	Update(expression string, item interface{}, key interface{}) error
	Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
	Scan(tableName string) (*dynamodb.ScanOutput, error)
	Delete(item *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
}

type DynamoDatabase struct {
	TableName string

	client *dynamodb.DynamoDB
}

func NewDynamoDatabaseWithSession(tableName string, sess *session.Session) *DynamoDatabase {
	client := dynamodb.New(sess)

	return &DynamoDatabase{
		TableName: tableName,
		client:    client,
	}
}

func (d *DynamoDatabase) getClient() *dynamodb.DynamoDB {
	if d.client != nil {
		return d.client
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	d.client = dynamodb.New(sess)
	return d.client
}

func (d *DynamoDatabase) Create(item interface{}) error {
	svc := d.getClient()

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Printf("Error marshaling item. %+v", item)
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(d.TableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Printf("Error on putting item in table: (%s), item: error (%s)", d.TableName, err)
		return err
	}

	return nil
}

func (d *DynamoDatabase) Update(expression string, item, key interface{}) error {
	svc := d.getClient()

	expr, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Printf("Got error marshalling item: %s", err)
		return err
	}

	k, err := dynamodbattribute.MarshalMap(key)
	if err != nil {
		log.Printf("Got error marshalling key: %s", err)
		return err
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: expr,
		TableName:                 aws.String(d.TableName),
		Key:                       k,
		ReturnValues:              aws.String("UPDATED_NEW"),
		UpdateExpression:          aws.String(expression),
	}

	_, err = svc.UpdateItem(input)
	if err != nil {
		log.Printf("Got error calling UpdateItem: %s", err)
	}

	return nil
}

func (d *DynamoDatabase) Get(item *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	svc := d.getClient()
	return svc.GetItem(item)
}

func (d *DynamoDatabase) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	svc := d.getClient()
	return svc.Query(input)
}

func (d *DynamoDatabase) Scan(tableName string) (*dynamodb.ScanOutput, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	return svc.Scan(input)
}

func (d *DynamoDatabase) Delete(item *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	svc := d.getClient()
	return svc.DeleteItem(item)
}
