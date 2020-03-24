package client

import (
	"errors"
	"reflect"

	"github.com/stripe/stripe-go"
)

type MockCreatePaymentMethodExpectation struct {
	PaymentMethodParams *stripe.PaymentMethodParams
}

type MockCreateCustomerExpectation struct {
	CustomerMethodParams *stripe.CustomerParams
}

type MockCreateSubscriptionExpectation struct {
	SubscriptionParams *stripe.SubscriptionParams
}

type MockStripeClient struct {
	MockCreatePaymentMethodExpectation
	MockCreateCustomerExpectation
	MockCreateSubscriptionExpectation
	PaymentMethodResult   *stripe.PaymentMethod
	CustomerResult        *stripe.Customer
	SubscriptionResult    *stripe.Subscription
	CheckoutSessionResult *stripe.CheckoutSession
	ErrorResult           error
}

func NewMockStripeClient() *MockStripeClient {
	return &MockStripeClient{}
}

func (m MockStripeClient) CreatePaymentMethod(paymentMethodParams *stripe.PaymentMethodParams) (*stripe.PaymentMethod, error) {
	if !reflect.DeepEqual(paymentMethodParams, m.MockCreatePaymentMethodExpectation.PaymentMethodParams) {
		message := "payment method expectation are different from params"
		return nil, errors.New(message)
	}

	return m.PaymentMethodResult, m.ErrorResult
}

func (m MockStripeClient) CreateCustomer(customerParams *stripe.CustomerParams) (*stripe.Customer, error) {
	if !reflect.DeepEqual(customerParams, m.MockCreateCustomerExpectation.CustomerMethodParams) {
		message := "customer expectation are different from params"
		return nil, errors.New(message)
	}

	return m.CustomerResult, m.ErrorResult
}

func (m MockStripeClient) CreateSubscription(subscriptionParams *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	if !reflect.DeepEqual(subscriptionParams, m.MockCreateSubscriptionExpectation.SubscriptionParams) {
		message := "subscription expectation are different from params"
		return nil, errors.New(message)
	}

	return m.SubscriptionResult, m.ErrorResult
}

func (m MockStripeClient) CreateCheckoutSession(checkoutSessionParams *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error) {
	return m.CheckoutSessionResult, m.ErrorResult
}

func (m MockStripeClient) GetCheckoutSession(sessionId string) (*stripe.CheckoutSession, error) {
	return m.CheckoutSessionResult, m.ErrorResult
}

func (m MockStripeClient) CancelSubscription(subscriptionId string, subscriptionParams *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	return nil, nil
}
