package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"os"
	"time"
)

func Handler(_ context.Context) (string, error) {
	queueUrl := os.Getenv("QUEUE_URL")
	queueEndpoint := os.Getenv("QUEUE_ENDPOINT")
	localstackHostname := os.Getenv("LOCALSTACK_HOSTNAME")
	localstackPort := os.Getenv("EDGE_PORT")
	if len(localstackHostname) > 0 && len(localstackPort) > 0 {
		queueEndpoint = fmt.Sprintf("http://%s:%s", localstackHostname, localstackPort)
	}
	fmt.Println("QUEUE:", queueUrl)
	fmt.Println("ENDPOINT:", queueEndpoint)
	fmt.Println("The time is", time.Now().Format(time.RFC1123))
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("us-east-1"),
		Credentials:      credentials.NewStaticCredentials("test", "test", ""),
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String(queueEndpoint),
	})

	// Create an SQS client
	svc := sqs.New(sess)

	// Send a message to the queue
	_, err = svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String("Hello, World!"),
		QueueUrl:    aws.String(queueUrl),
	})
	if err != nil {
		return "", fmt.Errorf("error sending message to SQS: %v", err)
	} else {
		fmt.Println("Successfully sent message to SQS")
	}

	return "Message sent successfully!", nil
}

func main() {
	// Start the Lambda handler
	lambda.Start(Handler)
}
