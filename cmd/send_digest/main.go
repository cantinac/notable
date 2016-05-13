package main

import (
	notable "github.com/cantinac/notable"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(notable.Notes()) > 0 {
		notable.SendEmail(
			os.Getenv("SPARKPOST_API_KEY"),
			os.Getenv("TO_EMAIL"),
			os.Getenv("FROM_EMAIL"),
		)
		if os.Getenv("NO_RESET") == "" {
			notable.Reset()
		} else {
			log.Print("Not resetting notes")
		}
	}
}
