package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/net/html"
)

const (
	server   = "irc.libera.chat:6697"
	nickname = "Sofiaaa"
	username = "SofiaPertama"
	channel  = "##sofia"
	realname = "Ratu Sofia"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Cant load .env file")
	}

	saslUser := os.Getenv("SASL_USER")
	saslPass := os.Getenv("SASL_PASS")

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

		if strings.Contains(line, " 903 ") {
			fmt.Println("SASL sucess")
			fmt.Fprintf(conn, "CAP END\r\n")
			fmt.Fprintf(conn, "NICK %s\r\n", nickname)
			fmt.Fprintf(conn, "USER %s 8 * :%s\r\n", username, realname)
			continue
		}

		if strings.Contains(line, " 904 ") {
			fmt.Println("SASL gagal")
			fmt.Fprintf(conn, "CAP END\r\n")
			break
		}

		if !joined && strings.Contains(line, " 001 ") {
			fmt.Fprintf(conn, "JOIN %s\r\n", channel)
			joined = true
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
					reply := fmt.Sprintf("PRIVMSG %s :[Title] %s (sent by %s)\r\n", channel, title, sender)
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
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)
	for {
		title := z.Next()
		switch title {

		case html.ErrorToken:
			return "", fmt.Errorf("%s", "Title not found")

		case html.StartTagToken:
			t := z.Token()
			if t.Data == "title" {
				z.Next()
				return strings.TrimSpace(z.Token().Data), nil
			}

		}

	}
}
