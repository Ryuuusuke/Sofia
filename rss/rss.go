package rss

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
	Title string `xml:"title"`
}

type Item struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

type RDFChannel struct {
	Title string `xml:"title"`
}

type RDF struct {
	Items   []Item     `xml:"item"`
	Channel RDFChannel `xml:"channel"`
}

func FetchRSS(feedURL string) ([]Item, string, error) {
	resp, err := http.Get(feedURL)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	var rss RSS
	if err := xml.NewDecoder(resp.Body).Decode(&rss); err == nil && len(rss.Channel.Items) > 0 {
		return rss.Channel.Items, rss.Channel.Title, nil
	}

	resp.Body.Close()
	resp, err = http.Get(feedURL)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	var rdf RDF
	if err := xml.NewDecoder(resp.Body).Decode(&rdf); err == nil && len(rdf.Items) > 0 {
		source := rdf.Channel.Title
		return rdf.Items, source, nil
	}

	return nil, "", fmt.Errorf("unable to parse items from feed: %s", feedURL)
}

func removeQueryParams(rawUrl string) string {
	parsed, err := url.Parse(rawUrl)
	if err != nil {
		return rawUrl
	}

	parsed.RawQuery = ""
	return parsed.String()
}

func StartRSSLoop(cfg *ini.File, conn io.Writer, channel string) {
	rssSec := cfg.Section("rss")
	rawUrls := rssSec.Key("urls").Strings(",")

	if len(rawUrls) == 0 {
		log.Println("> RSS: No URL feed has been set")
		return
	}

	interval := rssSec.Key("interval").MustInt(10)

	lastSeen := make(map[string]string)

	pendingItems := make(map[string][]Item)

	go func() {
		for {
			for _, feedUrl := range rawUrls {
				items, source, err := FetchRSS(feedUrl)
				if err != nil {
					log.Printf("> RSS: Failed to fetch from %s: %v\n", feedUrl, err)
					continue
				}

				if len(items) == 0 {
					continue
				}

				newItems := []Item{}
				for _, item := range items {
					if item.Link == lastSeen[feedUrl] {
						break
					}
					newItems = append(newItems, item)
				}

				if len(newItems) > 0 {
					lastSeen[feedUrl] = items[0].Link
					pendingItems[source] = append(newItems, pendingItems[source]...)
				}
			}
			time.Sleep(time.Duration(interval) * time.Minute)
		}
	}()

	go func() {
		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			log.Printf("Failed to load timezone: %v", err)
			loc = time.Local // fallback
		}

		for {
			now := time.Now().In(loc)
			next10AM := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location())
			if now.After(next10AM) {
				next10AM = next10AM.Add(24 * time.Hour)
			}
			sleepDuration := time.Until(next10AM)
			log.Printf("> RSS: Sleeping until %v to send feed\n", next10AM)
			time.Sleep(sleepDuration)

			sentSources := make(map[string]bool)

			fmt.Fprintf(conn, "PRIVMSG %s : 🍡 Anime", channel)
			for source, items := range pendingItems {
				for i := len(items) - 1; i >= 0; i-- {
					item := items[i]
					cleanURL := removeQueryParams(item.Link)
					if strings.Contains(source, "Anime") {
						fmt.Fprintf(conn, "PRIVMSG %s : [\x0311%s\x0F] \x0300,02%s\x0F - %s\r\n", channel, source, item.Title, cleanURL)
						sentSources[source] = true
					}
				}
			}

			fmt.Fprintf(conn, "PRIVMSG %s : 🌐 World News", channel)
			for source, items := range pendingItems {
				for i := len(items) - 1; i >= 0; i-- {
					item := items[i]
					cleanURL := removeQueryParams(item.Link)
					if strings.Contains(source, "Stories") {
						fmt.Fprintf(conn, "PRIVMSG %s : [\x0311%s\x0F] \x0300,02%s\x0F - %s\r\n", channel, source, item.Title, cleanURL)
						sentSources[source] = true
					}
				}
			}

			for source := range sentSources {
				pendingItems[source] = nil
			}
		}
	}()
}
