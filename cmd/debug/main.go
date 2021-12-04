package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	_ "github.com/lib/pq"
)

const (
	dbURL = "postgres://prixfixe_api:password@api-database.cluster-ctj4wxgujo7g.us-east-1.rds.amazonaws.com:5432/prixfixe"
	// queueURL = "https://sqs.us-east-1.amazonaws.com/966107642521/writes.fifo"
	queueArn = "arn:aws:sqs:us-east-1:966107642521:writes.fifo"
)

var (
	db    *sql.DB
	queue *sqs.SQS
)

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server started at port 80")

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	queue = sqs.New(sess)

	log.Fatal(http.ListenAndServe(":80", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handler invoked for request: %s %s", r.Method, r.URL.String())

	var version string
	row := db.QueryRow("SELECT VERSION()")
	if err := row.Scan(&version); err != nil {
		log.Printf("error querying database: %v", err)
		fmt.Fprintf(w, "Error: %s", err)
	}

	if _, err := queue.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Things": {
				DataType:    aws.String("String"),
				StringValue: aws.String("Stuff"),
			},
		},
		MessageBody: aws.String("just testin'"),
		QueueUrl:    aws.String(queueArn),
	}); err != nil {
		log.Printf("error writing message to queue: %v", err)
		fmt.Fprintf(w, "Error: %s", err)
	}

	fmt.Fprintf(w, "Hello, there from version %s!", version)
}
