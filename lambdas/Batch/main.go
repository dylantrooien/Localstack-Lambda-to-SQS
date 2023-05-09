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
	//queueUrl := "http://localhost:4566/000000000000/integration-marketo-batch-queue-local"
	queueUrl := os.Getenv("QUEUE_URL")
	fmt.Println("QUEUE:", queueUrl)
	// Print the current time
	fmt.Println("The time is", time.Now().Format(time.RFC1123))
	//sess, err := session.NewSessionWithOptions(session.Options{
	//	Config: aws.Config{
	//		Region:      aws.String("us-east-1"),
	//		Credentials: credentials.NewStaticCredentials("test", "test", "test"),
	//		Endpoint:    aws.String(queueUrl),
	//	},
	//})
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String("us-east-1"),
		Credentials:      credentials.NewStaticCredentials("test", "test", ""),
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String("http://localhost:4566"),
	})
	//sess, err := session.NewSession()

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
