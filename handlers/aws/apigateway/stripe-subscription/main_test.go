package main

import (
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stripe/stripe-go"
	"gitlab.com/ncent/stripe-integration/services/stripe/client"
	stripeClient "gitlab.com/ncent/stripe-integration/services/stripe/client"
)

var _ = Describe("Stripe Subscription Suite", func() {
	var mockStripeService *client.MockStripeClient
	JustBeforeEach(func() {
		os.Setenv("DEFAULT_PLAN", "plan_abc123")
		os.Setenv("STAGE", "test")
		mockStripeService = stripeClient.NewMockStripeClient()
		stripeService = mockStripeService
	})

	Context("Given new stripe subscription", func() {
		It("It should create the subscription correctly", func() {
			mockStripeService.PaymentMethodResult = &stripe.PaymentMethod{
				ID: "card_abc123",
			}

			mockStripeService.MockCreatePaymentMethodExpectation.PaymentMethodParams = &stripe.PaymentMethodParams{
				Type: stripe.String("card"),
				Card: &stripe.PaymentMethodCardParams{
					Number:   stripe.String("4000000760000002"),
					ExpMonth: stripe.String("11"),
					ExpYear:  stripe.String("27"),
					CVC:      stripe.String("123"),
				},
			}

			mockStripeService.CustomerResult = &stripe.Customer{
				ID: "customer_abc123",
			}

			mockStripeService.MockCreateCustomerExpectation.CustomerMethodParams = &stripe.CustomerParams{
				PaymentMethod: stripe.String("card_abc123"),
				Email:         stripe.String("foobarman@gmail.com"),
				InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
					DefaultPaymentMethod: stripe.String("card_abc123"),
				},
				Params: stripe.Params{
					Metadata: map[string]string{
						userIdentifier: "123",
					},
				},
			}

			mockStripeService.SubscriptionResult = &stripe.Subscription{
				ID: "subscription_abc123",
			}

			items := []*stripe.SubscriptionItemsParams{
				{
					Plan: stripe.String("plan_abc123"),
					Params: stripe.Params{
						Metadata: map[string]string{
							planIdentifier: "321",
						},
					},
				},
			}

			mockStripeService.MockCreateSubscriptionExpectation.SubscriptionParams = &stripe.SubscriptionParams{
				Customer: stripe.String("customer_abc123"),
				Items:    items,
				Params: stripe.Params{
					Metadata: map[string]string{
						planIdentifier: "321",
					},
				},
			}

			_, err := handler(events.APIGatewayProxyRequest{
				Body: `{
					"user_uuid": "123",
					"plan_uuid": "321",
					"email": "foobarman@gmail.com",
					"card": {
						"number": "4000000760000002",
						"exp_month": "11",
						"exp_year": "27",
						"cvc":
						"123"
					}
				}`,
			})

			Expect(err).To(BeNil())
		})
	})
})

func TestHandlerStripeSubscription(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stripe Subscription Suite")
}
