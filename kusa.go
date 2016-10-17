package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool `json:"short"`
}

type Attachment struct {
	Color string `json:"color"`
	AutherName string `json:"author_name"`
	AutherLink string `json:"author_link"`
	AutherIcon string `json:"author_icon"`
	Title string `json:"title"`
	TitleLink string `json:"title_link"`
	Text string `json:"text"`
	Fields []Field `json:"fields"`
}

type RequestValue struct {
	Channel string `json:"channel"`
	Username string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
	Text string `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

func main() {
	users := strings.Split(getenvOrExit("GITHUB_USERS"), ":")

	for _, user := range users {
		attachments := []Attachment{}

		url := githubUrl(user)
		doc, err := goquery.NewDocument(url)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		count, streak, maxStreak, totalCnt, weekCnt := contributionsOn(doc)
		color := ""
		icon := ""
		text := fmt.Sprintf("%s's contributions", user)
		fields := []Field{}

		if count > 0 {
			color = "good"
			icon = getenvOrDefault("ICON_EMOJI", ":seedling:")
			buf := make([]byte, 0)
			for i := 0; i < count; i++ {
				buf = append(buf, ":cherry_blossom:"...)
			}
			fields = append(fields, Field{Title: "Contributions", Value: string(buf), Short: false})
		} else {
			color = "danger"
			icon = getenvOrDefault("ICON_EMOJI_NO_CONTRIBUTION", ":japanese_goblin:")
			value := getenvOrDefault("KUSA_MSG_NO_CONTRIBUTIONS", ":warning: There are no contributions today ! :warning:")
			fields = append(fields, Field{Title: "Contributions", Value: value, Short: false})
		}

		fields = append(fields, Field{Title: "Current streak", Value: fmt.Sprintf("%d days", streak), Short: true})
		fields = append(fields, Field{Title: "Longest streak", Value: fmt.Sprintf("%d days", maxStreak), Short: true})
		fields = append(fields, Field{Title: "In the last year", Value: fmt.Sprintf("%d contributions", totalCnt), Short: true})
		fields = append(fields, Field{Title: "In the last week", Value: fmt.Sprintf("%d contributions", weekCnt), Short: true})

		attachments = append(attachments, Attachment{Color: color, Text: text, Fields: fields})

		httpPost(createValue(icon, attachments))
	}
}

func getenvOrDefault(key string, defvalue string) string {
	value, exists := os.LookupEnv(key)
	if (exists) {
		return value
	} else {
		return defvalue
	}
}

func getenvOrExit(key string) string {
	value, exists := os.LookupEnv(key)
	if (!exists) {
		fmt.Println(fmt.Sprintf("%s is required.", key))
		os.Exit(1)
	}
	return value
}

func githubUrl(user string) string {
	return fmt.Sprintf("https://github.com/%s", user)
}

func contributionsOn(doc *goquery.Document) (int, int, int, int, int) {
	cnt := 0
	streak := 0
	maxStreak := 0
	week := []int{}
	totalCnt := 0
	doc.Find(".js-calendar-graph-svg .day").Each(func(i int, s *goquery.Selection) {
		str, _ := s.Attr("data-count")
		cnt, _ = strconv.Atoi(str)
		totalCnt += cnt
		week = append(week, cnt)
		if len(week) > 7 {
			week = week[1:]
		}
		if cnt > 0 {
			streak += 1
		} else {
			streak = 0
		}
		if streak > maxStreak {
			maxStreak = streak
		}
	})
	weekCnt := 0
	for _, c := range week {
		weekCnt += c
	}
	return cnt, streak, maxStreak, totalCnt, weekCnt
}

func httpPost(req RequestValue) {
	value, err := json.Marshal(req)

	resp, err := http.Post(getenvOrExit("SLACK_WEBHOOK_URL"), "application/json", bytes.NewBuffer(value))

	if err != nil {
		fmt.Println(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	defer resp.Body.Close()
}

func createValue(icon string, attachments []Attachment) RequestValue {
	channel := getenvOrExit("SLACK_CHANNEL")
	username := getenvOrDefault("SLACK_USERNAME", "kusabot")
	return RequestValue{Channel: channel, Username: username, IconEmoji: icon, Attachments: attachments}
}
