package main

import (
	"testing"

	"github.com/stripe/stripe-go"

	"github.com/aws/aws-lambda-go/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/ncent/stripe-integration/services/stripe/client"
	stripeClient "gitlab.com/ncent/stripe-integration/services/stripe/client"
)

var _ = Describe("Create Checkout Session Suite", func() {
	var mockStripeClient *client.MockStripeClient
	var queryStringParameters = map[string]string{
		"sessionId": "session_001",
	}
	JustBeforeEach(func() {
		mockStripeClient = stripeClient.NewMockStripeClient()
		mockStripeClient.CheckoutSessionResult = &stripe.CheckoutSession{
			ID: "session_001",
		}
		stripeService = mockStripeClient
	})

	Context("Given the request on webhook arrived", func() {
		It("Then it will create a message on event bridge", func() {
			result, err := handler(events.APIGatewayProxyRequest{
				QueryStringParameters: queryStringParameters,
				Body:                  `{}`,
			})

			Expect(err).To(BeNil())
			Expect(result.Body).To(Equal("{\"cancel_url\":\"\",\"client_reference_id\":\"\",\"customer\":null,\"customer_email\":\"\",\"deleted\":false,\"display_items\":null,\"id\":\"session_001\",\"livemode\":false,\"locale\":\"\",\"mode\":\"\",\"object\":\"\",\"payment_intent\":null,\"payment_method_types\":null,\"setup_intent\":null,\"subscription\":null,\"submit_type\":\"\",\"success_url\":\"\"}"))
		})
	})
})

func TestHandlerCheckoutSessionSuit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create Checkout Session Suite")
}
