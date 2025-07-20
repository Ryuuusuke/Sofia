package command

import (
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"
)

func init() {
	Register("chances", chanceCommand)
	Register("chance", chanceCommand)
}

func chanceCommand(conn io.Writer, channel, sender string, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(conn, "PRIVMSG %s: Usage: .chances <statement>\r\n", channel)
		return
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	percent := r.Intn(101)

	statement := strings.Join(args, " ")

	var response string
	switch {
	case percent >= 90:
		response = fmt.Sprintf("AMAZING! There's a \x0306%d%%\x0F chance %s", percent, statement)
	case percent >= 70:
		response = fmt.Sprintf("Pretty likely. I'd say \x0306%d%%\x0F chance %s", percent, statement)
	case percent >= 40:
		response = fmt.Sprintf("Hmm, maybe. About \x0306%d%%\x0F chance %s", percent, statement)
	case percent >= 20:
		response = fmt.Sprintf("Unlikely. Only \x0306%d%%\x0F chance %s", percent, statement)
	default:
		response = fmt.Sprintf("Nah, almost impossible. Just \x0306%d%%\x0F chance %s", percent, statement)
	}

	fmt.Fprintf(conn, "PRIVMSG %s :Let me think..\r\n", channel)
	time.Sleep(3 * time.Second)
	fmt.Fprintf(conn, "PRIVMSG %s :%s\r\n", channel, response)
}
