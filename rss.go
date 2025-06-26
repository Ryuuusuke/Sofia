package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

const rssURL = "https://myanimelist.net/rss/news.xml"

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"PubDate"`
}

func FetchRSS() ([]Item, error) {
	resp, err := http.Get(rssURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rss RSS
	err = xml.NewDecoder(resp.Body).Decode(&rss)
	if err != nil {
		return nil, err
	}

	return rss.Channel.Items, nil
}

func removeQueryParams(rawUrl string) string {
	parsed, err := url.Parse(rawUrl)
	if err != nil {
		return rawUrl
	}

	parsed.RawQuery = ""
	return parsed.String()
}

func StartRSSLoop(conn io.Writer, channel string) {
	var lastSent string

	for {
		items, err := FetchRSS()
		if err != nil {
			log.Println("Gagal mengambil RSS: ", err)
			time.Sleep(10 * time.Minute)
			continue
		}

		if len(items) > 0 && items[0].Link != lastSent {
			lastSent = items[0].Link
			cleanLink := removeQueryParams(items[0].Link)
			msg := fmt.Sprintf("PRIVMSG %s : [\x0302MAL News\x0F] %s - %s\r\n", channel, items[0].Title, cleanLink)
			fmt.Fprint(conn, msg)
		}

		time.Sleep(10 * time.Minute)
	}
}
