package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

var (
	api *anaconda.TwitterApi

	port = os.Getenv("PORT")

	tweetDateLayout = "Mon Jan 2 15:04:05 -0700 2006"

	accessToken       = os.Getenv("ACCESS_TOKEN")
	accessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")
	consumerKey       = os.Getenv("CONSUMER_KEY")
	consumerSecret    = os.Getenv("CONSUMER_SECRET")
)

func main() {
	if accessToken == "" || accessTokenSecret == "" || consumerKey == "" || consumerSecret == "" {
		panic("Missing Twitter credentials")
	}

	api = anaconda.NewTwitterApiWithCredentials(accessToken, accessTokenSecret, consumerKey, consumerSecret)

	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", getTweet)
	fmt.Println("Server started at port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getTweet(w http.ResponseWriter, r *http.Request) {
	id := strings.ReplaceAll(r.URL.Path, "/", "")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	tweet, err := api.GetTweet(i, nil)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	templateFuncs := template.FuncMap{
		"base64": func(url string) string {
			res, err := http.Get(url)
			if err != nil {
				return ""
			}
			buf := new(bytes.Buffer)
			buf.ReadFrom(res.Body)
			return base64.StdEncoding.EncodeToString(buf.Bytes())
		},
		"isoDate": func(date string) string {
			t, _ := time.Parse(tweetDateLayout, date)
			return t.Format(time.RFC3339)
		},
		"humanDate": func(date string) string {
			t, _ := time.Parse(tweetDateLayout, date)
			return t.Format("3:04 PM · Jan 2, 2006")
		},
	}

	t := template.Must(
		template.New("tweet.svg.tmpl").
			Funcs(templateFuncs).
			ParseFiles("tweet.svg.tmpl"))

	w.Header().Set("Content-type", "image/svg+xml")
	err = t.Execute(w, tweet)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
}
