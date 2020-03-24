package main

import (
	"os"
	"testing"

	"github.com/stripe/stripe-go"

	"github.com/aws/aws-lambda-go/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/ncent/stripe-integration/services/stripe/client"
	stripeClient "gitlab.com/ncent/stripe-integration/services/stripe/client"
)

const authorizationHeader = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFyYmVyZG1uQGdtYWlsLmNvbSIsInB1YmxpY0tleSI6IjA0OGE3NWMwYjk0NWUzZjA1NTQ2MmMxY2RiM2I2ODJlZTA3YjRkZjVkMjE4MzlkODg3Nzg3MjlkNDg3MGIwZTU2ZWNiMWY3N2QxOGQyODFjYTA0MDVhYWEzMzVlOTMzMDg4ZTExYTM0Mjg4ZmMyMDJmYTc0ZGIxOTdmYWU2YTk2NjIifQ.vQ4ra3BPQV4OenhiaIqFMqpfrnywHlwudPW7eg69J-E"
const sessionID = "session_001"

var _ = Describe("Create Checkout Session Suite", func() {
	var mockStripeClient *client.MockStripeClient
	JustBeforeEach(func() {
		os.Setenv("SUCCESS_URL", "http://example.org/success")
		os.Setenv("CANCEL_URL", "http://example.org/cancel")
		mockStripeClient = stripeClient.NewMockStripeClient()
		mockStripeClient.CheckoutSessionResult = &stripe.CheckoutSession{
			ID: sessionID,
		}
		stripeService = mockStripeClient
	})

	Context("Given the request on webhook arrived", func() {
		It("Then it will create a message on event bridge", func() {
			result, err := handler(events.APIGatewayProxyRequest{
				Headers: map[string]string{
					"Authorization": authorizationHeader,
				},
				Body: `{}`,
			})

			Expect(err).To(BeNil())
			Expect(result.Body).To(Equal("{\"sessionId\":\"session_001\"}"))
		})
	})
})

func TestHandlerCheckoutSessionSuit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create Checkout Session Suite")
}
