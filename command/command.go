package command

import (
	"fmt"
	"io"
	"strings"
)

type CommandFunc func(conn io.Writer, channel, sender string, args []string)

var Commands = map[string]CommandFunc{}

func Register(name string, handler CommandFunc) {
	Commands[strings.ToLower(name)] = handler
}

func Handle(conn io.Writer, channel, sender, content string) {
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
