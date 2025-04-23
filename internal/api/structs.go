package api

import "net/http"

type Client struct {
	csrf        string
	_csrf       string
	SessionId   string
	RoomCode    string
	TaskID      string
	QuestionNo  int
	Answer      string
	Http        *http.Client
	UserName    string
	StreakCount int
}

type RoomCompleteReq struct {
	Answer     string `json:"answer"`
	QuestionNo int    `json:"questionNo"`
	RoomCode   string `json:"roomCode"`
	TaskID     string `json:"taskId"`
}

type RoomCompleteResp struct {
	Status string `json:"status"`
	Data   struct {
		IsCorrect             bool `json:"isCorrect"`
		IsRoomCompleted       bool `json:"isRoomCompleted"`
		IsStreakFreezeAwarded bool `json:"isStreakFreezeAwarded"`
		IsStreakIncreased     bool `json:"isStreakIncreased"`
		CurrentStreak         int  `json:"currentStreak"`
	} `json:"data"`
}
