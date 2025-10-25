package main

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"sofia/command"
	"sofia/modules/github"
	"sofia/rss"
	"sofia/stdin"
	"sofia/utils"
	"strings"
	"time"

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
		command.HandlePong(line)
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

			go rss.StartRSSLoop(cfg, conn, channel)
			go github.StartGithubWatcher(cfg, conn, channel)

			continue
		}

		if strings.HasSuffix(line, "JOIN "+channel) || strings.HasSuffix(line, "JOIN :"+channel) {
			prefix := strings.SplitN(line, "!", 2)
			if len(prefix) > 0 && strings.HasPrefix(prefix[0], ":") {
				nick := strings.TrimPrefix(prefix[0], ":")

				if nick != nickname {
					message := fmt.Sprintf("NOTICE %s :Selamat datang, %s-sama! Mohon tunggu sebentar, Ramu-neechan akan memberikanmu akses voice agar kamu bisa ikut berbicara disini\r\n", channel, nick)
					fmt.Fprint(conn, message)
				}
			}
		}

		if strings.Contains(line, "PRIVMSG "+channel) {
			message := line
			sender := utils.GetNickname(line)
			parts := strings.SplitN(line, " :", 2)
			if len(parts) >= 2 {
				content := parts[1]
				command.Handle(conn, channel, sender, content)
			}
			if strings.Contains(message, "(re") {
				continue
			}
			if strings.Contains(message, "http") {
				sender := utils.GetNickname(line)
				urls := utils.GetUrl(message)
				for _, url := range urls {
					title, err := utils.GetUrlTitle(url)
					if err != nil {
						fmt.Printf("failed to fetch web title from %s: %s\n", url, err)
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
