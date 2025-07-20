package command

import (
	"fmt"
	"io"
	"math/rand"
	"strings"
)

type CommandFunc func(conn io.Writer, channel, sender string, args []string)

var Commands = map[string]CommandFunc{}

func Register(name string, handler CommandFunc) {
	Commands[strings.ToLower(name)] = handler
}

func Handle(conn io.Writer, channel, sender, content string) {
	if !strings.HasPrefix(content, ",") {
		return
	}

	fields := strings.Fields(content)
	if len(fields) == 0 {
		return
	}

	cmd := strings.ToLower(strings.TrimPrefix(fields[0], ","))
	args := fields[1:]

	if handler, ok := Commands[cmd]; ok {
		handler(conn, channel, sender, args)
	} else {
		fmt.Fprintf(conn, "PRIVMSG %s :Unknown command: %s\r\n", channel, cmd)
	}
}

// ping stuff
var pongCallbacks = map[string]func(){}

func RegisterPongCallback(id string, fn func()) {
	pongCallbacks[id] = fn
}

func HandlePong(line string) {
	if !strings.HasPrefix(line, ":") || !strings.Contains(line, "PONG") {
		return
	}
	parts := strings.Split(line, ":")
	if len(parts) < 3 {
		return
	}
	id := strings.TrimSpace(parts[2])
	if cb, ok := pongCallbacks[id]; ok {
		go cb()
		delete(pongCallbacks, id)
	}
}

func GeneratePingID() string {
	return fmt.Sprintf("latency-%d", rand.Int())
}
