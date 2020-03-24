package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	stripeClient "gitlab.com/ncent/stripe-integration/services/stripe/client"
)

type (
	userRequest struct {
		StripeSubscriptionID string `json:"stripe_subscription_id"`
	}
)

func init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	log = logger.WithFields(logrus.Fields{
		"micro-service": "stripe-integration",
		"handler":       "aws.apigateway.stripe_cancel",
	})
}

var (
	stripeService stripeClient.IStripeClient
	log           *logrus.Entry
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{
		Headers:    make(map[string]string),
		StatusCode: http.StatusOK,
	}

	resp.Headers["Access-Control-Allow-Origin"] = "*"

	var userReq userRequest
	err := json.Unmarshal([]byte(request.Body), &userReq)

	if err != nil {
		log.Errorf("Failed to unmarshal request %v", err)
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	// Do not cancel immediatelly only on the end of the next billing cycle
	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}

	subscription, err := stripeService.CancelSubscription(userReq.StripeSubscriptionID, params)

	if err != nil {
		log.Errorf("Failed to cancel subscription %v", err)
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	log.Infof("Subscription with ID %s cancelled", subscription.ID)

	return resp, nil
}

func main() {
	stripeService = stripeClient.NewStripeClient()
	lambda.Start(handler)
}
