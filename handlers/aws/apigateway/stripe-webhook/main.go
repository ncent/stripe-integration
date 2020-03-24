package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/webhook"
	"gitlab.com/ncent/stripe-integration/services/aws/eventbridge/client"
)

type Request struct {
	Type string `json:"type"`
}

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
		"handler":       "aws.apigateway.stripe_webhook",
	})
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{
		Headers:    make(map[string]string),
		StatusCode: http.StatusOK,
	}

	log.Infof("Request body arrived (%s)", request.Body)

	resp.Headers["Access-Control-Allow-Origin"] = "*"

	stage, ok := os.LookupEnv("STAGE")

	if !ok {
		log.Errorf("Stage not set")
		resp.StatusCode = http.StatusInternalServerError
		return resp, errors.New("Stage not set")
	}

	eventBus, ok := os.LookupEnv("EVENT_BUS_NAME")

	if !ok {
		log.Errorf("Event but not set, using default bus")
		eventBus = "default"
	}

	webhookSecret, ok := os.LookupEnv("WEBHOOK_SECRET")

	if !ok {
		log.Errorf("Endpoint stripe secret not set")
		resp.StatusCode = http.StatusInternalServerError
		return resp, errors.New("Stripe webhook secret not set")
	}

	stripeSignature := request.Headers["Stripe-Signature"]
	_, err := webhook.ConstructEvent([]byte(request.Body), stripeSignature, webhookSecret)

	if err != nil {
		log.Errorf("Error verifying webhook signature: %v", err)
		resp.StatusCode = http.StatusForbidden
		return resp, err
	}

	var req Request
	err = json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		log.Errorf("Failed to unmarshal request body %v", err)
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	err = eventBridgeService.PutEvent(
		eventBus,
		stage+".webhook.stripe",
		req.Type,
		request.Body,
	)

	if err != nil {
		log.Errorf("Failed to put event %v", err)
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	resp.Body = request.Body
	return resp, nil
}

func main() {
	eventBridgeService = client.NewEventBridgeService()
	lambda.Start(handler)
}
