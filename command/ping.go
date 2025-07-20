package command

import (
	"fmt"
	"io"
	"time"
)

func init() {
	Register("ping", pingCommand)
}

func pingCommand(conn io.Writer, channel, sender string, args []string) {
	pingID := GeneratePingID()
	start := time.Now()

	RegisterPongCallback(pingID, func() {
		latency := time.Since(start)
		latencyMs := fmt.Sprintf("%.1f", latency.Seconds()*1000)
		fmt.Fprintf(conn, "PRIVMSG %s :\x0300,09 PONG! \x0F Latency is \x0306%sms\x0F\r\n", channel, latencyMs)
	})

	fmt.Fprintf(conn, "PING :%s\r\n", pingID)
}
