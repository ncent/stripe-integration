package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"

	"gitlab.com/ncent/stripe-integration/services"
	stripeClient "gitlab.com/ncent/stripe-integration/services/stripe/client"
)

const paymentMethod = "card"

var (
	stripeService stripeClient.IStripeClient
	log           *logrus.Entry
)

type userRequest struct {
	PlanUUID string `json:"plan_uuid"`
}

type response struct {
	SessionID string `json:"sessionId"`
}

func init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	log = logger.WithFields(logrus.Fields{
		"micro-service": "stripe-integration",
		"handler":       "aws.apigateway.checkout_session",
	})
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	successURL, ok := os.LookupEnv("SUCCESS_URL")

	if !ok {
		log.Fatal("Success URL not found on env vars")
	}

	cancelURL, ok := os.LookupEnv("CANCEL_URL")

	if !ok {
		log.Fatal("Cancel URL not found on env vars")
	}

	resp := events.APIGatewayProxyResponse{
		Headers:    make(map[string]string),
		StatusCode: http.StatusOK,
	}

	var urerReq userRequest
	err := json.Unmarshal([]byte(request.Body), &urerReq)

	if err != nil {
		log.Errorf("Failed to parse json request %v", err)
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	resp.Headers["Access-Control-Allow-Origin"] = "*"
	resp.Headers["Access-Control-Allow-Credentials"] = "true"

	token, err := services.GetTokenFromHeader(request)

	if err != nil {
		log.Errorf("Error on retrieve token header")
		resp.StatusCode = http.StatusBadRequest
		return resp, err
	}

	userUUID, err := services.GetJWTPublicKey(token)
	userEmail, err := services.GetJWTEmail(token)

	if err != nil {
		log.Errorf("Failed to get publickey %v", err)
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	log.Infof("Checkout session for user %s", userUUID)
	log.Infof("Checkout session for urerReq %+v", urerReq)

	checkoutSessionParams := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			paymentMethod,
		}),
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Items: []*stripe.CheckoutSessionSubscriptionDataItemsParams{
				&stripe.CheckoutSessionSubscriptionDataItemsParams{
					Plan: stripe.String(urerReq.PlanUUID),
				},
			},
		},
		SuccessURL:        stripe.String(successURL),
		CancelURL:         stripe.String(cancelURL),
		ClientReferenceID: stripe.String(userUUID),
		CustomerEmail:     stripe.String(userEmail),
	}

	session, err := stripeService.CreateCheckoutSession(checkoutSessionParams)

	if err != nil {
		log.Errorf("Failed to create subscription %+v", err)
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	log.Infof("Session with ID %s created", session.ID)

	// Create response with sessionID
	response, err := json.Marshal(response{SessionID: session.ID})
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
