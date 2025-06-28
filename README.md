# Sofia IRC Bot

**Sofia** is a lightweight IRC bot written in Go. It supports:

* SASL authentication
* Link title fetching (including YouTube)
* RSS looping (modular)
* Sending messages directly from terminal (stdin)
* Configuration via `.ini` file

<br>

## 🔧 Configuration

Create a file named `config.ini` in the root directory:

```ini
[sasl]
sasl = true
user = your-sasl-username
password = your-sasl-password

[irc]
server = irc.libera.chat:6697
nickname = sofiaaa
username = SofiaPertama
realname = Ratu Sofia
channel = `##sofia`
```

> Make sure to wrap the `channel` value in backticks (`) so that the `#\` character is not treated as a comment by the parser.

<br>

## 🚀 Running the Bot

Run directly with:

```bash
go run main.go
```

Or build it first:

```bash
go build -o sofia .
./sofia
```

<br>

## 🖥️ Sending Messages from Terminal

Type directly into the terminal where the bot is running to send a message to the configured IRC channel.

<br>

## 🌐 Link Preview Feature

* When someone posts a link in the channel, the bot will try to fetch the **page title** automatically.
* YouTube links are fetched using the [YouTube oEmbed API](https://www.youtube.com/oembed).
* For regular webpages, the bot uses `chromedp` (headless Chrome via Go) to extract the `<title>`.

<br>

## 📰 RSS Feed

You can add your own RSS modules under the `rss/` folder.
The bot is modular and supports custom RSS loops tailored to your needs.

<br>

## 🧱 Directory Structure

```
.
├── main.go          # Entry point
├── config.ini       # User-provided config file
├── stdin/           # Module for stdin input
├── rss/             # RSS handler module
└── go.mod           # Go module file
```

<br>

## 📦 Dependencies

* [`go-ini/ini`](https://github.com/go-ini/ini) - For parsing `.ini` configuration files
* [`chromedp`](https://github.com/chromedp/chromedp) - For fetching web page titles via headless Chrome
* Go standard libraries (`net`, `bufio`, `tls`, `regexp`, etc.)

<br>

## 📄 License

This project is licensed under the **BSD 3-Clause License**.
See the [`LICENSE`](./LICENSE) file for full details.
