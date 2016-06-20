package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Kusa struct {
	ContributionCount int `json:"contribution_count"`
}

func kusaHandler(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	body := Kusa{}

	defer func() {
		outjson, err := json.Marshal(body)
		if err != nil {
			// TODO: Output to log
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "applcation/json")
		w.WriteHeader(status)
		fmt.Fprint(w, string(outjson))
	}()

	if r.Method != "GET" {
		status = http.StatusMethodNotAllowed
		return
	}

	username := r.URL.Path[len("/siba/"):]

	doc, err := goquery.NewDocument(githubUrl(username))
	if err != nil {
		// TODO: Output to log
		fmt.Println(err)
	}

	count := contributionsOn(doc, today())

	body = Kusa{ContributionCount: count}
}

func main() {
	http.HandleFunc("/kusa/", kusaHandler)
	http.ListenAndServe(getenvOrDefault("PORT", ":8080"), nil)
}

func getenvOrDefault(key string, defvalue string) string {
	value, exists := os.LookupEnv(key)
	if (exists) {
		return value
	} else {
		return defvalue
	}
}

func githubUrl(user string) string {
	return fmt.Sprintf("https://github.com/%s", user)
}

func contributionsOn(doc *goquery.Document, date string) int {
	str, _ := doc.Find(query(date)).Attr("data-count")
	cnt, _ := strconv.Atoi(str)
	return cnt
}

func query(date string) string {
	return fmt.Sprintf(".js-calendar-graph-svg .day[data-date='%s']", date)
}

func today() string {
	return time.Now().Format("2006-01-02")
}
