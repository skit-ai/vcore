package events

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var WAREHOUSE_QUEUE_URL string = "https://sqs.ap-south-1.amazonaws.com/536612919621/warehouse-data-transfer"
var AWS_SESSION *session.Session = nil


func getSQSSession() *session.Session {
	if AWS_SESSION != nil {
		return AWS_SESSION
	}

	return session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}


func SendCostEvent(costEvent CostEvent) error {
	session := getSQSSession()
	svc := sqs.New(session)

	body, jsonErr := json.Marshal(costEvent)
	if jsonErr != nil {
		return jsonErr
	}

	_, err := svc.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"EventType": {
				DataType:    aws.String("String"),
				StringValue: aws.String(string(WAREHOUSE_COST_TRACKER)),
			},
		},
		MessageBody: aws.String(string(body)),
		QueueUrl:    &WAREHOUSE_QUEUE_URL,
	})

	return err
}
