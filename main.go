package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"streakmaintain/internal/api"
	"streakmaintain/internal/webhook"
	"time"
)

type Config struct {
	RoomCode   string   `json:"roomCode"`
	QuestionNo int      `json:"questionNo"`
	TaskID     string   `json:"taskID"`
	Answer     string   `json:"answer"`
	Sessions   []string `json:"sessions"`
	Webhook    string   `json:"webhook"`
	TimeZone   string   `json:"timeZone"`
}

func main() {
	/* Inbuilt Timer */
	inbuiltTimer := flag.Bool("disabletimer", false, "Disable In-Built Timer")
	flag.Parse()

	var configuration Config
	config, err := os.ReadFile("./config/settings.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(config, &configuration); err != nil {
		log.Fatal(err)
	}
	webhook.WebHookURL = configuration.Webhook

	fmt.Printf("[!] No Of Sessions Found: %d\n", len(configuration.Sessions))

	fmt.Printf("-------------\nRoomCode: %s\nTaskId: %s\nQuestionNo: %d\nAnswer: %s\nWebhook: %s.....\nTimeZone: %s\n-------------\n", configuration.RoomCode, configuration.TaskID, configuration.QuestionNo, configuration.Answer, webhook.WebHookURL[:50], configuration.TimeZone)
	if *inbuiltTimer {
		for _, s := range configuration.Sessions {
			client := api.Client{
				SessionId:  s,
				RoomCode:   configuration.RoomCode,
				TaskID:     configuration.TaskID,
				QuestionNo: configuration.QuestionNo,
				Answer:     configuration.Answer,
				Http:       &http.Client{},
			}
			if err := client.Streak(); err != nil {
				if err := webhook.SendError(client.UserName, err); err != nil {
					fmt.Println("error ocurred sending webhook message!")
				}
			} else {
				if err := webhook.SendInfo("Username: " + client.UserName + "| StreakCount: " + strconv.Itoa(client.StreakCount+1)); err != nil {
					fmt.Println("error occurred sending webhook message!")
				}
			}
		}
		return
	}

	/* Define Your Timezone */
	location, err := time.LoadLocation(configuration.TimeZone)
	if err != nil {
		fmt.Println("Error loading time zone:", err)
		return
	}

	fmt.Println("[!] Using Inbuilt Timer!")
	fmt.Println("[!] Checking all the session cookies")
	for _, s := range configuration.Sessions {
		client := api.Client{
			SessionId: s,
			Http:      &http.Client{},
		}
		if err := client.ValidateCookie(); err != nil {
			log.Fatal("Expired/Invalid Cookie Provided: ", err)
		}
	}
	fmt.Println("[!] All valid cookies found!")

	for {
		timeNow := time.Now().In(location)
		durationUntilNextMidnight := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day()+1, 0, 0, 0, 0, location).Sub(timeNow)
		fmt.Printf("Time remaining until next trigger: %02d Hour %02d Minute %02d Seconds\n", int(durationUntilNextMidnight.Hours()), int(durationUntilNextMidnight.Minutes())%60, int(durationUntilNextMidnight.Seconds())%60)
		time.Sleep(durationUntilNextMidnight)
		for _, s := range configuration.Sessions {
			client := api.Client{
				SessionId:  s,
				RoomCode:   configuration.RoomCode,
				TaskID:     configuration.TaskID,
				QuestionNo: configuration.QuestionNo,
				Answer:     configuration.Answer,
				Http:       &http.Client{},
			}
			if err := client.Streak(); err != nil {
				if err := webhook.SendError(client.UserName, err); err != nil {
					fmt.Println("error ocurred sending webhook message!")
				}
			} else {
				if err := webhook.SendInfo("Username: " + client.UserName + "| StreakCount: " + strconv.Itoa(client.StreakCount+1)); err != nil {
					fmt.Println("error occurred sending webhook message!")
				}
			}
		}
	}
}
