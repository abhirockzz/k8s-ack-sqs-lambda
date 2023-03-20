package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var table string
var client *dynamodb.Client

func init() {
	table = os.Getenv("TABLE_NAME")
	if table == "" {
		log.Fatal("missing environment variable TABLE_NAME")
	}
	cfg, _ := config.LoadDefaultConfig(context.Background())
	client = dynamodb.NewFromConfig(cfg)

}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		fmt.Println("received message from sqs", message.MessageId, "with body", message.Body)
		fmt.Println("storing message info to dynamodb table", table)

		item := make(map[string]types.AttributeValue)

		item["email"] = &types.AttributeValueMemberS{Value: message.Body}

		for attr, val := range message.Attributes {
			fmt.Println(attr, "=", val)
			item[attr] = &types.AttributeValueMemberS{Value: val}
		}

		for attr, val := range message.MessageAttributes {
			fmt.Println(attr, "=", val)

			switch val.DataType {
			case "String":
				item[attr] = &types.AttributeValueMemberS{Value: *val.StringValue}
			case "StringList":
				//item[attr] = &types.AttributeValueMemberN{Value: val.StringListValues[]}
			case "Binary":
				item[attr] = &types.AttributeValueMemberB{Value: val.BinaryValue}
			case "BinaryList":
				//item[attr] = &types.AttributeValueMemberSS{Value: val.StringListValues}
			}

		}

		_, err := client.PutItem(context.Background(), &dynamodb.PutItemInput{
			TableName: aws.String(table),
			Item:      item,
		})

		if err != nil {
			return err
		}

		fmt.Println("item added to table")
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
