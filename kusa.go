package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Kusa struct {
}

type KusaStatus struct {
	User                                        string
	Count, Streak, MaxStreak, TotalCnt, WeekCnt int
}

func (k *Kusa) Fetch(users []string) ([]KusaStatus, error) {
	var statuses []KusaStatus

	for _, user := range users {
		url := k.githubUrl(user)
		doc, err := goquery.NewDocument(url)

		if err != nil {
			return nil, err
		}

		count, streak, maxStreak, totalCnt, weekCnt := k.contributionsOn(doc)
		statuses = append(statuses, KusaStatus{
			User:      user,
			Count:     count,
			Streak:    streak,
			MaxStreak: maxStreak,
			TotalCnt:  totalCnt,
			WeekCnt:   weekCnt,
		})
	}

	return statuses, nil
}

func (k *Kusa) githubUrl(user string) string {
	return fmt.Sprintf("https://github.com/%s", user)
}

func (k *Kusa) contributionsOn(doc *goquery.Document) (int, int, int, int, int) {
	cnt := 0
	streak := 0
	maxStreak := 0
	week := []int{}
	totalCnt := 0
	now := time.Now()

	doc.Find(".js-calendar-graph-svg .day").Each(func(i int, s *goquery.Selection) {
		dateStr, _ := s.Attr("data-date")
		date, _ := time.ParseInLocation("2006-01-02", dateStr, now.Location())

		if !date.After(now) {
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
		}
	})
	weekCnt := 0
	for _, c := range week {
		weekCnt += c
	}
	return cnt, streak, maxStreak, totalCnt, weekCnt
}
