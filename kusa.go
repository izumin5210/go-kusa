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
	"time"

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

		count := contributionsOn(doc, today())
		color := ""
		icon := ""
		text := ""

		if count > 0 {
			color = "good"
			icon = getenvOrDefault("ICON_EMOJI", ":seedling:")
			buf := make([]byte, 0)
			for i := 0; i < count; i++ {
				buf = append(buf, ":cherry_blossom:"...)
			}
			text = string(buf)
		} else {
			color = "danger"
			icon = getenvOrDefault("ICON_EMOJI_NO_CONTRIBUTION", ":japanese_goblin:")
			text = getenvOrDefault("KUSA_MSG_NO_CONTRIBUTIONS", ":warning: There are no contributions today ! :warning:")
		}

		attachments = append(attachments, Attachment{Color: color, Title: user, TitleLink: url, Text: text})

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

func today() string {
	return time.Now().Format("2006-01-02")
}

func contributionsOn(doc *goquery.Document, date string) int {
	str, _ := doc.Find(query(date)).Attr("data-count")
	cnt, _ := strconv.Atoi(str)
	return cnt
}

func query(date string) string {
	return fmt.Sprintf(".js-calendar-graph-svg .day[data-date='%s']", date)
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
