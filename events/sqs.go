package events

import (
	"fmt"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var AWS_SESSION *awsSession.Session = nil
var WAREHOUSE_QUEUE_NAME string = "warehouse-data-transfer"
var WAREHOUSE_QUEUE_URL *string


func getSQSSession() (*awsSession.Session, error) {
	if AWS_SESSION != nil {
		return AWS_SESSION, nil
	}
	// create session from default env variables
	session, err := awsSession.NewSession()
	AWS_SESSION = awsSession.Must(session, err)

	return AWS_SESSION, nil
}


func getQueueURL(svc *sqs.SQS, queue *string) (*string, error) {
    result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
        QueueName: queue,
    })
    if err != nil {
        return nil, err
    }

    return result.QueueUrl, nil
}


func SetAWSCredentials(awsAccessKey, awsSecretKey, awsRegion string) error {
	session, err := awsSession.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKey,awsSecretKey, ""),
	})
	AWS_SESSION = awsSession.Must(session, err)

	return err
}


func SendCostEvent(costEvent CostEvent) {
	go func() {
		session, err := getSQSSession()
		if err != nil {
			fmt.Println("SendCostEvent: ", err)
			return
		}

		svc := sqs.New(session)
		if WAREHOUSE_QUEUE_URL == nil {
			WAREHOUSE_QUEUE_URL, err = getQueueURL(svc, &WAREHOUSE_QUEUE_NAME)
		}
		if err != nil {
			fmt.Println("SendCostEvent: ", err)
			return
		}

		body, jsonErr := json.Marshal(costEvent)
		if jsonErr != nil {
			fmt.Println("SendCostEvent: ", err)
			return
		}

		_, err = svc.SendMessage(&sqs.SendMessageInput{
			MessageAttributes: map[string]*sqs.MessageAttributeValue{
				"EventType": {
					DataType:    aws.String("String"),
					StringValue: aws.String(string(WAREHOUSE_COST_TRACKER)),
				},
			},
			MessageBody: aws.String(string(body)),
			QueueUrl:    WAREHOUSE_QUEUE_URL,
		})
		fmt.Println("SendCostEvent: ", err)
	}()
}
