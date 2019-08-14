package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"strconv"
	"time"
)

var DB *dynamodb.DynamoDB

func init() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	DB = dynamodb.New(sess)
}

func Delete(tableName, key, value string) (err error) {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			key: {
				S: aws.String(value),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err = DB.DeleteItem(input)

	return
}

func Put(tableName string, item interface{}, ttl time.Duration) (err error) {
	i, err := dynamodbattribute.MarshalMap(item)

	input := &dynamodb.PutItemInput{
		Item:      i,
		TableName: aws.String(tableName),
	}

	duration := int(ttl.Seconds())
	currentTime := int(time.Now().Unix())
	currentTime += duration
	formattedTtl := strconv.Itoa(currentTime)

	if ttl != 0 {
		input.Item["ttl"] = &dynamodb.AttributeValue{N: aws.String(formattedTtl)}
	}

	_, err = DB.PutItem(input)
	return
}

func GetItem(tableName string, keys map[string]interface{}, object interface{}) (err error) {

	builder := expression.NewBuilder()

	filter := expression.ConditionBuilder{}
	firstKey := true

	for key, value := range keys {
		if firstKey {
			filter = expression.Name(key).Equal(expression.Value(value))
			firstKey = false
		} else {
			filter = filter.And(expression.Name(key).Equal(expression.Value(value)))
		}
	}

	exp := builder.WithFilter(filter)

	expr, err := exp.Build()

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	out, err := DB.Scan(params)

	if len(out.Items) > 0 {
		dynamodbattribute.UnmarshalMap(out.Items[0], &object)
	}

	return
}

func GetItemList(tableName string, keys map[string]interface{}, object interface{}) (err error) {

	builder := expression.NewBuilder()

	filter := expression.ConditionBuilder{}
	firstKey := true

	for key, value := range keys {
		if firstKey {
			filter = expression.Name(key).Equal(expression.Value(value))
			firstKey = false
		} else {
			filter = filter.And(expression.Name(key).Equal(expression.Value(value)))
		}
	}

	exp := builder.WithFilter(filter)

	expr, err := exp.Build()

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	out, err := DB.Scan(params)

	err = dynamodbattribute.UnmarshalListOfMaps(out.Items, &object)

	return
}
