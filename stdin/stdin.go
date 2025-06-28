package stdin

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func HandleStdin(conn io.Writer) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		fmt.Fprintf(conn, "%s\r\n", line)
	}
	if err := scanner.Err(); err != nil {
		log.Printf("error when reading stdin : %v", err)
	}
}
