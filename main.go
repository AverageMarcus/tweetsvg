package main

import (
	"bytes"
	"embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/joho/godotenv"
	"github.com/rivo/uniseg"
	emoji "github.com/tmdvs/Go-Emoji-Utils"
)

//go:embed index.html tweet.svg.tmpl suspendedTweet.svg

var content embed.FS

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) > 1 {
			getTweet(w, r)
		} else {
			body, _ := content.ReadFile("index.html")
			w.Write(body)
		}
	})
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
		switch err := err.(type) {
		case *anaconda.ApiError:
			switch err.Decoded.Errors[0].Code {
			case 63:
				fmt.Printf("Generating suspended tweet image for %s\n", id)
				suspendedTweet(w)
				return
			}
		}
		fmt.Println(err)
		w.WriteHeader(404)
		return
	}

	gr := uniseg.NewGraphemes(tweet.FullText)
	count := 0
	displayText := ""
	for gr.Next() {
		if count >= tweet.DisplayTextRange[0] && count < tweet.DisplayTextRange[1] {
			displayText += gr.Str()
		}
		count += 1
	}
	tweet.FullText = displayText

	for _, user := range tweet.Entities.User_mentions {
		tweet.FullText = strings.ReplaceAll(tweet.FullText, "@"+user.Screen_name, fmt.Sprintf("<a rel=\"noopener\" target=\"_blank\" href=\"https://twitter.com/%s/\">@%s</a>", user.Screen_name, user.Screen_name))
	}
	for _, url := range tweet.Entities.Urls {
		tweet.FullText = strings.ReplaceAll(tweet.FullText, url.Url, fmt.Sprintf("<a rel=\"noopener\" target=\"_blank\" href=\"https://twitter.com/%s/\">%s</a>", url.Expanded_url, url.Display_url))
	}
	for _, hashtag := range tweet.Entities.Hashtags {
		tweet.FullText = strings.ReplaceAll(tweet.FullText, "#"+hashtag.Text, fmt.Sprintf("<a rel=\"noopener\" target=\"_blank\" href=\"https://twitter.com/hashtag/%s\">#%s</a>", hashtag.Text, hashtag.Text))
	}

	tweet.FullText = strings.ReplaceAll(tweet.FullText, "\n", "<br />")

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
		"calculateHeight": func(tweet anaconda.Tweet) string {
			height := 64.0 /* Avatar */ + 20 /* footer */ + 46 /* text margin */ + 22 /* margin */

			lineWidth := 0.0
			lineHeight := 28.0
			tweetText := strings.ReplaceAll(tweet.FullText, "<br />", " \n")
			tweetText = strip.StripTags(tweetText)
			words := regexp.MustCompile(`[ |-]`).Split(tweetText, -1)
			for _, word := range words {
				if len(emoji.FindAll(word)) > 0 {
					lineHeight = 32.0
				}

				if strings.Contains(word, "\n") {
					height += lineHeight
					lineHeight = 28.0
					lineWidth = 0
					continue
				}

				chars := strings.Split(word, "")
				wordWidth := 0.0
				for _, char := range chars {
					wordWidth += getCharWidth(char)
				}

				if wordWidth > 435 {
					height += (lineHeight * (math.Ceil(wordWidth/435) + 1))
					lineHeight = 28.0
					lineWidth = 0
				} else if lineWidth+getCharWidth(" ")+wordWidth > 435 {
					height += lineHeight
					lineHeight = 28.0
					lineWidth = wordWidth
				} else {
					lineWidth += wordWidth
				}
			}
			if lineWidth > 0 {
				height += lineHeight
			}

			if tweet.InReplyToScreenName != "" {
				height += 42
			}

			for i, img := range tweet.ExtendedEntities.Media {
				ratio := float64(img.Sizes.Small.W) / 468
				tweet.ExtendedEntities.Media[i].Sizes.Small.W = 468
				tweet.ExtendedEntities.Media[i].Sizes.Small.H = int((float64(img.Sizes.Small.H) / ratio) + 5.0)
			}

			if len(tweet.ExtendedEntities.Media) > 1 {
				for i, img := range tweet.ExtendedEntities.Media {
					tweet.ExtendedEntities.Media[i].Sizes.Small.W = (img.Sizes.Small.W / 2) - 20
					tweet.ExtendedEntities.Media[i].Sizes.Small.H = (img.Sizes.Small.H / 2) - 20
				}
			}

			switch len(tweet.ExtendedEntities.Media) {
			case 1:
				height += float64(tweet.ExtendedEntities.Media[0].Sizes.Small.H)
			case 2:
				height += math.Max(float64(tweet.ExtendedEntities.Media[0].Sizes.Small.H), float64(tweet.ExtendedEntities.Media[1].Sizes.Small.H)) + 5
			case 3:
				height += math.Max(float64(tweet.ExtendedEntities.Media[0].Sizes.Small.H), float64(tweet.ExtendedEntities.Media[1].Sizes.Small.H)) + 5
				height += float64(tweet.ExtendedEntities.Media[2].Sizes.Small.H) + 35
			case 4:
				height += math.Max(float64(tweet.ExtendedEntities.Media[0].Sizes.Small.H), float64(tweet.ExtendedEntities.Media[1].Sizes.Small.H)) + 10
				height += math.Max(float64(tweet.ExtendedEntities.Media[2].Sizes.Small.H), float64(tweet.ExtendedEntities.Media[3].Sizes.Small.H)) + 10
				height += 7
			}

			return fmt.Sprintf("%dpx", int64(height))
		},
	}

	t := template.Must(
		template.New("tweet.svg.tmpl").
			Funcs(templateFuncs).
			ParseFS(content, "tweet.svg.tmpl"))

	w.Header().Set("Content-type", "image/svg+xml")
	err = t.Execute(w, tweet)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
}

func suspendedTweet(w http.ResponseWriter) {
	w.Header().Set("Content-type", "image/svg+xml")
	tweet, _ := content.ReadFile("suspendedTweet.svg")
	w.Write(tweet)
}
