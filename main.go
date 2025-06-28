package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"sofia/rss"
	"sofia/stdin"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"gopkg.in/ini.v1"
)

type YtOembedResp struct {
	Title string `json:"title"`
}

func main() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("Can't load config: %v", err)
	}

	server := cfg.Section("irc").Key("server").String()
	channel := cfg.Section("irc").Key("channel").String()
	nickname := cfg.Section("irc").Key("nickname").String()
	username := cfg.Section("irc").Key("username").String()
	realname := cfg.Section("irc").Key("realname").String()

	useSasl, err := cfg.Section("sasl").Key("sasl").Bool()
	if err != nil {
		log.Fatalf("Can't not read useSasl %v", err)
	}

	conn, err := tls.Dial("tcp", server, &tls.Config{
		ServerName: "irc.libera.chat",
	})
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	joined := false

	fmt.Fprintf(conn, "CAP REQ :sasl\r\n")

	go stdin.HandleStdin(conn)

	// core
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Koneksi ditutup oleh server")
				fmt.Println(err)
				break
			}
			log.Fatalf("Error dari server: %s", err)
		}
		line = strings.TrimSpace(line)
		fmt.Println(">>", line)

		if strings.HasPrefix(line, "PING") {
			pongMsg := strings.Replace(line, "PING", "PONG", 1)
			fmt.Fprintf(conn, "%s\r\n", pongMsg)
			continue
		}

		if useSasl {

			saslUser := cfg.Section("sasl").Key("user").String()
			saslPass := cfg.Section("sasl").Key("password").String()

			if strings.Contains(line, "CAP") && strings.Contains(line, "ACK") {
				fmt.Fprintf(conn, "AUTHENTICATE PLAIN\r\n")
				continue
			}

			if line == "AUTHENTICATE +" {
				payload := fmt.Sprintf("\x00%s\x00%s", saslUser, saslPass)
				encoded := base64.StdEncoding.EncodeToString([]byte(payload))
				fmt.Fprintf(conn, "AUTHENTICATE %s\r\n", encoded)
				continue
			}

			if strings.Contains(line, " 904 ") {
				fmt.Println("SASL gagal")
				fmt.Fprintf(conn, "CAP END\r\n")
				break
			}

		}
		if strings.Contains(line, " 903 ") {
			fmt.Println("Login sucess")
			fmt.Fprintf(conn, "CAP END\r\n")
			fmt.Fprintf(conn, "NICK %s\r\n", nickname)
			fmt.Fprintf(conn, "USER %s 8 * :%s\r\n", username, realname)
			continue
		}

		if !joined && strings.Contains(line, " 001 ") {
			fmt.Fprintf(conn, "JOIN %s\r\n", channel)
			joined = true

			go rss.StartRSSLoop(conn, channel)

			continue
		}

		if strings.HasSuffix(line, "JOIN "+channel) || strings.HasSuffix(line, "JOIN :"+channel) {
			prefix := strings.SplitN(line, "!", 2)
			if len(prefix) > 0 && strings.HasPrefix(prefix[0], ":") {
				nick := strings.TrimPrefix(prefix[0], ":")

				if nick != nickname {
					message := fmt.Sprintf("PRIVMSG %s :Selamat datang %s!\r\n", channel, nick)
					fmt.Fprint(conn, message)
				}
			}
		}

		if strings.Contains(line, "PRIVMSG "+channel) {
			message := line
			if strings.Contains(message, "http") {
				sender := getNickname(line)
				urls := getUrl(message)
				for _, url := range urls {
					title, err := getUrlTitle(url)
					if err != nil {
						fmt.Printf("Gagal fetch title dari %s: %s\n", url, err)
						continue
					}
					reply := fmt.Sprintf("PRIVMSG %s :[\x0309Title\x0F] %s (sent by %s)\r\n", channel, title, sender)
					fmt.Fprint(conn, reply)
				}

			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func getNickname(line string) string {
	prefix := strings.SplitN(line, "!", 2)
	return strings.TrimPrefix(prefix[0], ":")
}

func getUrl(message string) []string {
	urlRegex := regexp.MustCompile(`https?://[^\s]+`)
	urls := urlRegex.FindAllString(message, -1)
	return urls
}

func getUrlTitle(url string) (string, error) {
	if isYt(url) {
		oembedURL := "https://www.youtube.com/oembed?url=" + url

		resp, err := http.Get(oembedURL)
		if err != nil {
			return "", fmt.Errorf("gagal mengambil data dari oEmbed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("status code bukan 200 dari oEmbed: %d", resp.StatusCode)
		}

		var oembedResp YtOembedResp
		err = json.NewDecoder(resp.Body).Decode(&oembedResp)
		if err != nil {
			return "", fmt.Errorf("gagal decode JSON dari oEmbed: %w ", err)
		}

		return oembedResp.Title, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	ctx, cancelBrowser := chromedp.NewContext(ctx)
	defer cancelBrowser()

	var title string
	err := chromedp.Run(ctx, chromedp.Navigate(url), chromedp.Title(&title))
	if err != nil {
		return "", err
	}
	return title, nil
}

func isYt(url string) bool {
	lower := strings.ToLower(url)
	return strings.Contains(lower, "youtube.com/watch") || strings.Contains(lower, "youtu.be/")
}
