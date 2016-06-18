package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type RequestValue struct {
	Channel string `json:"channel"`
	Username string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
	Text string `json:"text"`
}

func main() {
	user := getenvOrExit("GITHUB_USER")

	doc, err := goquery.NewDocument(githubUrl(user))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	httpPost(contributionsOn(doc, today()))
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

func httpPost(cnt int) {
	value, err := json.Marshal(createValues(cnt))

	resp, err := http.Post(getenvOrExit("SLACK_WEBHOOK_URL"), "application/json", bytes.NewBuffer(value))

	if err != nil {
		fmt.Println(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	defer resp.Body.Close()
}

func createValues(cnt int) RequestValue {
	channel := getenvOrExit("SLACK_CHANNEL")
	username := getenvOrDefault("SLACK_USERNAME", "kusabot")
	icon := ""
	text := ""
	if cnt > 0 {
		exists := false
		icon = getenvOrDefault("ICON_EMOJI", ":seedling:")
		text, exists = os.LookupEnv("KUSA_MSG_DEFAULT")
		if !exists {
			buf := make([]byte, 0)
			for i := 0; i < cnt; i++ {
				buf = append(buf, ":cherry_blossom:"...)
			}
			text = string(buf)
		}
	} else {
		icon = getenvOrDefault("ICON_EMOJI_NO_CONTRIBUTION", ":japanese_goblin:")
		text = getenvOrDefault("KUSA_MSG_NO_CONTRIBUTIONS", ":warning: There are no contributions today ! :warning:")
	}
	return RequestValue{Channel: channel, Username: username, IconEmoji: icon, Text: text}
}
