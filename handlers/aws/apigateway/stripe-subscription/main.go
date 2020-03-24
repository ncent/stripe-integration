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

const (
	planIdentifier = "plan_uuid"
	userIdentifier = "user_uuid"
)

type (
	card struct {
		Number   string `json:"number"`
		ExpMonth string `json:"exp_month"`
		ExpYear  string `json:"exp_year"`
		CVC      string `json:"cvc"`
	}

	userRequest struct {
		UserUUID  string `json:"user_uuid"`
		PlanUUID  string `json:"plan_uuid"`
		UserEmail string `json:"email"`
		Card      card   `json:"card"`
	}
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
		"handler":       "aws.apigateway.stripe_subscription",
	})
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	plan, ok := os.LookupEnv("DEFAULT_PLAN")

	if !ok {
		log.Fatal("Plan Key not found on env vars")
	}

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

	params := &stripe.PaymentMethodParams{
		Type: stripe.String("card"),
		Card: &stripe.PaymentMethodCardParams{
			Number:   stripe.String(userReq.Card.Number),
			ExpMonth: stripe.String(userReq.Card.ExpMonth),
			ExpYear:  stripe.String(userReq.Card.ExpYear),
			CVC:      stripe.String(userReq.Card.CVC),
		},
	}

	paymentMethod, err := stripeService.CreatePaymentMethod(params)

	if err != nil {
		log.Errorf("Failed to create payment methods %v", err)
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	customerParams := &stripe.CustomerParams{
		PaymentMethod: stripe.String(paymentMethod.ID),
		Email:         stripe.String(userReq.UserEmail),
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(paymentMethod.ID),
		},
		Params: stripe.Params{
			Metadata: map[string]string{
				userIdentifier: userReq.UserUUID,
			},
		},
	}

	customer, err := stripeService.CreateCustomer(customerParams)

	if err != nil {
		log.Errorf("Failed to create customer %v", err)
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	items := []*stripe.SubscriptionItemsParams{
		{
			Plan: stripe.String(plan),
			Params: stripe.Params{
				Metadata: map[string]string{
					planIdentifier: userReq.PlanUUID,
				},
			},
		},
	}

	subscriptionParams := &stripe.SubscriptionParams{
		Customer: stripe.String(customer.ID),
		Items:    items,
		Params: stripe.Params{
			Metadata: map[string]string{
				planIdentifier: userReq.PlanUUID,
			},
		},
	}

	subscription, err := stripeService.CreateSubscription(subscriptionParams)

	if err != nil {
		log.Errorf("Failed to create subscription %v", err)
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	log.Infof("Subscription with ID %s created", subscription.ID)

	return resp, nil
}

func main() {
	stripeService = stripeClient.NewStripeClient()
	lambda.Start(handler)
}
