package client

import (
	"os"

	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

var log *logrus.Entry

func init() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	log = logger.WithFields(logrus.Fields{
		"service": "service.stripe.client",
	})
}

type IStripeClient interface {
	CreatePaymentMethod(paymentMethodParams *stripe.PaymentMethodParams) (*stripe.PaymentMethod, error)
	CreateCustomer(customerParams *stripe.CustomerParams) (*stripe.Customer, error)
	CreateSubscription(subscriptionParams *stripe.SubscriptionParams) (*stripe.Subscription, error)
	CancelSubscription(subscriptionId string, subscriptionParams *stripe.SubscriptionParams) (*stripe.Subscription, error)
	CreateCheckoutSession(checkoutSessionParams *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error)
	GetCheckoutSession(sessionId string) (*stripe.CheckoutSession, error)
}

type StripeClient struct {
	stripeAPI *client.API
}

func NewStripeClient() *StripeClient {
	stripeKey, ok := os.LookupEnv("STRIPE_KEY")

	log.Infof("Stripe Key %s", stripeKey)

	if !ok {
		log.Fatal("Stripe Key not found on env vars")
	}

	stripeAPI := &client.API{}
	stripeAPI.Init(stripeKey, nil)
	return &StripeClient{stripeAPI: stripeAPI}
}

func (sc StripeClient) CreatePaymentMethod(paymentMethodParams *stripe.PaymentMethodParams) (*stripe.PaymentMethod, error) {
	log.Infof("Creating payment method within Stripe %s", *paymentMethodParams.Type)
	return sc.stripeAPI.PaymentMethods.New(paymentMethodParams)
}

func (sc StripeClient) CreateCustomer(customerParams *stripe.CustomerParams) (*stripe.Customer, error) {
	log.Infof("Creating customer within Stripe %s", *customerParams.Email)
	return sc.stripeAPI.Customers.New(customerParams)
}

func (sc StripeClient) CreateSubscription(subscriptionParams *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	log.Infof("Creating a subscription within Stripe %s", *subscriptionParams.Customer)
	return sc.stripeAPI.Subscriptions.New(subscriptionParams)
}

func (sc StripeClient) CreateCheckoutSession(checkoutSessionParams *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error) {
	json, _ := json.Marshal(checkoutSessionParams)

	log.Infof("Creating a checkout session within Stripe %s", string(json))
	return sc.stripeAPI.CheckoutSessions.New(checkoutSessionParams)
}

func (sc StripeClient) GetCheckoutSession(sessionId string) (*stripe.CheckoutSession, error) {
	log.Infof("Get a checkout session within session id %s", sessionId)
	return sc.stripeAPI.CheckoutSessions.Get(sessionId, nil)
}

func (sc StripeClient) CancelSubscription(subscriptionId string, subscriptionParams *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	log.Infof("Cancel a subscription within Stripe %s [%s]", *subscriptionParams.Customer, subscriptionId)
	return sc.stripeAPI.Subscriptions.Update(subscriptionId, subscriptionParams)
}
