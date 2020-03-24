package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
	"gitlab.com/ncent/stripe-integration/services/aws/eventbridge/client"
)

// Should match event detail
// Source: prod.webhook.stripe
// DetailType: checkout.session.completed

const (
	providerName = "stripe"
	detailType   = "subscription-created"
)

type (
	stripeEvent struct {
		Data struct {
			Object struct {
				ClientReferenceID string `json:"client_reference_id"`
				CustomerEmail     string `json:"customer_email"`
				DisplayItems      []struct {
					Plan struct {
						ID      string `json:"id"`
						Created int64  `json:"created"`
					} `json:"plan"`
				} `json:"display_items"`
			} `json:"object"`
		} `json:"data"`
	}

	detailEvent struct {
		PlanUUID  string `json:"plan_uuid"`
		UserUUID  string `json:"user_uuid"`
		StartDate string `json:"start"`
		EndDate   string `json:"end"`
		Provider  string `json:"provider"`
		Email     string `json:"email"`
	}

	eventRequest struct {
		Detail       string
		DetailType   string
		EventBusName string
		Source       string
	}
)

var (
	eventBridgeService client.IEventBridgeService
	log                *logrus.Entry
)

func init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	log = logger.WithFields(logrus.Fields{
		"micro-service": "stripe-integration",
		"handler":       "aws.eventbridge.create_subscription",
	})
}

func putEvent(eventRequest *eventRequest) error {
	err := eventBridgeService.PutEvent(
		eventRequest.EventBusName,
		eventRequest.Source,
		eventRequest.DetailType,
		eventRequest.Detail,
	)

	if err != nil {
		log.Printf("Failed to put event %v", err)
		return err
	}

	return nil
}

func handler(ctx context.Context, event events.CloudWatchEvent) error {
	stage, ok := os.LookupEnv("STAGE")

	if !ok {
		log.Fatal("Stage var not found on env vars")
	}

	eventBusName, ok := os.LookupEnv("EVENT_BUS_NAME")

	if !ok {
		log.Fatal("EventBusName not found on env vars")
	}

	var stripeEv stripeEvent
	err := json.Unmarshal([]byte(event.Detail), &stripeEv)

	if err != nil {
		log.Errorf("Failed to unmarshal strip event %v", err)
		return err
	}

	tm := time.Unix(stripeEv.Data.Object.DisplayItems[0].Plan.Created, 0)
	date := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		tm.UTC().Year(), tm.UTC().Month(), tm.UTC().Day(),
		tm.UTC().Hour(), tm.UTC().Minute(), tm.UTC().Second())

	detailEv := &detailEvent{
		PlanUUID:  stripeEv.Data.Object.DisplayItems[0].Plan.ID,
		UserUUID:  stripeEv.Data.Object.ClientReferenceID,
		StartDate: date,
		Provider:  providerName,
		Email:     stripeEv.Data.Object.CustomerEmail,
	}

	json, err := json.Marshal(detailEv)

	if err != nil {
		log.Infof("Failed to marshal subscription event %v", err)
		return err
	}

	log.Infof("Create subscription for (%s)", json)

	err = putEvent(&eventRequest{
		Source:       stage + ".monetization",
		DetailType:   detailType,
		EventBusName: eventBusName,
		Detail:       string(json),
	})

	if err != nil {
		log.Infof("Failed to put event on queue %v", err)
		return err
	}

	return nil
}

func main() {
	eventBridgeService = client.NewEventBridgeService()
	lambda.Start(handler)
}
