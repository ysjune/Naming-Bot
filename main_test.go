package main

import (
	"testing"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/sashabaranov/go-openai"
)

func TestMain(t *testing.T) {
	// Call the main function and check that it runs without errors
	main()
}

func TestHandleEventMessage(t *testing.T) {
	// Create a mock EventsAPIEvent and Client
	event := slackevents.EventsAPIEvent{}
	client := slack.New("mock-token")

	// Call the HandleEventMessage function with the mock event and client
	err := HandleEventMessage(event, client)

	// Check that the function returns no error
	if err != nil {
		t.Errorf("HandleEventMessage returned an error: %v", err)
	}
}

func TestHandleAppMentionEventToBot(t *testing.T) {
	// Create a mock AppMentionEvent and Client
	event := slackevents.AppMentionEvent{}
	client := slack.New("mock-token")

	// Call the HandleAppMentionEventToBot function with the mock event and client
	err := HandleAppMentionEventToBot(&event, client)

	// Check that the function returns no error
	if err != nil {
		t.Errorf("HandleAppMentionEventToBot returned an error: %v", err)
	}
}

func TestGetGptConnect(t *testing.T) {
	// Call the getGptConnect function with a mock AI token and full text
	response, err := getGptConnect("mock-token", "mock-text")

	// Check that the function returns a valid response and no error
	if response == nil {
		t.Errorf("getGptConnect returned a nil response")
	}
	if err != nil {
		t.Errorf("getGptConnect returned an error: %v", err)
	}
}

