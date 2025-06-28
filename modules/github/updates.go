package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gopkg.in/ini.v1"
)

type Commit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Message   string `json:"message"`
		Committer struct {
			Name string    `json:"name"`
			Date time.Time `json:"date"`
		} `json:"committer"`
	} `json:"commit"`
	HTMLURL string `json:"html_url"`
}

func FetchLatestCommit(owner, repo string) (*Commit, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", owner, repo)
	req, _ := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	req.Header.Set("User-Agent", "sofia-bot")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var commits []Commit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return nil, err
	}

	if len(commits) == 0 {
		return nil, fmt.Errorf("no commits found")
	}

	return &commits[0], nil
}

func StartGithubWatcher(cfg *ini.File, conn io.Writer, channel string) {
	githubSec := cfg.Section("github")
	if !githubSec.Key("enabled").MustBool(false) {
		return
	}

	owner := githubSec.Key("owner").String()
	repo := githubSec.Key("repo").String()
	interval := githubSec.Key("interval").MustInt(5)

	var lastSHA string

	for {
		commit, err := FetchLatestCommit(owner, repo)
		if err != nil {
			fmt.Println("failed to fetch commit: ", err)
			time.Sleep(time.Duration(interval) * time.Minute)
			continue
		}

		if commit.SHA != lastSHA {
			fmt.Fprintf(conn, "PRIVMSG %s :[\x0309GitHub\x0F] New commit in %s by %s: %s (%s)\r\n", channel, repo, commit.Commit.Committer.Name, commit.Commit.Message, commit.HTMLURL)
			lastSHA = commit.SHA
		}

		time.Sleep(time.Duration(interval) * time.Minute)
	}
}
