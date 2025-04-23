package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) csrfToken() error {
	req, err := http.NewRequest("GET", "https://tryhackme.com/api/v2/auth/csrf", nil)
	if err != nil {
		return err
	}
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("Sec-Ch-Ua-Full-Version-List", `"Google Chrome";v="135.0.7049.52", "Not-A.Brand";v="8.0.0.0", "Chromium";v="135.0.7049.52"`)
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua", `"Google Chrome";v="135", "Not-A.Brand";v="8", "Chromium";v="135"`)
	req.Header.Set("Sec-Ch-Ua-Bitness", `"64"`)
	req.Header.Set("Sec-Ch-Ua-Model", `""`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Arch", `"x86"`)
	req.Header.Set("Sec-Ch-Ua-Full-Version", `"135.0.7049.52"`)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua-Platform-Version", `"6.13.5"`)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://tryhackme.com/room/"+c.RoomCode)
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("If-None-Match", `W/"4c-FeXdzyT4Y5RCimEZAYusKlhr3do"`)
	req.Header.Set("Priority", "u=1, i")
	resp, err := (*c.Http).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("error status code recieved while fetching csrf")
	}

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response struct {
		Status string `json:"status"`
		Data   struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(bodyText, &response); err != nil {
		return err
	}

	if response.Status != "success" {
		fmt.Println("[!] Unknown response from csrf endpoint: ", response)
	}
	c.csrf = response.Data.Token
	c._csrf = resp.Cookies()[1].Value
	fmt.Println("csrf token: ", c.csrf)
	fmt.Println("csrf cookie: ", c._csrf)
	return nil
}

func (c *Client) ValidateCookie() error {
	if err := c.csrfToken(); err != nil {
		return err
	}
	req, err := http.NewRequest("GET", "https://tryhackme.com/api/v2/users/self", nil)
	if err != nil {
		return err
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("csrf-token", c.csrf)
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://tryhackme.com/dashboard")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="135", "Not-A.Brand";v="8", "Chromium";v="135"`)
	req.Header.Set("sec-ch-ua-arch", `"x86"`)
	req.Header.Set("sec-ch-ua-bitness", `"64"`)
	req.Header.Set("sec-ch-ua-full-version", `"135.0.7049.52"`)
	req.Header.Set("sec-ch-ua-full-version-list", `"Google Chrome";v="135.0.7049.52", "Not-A.Brand";v="8.0.0.0", "Chromium";v="135.0.7049.52"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-model", `""`)
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)
	req.Header.Set("sec-ch-ua-platform-version", `"6.13.5"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", "connect.sid="+c.SessionId)
	resp, err := (*c.Http).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("error status code recieved while validating cookie")
	}

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var response struct {
		Status string `json:"status"`
		Data   struct {
			User struct {
				Username string `json:"username"`
				Streak   struct {
					Streak int `json:"streak"`
				} `json:"streak"`
			} `json:"user"`
		} `json:"data"`
	}
	if err := json.Unmarshal(bodyText, &response); err != nil {
		return err
	}

	if response.Status != "success" {
		fmt.Println("[!] Unknown response from validatecookie endpoint: ", response)
	}
	c.UserName = response.Data.User.Username
	c.StreakCount = response.Data.User.Streak.Streak
	fmt.Printf("User: %s | Streak: %d\n", c.UserName, c.StreakCount)
	return nil
}

func (c *Client) resetProgress() error {
	var data = strings.NewReader(`{"roomCode":"` + c.RoomCode + `"}`)
	req, err := http.NewRequest("POST", "https://tryhackme.com/api/v2/rooms/reset-progress", data)
	if err != nil {
		return err
	}

	req.Header.Set("Sec-Ch-Ua-Full-Version-List", `"Google Chrome";v="135.0.7049.52", "Not-A.Brand";v="8.0.0.0", "Chromium";v="135.0.7049.52"`)
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Ch-Ua", `"Google Chrome";v="135", "Not-A.Brand";v="8", "Chromium";v="135"`)
	req.Header.Set("Sec-Ch-Ua-Bitness", `"64"`)
	req.Header.Set("Csrf-Token", c.csrf)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Model", `""`)
	req.Header.Set("Sec-Ch-Ua-Arch", `"x86"`)
	req.Header.Set("Sec-Ch-Ua-Full-Version", `"135.0.7049.52"`)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Sec-Ch-Ua-Platform-Version", `"6.13.5"`)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", "https://tryhackme.com")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://tryhackme.com/room/"+c.RoomCode)
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Add("Cookie", "connect.sid="+c.SessionId+";_csrf="+c._csrf)
	resp, err := (*c.Http).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("error status code recieved while reseting progress")
	}

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(bodyText))
	var response struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	}

	if err := json.Unmarshal(bodyText, &response); err != nil {
		return err
	}

	if response.Status != "success" {
		fmt.Println("[!] Unknown response from resetprogress endpoint: ", response)
	}
	fmt.Println("[*] Reset Room Progress Done")
	return nil
}

func (c *Client) completeRoom() error {
	var rawData = &RoomCompleteReq{
		Answer:     c.Answer,
		QuestionNo: c.QuestionNo,
		RoomCode:   c.RoomCode,
		TaskID:     c.TaskID,
	}
	data, err := json.Marshal(rawData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://tryhackme.com/api/v2/rooms/answer", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("csrf-token", c.csrf)
	req.Header.Set("origin", "https://tryhackme.com")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://tryhackme.com/room/"+c.RoomCode)
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="135", "Not-A.Brand";v="8", "Chromium";v="135"`)
	req.Header.Set("sec-ch-ua-arch", `"x86"`)
	req.Header.Set("sec-ch-ua-bitness", `"64"`)
	req.Header.Set("sec-ch-ua-full-version", `"135.0.7049.52"`)
	req.Header.Set("sec-ch-ua-full-version-list", `"Google Chrome";v="135.0.7049.52", "Not-A.Brand";v="8.0.0.0", "Chromium";v="135.0.7049.52"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-model", `""`)
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)
	req.Header.Set("sec-ch-ua-platform-version", `"6.13.5"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", "connect.sid="+c.SessionId+";_csrf="+c._csrf)
	resp, err := (*c.Http).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response RoomCompleteResp
	if err := json.Unmarshal(bodyText, &response); err != nil {
		return err
	}

	if response.Status != "success" {
		fmt.Println("[!] Unknown response from complete endpoint: ", response)
	}
	if !response.Data.IsStreakIncreased {
		fmt.Println("[!] Streak Already Maintained!")
		return fmt.Errorf("streak already maintained")
	}
	fmt.Println("[!] Streak Maintained!")
	return nil
}

func (c *Client) Streak() error {
	if err := c.ValidateCookie(); err != nil {
		return err
	}
	if err := c.resetProgress(); err != nil {
		return err
	}
	if err := c.completeRoom(); err != nil {
		return err
	}
	return nil
}
