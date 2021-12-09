package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	dbURL = "postgres://prixfixe_api:F9F7[#v=hSWL6f-KT#z6b2[2RIjWMlW_NhUxVBCu5GicS1_Rj]dkia_NXUq]KC=d@api-database.cluster-ctj4wxgujo7g.us-east-1.rds.amazonaws.com:5432/prixfixe"
	// queueURL = "https://sqs.us-east-1.amazonaws.com/966107642521/writes.fifo"
)

var (
	db *sql.DB
	// queue *sqs.SQS
)

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server started at port 80")

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	// sess := session.Must(session.NewSessionWithOptions(session.Options{
	// 	SharedConfigState: session.SharedConfigEnable,
	// }))
	// queue = sqs.New(sess)

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

	// if _, err := queue.SendMessage(&sqs.SendMessageInput{
	// 	MessageAttributes: map[string]*sqs.MessageAttributeValue{
	// 		"Things": {
	// 			DataType:    aws.String("String"),
	// 			StringValue: aws.String("Stuff"),
	// 		},
	// 	},
	// 	MessageDeduplicationId: aws.String(ksuid.New().String()),
	// 	MessageGroupId:         aws.String("writes"),
	// 	MessageBody:            aws.String("just testin'"),
	// 	QueueUrl:               aws.String(queueURL),
	// }); err != nil {
	// 	log.Printf("error writing message to queue: %v\n", err)
	// 	fmt.Fprintf(w, "Error: %s", err)
	// }

	fmt.Fprintf(w, "Hello, there from version %s!", version)
}
