package client

import (
	"errors"
	"fmt"
)

type MockEventBridgePutEventExpectation struct {
	EventBusName string
	Source       string
	DetailType   string
	JsonDetail   string
}

type MockEventBridgeService struct {
	MockEventBridgePutEventExpectation
	ErrorReturned error
}

func NewMockEventBridgeService() *MockEventBridgeService {
	return &MockEventBridgeService{}
}

func (m *MockEventBridgeService) PutEvent(eventBusName, source, detailType, jsonDetail string) error {
	if m.MockEventBridgePutEventExpectation.EventBusName != eventBusName {
		message := fmt.Sprintf("eventBusName expectation %s is different from %s", eventBusName, m.MockEventBridgePutEventExpectation.EventBusName)
		return errors.New(message)
	}

	if m.MockEventBridgePutEventExpectation.Source != source {
		message := fmt.Sprintf("source expectation %s is different from %s", source, m.MockEventBridgePutEventExpectation.Source)
		return errors.New(message)
	}

	if m.MockEventBridgePutEventExpectation.DetailType != detailType {
		message := fmt.Sprintf("detailType expectation %s is different from %s", detailType, m.MockEventBridgePutEventExpectation.DetailType)
		return errors.New(message)
	}

	if m.MockEventBridgePutEventExpectation.JsonDetail != jsonDetail {
		message := fmt.Sprintf("jsonDetail expectation %s is different from %s", jsonDetail, m.MockEventBridgePutEventExpectation.JsonDetail)
		return errors.New(message)
	}
	return m.ErrorReturned
}
