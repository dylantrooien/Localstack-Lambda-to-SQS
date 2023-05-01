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
)

func Handler(_ context.Context) (string, error) {
	//queueUrl := "http://localhost:4566/000000000000/integration-marketo-batch-queue-local"
	queueUrl := os.Getenv("QUEUE_URL")
	fmt.Println("queueUrl:", queueUrl)
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String("us-east-1"),
			Credentials: credentials.NewStaticCredentials("test", "test", "test"),
			Endpoint:    aws.String(queueUrl),
		},
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

//package main
//
//import (
//"context"
//"encoding/json"
//"fmt"
//
//"github.com/aws/aws-lambda-go/events"
//"github.com/aws/aws-lambda-go/lambda"
//"github.com/aws/aws-sdk-go/aws"
//"github.com/aws/aws-sdk-go/aws/session"
//"github.com/aws/aws-sdk-go/service/sqs"
//)
//
//type Message struct {
//	Greeting string `json:"greeting"`
//	Name     string `json:"name"`
//}
//
//func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
//	// Create an SQS service client
//	svc := sqs.New(session.Must(session.NewSession()))
//
//	// Create a message
//	message := Message{
//		Greeting: "Hello",
//		Name:     "World",
//	}
//	body, err := json.Marshal(message)
//	if err != nil {
//		return events.APIGatewayProxyResponse{}, fmt.Errorf("failed to marshal message: %v", err)
//	} else {
//		fmt.Println("Message marshalled")
//	}
//
//	// Send the message to the SQS queue
//	queueURL := "http://localhost:4566/000000000000/integration-marketo-batch-queue-local"
//	_, err = svc.SendMessage(&sqs.SendMessageInput{
//		MessageBody: aws.String(string(body)),
//		QueueUrl:    aws.String(queueURL),
//	})
//	if err != nil {
//		return events.APIGatewayProxyResponse{}, fmt.Errorf("failed to send message: %v", err)
//	} else {
//		fmt.Println("Message has been sent")
//	}
//
//	return events.APIGatewayProxyResponse{
//		StatusCode: 200,
//		Body:       "Message sent",
//	}, nil
//}
//
//func main() {
//	lambda.Start(handler)
//}

//package main
//
//import (
//	"database/sql"
//	"encoding/json"
//	"fmt"
//	"github.com/aws/aws-lambda-go/events"
//	"github.com/aws/aws-lambda-go/lambda"
//	"os"
//
//	"github.com/aws/aws-sdk-go/aws"
//	"github.com/aws/aws-sdk-go/aws/session"
//	"github.com/aws/aws-sdk-go/service/sqs"
//	_ "github.com/lib/pq"
//)
//
//const (
//	host     = "database"
//	port     = 5432
//	user     = "swoogo"
//	password = "swoogo"
//	dbname   = "integrations"
//	queueURL = "http://localhost:4566/000000000000/integration-marketo-batch-queue-local"
//)
//
//func main() {
//	// Set the Postgres database connection parameters
//	err := os.Setenv("HOST", host)
//	if err != nil {
//		return
//	}
//	err = os.Setenv("PORT", fmt.Sprintf("%d", port))
//	if err != nil {
//		return
//	}
//	err = os.Setenv("DATABASE", dbname)
//	if err != nil {
//		return
//	}
//	err = os.Setenv("USER", user)
//	if err != nil {
//		return
//	}
//	err = os.Setenv("PASSWORD", password)
//	if err != nil {
//		return
//	}
//
//	// Open a connection to the database
//	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
//	if err != nil {
//		fmt.Println("Error opening database connection:", err)
//		return
//	} else {
//		fmt.Println("Successfully connected to database")
//	}
//	defer func(db *sql.DB) {
//		err := db.Close()
//		if err != nil {
//			fmt.Println("Error closing database connection:", err)
//		} else {
//			fmt.Println("Successfully closed database connection")
//		}
//	}(db)
//
//	// Create an AWS session
//	sess := session.Must(session.NewSessionWithOptions(session.Options{
//		Config: aws.Config{
//			Region: aws.String("us-west-2"),
//		},
//	}))
//
//	// Create an SQS service client
//	svc := sqs.New(sess)
//
//	// Query the database for messages in the marketo_batch table
//	rows, err := db.Query("SELECT id, program_id, \"messageIds\", batch, created FROM marketo_batch WHERE status IS NULL")
//	if err != nil {
//		fmt.Println("Error querying database:", err)
//		return
//	} else {
//		fmt.Println("Successfully queried database")
//	}
//	defer func(rows *sql.Rows) {
//		err := rows.Close()
//		if err != nil {
//			fmt.Println("Error closing database rows:", err)
//		} else {
//			fmt.Println("Successfully closed database rows")
//		}
//	}(rows)
//
//	// Loop through each row in the result set
//	for rows.Next() {
//		var id int
//		var programId string
//		var messageIds []byte
//		var batch []byte
//		var created string
//
//		// Scan the row into variables
//		err := rows.Scan(&id, &programId, &messageIds, &batch, &created)
//		if err != nil {
//			fmt.Println("Error scanning marketo_batch row:", err)
//			continue
//		} else {
//			fmt.Printf("Scanned marketo_batch row with id %d\n", id)
//		}
//
//		type Message struct {
//			ID      string `json:"id"`
//			Content string `json:"content"`
//		}
//
//		var messages []Message
//		err = json.Unmarshal(batch, &messages)
//		if err != nil {
//			fmt.Println("Error unmarshalling batch:", err)
//			continue
//		} else {
//			fmt.Printf("Unmarshalled %d messages from batch\n", len(messages))
//		}
//
//		var payloads []*sqs.SendMessageBatchRequestEntry
//
//		for _, msg := range messages {
//			payload := &sqs.SendMessageBatchRequestEntry{
//				Id:          aws.String(msg.ID),
//				MessageBody: aws.String(msg.Content),
//			}
//			payloads = append(payloads, payload)
//		}
//
//
//
//
//
//		func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
//			for _, message := range sqsEvent.Records {
//			fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
//		}
//
//			return nil
//		}
//
//		func main() {
//			lambda.Start(handler)
//		}
//
//
//
//
//
//		// Send the message batch to the SQS queue
//		result, err := svc.SendMessageBatch(&sqs.SendMessageBatchInput{
//			Entries:  payloads,
//			QueueUrl: aws.String(queueURL),
//		})
//
//		if err != nil {
//			fmt.Println("Error sending message batch:", err)
//			continue
//		} else {
//			fmt.Printf("Sent %d messages to queue\n", len(payloads))
//		}
//
//		// Update the marketo_batch status to 'in progress'
//		_, err = db.Exec("UPDATE marketo_batch SET status='queued' WHERE id=$1", id)
//		if err != nil {
//			fmt.Println("Error updating marketo_batch status:", err)
//			continue
//		} else {
//			fmt.Printf("Updated marketo_batch status to 'in progress' for id %d\n", id)
//		}
//
//		fmt.Printf("Sent %d messages to queue. Result: %v\n", len(payloads), result)
//	}
//}
