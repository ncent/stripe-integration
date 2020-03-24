package services

import "github.com/aws/aws-sdk-go/service/cloudwatchevents"

type ICloudWatchEvents interface {
	PutEvents(putEventsInput *cloudwatchevents.PutEventsInput) (*cloudwatchevents.PutEventsOutput, error)
}
