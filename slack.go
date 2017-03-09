package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Attachment struct {
	Color      string  `json:"color"`
	AutherName string  `json:"author_name"`
	AutherLink string  `json:"author_link"`
	AutherIcon string  `json:"author_icon"`
	Title      string  `json:"title"`
	TitleLink  string  `json:"title_link"`
	Text       string  `json:"text"`
	Fields     []Field `json:"fields"`
}

type RequestValue struct {
	Channel     string       `json:"channel"`
	Username    string       `json:"username"`
	IconEmoji   string       `json:"icon_emoji"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Slack struct {
	WebhookUrl, Channel, User           string
	IconEmoji, IconEmojiNoContributions string
	MsgNoContributions                  string
}

func (s *Slack) Post(status KusaStatus) error {
	color := ""
	icon := ""
	text := fmt.Sprintf("%s's contributions", status.User)
	fields := []Field{}

	if status.Count > 0 {
		color = "good"
		icon = s.IconEmoji
		buf := make([]byte, 0)
		for i := 0; i < status.Count; i++ {
			buf = append(buf, ":cherry_blossom:"...)
		}
		fields = append(fields, Field{Title: "Contributions", Value: string(buf), Short: false})
	} else {
		color = "danger"
		icon = s.IconEmojiNoContributions
		value := s.MsgNoContributions
		fields = append(fields, Field{Title: "Contributions", Value: value, Short: false})
	}

	fields = append(fields, Field{Title: "Current streak", Value: fmt.Sprintf("%d days", status.Streak), Short: true})
	fields = append(fields, Field{Title: "Longest streak", Value: fmt.Sprintf("%d days", status.MaxStreak), Short: true})
	fields = append(fields, Field{Title: "In the last year", Value: fmt.Sprintf("%d contributions", status.TotalCnt), Short: true})
	fields = append(fields, Field{Title: "In the last week", Value: fmt.Sprintf("%d contributions", status.WeekCnt), Short: true})

	attachments := []Attachment{}
	attachments = append(attachments, Attachment{Color: color, Text: text, Fields: fields})

	return s.httpPost(s.createValue(status.User, icon, attachments))
}

func (s *Slack) httpPost(req RequestValue) error {
	value, err := json.Marshal(req)

	resp, err := http.Post(s.WebhookUrl, "application/json", bytes.NewBuffer(value))

	if err != nil {
		return err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	defer resp.Body.Close()

	return nil
}

func (s *Slack) createValue(user, icon string, attachments []Attachment) RequestValue {
	return RequestValue{Channel: s.Channel, Username: user, IconEmoji: icon, Attachments: attachments}
}
