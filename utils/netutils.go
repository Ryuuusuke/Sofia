package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type YTOembedResp struct {
	Title string `json:"title"`
}

func GetNickname(line string) string {
	prefix := strings.SplitN(line, "!", 2)
	return strings.TrimPrefix(prefix[0], ":")
}

func GetUrl(message string) []string {
	urlRegex := regexp.MustCompile(`<?(https?://[^<>\s]+)>?`)
	matches := urlRegex.FindAllStringSubmatch(message, -1)

	urls := make([]string, 0) // gg ga tuh deklarasi nya wkwokwokokok
	for _, match := range matches {
		if len(match) > 1 {
			urls = append(urls, match[1])
		}
	}
	return urls
}

func IsYoutube(url string) bool {
	lower := strings.ToLower(url)
	return strings.Contains(lower, "youtube.com/watch") || strings.Contains(lower, "youtu.be/")
}

func GetUrlTitle(url string) (string, error) {
	if IsYoutube(url) {
		oembedUrl := "https://www.youtube.com/oembed?url=" + url
		resp, err := http.Get(oembedUrl)
		if err != nil {
			return "", fmt.Errorf("failed to fetch Youtube oEmbed: %w", err)
		}
		defer resp.Body.Close()

		var oembedResp YTOembedResp
		if err := json.NewDecoder(resp.Body).Decode(&oembedResp); err != nil {
			return "", fmt.Errorf("failed to parse oEmbed JSON: %w", err)
		}
		return oembedResp.Title, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ctx, cancelBrowser := chromedp.NewContext(ctx)
	defer cancelBrowser()

	var title string
	if err := chromedp.Run(ctx, chromedp.Navigate(url), chromedp.Title(&title)); err != nil {
		return "", err
	}
	return title, nil
}
