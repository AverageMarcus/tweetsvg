package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
)

var (
	api *anaconda.TwitterApi

	tweetDateLayout = "Mon Jan 2 15:04:05 -0700 2006"

	port              string
	accessToken       string
	accessTokenSecret string
	consumerKey       string
	consumerSecret    string
)

func init() {
	godotenv.Load(os.Getenv("DOTENV_DIR") + ".env")

	port = os.Getenv("PORT")
	accessToken = os.Getenv("ACCESS_TOKEN")
	accessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")
	consumerKey = os.Getenv("CONSUMER_KEY")
	consumerSecret = os.Getenv("CONSUMER_SECRET")
}

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

	re := regexp.MustCompile(`[\x{1F300}-\x{1F6FF}]`)
	emojis := re.FindAllString(tweet.FullText, -1)

	emojiCount := 0
	for _, emoji := range emojis {
		emojiCount += len([]byte(emoji)) - 1
	}
	tweet.FullText = tweet.FullText[tweet.DisplayTextRange[0] : tweet.DisplayTextRange[1]+emojiCount]

	for _, user := range tweet.Entities.User_mentions {
		tweet.FullText = strings.ReplaceAll(tweet.FullText, "@"+user.Screen_name, fmt.Sprintf("<a rel=\"noopener\" target=\"_blank\" href=\"https://twitter.com/%s/\">@%s</a>", user.Screen_name, user.Screen_name))
	}
	for _, url := range tweet.Entities.Urls {
		tweet.FullText = strings.ReplaceAll(tweet.FullText, url.Url, fmt.Sprintf("<a rel=\"noopener\" target=\"_blank\" href=\"https://twitter.com/%s/\">%s</a>", url.Expanded_url, url.Display_url))
	}
	for _, hashtag := range tweet.Entities.Hashtags {
		tweet.FullText = strings.ReplaceAll(tweet.FullText, "#"+hashtag.Text, fmt.Sprintf("<a rel=\"noopener\" target=\"_blank\" href=\"https://twitter.com/hashtag/%s\">#%s</a>", hashtag.Text, hashtag.Text))
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
		"html": func(in string) template.HTML {
			return template.HTML(in)
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
