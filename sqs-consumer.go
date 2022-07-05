package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/nickzelei/go-sqs-consumer/config"
)

func main() {
	consume()
}

func consume() {
	conf, err := config.ReadConfig()

	if err != nil {
		panic(fmt.Errorf("fatal error with config file: %s", err))
	}

	queueURL := conf.GetString("SQS_QUEUE_URL")
	poolSize := conf.GetInt("MAX_WORKERS")

	fmt.Printf("Queue url: %s\n", queueURL)

	sqsAPI := getSqs()

	fmt.Printf("Starting %d worker(s)...\n", poolSize)
	var wg sync.WaitGroup
	for w := 0; w < poolSize; w++ {
		wg.Add(w + 1)
		go func(workerId int) {
			defer wg.Done()
			worker(sqsAPI, workerId, queueURL)
		}(w + 1)
	}
	wg.Wait()
}

func worker(sqsAPI sqsiface.SQSAPI, id int, queueURL string) {
	for {
		messages, err := retrieveMessages(sqsAPI, queueURL, 10)
		if err != nil {
			log.Fatal(err)
			continue
		}

		fmt.Printf("Worker %d found %d messages\n", id, len(messages))

		var wg sync.WaitGroup
		for _, message := range messages {
			wg.Add(1)
			go handleMessage(sqsAPI, queueURL, message, &wg)
		}
		wg.Wait()
	}
}

func handleMessage(sqsAPI sqsiface.SQSAPI, queueURL string, message *sqs.Message, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println(message)

	deleteMessage(sqsAPI, queueURL, *message.ReceiptHandle)
}

func retrieveMessages(sqsAPI sqsiface.SQSAPI, queueURL string, maxMessages int64) ([]*sqs.Message, error) {
	log.Println("Receiving Messages from SQS...")
	output, err := sqsAPI.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueURL),
		MaxNumberOfMessages: aws.Int64(maxMessages),
		WaitTimeSeconds:     aws.Int64(20),
	})

	if err != nil {
		return nil, err
	}
	return output.Messages, nil
}

func deleteMessage(sqsAPI sqsiface.SQSAPI, queueURL, receiptHandle string) error {
	_, err := sqsAPI.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})

	return err
}

func getSqs() sqsiface.SQSAPI {
	return sqs.New(session.Must(session.NewSession()))
}
