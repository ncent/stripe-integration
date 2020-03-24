package main

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/ncent/stripe-integration/services/aws/eventbridge/client"
)

var stripeWebhook = `{
	"data": {
		"object": {
			"client_reference_id": "048a75c0b945e3f055462c1cdb3b682ee07b4df5d21839d88778729d4870b0e56ecb1f77d18d281ca0405aaa335e933088e11a34288fc202fa74db197fae6a9662",
			"customer_email": "foobarexample@test.com",
			"display_items": [
				{
					"plan": {
						"id": "plan_GBAiQM9TVNFlSa",
						"created": 1573761180
					}
				}
			]
		}
	}
}`

var _ = Describe("Stripe Integration Suite", func() {
	var mockEventBridge *client.MockEventBridgeService

	JustBeforeEach(func() {
		os.Setenv("STAGE", "test")
		os.Setenv("STRIPE_KEY", "123")
		os.Setenv("EVENT_BUS_NAME", "testbus")
		mockEventBridge = client.NewMockEventBridgeService()
		eventBridgeService = mockEventBridge
	})

	Context("Given new stripe subscription event", func() {
		It("It should create the subscription event correctly", func() {
			mockEventBridge.MockEventBridgePutEventExpectation = client.MockEventBridgePutEventExpectation{
				EventBusName: "testbus",
				Source:       "test.monetization",
				DetailType:   detailType,
				JsonDetail:   `{"plan_uuid":"plan_GBAiQM9TVNFlSa","user_uuid":"048a75c0b945e3f055462c1cdb3b682ee07b4df5d21839d88778729d4870b0e56ecb1f77d18d281ca0405aaa335e933088e11a34288fc202fa74db197fae6a9662","start":"2019-11-14T19:53:00","end":"","provider":"stripe","email":"foobarexample@test.com"}`,
			}

			err := handler(context.Background(), events.CloudWatchEvent{
				Detail: []byte(stripeWebhook),
			})

			Expect(err).To(BeNil())
		})

		It("It should fail if the subscription event is missing detail", func() {
			mockEventBridge.MockEventBridgePutEventExpectation = client.MockEventBridgePutEventExpectation{
				EventBusName: "testbus",
				DetailType:   detailType,
				JsonDetail:   `{"plan_uuid":"plan_GBAiQM9TVNFlSa", "start":"2019-11-14T19:53:00","end":""}`,
			}

			err := handler(context.Background(), events.CloudWatchEvent{
				Detail: []byte(stripeWebhook),
			})

			Expect(err).To(Not(BeNil()))
		})

		It("It should fail if the subscription can't put the message on queue", func() {
			mockEventBridge.MockEventBridgePutEventExpectation = client.MockEventBridgePutEventExpectation{
				EventBusName: "testbus",
				Source:       "test.monetization",
				DetailType:   detailType,
				JsonDetail:   `{"plan_uuid":"plan_GBAiQM9TVNFlSa","user_uuid":"048a75c0b945e3f055462c1cdb3b682ee07b4df5d21839d88778729d4870b0e56ecb1f77d18d281ca0405aaa335e933088e11a34288fc202fa74db197fae6a9662","start":"2019-11-14T19:53:00","end_date":"","provider":"stripe"}`,
			}

			mockEventBridge.ErrorReturned = errors.New("Unexpected error")

			err := handler(context.Background(), events.CloudWatchEvent{
				Detail: []byte(stripeWebhook),
			})

			Expect(err).To(Not(BeNil()))
		})
	})
})

func TestHandlerStripeIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stripe Integration Suite")
}
