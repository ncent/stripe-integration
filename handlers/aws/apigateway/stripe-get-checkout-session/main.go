package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"

	stripeClient "gitlab.com/ncent/stripe-integration/services/stripe/client"
)

var (
	stripeService stripeClient.IStripeClient
	log           *logrus.Entry
)

func init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	log = logger.WithFields(logrus.Fields{
		"micro-service": "stripe-integration",
		"handler":       "aws.apigateway.get_checkout_session",
	})
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	resp := events.APIGatewayProxyResponse{
		Headers:    make(map[string]string),
		StatusCode: http.StatusOK,
	}

	resp.Headers["Access-Control-Allow-Origin"] = "*"
	resp.Headers["Access-Control-Allow-Credentials"] = "true"

	sessionID, ok := request.QueryStringParameters["sessionId"]

	// Return Bad Request if session id is not passed
	if !ok || len(sessionID) == 0 {
		resp.StatusCode = http.StatusBadRequest
		return resp, nil
	}

	session, err := stripeService.GetCheckoutSession(sessionID)

	if err != nil {
		log.Errorf("Failed to get a  session %+v", err)
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	log.Infof("Session with ID %s created", session.ID)

	// Create response with session
	response, err := json.Marshal(session)
	if err != nil {
		return resp, err
	}

	resp.StatusCode = http.StatusOK
	resp.Body = string(response)

	return resp, nil
}

func main() {
	stripeService = stripeClient.NewStripeClient()
	lambda.Start(handler)
}
