package main

import (
	"flag"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

func getAwsSession(endpoint string) *session.Session {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Endpoint: aws.String(endpoint)},
	}))

	return sess
}

func pollSqs(sqsEndpoint string, sqsUrl string, chn chan<- *sqs.Message) {

	sqsSvc := sqs.New(getAwsSession(sqsEndpoint))
	for {
		output, err := sqsSvc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(sqsUrl),
			MaxNumberOfMessages: aws.Int64(5),
			WaitTimeSeconds:     aws.Int64(1),
		})

		if err != nil {
			log.Fatalf("failed to fetch sqs message %v", err)
		}

		for _, message := range output.Messages {
			chn <- message
			_, _ = sqsSvc.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      aws.String(sqsUrl),
				ReceiptHandle: message.ReceiptHandle,
			})
		}

	}

}

func handleMessage(message *sqs.Message) {
	log.Printf("Received this message: %s", message.GoString())
}

func init() {
	flag.String("endpoint", "", "aws endpoint to use")
	flag.String("url", "", "sqs url to poll")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	viper.SetConfigName("config")

	viper.AutomaticEnv()
}
func main() {

	awsEndpoint := viper.GetString("endpoint")
	sqsUrl := viper.GetString("url")

	chnMessages := make(chan *sqs.Message, 5)
	go pollSqs(awsEndpoint, sqsUrl, chnMessages)

	log.Printf("Listening on stack queue: %s", sqsUrl)

	for message := range chnMessages {
		handleMessage(message)
	}
}
