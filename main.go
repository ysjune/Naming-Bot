package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"log"
	"os"
	"strings"
)

func main() {

	godotenv.Load("local.env")

	token := os.Getenv("SLACK_AUTH_TOKEN")
	//channelID := os.Getenv("SLACK_CHANNEL_ID")
	appToken := os.Getenv("SLACK_APP_TOKEN")

	client := slack.New(token, slack.OptionDebug(true), slack.OptionAppLevelToken(appToken))

	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(false),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	go func(ctx context.Context, client *slack.Client, socketClient *socketmode.Client) {
		for {
			select {
			case <-ctx.Done():
				log.Println("Shutting down socketmode listener")
				return
			case event := <-socketClient.Events:

				switch event.Type {
				case socketmode.EventTypeEventsAPI:
					eventsAPI, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						log.Printf("Could not type case the event to the EventsAPI: %v\n", event)
						continue
					}
					socketClient.Ack(*event.Request)
					err := HandleEventMessage(eventsAPI, client)
					if err != nil {
						log.Print(err)
					}
				}
			}
		}
	}(ctx, client, socketClient)

	socketClient.Run()
}

func HandleEventMessage(event slackevents.EventsAPIEvent, client *slack.Client) error {
	switch event.Type {

	case slackevents.CallbackEvent:

		innerEvent := event.InnerEvent

		switch evnt := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			err := HandleAppMentionEventToBot(evnt, client)
			if err != nil {
				return err
			}
		case *slackevents.MessageAction:

		}
	default:
		return errors.New("unsupported event type")
	}
	return nil
}

func HandleAppMentionEventToBot(event *slackevents.AppMentionEvent, client *slack.Client) error {
	aiToken := os.Getenv("OPEN_API_TOKEN")

	attachment := slack.Attachment{}

	text := strings.ToLower(event.Text)
	if strings.Contains(text, ".변수명") {
		removedMention := strings.ReplaceAll(text, "<@u052tfxnt0s>", "")
		removedCommand := strings.ReplaceAll(removedMention, ".변수명", "")
		fullText := fmt.Sprint("Suggest an English variable for '" + removedCommand + "' in camelCase")

		response, err := getGptConnect(aiToken, fullText)
		if err != nil {
			attachment.Text = fmt.Sprint("현재 GPT서버 상태가 :zany_face: 합니다. \n 나중에 이용해 주세요.")
			log.Printf("ChatCompletion error: %v \n", err)
		} else {
			attachment.Pretext = fmt.Sprint("*'" + removedCommand + "' 에 대한 추천 변수명입니다* :sunglasses:")
			content := response.Choices[0].Message.Content
			attachment.Text = fmt.Sprint(":camel: camelCase : *" + content + "* \n\n :pray: 단어가 제대로 나오지 않는다면 영어로 입력해보세요. \n ex) .변수명 cash지급 내역")
			attachment.Color = "#4af030"
		}
	} else if strings.Contains(text, ".마스터키") {
		removedMention := strings.ReplaceAll(text, "<@u052tfxnt0s>", "")
		removedCommand := strings.ReplaceAll(removedMention, ".마스터키", "")

		response, err := getGptConnect(aiToken, removedCommand)
		if err != nil {
			attachment.Text = fmt.Sprint("현재 GPT서버 상태가 :zany_face: 합니다. \n 나중에 이용해 주세요.")
			log.Printf("ChatCompletion error: %v \n", err)
		} else {
			attachment.Text = response.Choices[0].Message.Content
		}
	} else if strings.Contains(text, ".메소드명") {
		removedMention := strings.ReplaceAll(text, "<@u052tfxnt0s>", "")
		removedCommand := strings.ReplaceAll(removedMention, ".메소드명", "")
		fullText := fmt.Sprint("Recommend an English function name for '" + removedCommand + "' in camelCase. The function should consist of a verb + noun. for example, saveMoney")

		response, err := getGptConnect(aiToken, fullText)
		if err != nil {
			attachment.Text = fmt.Sprint("현재 GPT서버 상태가 :zany_face: 합니다. \n 나중에 이용해 주세요.")
			log.Printf("ChatCompletion error: %v \n", err)
		} else {
			attachment.Pretext = fmt.Sprint("*'" + removedCommand + "' 에 대한 추천 메소드명입니다* :partying_face:")
			content := response.Choices[0].Message.Content
			attachment.Text = fmt.Sprint(":camel: camelCase : *" + content + "* \n\n :pray: 단어가 제대로 나오지 않는다면 영어로 입력해보세요. \n ex) .메소드명 cash 지급 요청")
			attachment.Color = "#309df0"
		}
	} else {
		attachment.Text = "변수명, 메소드명 추천만 입력이 가능합니다 :face_exhaling:"
	}

	_, _, err := client.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		return fmt.Errorf("failed to post message : %w", err)
	}
	return nil
}

func getGptConnect(aiToken string, fullText string) (openai.ChatCompletionResponse, error) {
	aiClient := openai.NewClient(aiToken)
	response, err := aiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo0301,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fullText,
				},
			},
		})
	return response, err
}
