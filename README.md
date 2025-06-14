# Sofiaaa IRC Bot

**Sofiaaa** is a simple IRC bot written in Go that can:

- Send a welcome message when someone joins the channel.
- Detect and fetch webpage titles from URLs mentioned in messages.
- Authenticate using SASL to NickServ on the Libera.Chat network.

## Features

- **SASL authentication** to NickServ (required by Libera.Chat).
- **Automatic greetings** to users who join the channel.
- **URL title fetcher** – detects links in messages and posts their `<title>`.

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/yourusername/sofiaaa-bot.git
cd sofiaaa-bot
````

### 2. Create a `.env` file

Create a `.env` file in the project root with your NickServ credentials:

```env
SASL_USER=Sofiaaa
SASL_PASS=your_nickserv_password_here
```

> ⚠️ **Warning**: Do **not** commit `.env` to version control!

### 3. Run the bot

```bash
go run main.go
```

Make sure your system allows outbound TLS connections to `irc.libera.chat:6697`.

## Dependencies

* [`golang.org/x/net/html`](https://pkg.go.dev/golang.org/x/net/html) – for parsing HTML titles.
* [`github.com/joho/godotenv`](https://github.com/joho/godotenv) – to load `.env` variables.

Install with:

```bash
go get github.com/joho/godotenv
go get golang.org/x/net/html
```

## File Structure

```
.
├── main.go        # Main bot logic
├── .env           # Secret credentials (should be in .gitignore)
└── README.md      # This documentation
```

## Planned Features

* Logging to a file
* Basic IRC commands (e.g., `!ping`, `!help`)
* Spam or keyword filtering

## License

This project is licensed under the BSD 3-Clause License.  
See the [LICENSE](./LICENSE) file for details.
